package abi

import (
	"math/big"
	"strings"

	"github.com/Gessiux/neatchain/accounts/abi"
	"github.com/Gessiux/neatchain/common"
)

type FunctionType struct {
	id    int
	cross bool // Tx type, cross chain / non cross chain
	main  bool // allow to be execute on main chain or not
	child bool // allow to be execute on child chain or not
}

var (
	// Cross Chain Function
	CreateChildChain       = FunctionType{0, true, true, false}
	JoinChildChain         = FunctionType{1, true, true, false}
	DepositInMainChain     = FunctionType{2, true, true, false}
	DepositInChildChain    = FunctionType{3, true, false, true}
	WithdrawFromChildChain = FunctionType{4, true, false, true}
	WithdrawFromMainChain  = FunctionType{5, true, true, false}
	SaveDataToMainChain    = FunctionType{6, true, true, false}
	SetBlockReward         = FunctionType{7, true, false, true}
	// Non-Cross Chain Function
	VoteNextEpoch  = FunctionType{10, false, true, true}
	RevealVote     = FunctionType{11, false, true, true}
	Delegate       = FunctionType{12, false, true, true}
	UnDelegate     = FunctionType{13, false, true, true}
	Register       = FunctionType{14, false, true, true}
	UnRegister     = FunctionType{15, false, true, true}
	EditValidator  = FunctionType{16, false, true, true}
	WithdrawReward = FunctionType{17, false, true, true}
	UnForbidden    = FunctionType{18, false, true, true}
	SetCommission  = FunctionType{19, false, true, true}
	// Unknown
	Unknown = FunctionType{-1, false, false, false}
)

func (t FunctionType) IsCrossChainType() bool {
	return t.cross
}

func (t FunctionType) AllowInMainChain() bool {
	return t.main
}

func (t FunctionType) AllowInChildChain() bool {
	return t.child
}

func (t FunctionType) RequiredGas() uint64 {
	switch t {
	case CreateChildChain:
		return 42000
	case JoinChildChain:
		return 21000
	case DepositInMainChain:
		return 42000
	case DepositInChildChain:
		return 0
	case WithdrawFromChildChain:
		return 42000
	case WithdrawFromMainChain:
		return 0
	case SaveDataToMainChain:
		return 0
	case VoteNextEpoch:
		return 21000
	case RevealVote:
		return 21000
	case Delegate, UnDelegate, Register, UnRegister:
		return 21000
	case SetBlockReward:
		return 21000
	case EditValidator:
		return 21000
	case WithdrawReward:
		return 21000
	case UnForbidden:
		return 21000
	case SetCommission:
		return 21000
	default:
		return 0
	}
}

func (t FunctionType) String() string {
	switch t {
	case CreateChildChain:
		return "CreateChildChain"
	case JoinChildChain:
		return "JoinChildChain"
	case DepositInMainChain:
		return "DepositInMainChain"
	case DepositInChildChain:
		return "DepositInChildChain"
	case WithdrawFromChildChain:
		return "WithdrawFromChildChain"
	case WithdrawFromMainChain:
		return "WithdrawFromMainChain"
	case SaveDataToMainChain:
		return "SaveDataToMainChain"
	case VoteNextEpoch:
		return "VoteNextEpoch"
	case RevealVote:
		return "RevealVote"
	case Delegate:
		return "Delegate"
	case UnDelegate:
		return "UnDelegate"
	case Register:
		return "Register"
	case UnRegister:
		return "UnRegister"
	case SetBlockReward:
		return "SetBlockReward"
	case EditValidator:
		return "EditValidator"
	case WithdrawReward:
		return "WithdrawReward"
	case UnForbidden:
		return "UnForbidden"
	case SetCommission:
		return "SetCommission"
	default:
		return "UnKnown"
	}
}

func StringToFunctionType(s string) FunctionType {
	switch s {
	case "CreateChildChain":
		return CreateChildChain
	case "JoinChildChain":
		return JoinChildChain
	case "DepositInMainChain":
		return DepositInMainChain
	case "DepositInChildChain":
		return DepositInChildChain
	case "WithdrawFromChildChain":
		return WithdrawFromChildChain
	case "WithdrawFromMainChain":
		return WithdrawFromMainChain
	case "SaveDataToMainChain":
		return SaveDataToMainChain
	case "VoteNextEpoch":
		return VoteNextEpoch
	case "RevealVote":
		return RevealVote
	case "Delegate":
		return Delegate
	case "UnDelegate":
		return UnDelegate
	case "Register":
		return Register
	case "UnRegister":
		return UnRegister
	case "SetBlockReward":
		return SetBlockReward
	case "EditValidator":
		return EditValidator
	case "WithdrawReward":
		return WithdrawReward
	case "UnForbidden":
		return UnForbidden
	case "SetCommission":
		return SetCommission
	default:
		return Unknown
	}
}

type CreateChildChainArgs struct {
	ChainId          string
	MinValidators    uint16
	MinDepositAmount *big.Int
	StartBlock       *big.Int
	EndBlock         *big.Int
}

type JoinChildChainArgs struct {
	PubKey    []byte
	ChainId   string
	Signature []byte
}

type DepositInMainChainArgs struct {
	ChainId string
}

type DepositInChildChainArgs struct {
	ChainId string
	TxHash  common.Hash
}

type WithdrawFromChildChainArgs struct {
	ChainId string
}

type WithdrawFromMainChainArgs struct {
	ChainId string
	Amount  *big.Int
	TxHash  common.Hash
}

type VoteNextEpochArgs struct {
	VoteHash common.Hash
}

type RevealVoteArgs struct {
	PubKey    []byte
	Amount    *big.Int
	Salt      string
	Signature []byte
}

type DelegateArgs struct {
	Candidate common.Address
}

type UnDelegateArgs struct {
	Candidate common.Address
	Amount    *big.Int
}

type RegisterArgs struct {
	Pubkey     []byte
	Signature  []byte
	Commission uint8
}

type SetBlockRewardArgs struct {
	ChainId string
	Reward  *big.Int
}

type EditValidatorArgs struct {
	Moniker  string
	Website  string
	Identity string
	Details  string
}

type WithdrawRewardArgs struct {
	DelegateAddress common.Address
}

type UnForbiddenArgs struct {
}

type SetCommissionArgs struct {
	Commission uint8
}

const jsonChainABI = `
[
	{
		"type": "function",
		"name": "CreateChildChain",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			},
			{
				"name": "minValidators",
				"type": "uint16"
			},
			{
				"name": "minDepositAmount",
				"type": "uint256"
			},
			{
				"name": "startBlock",
				"type": "uint256"
			},
			{
				"name": "endBlock",
				"type": "uint256"
			}
		]
	},
	{
		"type": "function",
		"name": "JoinChildChain",
		"constant": false,
		"inputs": [
			{
				"name": "pubKey",
				"type": "bytes"
			},
			{
				"name": "chainId",
				"type": "string"
			},
			{
				"name": "signature",
				"type": "bytes"
			}
		]
	},
	{
		"type": "function",
		"name": "DepositInMainChain",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			}
		]
	},
	{
		"type": "function",
		"name": "DepositInChildChain",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			},
			{
				"name": "txHash",
				"type": "bytes32"
			}
		]
	},
	{
		"type": "function",
		"name": "WithdrawFromChildChain",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			}
		]
	},
	{
		"type": "function",
		"name": "WithdrawFromMainChain",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			},
			{
				"name": "amount",
				"type": "uint256"
			},
			{
				"name": "txHash",
				"type": "bytes32"
			}
		]
	},
	{
		"type": "function",
		"name": "SaveDataToMainChain",
		"constant": false,
		"inputs": [
			{
				"name": "data",
				"type": "bytes"
			}
		]
	},
	{
		"type": "function",
		"name": "VoteNextEpoch",
		"constant": false,
		"inputs": [
			{
				"name": "voteHash",
				"type": "bytes32"
			}
		]
	},
	{
		"type": "function",
		"name": "RevealVote",
		"constant": false,
		"inputs": [
			{
				"name": "pubKey",
				"type": "bytes"
			},
			{
				"name": "amount",
				"type": "uint256"
			},
			{
				"name": "salt",
				"type": "string"
			},
			{
				"name": "signature",
				"type": "bytes"
			}
		]
	},
	{
		"type": "function",
		"name": "Delegate",
		"constant": false,
		"inputs": [
			{
				"name": "candidate",
				"type": "address"
			}
		]
	},
	{
		"type": "function",
		"name": "UnDelegate",
		"constant": false,
		"inputs": [
			{
				"name": "candidate",
				"type": "address"
			},
			{
				"name": "amount",
				"type": "uint256"
			}
		]
	},
	{
		"type": "function",
		"name": "Register",
		"constant": false,
		"inputs": [
			{
				"name": "pubkey",
				"type": "bytes"
			},
            {
				"name": "signature",
				"type": "bytes"
			},
			{
				"name": "commission",
				"type": "uint8"
			}
		]
	},
	{
		"type": "function",
		"name": "UnRegister",
		"constant": false,
		"inputs": []
	},
	{
		"type": "function",
		"name": "SetBlockReward",
		"constant": false,
		"inputs": [
			{
				"name": "chainId",
				"type": "string"
			},
			{
				"name": "reward",
				"type": "uint256"
			}
		]
	},
	{
		"type": "function",
		"name": "EditValidator",
		"constant": false,
		"inputs": [
			{
				"name": "moniker",
				"type": "string"
			},
			{
				"name": "website",
				"type": "string"
			},
			{
				"name": "identity",
				"type": "string"
			},
			{
				"name": "details",
				"type": "string"
			}
		]
	},
	{
		"type": "function",
		"name": "WithdrawReward",
		"constant": false,
		"inputs": [
			{
				"name": "delegateAddress",
				"type": "address"
			}
		]
	},
	{
		"type": "function",
		"name": "UnForbidden",
		"constant": false,
		"inputs": []
	},
	{
		"type": "function",
		"name": "SetCommission",
		"constant": false,
		"inputs": [
			{
				"name": "commission",
				"type": "uint8"
			}
		]
	}
]`

// NeatChain Child Chain Token Incentive Address
var ChildChainTokenIncentiveAddr = common.StringToAddress("INT3EEEEEEEEEEEEEEEEEEEEEEEEEEEE")

// NeatChain Internal Contract Address
var ChainContractMagicAddr = common.StringToAddress("INT3FFFFFFFFFFFFFFFFFFFFFFFFFFFF") // don't conflict with neatchain/core/vm/contracts.go

var ChainABI abi.ABI

func init() {
	var err error
	ChainABI, err = abi.JSON(strings.NewReader(jsonChainABI))
	if err != nil {
		panic("fail to create the chain ABI: " + err.Error())
	}
}

func IsNeatChainContractAddr(addr *common.Address) bool {
	return addr != nil && *addr == ChainContractMagicAddr
}

func FunctionTypeFromId(sigdata []byte) (FunctionType, error) {
	m, err := ChainABI.MethodById(sigdata)
	if err != nil {
		return Unknown, err
	}

	return StringToFunctionType(m.Name), nil
}
