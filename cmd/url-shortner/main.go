package main

import (
	"fmt"
	"url-shortner/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Printf("%+v\n", cfg)
	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, "chi render"

	// TODO: init server: run server
}
