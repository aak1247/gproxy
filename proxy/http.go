package proxy

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// NewHttpProxy 启动代理服务
func NewHttpProxy(url string, localPort int) {
	r := gin.Default()
	r.Any("/*path", httpProxy(url))
	r.Run("0.0.0.0:" + strconv.Itoa(localPort))
}

func NewHttpsProxy(url string, localPort int, cert, key string) {
	r := gin.Default()
	r.Any("/*path", httpProxy(url))
	r.RunTLS("0.0.0.0:"+strconv.Itoa(localPort), cert, key)
}

func httpProxy(url string) func(c *gin.Context) {
	return func(c *gin.Context) {
		// 读取请求头
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			headers[k] = v[0]
		}
		// 覆盖请求头
		headers["Host"] = GetDomain(url)
		if Anonymous {
			// 删除请求源信息
			headers["X-Real-Ip"] = RealIP
			headers["X-Forwarded-For"] = RealIP
			headers["User-Agent"] = UserAgent
			delete(headers, "Referer")
		}
		// 读取请求体
		body, _ := c.GetRawData()
		// 发送请求
		resp, err := SendRequest(c.Request.Method, url+c.Request.RequestURI, headers, body)
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
		} else {
			// 设置响应头
			for k, v := range resp.Header {
				c.Header(k, v[0])
			}
			// 设置响应码
			c.Status(resp.StatusCode)
			// 设置响应体
			// 读取响应体
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(resp.Body)
			if err != nil {
				c.JSON(500, gin.H{
					"message": err.Error(),
				})
			}
			_, err = c.Writer.Write(buf.Bytes())
			if err != nil {
				c.JSON(500, gin.H{
					"message": err.Error(),
				})
			}
		}
	}
}

func SendRequest(method string, url string, headers map[string]string, body []byte) (*http.Response, error) {
	// 创建请求
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	log.Println(req.Header)
	log.Println(url)
	// 发送请求
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetDomain(url string) string {
	// 从url中获取域名
	v1 := strings.Split(url, "//")[1]
	v2 := strings.Split(v1, "/")[0]
	return v2
}
