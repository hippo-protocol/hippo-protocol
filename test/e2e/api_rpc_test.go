package test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const apiBaseUrl = "http://localhost:1317"
const rpcBaseUrl = "http://localhost:26657"

func TestApi(t *testing.T) {
	client := &http.Client{}

	tests := []struct {
		path   string
		expect string
	}{
		{path: "/cosmos/auth/v1beta1/accounts", expect: `"accounts":`},

		{path: "/cosmos/base/tendermint/v1beta1/blocks/latest", expect: `"height":`},
		{path: "/cosmos/base/tendermint/v1beta1/node_info", expect: `"network":"hippo-protocol-testnet-1"`},

		{path: "/cosmos/bank/v1beta1/supply/by_denom?denom=ahp", expect: `"denom":"ahp"`},

		{path: "/cosmos/distribution/v1beta1/community_pool", expect: `"denom":"ahp"`},

		{path: "/cosmos/staking/v1beta1/pool", expect: `"bonded_tokens":`},
		{path: "/cosmos/staking/v1beta1/validators", expect: `"operator_address":`},
		{path: "/cosmos/staking/v1beta1/params", expect: `"min_commission_rate":`},

		{path: "/cosmos/mint/v1beta1/inflation", expect: `"inflation":`},

		{path: "/cosmos/gov/v1/proposals", expect: `"proposals":`},
		{path: "/cosmos/gov/v1/params/voting", expect: `"min_deposit":`},
	}

	for _, test := range tests {
		response, err := client.Get(apiBaseUrl + test.path)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}
		assert.Equal(t, http.StatusOK, response.StatusCode, fmt.Sprintf("Expected status code 200, got %d", response.StatusCode))
		assert.Contains(t, string(body), test.expect, fmt.Sprintf("Expected %s in response", test.expect))
	}

}

func TestRpc(t *testing.T) {
	client := &http.Client{}

	tests := []struct {
		path   string
		expect string
	}{
		{path: "/blockchain?minHeight=1&maxHeight=2", expect: `"last_height":`},
		{path: `/tx_search?query="tx.height>=1"&order_by="desc"&per_page=20&page=1`, expect: `"txs":`},
	}

	for _, test := range tests {
		response, err := client.Get(rpcBaseUrl + test.path)
		if err != nil {
			t.Fatalf("Error sending request: %v", err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
			return
		}
		assert.Equal(t, http.StatusOK, response.StatusCode, fmt.Sprintf("Expected status code 200, got %d", response.StatusCode))
		assert.Contains(t, string(body), test.expect, fmt.Sprintf("Expected %s in response", test.expect))
	}
}
