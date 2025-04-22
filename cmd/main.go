package main

import (
	"notification/internal/app"
)

func main() {
	r := app.App()

	// Здесь можно добавить логгер, сигналы завершения и пр.
	r.Run(":8079")
}
