package cmd_test

import (
	"testing"

	"github.com/hippocrat-dao/hippo-protocol/hippod/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd()
	require.NotNil(t, rootCmd, "rootCmd should not be nil")
	require.IsType(t, &cobra.Command{}, rootCmd, "rootCmd should be of type *cobra.Command")
	require.Equal(t, "hippod", rootCmd.Use, "Command name should be 'hippod'")
	require.NotEmpty(t, rootCmd.Commands(), "rootCmd should have subcommands")
}
