package proxy

import (
	"fmt"
	ws2 "github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NewWSProxy(url string, localPort int) {
	//http.Handle("/proxyWsConn", )
	http.HandleFunc("/", wsProxyUrl(url))
	if err := http.ListenAndServe(":"+strconv.Itoa(localPort), nil); err != nil {
		log.Fatalf("failed to start ws server %v\n", err)
	}
}

func NewWSSProxy(url string, localPort int, cert, key string) {
	//http.Handle("/proxyWsConn", )
	http.HandleFunc("/", wsProxyUrl(url))
	if err := http.ListenAndServeTLS(":"+strconv.Itoa(localPort), cert, key, nil); err != nil {
		log.Fatalf("failed to start ws server %v\n", err)
	}
}

func wsProxyUrl(url string) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		if ws2.IsWebSocketUpgrade(request) {
			conn := initWSConn(writer, request)
			// 拿到所有的header和参数
			header := request.Header
			requestUrl := request.URL
			log.Printf("header %v", header)
			log.Printf("url %v", requestUrl)
			if url[len(url)-1] == '/' {
				url = url[:len(url)-1]
			}
			targetUrl := url + requestUrl.Path
			if len(requestUrl.RawQuery) != 0 {
				targetUrl += "?" + requestUrl.RawQuery
			}
			proxyWsConn(conn, targetUrl, header)
		} else {
			writer.Write([]byte("ok"))
		}
	}
}

func initWSConn(writer http.ResponseWriter, request *http.Request) *ws2.Conn {
	var upgrader = ws2.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			log.Println("升级协议", r.Header["User-Agent"])
			return true
		},
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	return conn
}

func proxyWsConn(ws *ws2.Conn, url string, headers http.Header) {
	defer func() {
		f := recover()
		if f != nil {
			fmt.Printf("fatal error %v\n", f)
		}
	}()
	// WSClient init
	server := &WSClient{}
	if err := server.Connect(url, headers); err != nil {
		log.Printf("connect error %v\n", err)
		ws.Close()
		return
	}

	// 上行 goroutine
	go func() {
		defer ws.Close()
		defer server.Close()
		defer func() {
			if f := recover(); f != nil {
				fmt.Printf("panic %v\n", f)
			}
		}()
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				if ce, ok := err.(*ws2.CloseError); ok && ws2.IsCloseError(ce, ce.Code) {
					log.Printf("close error %v\n", err)
					return
				} else if strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("close error %v\n", err)
					return
				} else {
					log.Printf("read client message failed %v\n", err)
					continue
				}
			}
			log.Println("获取客户端发送的消息:" + string(message))
			err = server.WriteMessage(mt, message)
			if err != nil {
				log.Printf("send message to server failed %v\n", err)
			}
			if mt == ws2.CloseMessage {
				log.Printf("close message %s\n", string(message))
				return
			}
		}
	}()

	// 下行 goroutine
	go func() {
		defer ws.Close()
		defer server.Close()
		defer func() {
			if f := recover(); f != nil {
				log.Printf("panic %v\n", f)
			}
		}()
		for {
			mt, message, err := server.ReadMessage()
			if err != nil {
				if ce, ok := err.(*ws2.CloseError); ok && ws2.IsCloseError(ce, ce.Code) {
					log.Printf("close error %v\n", err)
					return
				} else if strings.Contains(err.Error(), "use of closed network connection") {
					log.Printf("close error %v\n", err)
					return
				} else {
					log.Printf("read server message failed %v\n", err)
					continue
				}
			}
			log.Println("获取服务器发送的消息:" + string(message))
			err = ws.WriteMessage(mt, message)
			if err != nil {
				log.Printf("send message to client failed %v\n", err)
			}
			if mt == ws2.CloseMessage {
				log.Printf("close message %s\n", string(message))
				return
			}
		}
	}()
}

type WSClient struct {
	*ws2.Conn
}

func (cli *WSClient) Connect(url string, headers http.Header) error {
	// delete dup header
	headers.Del("Sec-WebSocket-Version")
	headers.Del("Sec-WebSocket-Key")
	headers.Del("Connection")
	headers.Del("Upgrade")
	headers.Del("Sec-Websocket-Extensions")
	headers.Del("Sec-Websocket-Protocol")
	conn, resp, err := ws2.DefaultDialer.Dial(url, headers)
	if err != nil {
		log.Printf("connect failed %v\n", err)
	} else {
		log.Printf("connect success %v\n", resp.Status)
		cli.Conn = conn
	}
	return err
}
