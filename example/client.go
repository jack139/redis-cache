package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"golang.org/x/net/http2"
)

const (
	post_url = "https://localhost:8443/redis/cache"
)

var (
	httpVersion = flag.Int("version", 2, "HTTP version")
	flagKey = flag.String("key", "", "cache key")
	flagShoot = flag.Int("num", 1, "num of goroutine")

	guard chan struct{}
)

func do_post(key string){
	defer func(){ // 释放 锁
		<-guard
	}()

	client := &http.Client{}

	// Create a pool with the server certificate since it is not signed
	// by a known CA
	caCert, err := ioutil.ReadFile("../cert/server.crt")
	if err != nil {
		log.Fatalf("Reading server certificate: %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create TLS configuration with the certificate of the server
	tlsConfig := &tls.Config{
		RootCAs: caCertPool,
	}

	// Use the proper transport in the client
	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	// Perform the request
	resp, err := client.PostForm(post_url, url.Values{"key": {key}})
	if err != nil {
		log.Fatalf("Failed POST: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body: %s", err)
	}
	fmt.Printf(
		"Got response %d: %s %s\n",
		resp.StatusCode, resp.Proto, string(body))
}

func main() {
	flag.Parse()

	guard = make(chan struct{}, *flagShoot)

	for i:=0;i<*flagShoot;i++ {
		guard <- struct{}{}
		go do_post(*flagKey)
	}

	// 都获取到才结束
	for i:=0;i<*flagShoot;i++ {
		guard <- struct{}{}
	}
}