package main

import (
	"flag"

	"github.com/Izumra/RefInatore/app/renamer"
)

func main() {
	var folderPath string

	flag.StringVar(
		&folderPath,
		"folderPath",
		"",
		"folderPath",
	)
	flag.Parse()

	renamer := renamer.New(folderPath)

	err := renamer.Rename()
	if err != nil {
		panic(err)
	}
}
