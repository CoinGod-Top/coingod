package consensus

import (
	"encoding/binary"
	"fmt"
	"math"
	"strings"

	log "github.com/sirupsen/logrus"

	"coingod/crypto/ed25519/chainkd"
	"coingod/protocol/bc"
)

// consensus variables
const (
	// Max gas that one block contains
	MaxBlockGas    = uint64(10000000)
	VMGasRate      = int64(200)
	StorageGasRate = int64(1)
	MaxGasAmount   = int64(300000)

	// These configs need add to casper config in elegant way
	MaxNumOfValidators = int(10)
	InitCGSupply       = 320000000000000
	RewardThreshold    = 0.5
	BlockReward        = uint64(200000000)

	// config parameter for coinbase reward
	CoinbasePendingBlockNumber = uint64(10)
	MinVoteOutputAmount        = uint64(100000000)

	PayToWitnessPubKeyHashDataSize = 20
	PayToWitnessScriptHashDataSize = 32
	BCRPContractHashDataSize       = 32
	CoinbaseArbitrarySizeLimit     = 128

	BCRPRequiredCGAmount = uint64(100000000)

	CGAlias               = "CG"
	COINGODAlias          = "CoinGod"
	defaultVotePendingNum = 202500
)

type CasperConfig struct {
	// BlockTimeInterval, milliseconds, the block time interval for producing a block
	BlockTimeInterval uint64

	// MaxTimeOffsetMs represent the max number of seconds a block time is allowed to be ahead of the current time
	MaxTimeOffsetMs uint64

	// BlocksOfEpoch represent the block num in one epoch
	BlocksOfEpoch uint64

	// MinValidatorVoteNum is the minimum vote number of become validator
	MinValidatorVoteNum uint64

	// VotePendingBlockNumber is the locked block number of vote utxo
	VotePendingBlockNums []VotePendingBlockNum

	FederationXpubs []chainkd.XPub
}

type VotePendingBlockNum struct {
	BeginBlock uint64
	EndBlock   uint64
	Num        uint64
}

// CGAssetID is CG's asset id, the soul asset of Coingod
var CGAssetID = &bc.AssetID{
	V0: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V1: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V2: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V3: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
}

// CGDefinitionMap is the ....
var CGDefinitionMap = map[string]interface{}{
	"name":        COINGODAlias,
	"symbol":      CGAlias,
	"decimals":    8,
	"description": `CoinGod`,
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string, params *Params) bool {
	prefix = strings.ToLower(prefix)
	return prefix == params.Bech32HRPSegwit+"1"
}

// Params store the config for different network
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name            string
	Bech32HRPSegwit string
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds []string

	// CasperConfig defines the casper consensus parameters
	CasperConfig
}

// ActiveNetParams is ...
var ActiveNetParams = MainNetParams

// NetParams is the correspondence between chain_id and Params
var NetParams = map[string]Params{
	"mainnet": MainNetParams,
	"wisdom":  TestNetParams,
	"solonet": SoloNetParams,
}

// MainNetParams is the config for production
var MainNetParams = Params{
	Name:            "main",
	Bech32HRPSegwit: "cg",
	DefaultPort:     "46657",
	DNSSeeds:        []string{},
	CasperConfig: CasperConfig{
		BlockTimeInterval:   6000,
		MaxTimeOffsetMs:     3000,
		BlocksOfEpoch:       10,
		MinValidatorVoteNum: 3e14,
		VotePendingBlockNums: []VotePendingBlockNum{
			{BeginBlock: 0, EndBlock: math.MaxUint64, Num: defaultVotePendingNum},
		},
		FederationXpubs: []chainkd.XPub{
			xpub("8c675cc0d0de07618dedd702fe54321f3dd0ab46b4b50deac4b87940ac0a974f79b9e33ca3161bf8cbd8d64b8214bd85db2e9bb04be0393f41041278278530c3"),
		},
	},
}

// TestNetParams is the config for test-net
var TestNetParams = Params{
	Name:            "test",
	Bech32HRPSegwit: "tc",
	DefaultPort:     "46656",
	DNSSeeds:        []string{},
	CasperConfig: CasperConfig{
		BlockTimeInterval:    6000,
		MaxTimeOffsetMs:      3000,
		BlocksOfEpoch:        100,
		MinValidatorVoteNum:  1e8,
		VotePendingBlockNums: []VotePendingBlockNum{{BeginBlock: 0, EndBlock: math.MaxUint64, Num: 10}},
		FederationXpubs: []chainkd.XPub{
			xpub("8c675cc0d0de07618dedd702fe54321f3dd0ab46b4b50deac4b87940ac0a974f79b9e33ca3161bf8cbd8d64b8214bd85db2e9bb04be0393f41041278278530c3"),
		},
	},
}

// SoloNetParams is the config for test-net
var SoloNetParams = Params{
	Name:            "solo",
	Bech32HRPSegwit: "sc",
	CasperConfig: CasperConfig{
		BlockTimeInterval:    6000,
		MaxTimeOffsetMs:      24000,
		BlocksOfEpoch:        100,
		MinValidatorVoteNum:  1e8,
		VotePendingBlockNums: []VotePendingBlockNum{{BeginBlock: 0, EndBlock: math.MaxUint64, Num: 10}},
		FederationXpubs:      []chainkd.XPub{},
	},
}

func VotePendingBlockNums(height uint64) uint64 {
	for _, pendingNum := range ActiveNetParams.VotePendingBlockNums {
		if height >= pendingNum.BeginBlock && height < pendingNum.EndBlock {
			return pendingNum.Num
		}
	}
	return defaultVotePendingNum
}

// InitActiveNetParams load the config by chain ID
func InitActiveNetParams(chainID string) error {
	var exist bool
	if ActiveNetParams, exist = NetParams[chainID]; !exist {
		return fmt.Errorf("chain_id[%v] don't exist", chainID)
	}
	return nil
}

func xpub(str string) (xpub chainkd.XPub) {
	if err := xpub.UnmarshalText([]byte(str)); err != nil {
		log.Panicf("Fail converts a string to xpub")
	}
	return xpub
}
