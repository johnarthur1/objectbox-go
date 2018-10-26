package main

import (
	"fmt"
	"os"
	"path"

	"github.com/objectbox/objectbox-go/internal/generator"
)

func main() {
	file, _, err := getArgs()
	stopOnError(err)

	fmt.Printf("Generating ObjectBox bindings for %s", file)

	err = generator.Process(file)
	stopOnError(err)
}

func stopOnError(err error) {
	fmt.Printf(err.Error())
	os.Exit(1)
}

func getArgs() (file string, line uint, err error) {
	line = 0

	// if the command is run by go:generate some environment variables are set
	// https://golang.org/pkg/cmd/go/internal/generate/
	if gofile, exists := os.LookupEnv("GOFILE"); exists {
		file = gofile
		// TODO if we want to create for just one struct
		//if goline, exists := os.LookupEnv("GOLINE"); exists {
		//	line, err := strconv.ParseUint(goline, 10, 0)
		//	if err != nil {
		//		err = fmt.Errorf("Error parsing GOLINE environment variable as int: %s", err)
		//		return
		//	}
		//}
	}

	if len(file) == 0 {
		if len(os.Args) <= 1 {
			err = fmt.Errorf("usage: %s file.go", path.Base(os.Args[0]))
			return
		} else {
			file = os.Args[1]
		}
	}

	return
}
