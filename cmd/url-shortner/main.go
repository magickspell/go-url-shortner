package main

import (
	"fmt"
	"go-url-shortner/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Printf("CONIFG LOADED:\n%+v\n", cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: init server: run server
}
