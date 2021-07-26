package main

import (
	"path/filepath"

	cfg "github.com/Gessiux/go-config"
	"github.com/Gessiux/neatchain/chain/accounts/keystore"
	ncTypes "github.com/Gessiux/neatchain/chain/consensus/neatbyft/types"
	"github.com/Gessiux/neatchain/chain/log"
	neatnode "github.com/Gessiux/neatchain/network/node"
	"github.com/Gessiux/neatchain/utilities/utils"
	"gopkg.in/urfave/cli.v1"
)

const (
	// Client identifier to advertise over the network
	MainChain    = "neatchain"
	TestnetChain = "testnet"
)

type Chain struct {
	Id       string
	Config   cfg.Config
	NeatNode *neatnode.Node
}

func LoadMainChain(ctx *cli.Context, chainId string) *Chain {

	chain := &Chain{Id: chainId}
	config := utils.GetNeatConConfig(chainId, ctx)
	chain.Config = config

	log.Info("Make full node")
	stack := makeFullNode(ctx, GetCMInstance(ctx).cch, chainId)
	chain.NeatNode = stack

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
	config := utils.GetNeatConConfig(chainId, ctx)
	chain.Config = config

	log.Infof("chainId: %s, makeFullNode", chainId)
	cch := GetCMInstance(ctx).cch
	stack := makeFullNode(ctx, cch, chainId)
	if stack == nil {
		return nil
	} else {
		chain.NeatNode = stack
		return chain
	}
}

func StartChain(ctx *cli.Context, chain *Chain, startDone chan<- struct{}) error {

	log.Infof("Start Chain: %s", chain.Id)
	go func() {
		utils.StartNode(ctx, chain.NeatNode)

		if startDone != nil {
			startDone <- struct{}{}
		}
	}()

	return nil
}

func CreateChildChain(ctx *cli.Context, chainId string, validator ncTypes.PrivValidator, keyJson []byte, validators []ncTypes.GenesisValidator) error {

	// Get NeatCon config base on chain id
	config := utils.GetNeatConConfig(chainId, ctx)

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

	// Init the NEAT Genesis
	err := initEthGenesisFromExistValidator(chainId, config, validators)
	if err != nil {
		return err
	}

	// Init the NEAT Blockchain
	init_neatchain(chainId, config.GetString("neat_genesis_file"), ctx)

	// Init the NeatCon Genesis
	init_em_files(config, chainId, config.GetString("neat_genesis_file"), validators)

	return nil
}
