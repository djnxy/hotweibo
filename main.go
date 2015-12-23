package main

import (
	"net/http"
	"runtime"
	"time"

	"gopkg.in/tylerb/graceful.v1"
	"weibo.com/hotweibo/lib/config"
	"weibo.com/hotweibo/service/mixflow"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		mixflow.Test()
	})
	port := config.Instance.MustValue("app", "port")
	delay := config.Instance.MustInt("app", "grace_delay")
	graceful.Run(":"+port, time.Duration(delay)*time.Second, mux)
}
