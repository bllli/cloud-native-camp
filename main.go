package main

import (
	"fmt"
	"github.com/urfave/negroni"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Logger interface {
	log(level int, msg string)
}

type SimpleLogger struct {
}

func (l *SimpleLogger) log(level int, msg string) {
	fmt.Print(time.Now().Format("2006-04-02 15-04-05"), "; ", msg)
}

var (
	logger Logger = nil
)

func init() {
	logger = &SimpleLogger{}
}


func logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()
		w := negroni.NewResponseWriter(writer)
		next.ServeHTTP(w, request)
		logger.log(0, fmt.Sprintf("Requset path: %s; client ip: %s; spent: %d; statusCode: %d\n", request.RequestURI, request.RemoteAddr, time.Now().Sub(start), w.Status()))
	})
}


func echoRequestHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			writer.Header().Set(k, strings.Join(v, ","))
		}
		next.ServeHTTP(writer, request)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Welcome!")
}

func main() {
	http.Handle("/favicon.ico", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(404)
	}))
	http.Handle("/healthz", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, "200")
	}))
	http.Handle("/env", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		key := "VERSION"
		env := os.Getenv(key)
		if env == "" {
			env = "not found"
		}
		writer.Header().Set(key, env)
	}))
	http.Handle("/", logRequestMiddleware(echoRequestHeaderMiddleware(http.HandlerFunc(index))))
	server := &http.Server{
		Addr: ":8888",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
