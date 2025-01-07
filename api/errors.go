package api

import (
	"context"

	"coingod/account"
	"coingod/asset"
	"coingod/blockchain/pseudohsm"
	"coingod/blockchain/rpc"
	"coingod/blockchain/signers"
	"coingod/blockchain/txbuilder"
	"coingod/contract"
	"coingod/errors"
	"coingod/net/http/httperror"
	"coingod/net/http/httpjson"
	"coingod/protocol/validation"
	"coingod/protocol/vm"
)

var (
	// ErrDefault is default Coingod API Error
	ErrDefault = errors.New("Coingod API Error")
)

func isTemporary(info httperror.Info, err error) bool {
	switch info.ChainCode {
	case "CG000": // internal server error
		return true
	case "CG001": // request timed out
		return true
	case "CG761": // outputs currently reserved
		return true
	case "CG706": // 1 or more action errors
		errs := errors.Data(err)["actions"].([]httperror.Response)
		temp := true
		for _, actionErr := range errs {
			temp = temp && isTemporary(actionErr.Info, nil)
		}
		return temp
	default:
		return false
	}
}

var respErrFormatter = map[error]httperror.Info{
	ErrDefault: {500, "CG000", "Coingod API Error"},

	// Signers error namespace (2xx)
	signers.ErrBadQuorum: {400, "CG200", "Quorum must be greater than or equal to 1, and must be less than or equal to the length of xpubs"},
	signers.ErrBadXPub:   {400, "CG201", "Invalid xpub format"},
	signers.ErrNoXPubs:   {400, "CG202", "At least one xpub is required"},
	signers.ErrDupeXPub:  {400, "CG203", "Root XPubs cannot contain the same key more than once"},

	// Contract error namespace (3xx)
	contract.ErrContractDuplicated: {400, "CG302", "Contract is duplicated"},
	contract.ErrContractNotFound:   {400, "CG303", "Contract not found"},

	// Transaction error namespace (7xx)
	// Build transaction error namespace (70x ~ 72x)
	account.ErrInsufficient:         {400, "CG700", "Funds of account are insufficient"},
	account.ErrImmature:             {400, "CG701", "Available funds of account are immature"},
	account.ErrReserved:             {400, "CG702", "Available UTXOs of account have been reserved"},
	account.ErrMatchUTXO:            {400, "CG703", "UTXO with given hash not found"},
	ErrBadActionType:                {400, "CG704", "Invalid action type"},
	ErrBadAction:                    {400, "CG705", "Invalid action object"},
	ErrBadActionConstruction:        {400, "CG706", "Invalid action construction"},
	txbuilder.ErrMissingFields:      {400, "CG707", "One or more fields are missing"},
	txbuilder.ErrBadAmount:          {400, "CG708", "Invalid asset amount"},
	account.ErrFindAccount:          {400, "CG709", "Account not found"},
	asset.ErrFindAsset:              {400, "CG710", "Asset not found"},
	txbuilder.ErrBadContractArgType: {400, "CG711", "Invalid contract argument type"},
	txbuilder.ErrOrphanTx:           {400, "CG712", "Transaction input UTXO not found"},
	txbuilder.ErrExtTxFee:           {400, "CG713", "Transaction fee exceeded max limit"},
	txbuilder.ErrNoGasInput:         {400, "CG714", "Transaction has no gas input"},

	// Submit transaction error namespace (73x ~ 79x)
	// Validation error (73x ~ 75x)
	validation.ErrTxVersion:                 {400, "CG730", "Invalid transaction version"},
	validation.ErrWrongTransactionSize:      {400, "CG731", "Invalid transaction size"},
	validation.ErrBadTimeRange:              {400, "CG732", "Invalid transaction time range"},
	validation.ErrNotStandardTx:             {400, "CG733", "Not standard transaction"},
	validation.ErrWrongCoinbaseTransaction:  {400, "CG734", "Invalid coinbase transaction"},
	validation.ErrWrongCoinbaseAsset:        {400, "CG735", "Invalid coinbase assetID"},
	validation.ErrCoinbaseArbitraryOversize: {400, "CG736", "Invalid coinbase arbitrary size"},
	validation.ErrEmptyResults:              {400, "CG737", "No results in the transaction"},
	validation.ErrMismatchedAssetID:         {400, "CG738", "Mismatched assetID"},
	validation.ErrMismatchedPosition:        {400, "CG739", "Mismatched value source/dest position"},
	validation.ErrMismatchedReference:       {400, "CG740", "Mismatched reference"},
	validation.ErrMismatchedValue:           {400, "CG741", "Mismatched value"},
	validation.ErrMissingField:              {400, "CG742", "Missing required field"},
	validation.ErrNoSource:                  {400, "CG743", "No source for value"},
	validation.ErrOverflow:                  {400, "CG744", "Arithmetic overflow/underflow"},
	validation.ErrPosition:                  {400, "CG745", "Invalid source or destination position"},
	validation.ErrUnbalanced:                {400, "CG746", "Unbalanced asset amount between input and output"},
	validation.ErrOverGasCredit:             {400, "CG747", "Gas credit has been spent"},
	validation.ErrGasCalculate:              {400, "CG748", "Gas usage calculate got a math error"},

	// VM error (76x ~ 78x)
	vm.ErrAltStackUnderflow:  {400, "CG760", "Alt stack underflow"},
	vm.ErrBadValue:           {400, "CG761", "Bad value"},
	vm.ErrContext:            {400, "CG762", "Wrong context"},
	vm.ErrDataStackUnderflow: {400, "CG763", "Data stack underflow"},
	vm.ErrDisallowedOpcode:   {400, "CG764", "Disallowed opcode"},
	vm.ErrDivZero:            {400, "CG765", "Division by zero"},
	vm.ErrFalseVMResult:      {400, "CG766", "False result for executing VM"},
	vm.ErrLongProgram:        {400, "CG767", "Program size exceeds max int32"},
	vm.ErrRange:              {400, "CG768", "Arithmetic range error"},
	vm.ErrReturn:             {400, "CG769", "RETURN executed"},
	vm.ErrRunLimitExceeded:   {400, "CG770", "Run limit exceeded because the CG Fee is insufficient"},
	vm.ErrShortProgram:       {400, "CG771", "Unexpected end of program"},
	vm.ErrToken:              {400, "CG772", "Unrecognized token"},
	vm.ErrUnexpected:         {400, "CG773", "Unexpected error"},
	vm.ErrUnsupportedVM:      {400, "CG774", "Unsupported VM because the version of VM is mismatched"},
	vm.ErrVerifyFailed:       {400, "CG775", "VERIFY failed"},

	// Mock HSM error namespace (8xx)
	pseudohsm.ErrDuplicateKeyAlias: {400, "CG800", "Key Alias already exists"},
	pseudohsm.ErrLoadKey:           {400, "CG801", "Key not found or wrong password"},
	pseudohsm.ErrDecrypt:           {400, "CG802", "Could not decrypt key with given passphrase"},
}

// Map error values to standard coingod error codes. Missing entries
// will map to internalErrInfo.
//
// TODO(jackson): Share one error table across Chain
// products/services so that errors are consistent.
var errorFormatter = httperror.Formatter{
	Default:     httperror.Info{500, "CG000", "Coingod API Error"},
	IsTemporary: isTemporary,
	Errors: map[error]httperror.Info{
		// General error namespace (0xx)
		context.DeadlineExceeded: {408, "CG001", "Request timed out"},
		httpjson.ErrBadRequest:   {400, "CG002", "Invalid request body"},
		rpc.ErrWrongNetwork:      {502, "CG103", "A peer core is operating on a different blockchain network"},

		//accesstoken authz err namespace (86x)
		errNotAuthenticated: {401, "CG860", "Request could not be authenticated"},
	},
}
