package contenttype

import (
	"fmt"
	"log"
	"net/http"

	"github.com/valyala/fasthttp"
)

type HTTPAdapter interface {
	GetAcceptHeader() string
	GetContentType() string
}

type NetHTTPAdapter struct {
	Request *http.Request
}

func (n NetHTTPAdapter) GetAcceptHeader() string {
	acceptHeaders := n.Request.Header.Values("Accept")
	if len(acceptHeaders) == 0 {
		return ""
	}
	return acceptHeaders[0]
}

func (n NetHTTPAdapter) GetContentType() string {
	contentTypeHeaders := n.Request.Header.Values("Content-Type")
	if len(contentTypeHeaders) == 0 {
		return ""
	}

	return contentTypeHeaders[0]
}

type FastHTTPAdapter struct {
	Request *fasthttp.Request
}

func (f FastHTTPAdapter) GetAcceptHeader() string {
	acceptHeaders := string(f.Request.Header.Peek("Accept"))
	return acceptHeaders
}
func (f FastHTTPAdapter) GetContentType() string {
	contentType := string(f.Request.Header.Peek("Content-Type"))
	return contentType
}

type NilAdapter struct {
}

func (n NilAdapter) GetAcceptHeader() string {
	panic("invalid call")
}

func (n NilAdapter) GetContentType() string {
	panic("invalid call")
}

func getAdapter(request interface{}) (HTTPAdapter, error) {
	var adapter HTTPAdapter
	switch request := request.(type) {
	case http.Request:
		log.Println("request")
		adapter = NetHTTPAdapter{
			Request: &request,
		}
	case *http.Request:
		log.Println("*request")
		adapter = NetHTTPAdapter{
			Request: request,
		}
	case *fasthttp.Request:
		log.Println("fastrequest")
		adapter = FastHTTPAdapter{
			Request: request,
		}
	default:
		return NilAdapter{}, fmt.Errorf("no valid type found for %T", request)
	}
	return adapter, nil
}
