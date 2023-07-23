package main

import "os"

func main() {
	// 解析命令行参数（proxy domain [-p port]）
	args := os.Args[1:]
	if len(args) < 1 {
		panic("invalid args")
	}
	domain := args[0]
	port := "80"
	if len(args) > 1 && args[1] == "-p" {
		port = args[2]
	}
	// 启动代理服务
	NewProxy(domain, port)
}
