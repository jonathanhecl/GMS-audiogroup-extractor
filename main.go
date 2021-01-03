package main

import (
	"fmt"
	"os"
	"path/filepath"
)

var version string = "1.0.0"

func main() {

	fmt.Println("GameMaker Studio AudioGroup Extractor v" + version)
	fmt.Println("USAGE: " + filepath.Base(os.Args[0]) + " audiogroup1.dat")
	fmt.Println()

}
