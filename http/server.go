package http

import (
	"time"
	"log"
	"net/url"
	"github.com/valyala/fasthttp"
	"github.com/dgrr/http2"
	"github.com/fasthttp/router"

	"redis-cache/helper"
)


func RunServer() {

	/* router */
	r := router.New()
	r.GET("/", index)
	r.POST("/redis/cache", shoot)

	s := &fasthttp.Server{
		ReadTimeout: time.Second * 5,
		Handler: combined(r.Handler),
		Name:    "HTTP2 redis-cache",
	}

	log.Printf("start HTTP/2 server at %s\n", helper.Settings.Server.HTTP2_LISTEN)

	http2.ConfigureServer(s, http2.ServerConfig{})

	log.Fatal(s.ListenAndServeTLS(
		helper.Settings.Server.HTTP2_LISTEN,
		helper.Settings.Server.SSL_CERT_PATH + "/server.crt",
		helper.Settings.Server.SSL_CERT_PATH + "/server.key",
	))
}


func index(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Hello world.")
}


func shoot(ctx *fasthttp.RequestCtx) {
	// POST 的数据
	content := ctx.PostBody()

	//log.Println(string(content))

	m, err := url.ParseQuery(string(content))
	if err != nil {
		respError(ctx, err.Error())
		return
	}

	if m.Get("key")=="" {
		respError(ctx, "Empty key.")
		return
	}

	v, err := helper.Redis_shoot(m.Get("key"))
	if err!=nil {
		respError(ctx, err.Error())
		return
	}

	// 返回字符串
	//doJSONWrite(ctx, fasthttp.StatusOK, v)

	// 返回json
	respJson(ctx, v)
}
