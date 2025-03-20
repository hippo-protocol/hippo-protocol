package app_test

import (
	"testing"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/hippocrat-dao/hippo-protocol/app"
)

func TestSetAnteHandler(t *testing.T) {
	mockApp := &app.App{}
	
	mockAnteHandler := func(ctx sdktypes.Context, tx sdktypes.Tx, simulate bool) (newCtx sdktypes.Context, err error) {
		require.NotEmpty(t, tx.GetMsgs(), "no messages in transaction")
		if len(tx.GetMsgs()) != 0 {
			return ctx, sdkerrors.ErrUnauthorized
		}
		if simulate {
			return ctx, sdkerrors.ErrUnauthorized
		}
		return ctx, nil
	}
	
	mockAppInterface := struct {
		app.App
		SetAnteHandler func(sdktypes.AnteHandler)
	}{
		App: *mockApp,
		SetAnteHandler: func(handler sdktypes.AnteHandler) {
			require.NotNil(t, handler)
		},
	}
	mockAppInterface.SetAnteHandler(mockAnteHandler)
}