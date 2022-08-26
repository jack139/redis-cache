package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/dgrr/http2"
	"github.com/fasthttp/router"
)

// openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.crt

func main() {
	r := router.New()
	r.GET("/", index)

    s := &fasthttp.Server{
        Handler: r.Handler,
        Name:    "HTTP2 test",
    }

	host := fmt.Sprintf("%s:%d", "", 8443)
	log.Printf("start HTTP server at %s\n", host)

    http2.ConfigureServer(s, http2.ServerConfig{})
    
    log.Fatal(s.ListenAndServeTLS(host, "server.crt", "server.key"))
}

func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}
