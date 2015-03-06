package main

import (
	"os"
)

func main() {
	app := initApp()
	app.Run(os.Args)
}
