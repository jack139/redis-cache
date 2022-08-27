package http

import (
	"log"
	"time"
	"os"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/ferluci/fast-realip"
)

var (
	/* 日志输出使用 */
	output  = log.New(os.Stdout, "", 0)

	/* 返回值的 content-type */
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)


/* 处理返回值，返回json */
func respJson(ctx *fasthttp.RequestCtx, value string) {
	respJson := map[string]interface{}{
		"success": true,
		"data": value,
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func respError(ctx *fasthttp.RequestCtx, msg string) {
	log.Println("Error: ", msg)
	respJson := map[string]interface{}{
		"success": false,
		"msg": msg,
	}
	doJSONWrite(ctx, fasthttp.StatusOK, respJson)
}

func doJSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) {
	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	ctx.Response.SetStatusCode(code)
	start := time.Now()
	if err := json.NewEncoder(ctx).Encode(obj); err != nil {
		elapsed := time.Since(start)
		log.Printf("", elapsed, err.Error(), obj)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}



// 日志格式处理

// "github.com/AubSs/fasthttplogger"
func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// Combined format:
// [<time>] <remote-addr> | <HTTP/http-version> | <method> <url> - <status> - <response-time us> | <user-agent>
// [2017/05/31 - 13:27:28] 127.0.0.1:54082 | HTTP/1.1 | GET /hello - 200 - 48.279µs | Paw/3.1.1 (Macintosh; OS X/10.12.5) GCDHTTPRequest
func combined(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v (%v) | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			realip.FromRequest(ctx),
			getHttp(ctx),
			ctx.Method(),
			ctx.RequestURI(),
			ctx.Response.Header.StatusCode(),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}


// 返回 map 所有 key
func getMapKeys(m map[string]interface{}) *[]string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	return &keys
}
