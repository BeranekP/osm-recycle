package main

import (
	"os"

)

func main() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "-U" {
			FetchData()
		}
	}

	if !FilesExist() {
		ConvertData()
		ValidateData()
	}
	ServeData()

}


