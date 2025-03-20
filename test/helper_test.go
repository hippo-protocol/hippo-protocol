package test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetApp(t *testing.T) {

	application := GetApp()

	require.NotNil(t, application, "The app instance should not be nil")
	
}