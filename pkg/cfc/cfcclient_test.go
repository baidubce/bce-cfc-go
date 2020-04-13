package cfc

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"syscall"
	"testing"
)

type testHandler struct {
}

func (h *testHandler) Handle(input io.Reader, output io.Writer, context InvokeContext) error {
	n, err := io.Copy(output, input)
	log.Printf("copy %d bytes\n", n)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func TestRegisterHandler(t *testing.T) {
	RegisterNamedHandler("index.handler", &testHandler{})

	handler, err := getInvokeHandler("index.handler")
	if err != nil {
		t.Fail()
		return
	}
	h1 := handler.(*testHandler)
	if h1 == nil {
		t.Fail()
		return
	}
	handler, err = getInvokeHandler("index.handler2") // 获取默认handler
	if err == nil {
		t.Fail()
		return
	}
}

func setupTestEnv(ivkPipe int, rspPipe int) {
	os.Setenv("BCE_CFC_INVOKE_PIPE", strconv.Itoa(ivkPipe))
	os.Setenv("BCE_CFC_RESPONSE_PIPE", strconv.Itoa(rspPipe))
	os.Setenv("BCE_USER_CODE_ROOT", "/var/task")
	os.Setenv("BCE_USER_FUNCTION_BRN", "function_brn")
	os.Setenv("BCE_USER_FUNCTION_NAME", "function_name")
	os.Setenv("BCE_USER_FUNCTION_MEMSIZE", "128")
	os.Setenv("BCE_USER_FUNCTION_VERSION", "$LATEST")
	os.Setenv("BCE_USER_FUNCTION_HANDLER", "index.handler")
}

func TestInvokeFunction(t *testing.T) {
	fds, err := syscall.Socketpair(syscall.AF_LOCAL, syscall.SOCK_STREAM, 0)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		syscall.Close(fds[0])
		syscall.Close(fds[1])
	}()

	setupTestEnv(fds[1], fds[1])
	cfg, err := NewRuntimeConfig()
	if err != nil {
		t.Error(err)
		return
	}

	cli := NewCfcClient(cfg, 0)
	go func() {
		cli.WaitInvoke()
	}()

	request := &InvokeRequest{
		RequestID:   "111",
		EventObject: "hello world",
	}
	file := os.NewFile(uintptr(fds[0]), "")
	encoder := json.NewEncoder(file)
	err = encoder.Encode(request)
	if err != nil {
		t.Error(err)
		return
	}

	response := &InvokeResponse{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(response)
	if err != nil {
		t.Error(err)
		return
	}
	if !response.Success {
		t.Error("invoke failed")
		return
	}
	if response.FuncResult != "hello world" {
		t.Error(response.FuncResult)
	}
}
