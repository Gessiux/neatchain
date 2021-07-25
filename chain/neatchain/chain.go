package main

import (
	"path/filepath"

	cfg "github.com/Gessiux/go-config"
	"github.com/Gessiux/neatchain/chain/accounts/keystore"
	tdmTypes "github.com/Gessiux/neatchain/chain/consensus/neatbyft/types"
	"github.com/Gessiux/neatchain/chain/log"
	intnode "github.com/Gessiux/neatchain/network/node"
	"github.com/Gessiux/neatchain/utilities/utils"
	"gopkg.in/urfave/cli.v1"
)

const (
	// Client identifier to advertise over the network
	MainChain    = "neatchain"
	TestnetChain = "testnet"
)

type Chain struct {
	Id      string
	Config  cfg.Config
	IntNode *intnode.Node
}

func LoadMainChain(ctx *cli.Context, chainId string) *Chain {

	chain := &Chain{Id: chainId}
	config := utils.GetTendermintConfig(chainId, ctx)
	chain.Config = config

	log.Info("Make full node")
	stack := makeFullNode(ctx, GetCMInstance(ctx).cch, chainId)
	chain.IntNode = stack

	return chain
}

func LoadChildChain(ctx *cli.Context, chainId string) *Chain {

	log.Infof("now load child: %s", chainId)

	//chainDir := ChainDir(ctx, chainId)
	//empty, err := cmn.IsDirEmpty(chainDir)
	//log.Infof("chainDir is : %s, empty is %v", chainDir, empty)
	//if empty || err != nil {
	//	log.Errorf("directory %s not exist or with error %v", chainDir, err)
	//	return nil
	//}
	chain := &Chain{Id: chainId}
	config := utils.GetTendermintConfig(chainId, ctx)
	chain.Config = config

	log.Infof("chainId: %s, makeFullNode", chainId)
	cch := GetCMInstance(ctx).cch
	stack := makeFullNode(ctx, cch, chainId)
	if stack == nil {
		return nil
	} else {
		chain.IntNode = stack
		return chain
	}
}

func StartChain(ctx *cli.Context, chain *Chain, startDone chan<- struct{}) error {

	log.Infof("Start Chain: %s", chain.Id)
	go func() {
		utils.StartNode(ctx, chain.IntNode)

		if startDone != nil {
			startDone <- struct{}{}
		}
	}()

	return nil
}

func CreateChildChain(ctx *cli.Context, chainId string, validator tdmTypes.PrivValidator, keyJson []byte, validators []tdmTypes.GenesisValidator) error {

	// Get Tendermint config base on chain id
	config := utils.GetTendermintConfig(chainId, ctx)

	// Save the KeyStore File (Optional)
	if len(keyJson) > 0 {
		keystoreDir := config.GetString("keystore")
		keyJsonFilePath := filepath.Join(keystoreDir, keystore.KeyFileName(validator.Address))
		saveKeyError := keystore.WriteKeyStore(keyJsonFilePath, keyJson)
		if saveKeyError != nil {
			return saveKeyError
		}
	}

	// Save the Validator Json File
	privValFile := config.GetString("priv_validator_file_root")
	validator.SetFile(privValFile + ".json")
	validator.Save()

	// Init the INT Genesis
	err := initEthGenesisFromExistValidator(chainId, config, validators)
	if err != nil {
		return err
	}

	// Init the INT Blockchain
	init_int_blockchain(chainId, config.GetString("int_genesis_file"), ctx)

	// Init the Tendermint Genesis
	init_em_files(config, chainId, config.GetString("int_genesis_file"), validators)

	return nil
}
