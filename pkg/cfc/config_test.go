package cfc

import (
	"testing"
)

func TestRuntimeConfig(t *testing.T) {
	setupTestEnv(3, 4)
	cfg, err := NewRuntimeConfig()
	if err != nil {
		t.Error(err)
		return
	}
	if cfg.InvokePipe() != 3 {
		t.Error(cfg.InvokePipe())
		return
	}
	if cfg.ResponsePipe() != 4 {
		t.Error(cfg.ResponsePipe())
		return
	}
	if cfg.UserCodeRoot() != "/var/task" {
		t.Error(cfg.UserCodeRoot())
		return
	}
	if cfg.FunctionName() != "function_name" {
		t.Error(cfg.FunctionName())
		return
	}
	if cfg.FunctionBrn() != "function_brn" {
		t.Error(cfg.FunctionBrn())
		return
	}
	if cfg.FunctionVersion() != "$LATEST" {
		t.Error(cfg.FunctionVersion())
		return
	}
	if cfg.MemorySize() != 128 {
		t.Error(cfg.MemorySize())
		return
	}
	if cfg.Handler() != "index.handler" {
		t.Error(cfg.Handler())
		return
	}
}
