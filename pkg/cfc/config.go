package cfc

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	BCE_USER_CODE_ROOT        = "BCE_USER_CODE_ROOT"
	BCE_CFC_INVOKE_PIPE       = "BCE_CFC_INVOKE_PIPE"
	BCE_CFC_RESPONSE_PIPE     = "BCE_CFC_RESPONSE_PIPE"
	BCE_USER_FUNCTION_BRN     = "BCE_USER_FUNCTION_BRN"
	BCE_USER_FUNCTION_NAME    = "BCE_USER_FUNCTION_NAME"
	BCE_USER_FUNCTION_MEMSIZE = "BCE_USER_FUNCTION_MEMSIZE"
	BCE_USER_FUNCTION_VERSION = "BCE_USER_FUNCTION_VERSION"
	BCE_USER_FUNCTION_HANDLER = "BCE_USER_FUNCTION_HANDLER"
)

type RuntimeConfig struct {
	startTime    time.Time // runtime启动时间
	invokePipe   int
	responsePipe int
	userCodeRoot string
	functionBrn  string
	functionName string
	functionVer  string
	memorySize   int // MB为单位
	handler      string
}

func NewRuntimeConfig() (*RuntimeConfig, error) {
	invokePipe, err := strconv.ParseInt(os.Getenv(BCE_CFC_INVOKE_PIPE), 10, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s error: %v", BCE_CFC_INVOKE_PIPE, err)
	}
	responsePipe, err := strconv.ParseInt(os.Getenv(BCE_CFC_RESPONSE_PIPE), 10, 0)
	if err != nil {
		responsePipe = invokePipe
	}
	memorySize, _ := strconv.ParseInt(os.Getenv(BCE_USER_FUNCTION_MEMSIZE), 10, 0)
	return &RuntimeConfig{
		startTime:    time.Now(),
		invokePipe:   int(invokePipe),
		responsePipe: int(responsePipe),
		userCodeRoot: os.Getenv(BCE_USER_CODE_ROOT),
		functionBrn:  os.Getenv(BCE_USER_FUNCTION_BRN),
		functionName: os.Getenv(BCE_USER_FUNCTION_NAME),
		functionVer:  os.Getenv(BCE_USER_FUNCTION_VERSION),
		handler:      os.Getenv(BCE_USER_FUNCTION_HANDLER),
		memorySize:   int(memorySize),
	}, nil
}

func (c *RuntimeConfig) StartTime() time.Time {
	return c.startTime
}

func (c *RuntimeConfig) InvokePipe() int {
	return c.invokePipe
}

func (c *RuntimeConfig) ResponsePipe() int {
	return c.responsePipe
}

func (c *RuntimeConfig) UserCodeRoot() string {
	return c.userCodeRoot
}

func (c *RuntimeConfig) FunctionName() string {
	return c.functionName
}

func (c *RuntimeConfig) FunctionBrn() string {
	return c.functionBrn
}

func (c *RuntimeConfig) FunctionVersion() string {
	return c.functionVer
}

func (c *RuntimeConfig) MemorySize() int {
	return c.memorySize
}

func (c *RuntimeConfig) Handler() string {
	return c.handler
}
