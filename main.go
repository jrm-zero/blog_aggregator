package main

import (
	"fmt"
	"github.com/jrm-zero/pokedexcli/internal/pokecache"
)

func main() {
	path, err := config.Read()
	fmt.Println(path)
	return
}