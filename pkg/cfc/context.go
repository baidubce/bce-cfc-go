package cfc

import "time"

type InvokeContext interface {
	GetRequestID() string
	GetFunctionBrn() string
	GetFunctionName() string
	GetFunctionVersion() string
	GetMemoryLimitMB() int
	GetCredential() *Credential
	GetClientContext() []byte
}

type cfcInvokeContext struct {
	requestId     string
	startTime     time.Time // 请求开始时间
	timeout       time.Duration
	functionBrn   string
	functionName  string
	functionVer   string
	handler       string
	memoryLimit   int
	credential    *Credential
	clientContext []byte
}

func NewInvokeContext(reqId string, config *RuntimeConfig, credential *Credential, clientCtx []byte) InvokeContext {
	return &cfcInvokeContext{
		requestId:     reqId,
		startTime:     time.Now(),
		functionBrn:   config.FunctionBrn(),
		functionName:  config.FunctionName(),
		functionVer:   config.FunctionVersion(),
		memoryLimit:   config.MemorySize(),
		handler:       config.Handler(),
		credential:    credential,
		clientContext: clientCtx,
	}
}

func (c *cfcInvokeContext) GetRequestID() string {
	return c.requestId
}

func (c *cfcInvokeContext) GetFunctionBrn() string {
	return c.functionBrn
}

func (c *cfcInvokeContext) GetFunctionName() string {
	return c.functionName
}

func (c *cfcInvokeContext) GetFunctionVersion() string {
	return c.functionVer
}

func (c *cfcInvokeContext) GetMemoryLimitMB() int {
	return c.memoryLimit
}

func (c *cfcInvokeContext) GetCredential() *Credential {
	return c.credential
}

func (c *cfcInvokeContext) GetClientContext() []byte {
	return c.clientContext
}
