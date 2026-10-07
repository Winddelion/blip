package main

import (
	"fmt"
	"os"

	"github.com/Winddelion/blip/internal/checker"
	"github.com/Winddelion/blip/internal/cli"
)

func main() {
	opts, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Recieved struct: %+v\n", opts)

	i := 0
	for i < len(opts.Urls) {
		result, err := checker.GetStatusCode(opts.Urls[i])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("Got:")
		fmt.Printf("Domain: %v\n", result.UsrURL)
		fmt.Printf("Status Code: %v\n", result.Resp)
		fmt.Printf("In time of: %v\n", result.ResponseTime)
		i++
	}
}
