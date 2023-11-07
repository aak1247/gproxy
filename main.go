package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	// 解析命令行参数（httpProxy domain [-p port]）
	args := os.Args[1:]
	if len(args) < 1 {
		panic("invalid args")
	}
	domain := args[0]
	var err error
	port := 80
	mode := "http"
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
		domain = arg
	}
	switch mode {
	case "http":
		// 启动HTTP代理服务
		NewHttpProxy(domain, port)
	case "tcp":
		NewTcpProxy(domain, port)
	}
}
