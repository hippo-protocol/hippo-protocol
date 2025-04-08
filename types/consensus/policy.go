package consensus

import "time"

const (
	MaxBlockSize    = 4194304   // 4MB, a single MsgSend Tx counts 200~500 bytes normally.
	MaxBlockGas     = 100000000 // 100 milion, a single MsgSend Tx consumes 50,000~100,000 gas.
	MinGasPrices    = "5000000000000" + DefaultHippoDenom
	BlockTimeSec    = 6
	UnbondingPeriod = 60 * 60 * 24 * 7 * 3 * time.Second
	// staking
	MaxValidators     = 22
	MinCommissionRate = 10
	// mint
	Minter              = 25
	InflationRateChange = 25
	InflationMin        = 0
	InflationMax        = 25
	BlocksPerYear       = uint64(60*60*24*365) / uint64(BlockTimeSec)
	// distr
	CommunityTax = 92
	// gov
	MinDepositTokens    = 50_000
	ExpeditedMinDeposit = 100_000
	MaxDepositPeriod    = 60 * 60 * 24 * 14 * time.Second
	VotingPeriod        = 60 * 60 * 24 * 14 * time.Second
	// slashing
	SignedBlocksWindow      = 10_000
	MinSignedPerWindow      = 75
	SlashFractionDoubleSign = 5
	SlashFractionDowntime   = 0
	// evidence
	MaxAgeDuration  = UnbondingPeriod * 30 / 21 // 30 days
	MaxAgeNumBlocks = BlocksPerYear * 30 / 365  // 30 days
)
