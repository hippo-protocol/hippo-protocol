package app

import (
	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	ibcante "github.com/cosmos/ibc-go/v8/modules/core/ante"

	corestoretypes "cosmossdk.io/core/store"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
)

// setAnteHandler Reference github.com/cosmos/cosmos-sdk/x/auth/ante/ante.go
func (app *App) setAnteHandler(txConfig client.TxConfig, nodeConfig wasmtypes.NodeConfig, txCounterStoreService corestoretypes.KVStoreService) {
	app.SetAnteHandler(
		sdktypes.ChainAnteDecorators(
			ante.NewSetUpContextDecorator(),                                          // outermost AnteDecorator. SetUpContext must be called first
			wasmkeeper.NewLimitSimulationGasDecorator(nodeConfig.SimulationGasLimit), // after setup context to enforce limits early
			wasmkeeper.NewCountTXDecorator(txCounterStoreService),
			wasmkeeper.NewGasRegisterDecorator(app.WasmKeeper.GetGasRegister()), // registers gas costs for wasm operations
			wasmkeeper.NewTxContractsDecorator(),                                 // handles contract transaction decorations
			// Note: circuit breaker decorator is not added as the circuit module is not integrated in this blockchain.
			// If circuit breaker functionality is needed in the future, the circuit module should be added first.
			ante.NewExtensionOptionsDecorator(nil),
			ante.NewValidateBasicDecorator(),
			ante.NewTxTimeoutHeightDecorator(),
			ante.NewValidateMemoDecorator(app.AccountKeeper),
			ante.NewConsumeGasForTxSizeDecorator(app.AccountKeeper),
			ante.NewDeductFeeDecorator(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper, nil),
			ante.NewSetPubKeyDecorator(app.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
			ante.NewValidateSigCountDecorator(app.AccountKeeper),
			ante.NewSigGasConsumeDecorator(app.AccountKeeper, ante.DefaultSigVerificationGasConsumer),
			ante.NewSigVerificationDecorator(app.AccountKeeper, txConfig.SignModeHandler()),
			ante.NewIncrementSequenceDecorator(app.AccountKeeper),
			ibcante.NewRedundantRelayDecorator(app.IBCKeeper),
		),
	)
}
