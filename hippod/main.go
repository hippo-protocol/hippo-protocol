package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/hippod/cmd"
)

func main() {

	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		server.NewDefaultContext().Logger.Error(err.Error())
		os.Exit(1)
	}
}
