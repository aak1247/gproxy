package main

import (
	"gproxy/proxy"
	"log"
	"os"
	"strconv"
)

func main() {
	// 解析命令行参数（gproxy [-p port] [-m mode] domain）
	args := os.Args[1:]
	if len(args) < 1 {
		panic("invalid args")
	}
	domain := args[0]
	var err error
	port := 80
	mode := "http"
	cert := ""
	key := ""
	for i, arg := range args {
		if arg == "-p" && len(args) > i+1 {
			port, err = strconv.Atoi(args[i+1])
			if err != nil {
				log.Fatalf("error occured: %v", err)
			}
			i++
			continue
		}
		if arg == "-m" && len(args) > i+1 {
			mode = args[i+1]
			i++
			continue
		}
		if arg == "--key" && len(args) > i+1 {
			key = args[i+1]
			i++
			continue
		}
		if arg == "--cert" && len(args) > i+1 {
			cert = args[i+1]
			i++
			continue
		}
		domain = arg
	}
	switch mode {
	case "http":
		// 启动HTTP代理服务
		proxy.NewHttpProxy(domain, port)
	case "https":
		// 启动HTTPS代理服务
		proxy.NewHttpsProxy(domain, port, cert, key)
	case "tcp":
		proxy.NewTcpProxy(domain, port)
	case "ws":
		proxy.NewWSProxy(domain, port)
	case "wss":
		proxy.NewWSSProxy(domain, port, cert, key)
	}
}
