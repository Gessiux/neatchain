package core

import (
	"errors"
	"math/big"
	"sync"

	"github.com/Gessiux/go-crypto"
	dbm "github.com/Gessiux/go-db"
	"github.com/Gessiux/neatchain/common"
	"github.com/Gessiux/neatchain/consensus/ipbft/epoch"
	"github.com/Gessiux/neatchain/core/state"
	"github.com/Gessiux/neatchain/core/types"
	neatAbi "github.com/Gessiux/neatchain/neatabi/abi"
	"github.com/Gessiux/neatchain/neatclient"
)

type TX3LocalCache interface {
	GetTX3(chainId string, txHash common.Hash) *types.Transaction
	DeleteTX3(chainId string, txHash common.Hash)

	WriteTX3ProofData(proofData *types.TX3ProofData) error

	GetTX3ProofData(chainId string, txHash common.Hash) *types.TX3ProofData
	GetAllTX3ProofData() []*types.TX3ProofData
}

type CrossChainHelper interface {
	GetMutex() *sync.Mutex
	GetClient() *neatclient.Client
	GetMainChainId() string
	GetChainInfoDB() dbm.DB

	CanCreateChildChain(from common.Address, chainId string, minValidators uint16, minDepositAmount, startupCost *big.Int, startBlock, endBlock *big.Int) error
	CreateChildChain(from common.Address, chainId string, minValidators uint16, minDepositAmount *big.Int, startBlock, endBlock *big.Int) error
	ValidateJoinChildChain(from common.Address, pubkey []byte, chainId string, depositAmount *big.Int, signature []byte) error
	JoinChildChain(from common.Address, pubkey crypto.PubKey, chainId string, depositAmount *big.Int) error
	ReadyForLaunchChildChain(height *big.Int, stateDB *state.StateDB) ([]string, []byte, []string)
	ProcessPostPendingData(newPendingIdxBytes []byte, deleteChildChainIds []string)

	VoteNextEpoch(ep *epoch.Epoch, from common.Address, voteHash common.Hash, txHash common.Hash) error
	RevealVote(ep *epoch.Epoch, from common.Address, pubkey crypto.PubKey, depositAmount *big.Int, salt string, txHash common.Hash) error
	UpdateNextEpoch(ep *epoch.Epoch, from common.Address, pubkey crypto.PubKey, depositAmount *big.Int, salt string, txHash common.Hash) error

	GetHeightFromMainChain() *big.Int
	GetEpochFromMainChain() (string, *epoch.Epoch)
	GetTxFromMainChain(txHash common.Hash) *types.Transaction

	ChangeValidators(chainId string)

	// for epoch only
	VerifyChildChainProofData(bs []byte) error
	SaveChildChainProofDataToMainChain(bs []byte) error

	TX3LocalCache
	ValidateTX3ProofData(proofData *types.TX3ProofData) error
	ValidateTX4WithInMemTX3ProofData(tx4 *types.Transaction, tx3ProofData *types.TX3ProofData) error

	////SaveDataToMainV1 acceps both epoch and tx3
	//VerifyChildChainProofDataV1(proofData *types.ChildChainProofDataV1) error
	//SaveChildChainProofDataToMainChainV1(proofData *types.ChildChainProofDataV1) error
}

// CrossChain Callback
type CrossChainValidateCb = func(tx *types.Transaction, state *state.StateDB, cch CrossChainHelper) error
type CrossChainApplyCb = func(tx *types.Transaction, state *state.StateDB, ops *types.PendingOps, cch CrossChainHelper, mining bool) error

// Non-CrossChain Callback
type NonCrossChainValidateCb = func(tx *types.Transaction, state *state.StateDB, bc *BlockChain) error
type NonCrossChainApplyCb = func(tx *types.Transaction, state *state.StateDB, bc *BlockChain, ops *types.PendingOps) error

type EtdInsertBlockCb func(bc *BlockChain, block *types.Block)

var validateCbMap = make(map[neatAbi.FunctionType]interface{})
var applyCbMap = make(map[neatAbi.FunctionType]interface{})
var insertBlockCbMap = make(map[string]EtdInsertBlockCb)

func RegisterValidateCb(function neatAbi.FunctionType, validateCb interface{}) error {

	_, ok := validateCbMap[function]
	if ok {
		return errors.New("the name has registered in validateCbMap")
	}

	validateCbMap[function] = validateCb
	return nil
}

func GetValidateCb(function neatAbi.FunctionType) interface{} {

	cb, ok := validateCbMap[function]
	if ok {
		return cb
	}

	return nil
}

func RegisterApplyCb(function neatAbi.FunctionType, applyCb interface{}) error {

	_, ok := applyCbMap[function]
	if ok {
		return errors.New("the name has registered in applyCbMap")
	}

	applyCbMap[function] = applyCb

	return nil
}

func GetApplyCb(function neatAbi.FunctionType) interface{} {

	cb, ok := applyCbMap[function]
	if ok {
		return cb
	}

	return nil
}

func RegisterInsertBlockCb(name string, insertBlockCb EtdInsertBlockCb) error {

	_, ok := insertBlockCbMap[name]
	if ok {
		return errors.New("the name has registered in insertBlockCbMap")
	}

	insertBlockCbMap[name] = insertBlockCb

	return nil
}

func GetInsertBlockCb(name string) EtdInsertBlockCb {

	cb, ok := insertBlockCbMap[name]
	if ok {
		return cb
	}

	return nil
}

func GetInsertBlockCbMap() map[string]EtdInsertBlockCb {

	return insertBlockCbMap
}
