package proxy

import (
	"bytes"
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"gproxy/auth"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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
		if !auth.IsHTTPAuthorized(c.Request) {
			log.Printf("unauthorized request to %s\n", c.Request.URL)
			c.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		for k, v := range c.Request.Header {
			headers[k] = v[0]
		}
		// 覆盖请求头
		nativeHost := headers["Host"]
		nativeOrigin := headers["Origin"]
		nativeReferer := headers["Referer"]
		headers["Host"] = GetDomain(url)
		if Anonymous {
			// 删除请求源信息
			headers["X-Real-Ip"] = RealIP
			headers["X-Forwarded-For"] = RealIP
			headers["User-Agent"] = UserAgent
			delete(headers, "Referer")
			// Origin
			headers["Origin"] = GetOrigin(url)
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
			if Anonymous {
				// 设置回Origin和Referer
				c.Header("Origin", nativeOrigin)
				c.Header("Referer", nativeReferer)
				c.Header("Host", nativeHost)
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

func SendRequest(method string, targetUrl string, headers map[string]string, body []byte) (*http.Response, error) {
	client := http.DefaultClient
	// 设置代理
	if ProxyPass != "" {
		proxy, _ := url.Parse(ProxyPass)
		tr := &http.Transport{
			Proxy:           http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{
			Transport: tr,
			Timeout:   time.Second * 5, //超时时间
		}
	}
	// 创建请求
	req, err := http.NewRequest(method, targetUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	log.Println(req.Header)
	log.Println(targetUrl)
	// 发送请求
	resp, err := client.Do(req)
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

func GetOrigin(url string) string {
	// 从url中获取完整域名
	s := strings.Split(url, "//")
	v2 := strings.Split(s[1], "/")[0]
	return s[0] + "//" + v2
}
