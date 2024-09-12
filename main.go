package main

import (
	"gproxy/auth"
	"gproxy/proxy"
	"log"
	"os"
	"strconv"
)

func main() {
	// 解析命令行参数（gproxy [-p port] [-m mode][--proxy proxypass][--anonymous][--realip realip][--useragent useragent] domain）
	args := os.Args[1:]
	domain := ""
	var err error
	port := 80
	mode := "http"
	cert := ""
	key := ""
	for i := 0; i < len(args); {
		arg := args[i]
		if arg == "-p" && len(args) > i+1 {
			port, err = strconv.Atoi(args[i+1])
			if err != nil {
				log.Fatalf("error occured: %v", err)
			}
			i += 2
			continue
		}
		if arg == "-m" && len(args) > i+1 {
			mode = args[i+1]
			i += 2
			continue
		}
		if arg == "--key" && len(args) > i+1 {
			key = args[i+1]
			i += 2
			continue
		}
		if arg == "--cert" && len(args) > i+1 {
			cert = args[i+1]
			i += 2
			continue
		}
		if arg == "--anonymous" {
			proxy.Anonymous = true
			i++
			continue
		}
		if arg == "--realip" && len(args) > i+1 {
			proxy.RealIP = args[i+1]
			i += 2
			continue
		}
		if arg == "--useragent" && len(args) > i+1 {
			proxy.UserAgent = args[i+1]
			i += 2
			continue
		}
		if arg == "--token" && len(args) > i+1 {
			auth.Token = args[i+1]
			i += 2
			continue
		}
		if arg == "--proxy" && len(args) > i+1 {
			proxy.ProxyPass = args[i+1]
			i += 2
			continue
		}
		domain = arg
		i++ // 防止死循环
	}
	// 检验逻辑
	if domain == "" {
		log.Fatal("domain is required")
	}
	// 仅http/https/tcp支持代理
	if (mode != "http" && mode != "https" && mode != "tcp") && proxy.ProxyPass != "" {
		log.Fatalf("proxy is only supported for http/https")
		return
	}
	// 仅http/https/ws/wss支持token
	if (mode != "http" && mode != "https" && mode != "ws" && mode != "wss") && auth.Token != "" {
		log.Fatalf("token is not supported for %s", mode)
		return
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
