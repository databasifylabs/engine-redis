package main

import (
	"os"

	sdk "github.com/databasifylabs/engine-sdk"
	sentinel "github.com/databasifylabs/engine-sentinel-family"
)

func main() {
	os.Exit(sdk.Run(sentinel.New(sentinel.Redis)))
}
