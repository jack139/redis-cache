package main

import (
	//"fmt"
	"log"
	"flag"

	"github.com/valyala/fasthttp"
	"github.com/dgrr/http2"
	"github.com/fasthttp/router"

	"redis-cache/helper"
)

// openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 3650 -out server.crt -addext "subjectAltName = DNS:localhost"


var configPath = flag.String("config", "config/settings.yaml", "YAML config file path")


func main() {
	flag.Parse()

	// 初始化
	helper.InitSettings(*configPath)


	ret, err := helper.Mssql_test()
	if err!=nil {
		log.Fatal(err)
	}

	log.Println(ret)


	r := router.New()
	r.GET("/", index)
	r.GET("/get", shoot)

    s := &fasthttp.Server{
        Handler: r.Handler,
        Name:    "HTTP2 test",
    }

	log.Printf("start HTTP server at %s\n", helper.Settings.Server.HTTP2_LISTEN)

    http2.ConfigureServer(s, http2.ServerConfig{})
    
    log.Fatal(s.ListenAndServeTLS(
    	helper.Settings.Server.HTTP2_LISTEN, 
    	helper.Settings.Server.SSL_CERT_PATH + "/server.crt", 
    	helper.Settings.Server.SSL_CERT_PATH + "/server.key",
    ))
}

func index(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	ctx.WriteString("Hello world.")
}

func shoot(ctx *fasthttp.RequestCtx) {
	log.Printf("%v", ctx.RemoteAddr())
	v, err := helper.Redis_shoot("123")
	if err!=nil {
		log.Println(err)
		ctx.WriteString("!fail")
	}
	ctx.WriteString(v)
}
