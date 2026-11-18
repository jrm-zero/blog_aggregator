package main

import (
	"fmt"
	"github.com/jrm-zero/blog_aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = cfg.SetUser("James"); err != nil {
		fmt.Println("fucked up")
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(cfg)

	return
}