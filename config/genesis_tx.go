package config

import (
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"coingod/consensus"
	"coingod/protocol/bc"
	"coingod/protocol/bc/types"
)

// AssetIssue asset issue params
type AssetIssue struct {
	NonceStr           string
	IssuanceProgramStr string
	AssetDefinitionStr string
	AssetIDStr         string
	Amount             uint64
}

func (a *AssetIssue) nonce() []byte {
	bs, err := hex.DecodeString(a.NonceStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) issuanceProgram() []byte {
	bs, err := hex.DecodeString(a.IssuanceProgramStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) assetDefinition() []byte {
	bs, err := hex.DecodeString(a.AssetDefinitionStr)
	if err != nil {
		panic(err)
	}

	return bs
}

func (a *AssetIssue) id() bc.AssetID {
	var bs [32]byte
	bytes, err := hex.DecodeString(a.AssetIDStr)
	if err != nil {
		panic(err)
	}

	copy(bs[:], bytes)
	return bc.NewAssetID(bs)
}

// GenesisTxs make genesis block txs
func GenesisTxs() []*types.Tx {
	contract, err := hex.DecodeString("0014f09582056ca6dea02c11c136a831fd25d090927e")
	if err != nil {
		log.Panicf("fail on decode genesis tx output control program")
	}

	var txs []*types.Tx
	firstTxData := types.TxData{
		Version: 1,
		Inputs:  []*types.TxInput{types.NewCoinbaseInput([]byte("Information is power. -- Jan/11/2013. Computing is power. -- Apr/24/2024."))},
		Outputs: []*types.TxOutput{types.NewOriginalTxOutput(*consensus.CGAssetID, consensus.InitCGSupply, contract, nil)},
	}
	txs = append(txs, types.NewTx(firstTxData))

	inputs := []*types.TxInput{}
	outputs := []*types.TxOutput{}
	for _, asset := range assetIssues {
		inputs = append(inputs, types.NewIssuanceInput(asset.nonce(), asset.Amount, asset.issuanceProgram(), nil, asset.assetDefinition()))
		outputs = append(outputs, types.NewOriginalTxOutput(asset.id(), asset.Amount, contract, nil))
	}

	secondTxData := types.TxData{Version: 1, Inputs: inputs, Outputs: outputs}
	txs = append(txs, types.NewTx(secondTxData))
	return txs
}

var assetIssues = []*AssetIssue{
	{
		NonceStr:           "8e972359c6441299",
		IssuanceProgramStr: "ae20d66ab117eca2bba6aefed569e52d6bf68a7a4ad7d775cbc01f7b2cfcd798f7b22031ab3c2147c330c5e360b4e585047d1dea5f529476ad5aff475ecdd55541923120851b4a24975df6dbeb4f8e5348542764f85bed67b763875325aa5e45116751b35253ad",
		AssetDefinitionStr: "7b22646563696d616c73223a382c226465736372697074696f6e223a225472756d702c2054686520476f64204f6620436f696e73222c226e616d65223a225472756d70222c2271756f72756d223a312c2272656973737565223a66616c73652c2273796d626f6c223a225452554d50227d",
		AssetIDStr:         "f211a4285e0679f23188e6f44052d6ebd727ce1df088208a36f5ad9843f244b2",
		Amount:             100000000000000,
	},
}
