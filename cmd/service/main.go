package main

import (
	"github.com/datalinkE/yet-another-chat/internal/config"
	"github.com/datalinkE/yet-another-chat/internal/service"
)

func main() {
	cfg := config.GetFromEnv()
	err := service.Run(cfg)
	if err != nil {
		panic(err)
	}
}
