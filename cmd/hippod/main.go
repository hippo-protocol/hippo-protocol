package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/cmd/hippod/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)
		default:
			os.Exit(1)
		}
	}
}
