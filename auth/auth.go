package auth

import "net/http"

var Token string = ""

const TokenHeader = "X-PROXY-Authorization"

func IsHTTPAuthorized(req *http.Request) bool {
	if Token == "" {
		return true
	}
	// 请求头鉴权
	token := req.Header.Get(TokenHeader)
	if token == Token {
		req.Header.Del(TokenHeader)
		return true
	}
	return false
}
