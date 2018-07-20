package server

import (
	"encoding/base64"
	"net"
	"net/http"
	"strconv"
	"strings"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/spaolacci/murmur3"
)

func (*Server) isMethodCacheable(m string) bool {
	switch strings.ToUpper(m) {
	case http.MethodHead:
		return true
	case http.MethodGet:
		return true
	default:
		return false
	}
}

func (*Server) isMethodAllowed(m string) bool {
	switch strings.ToUpper(m) {
	case http.MethodHead:
		return true
	case http.MethodGet:
		return true
	case http.MethodPost:
		return true
	case http.MethodPut:
		return true
	case http.MethodDelete:
		return true
	default:
		return false
	}
}

func (server *Server) checkAuthorization(r *http.Request) bool {
	// if server.checkStringIsNotEmpty(server.detectAPIKey(r)) {
	// 	return true
	// }
	return false
}

func (*Server) isStatusCodeCacheable(s int) bool {
	switch {
	case s == http.StatusOK:
		return true
	case s == http.StatusUnauthorized:
		return true
	default:
		return false
	}
}

func (server *Server) isClientProtocolSecure(r *http.Request) bool {
	if server.detectClientProtocol(r) == ProtocalHTTPS {
		return true
	}
	return false
}

func (server *Server) generateAPIKeyHash(apiKey string) string {
	return strconv.FormatUint(murmur3.Sum64([]byte(apiKey)), 16)
}

func (server *Server) generateRequestHash(r *http.Request) string {
	//Hash based on API Key + User Agent + Remote IP + URL Route
	key := server.detectAPIKey(r)
	key = key + "|" + server.detectUserAgent(r)
	key = key + "|" + server.detectClientIP(r)
	key = key + "|" + r.RequestURI
	key = key + "|" + r.Method
	hash := strconv.FormatUint(murmur3.Sum64([]byte(key)), 16)
	return hash
}

func (*Server) generateRequestID() string {
	uuid, _ := uuid.GenerateUUID()
	return uuid
}

func (*Server) detectContentType(c string) string {
	switch {
	case c == MIMEApplicationXML:
		return ContentTypeIsXML
	case c == MIMEApplicationXMLCharsetUTF8:
		return ContentTypeIsXML
	case c == MIMEApplicationJSON:
		return ContentTypeIsJSON
	case c == MIMEApplicationJSONCharsetUTF8:
		return ContentTypeIsJSON
	default:
		return ContentTypeIsJSON
	}
}

func (server *Server) detectClientIP(r *http.Request) string {
	result := r.Header[HeaderXForwardedFor]
	if len(result) >= 1 {
		// if server.checkStringIsNotEmpty(result[0]) {
		// 	return result[0]
		// }
	}
	return server.splitIPAddressFromPort(r.RemoteAddr)
}

func (server *Server) detectAPIKey(r *http.Request) string {
	username, _, ok := r.BasicAuth()
	if ok {
		return username
	}
	auth := r.Header.Get(HeaderAuthorization)
	noneStandardKey := server.getNonStandardBasicAuth(auth)
	if len(noneStandardKey) > 0 {
		return noneStandardKey
	}
	return server.getOAuthKey(r)
}

func (server *Server) getAPIKeyFromBasicAuth(r *http.Request) string {
	username, _, ok := r.BasicAuth()
	if ok {
		return username
	}
	auth := r.Header.Get(HeaderAuthorization)
	return server.getNonStandardBasicAuth(auth)
}

func (*Server) getNonStandardBasicAuth(auth string) string {

	authElements := strings.Fields(auth)
	if len(authElements) < 2 {
		return ""
	}
	if !strings.EqualFold(authElements[0], "Basic") {
		return ""
	}
	sanitizedAuth := authElements[1]
	c, err := base64.StdEncoding.DecodeString(sanitizedAuth)
	if err != nil {
		return ""
	}
	keyAndPwd := strings.Split(string(c), ":")
	return keyAndPwd[0]
}

func (*Server) getOAuthKey(r *http.Request) string {
	auth := r.Header.Get(HeaderAuthorization)
	authElements := strings.Fields(auth)
	if len(authElements) < 2 {
		// it is not OAuth
		return ""
	}
	if !strings.EqualFold(authElements[0], "Bearer") {
		return ""
	}
	return authElements[1]
}

func (*Server) detectUserAgent(r *http.Request) string {
	return r.Header.Get(HeaderUserAgent)
}

func (server *Server) detectClientProtocol(r *http.Request) string {
	result := r.Header[HeaderXForwardedProto]
	if len(result) >= 1 {
		// if server.checkStringIsNotEmpty(result[0]) {
		// 	if strings.ToLower(result[0]) == ProtocalHTTPS {
		// 		return ProtocalHTTPS
		// 	}
		// }
	}
	return ProtocalHTTP
}

func (server *Server) splitIPAddressFromPort(r string) string {
	ip, _, _ := net.SplitHostPort(r)
	return ip
}
