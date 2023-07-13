# Fine Snow | 细雪
<div align="center">
    <img alt="Go Badge" src="https://img.shields.io/badge/Go-2b7d9c?logo=go&logoColor=fff&style=flat"/>
    <img alt="License Badge" src="https://img.shields.io/github/license/fine-snow/finesnow"/>
</div>
<div align="center">
    <img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/fine-snow/finesnow?style=social">
    <img alt="GitHub forks" src="https://img.shields.io/github/forks/fine-snow/finesnow?style=social">
    <img alt="GitHub watchers" src="https://img.shields.io/github/watchers/fine-snow/finesnow?style=social">
</div>

### Feature Description | 功能描述
```text
A golang-based web framework that is easy to use | 一个易于使用的基于 Golang 的 Web 框架
```
### Simple Start | 简单开始
```go
/**
 * Precondition | 前提
 * 在项目根目录下进入命令行, 通过执行 go get github.com/fine-snow/finesnow 在项目中引入 fine-snow 的 mod 依赖
 * Go to the command line in the project root directory,
 * and introduce the fine-snow mod dependency into the project by executing 'go get github.com/fine-snow/finesnow'.
 */

// Sample Code | 示例代码
package main

import "github.com/fine-snow/finesnow/snow"


func sayHello(name string) string {
	return name + " Say Hello World"
}

func main() {
	
	// Register a GET HTTP request that returns the hello world string
	// 注册一个 GET HTTP 请求, 返回 hello world 字符串
	// Other types of HTTP request registration: snow. Post | snow. Put | snow. Delete
	// 其他类型 HTTP 请求注册方式: snow.Post | snow.Put | snow.Delete
	// After the project is launched, the browser accesses http://localhost:9801/hello to get hello world, and a simple get request is implemented
	// 项目启动后, 浏览器访问 http://localhost:9801/hello 得到 hello world, 一个简单的 GET HTTP 请求就实现了
	snow.Get("/hello", func() string {
		return "Hello World"
	})
	
	// One more GET HTTP request that brings in parameters
	// 再来一个带入参的 GET HTTP 请求
	// Note: Request functions with input parameters do not support anonymous writing
	// 提醒: 带有入参的请求函数不支持匿名写法
	// After the project is launched, the browser accesses http://localhost:9801/sayHello?name=Tom to get Tom Say Hello World, and a GET HTTP request that brings in parameters is implemented
	// 项目启动后, 浏览器访问 http://localhost:9801/sayHello?name=Tom 得到 Tom Say Hello World, 一个带入参的 GET HTTP 请求就实现了
	snow.Get("/sayHello", sayHello)
	
	// Run function startup framework: addr parameter is the IP and port to be started (null character, default port: 9801); The intercept parameter is a global interceptor
	// Run 函数启动框架: addr 参数为需要启动的ip和端口(传空字符, 默认端口: 9801); intercept 参数为全局拦截器
	snow.Run("", nil)
}
```