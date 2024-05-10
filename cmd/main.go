package main

import (
	"github.com/flpnascto/rate-limiter-go/audit"
	"github.com/flpnascto/rate-limiter-go/configs"
	"github.com/flpnascto/rate-limiter-go/internal/infra/web"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	go web.NewWebServer(configs)

	audit.LoadTest()

}
