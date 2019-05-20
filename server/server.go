package server

import (
	"fmt"
	"log"
	"net/http"
	"zlp-demo-golang/common"
)

type Server struct {
	*http.ServeMux
}

type Request struct {
	*http.Request
	PostData map[string]string
}

type HandlerFunc func(w http.ResponseWriter, r *Request) string

func NewServer() *Server {
	return &Server{
		ServeMux: http.NewServeMux(),
	}
}

func (sv *Server) withMiddlewares(handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData interface{}

		req := &Request{Request: r}

		if r.Method == http.MethodGet {
			r.ParseForm()
			requestData = r.Form
		} else {
			req.PostData = common.JSON.ParseReader(r.Body)
			requestData = req.PostData
		}

		log.Printf("[Request][%s][%s][%s] %+v", r.Method, r.Host, r.URL, requestData)
		resp := handler(w, req)
		log.Printf("[Response][%s][%s][%s] %+v", r.Method, r.Host, r.URL, resp)
		fmt.Fprint(w, resp)
	}
}

func (sv *Server) HandleFunc(pattern string, fn HandlerFunc) {
	sv.ServeMux.HandleFunc(pattern, sv.withMiddlewares(fn))
}
