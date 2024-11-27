package consensus

import "time"

const (
	MinGasPrices    = "5000000000000ahippo"
	BlockTimeSec    = 6
	UnbondingPeriod = 60 * 60 * 24 * 7 * 3 * time.Second
	// staking
	MaxValidators     = 21
	MinCommissionRate = 10
	// mint
	Minter              = 25
	InflationRateChange = 25
	InflationMin        = 0
	InflationMax        = 25
	BlocksPerYear       = uint64(60*60*24*365) / uint64(BlockTimeSec)
	// distr
	CommunityTax = 0
	// gov
	MinDepositTokens = 50_000
	MaxDepositPeriod = 60 * 60 * 24 * 14 * time.Second
	VotingPeriod     = 60 * 60 * 24 * 14 * time.Second
	// crisis
	ConstantFee = 100_000
	// slashing
	SignedBlocksWindow      = 10_000
	MinSignedPerWindow      = 5
	SlashFractionDoubleSign = 5
	SlashFractionDowntime   = 0.01
)
