package main

import (
	"fmt"
	"os"

	"github.com/Winddelion/blip/internal/cli"
)

func main() {
	opts, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", opts)

}
