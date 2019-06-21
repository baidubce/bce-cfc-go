package cfc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

type Credential struct {
	AccessKeyId     string
	AccessKeySecret string
	SessionToken    string
}

type InvokeRequest struct {
	RequestID       string
	AccessKeyID     string
	AccessKeySecret string
	SecurityToken   string
	FunctionBrn     string
	FunctionTimeout int
	ClientContext   string `json:",omitempty"`
	EventObject     string `json:",omitempty"`
}

type InvokeResponse struct {
	RequestID  string
	Success    bool
	FuncResult string `json:"result,omitempty"`
	FuncError  string `json:"error,omitempty"`
}

type InvokeHandler interface {
	Handle(input io.Reader, output io.Writer, context InvokeContext) error
}

var (
	handlerMap = map[string]InvokeHandler{}
)

const (
	DefaultHandlerName = "_default_"
)

func RegisterNamedHandler(name string, h InvokeHandler) {
	handlerMap[name] = h
}

func RegisterDefaultHandler(h InvokeHandler) {
	handlerMap[DefaultHandlerName] = h
}

func GetInvokeHandler(name string) InvokeHandler {
	if h, ok := handlerMap[name]; ok {
		return h
	}
	return handlerMap[DefaultHandlerName]
}

type CfcClient struct {
	config *RuntimeConfig
	reader io.ReadCloser
	bufrd  *bufio.Reader
	writer io.WriteCloser
}

const (
	MAX_INVOKE_EVENT_LENGTH = 128 * 1024
)

func NewCfcClient(c *RuntimeConfig, maxEventSize int) *CfcClient {
	invokePipe := os.NewFile(uintptr(c.InvokePipe()), "invoke-pipe")
	responsePipe := invokePipe
	if c.InvokePipe() != c.ResponsePipe() {
		responsePipe = os.NewFile(uintptr(c.ResponsePipe()), "response-pipe")
	}
	if maxEventSize < MAX_INVOKE_EVENT_LENGTH {
		maxEventSize = MAX_INVOKE_EVENT_LENGTH
	}
	return &CfcClient{
		config: c,
		reader: invokePipe,
		bufrd:  bufio.NewReaderSize(invokePipe, maxEventSize),
		writer: responsePipe,
	}
}

func (iv *CfcClient) recvRequest() (*InvokeRequest, error) {
	data, prefix, err := iv.bufrd.ReadLine()
	if prefix {
		return nil, fmt.Errorf("invoke message is too long")
	}
	if err != nil {
		return nil, err
	}
	request := &InvokeRequest{}
	err = json.Unmarshal(data, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func (iv *CfcClient) invokeFunction(request *InvokeRequest) (response *InvokeResponse) {
	credential := &Credential{
		AccessKeyId:     request.AccessKeyID,
		AccessKeySecret: request.AccessKeySecret,
		SessionToken:    request.SecurityToken,
	}
	invokeContext := NewInvokeContext(request.RequestID,
		iv.config, credential, []byte(request.ClientContext))
	input := bytes.NewBuffer([]byte(request.EventObject))
	output := new(bytes.Buffer)
	handler := GetInvokeHandler(iv.config.Handler())
	response = &InvokeResponse{
		RequestID: request.RequestID,
		Success:   true,
	}

	defer func() {
		if r := recover(); r != nil {
			response.Success = false
			response.FuncError = string(debug.Stack())
		}
	}()

	err := handler.Handle(input, output, invokeContext)
	if err != nil {
		response.Success = false
		response.FuncError = err.Error()
	} else {
		response.FuncResult = output.String()
	}
	return response
}

func (iv *CfcClient) sendResponse(resp *InvokeResponse) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	// 增加请求日志结束标记
	os.Stdout.Write([]byte("\000"))
	os.Stderr.Write([]byte("\000"))
	data = append(data, '\n')
	_, err = iv.writer.Write(data)
	return err
}

func (iv *CfcClient) WaitInvoke() error {
	for {
		request, err := iv.recvRequest()
		if err != nil {
			return err
		}

		response := iv.invokeFunction(request)
		err = iv.sendResponse(response)
		if err != nil {
			return err
		}
	}
}

func (iv *CfcClient) Close() {
	if iv.reader != nil {
		iv.reader.Close()
		iv.reader = nil
	}
	if iv.writer != nil {
		iv.writer.Close()
		iv.writer = nil
	}
}
