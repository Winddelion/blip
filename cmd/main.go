package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Winddelion/blip/internal/checker"
	"github.com/Winddelion/blip/internal/cli"
	"github.com/Winddelion/blip/internal/db"
)

func main() {
	opts, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Recieved struct: %+v\n", opts)
	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		fmt.Println("error connecting to db: ", err)
		os.Exit(1)
	}

	i := 0
	for i < len(opts.Urls) {

		if opts.Search {
			result, err := db.GetChecksByURL(ctx, pool, opts.Urls[i], 5)
			if err != nil {
				fmt.Println("error checking by url: ", err)
				os.Exit(1)
			}
			fmt.Println("found: ", result)
			i++
			continue
		}

		result, err := checker.GetStatusCode(opts.Urls[i])
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		id, err := db.InsertCheck(ctx, pool, result.UsrURL, int32(result.Resp), int32(result.ResponseTime))
		if err != nil {
			fmt.Println("error: ", err)
		}
		fmt.Println("Got:")
		fmt.Printf("Domain: %v\n", result.UsrURL)
		fmt.Printf("Status Code: %v\n", result.Resp)
		fmt.Printf("In time of: %v\n", result.ResponseTime)
		fmt.Printf("Written into db with id of: %v\n", id)
		i++
	}
}
