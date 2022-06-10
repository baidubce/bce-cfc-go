## 定义handler
用户代码需要实现``InvokeHandler``接口，该接口定义如下：
```
type InvokeHandler interface {
	Handle(input io.Reader, output io.Writer, context InvokeContext) error
}
```
**input** 用户输入的事件内容  
**output** 接收函数运行的输出结果  
**context** 函数运行时环境。保存了函数的相关信息，以及访问其他服务所需的AK/SK/Session Token  

## 注册handler
注册handler需要调用``RegisterNamedHandler``接口。
该函数第一个参数为handler名称，要与前端界面配置保持一致。

## 定义main
代码库中实现了函数调用的基础框架，需要用户定义的main函数中调用框架的``Main``开始事件处理。

## 完整示例代码
``demo.go``中定义并注册handler  
```
package main

import (
	"io"
	"log"

	"github.com/baidubce/bce-cfc-go/pkg/cfc"
)

func init() {
    // 注册handler
    cfc.RegisterNamedHandler("echo_handler", &EchoHandler{})
}

// 定义handler
type EchoHandler struct {
}

func (h *EchoHandler) Handle(input io.Reader, output io.Writer, context cfc.InvokeContext) error {
	n, err := io.Copy(output, input)
	log.Printf("copy %d bytes\n", n)
	if err != nil {
		log.Println(err)
	}
	return nil
}
```

``main.go``执行自定义初始化逻辑，并调用框架``Main``函数开始事件处理  
```
import (
	"flag"

	"github.com/baidubce/bce-cfc-go/pkg/cfc"
)

func main() {
	flag.Parse()
    // 进入框架，处理函数事件
	cfc.Main()
}
```

## 编译代码
编译代码时，需要设置``GOOS``和``GOARCH``两个环境变量来指定对应的运行时环境。  
```
export GOOS=linux
export GOARCH=amd64
go build
```

## 打包输出
用户可以将生成的二进制执行文件打包到zip文件中，并上传zip包即可。

## 自定义程序启动参数
可以通过在zip包中添加一个``bootstrap``脚本文件，实现自定义的程序启动参数。
该文件中定义一个名为``ExecStart``的数组。示例如下：
```
#!/bin/bash
 
ExecStart=(
    ./golang-runtime
    "--alsologtostderr"
    "-v7"
)
```
