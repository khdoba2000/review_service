package main

import (
	"fmt"
	"monorepo/src/review_service/app"
)

func main() {
	fmt.Println("main")

	app := app.New()

	app.Run()
}
