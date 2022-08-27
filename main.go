package main

import (
	"flag"

	"redis-cache/helper"
	"redis-cache/http"
)

// openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 3650 -out server.crt -addext "subjectAltName = DNS:localhost"


var configPath = flag.String("config", "config/settings.yaml", "YAML config file path")


func main() {
	flag.Parse()

	// 初始化
	helper.InitSettings(*configPath)

	http.RunServer()
}
