package main

import (
	"math/big"
	"os"

	"github.com/cosmos/cosmos-sdk/server"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/hippocrat-dao/hippo-protocol/app"
	"github.com/hippocrat-dao/hippo-protocol/hippod/cmd"
	"github.com/hippocrat-dao/hippo-protocol/types/consensus"
)

func main() {
	// apply custom power reduction for 'a' base denom unit 10^18
	sdk.DefaultPowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(consensus.DefaultHippoPrecision), nil))

	err := sdk.RegisterDenom(consensus.DefaultHippoDenom, sdk.NewDecWithPrec(1, consensus.DefaultHippoPrecision))
	if err != nil {
		panic(err)
	}

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
