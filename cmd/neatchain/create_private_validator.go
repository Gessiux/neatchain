package main

import (
	"fmt"
	"github.com/Gessiux/go-crypto"
	"github.com/Gessiux/go-wire"
	"github.com/Gessiux/neatchain/cmd/utils"
	"github.com/Gessiux/neatchain/common"
	"github.com/Gessiux/neatchain/consensus/ipbft/types"
	"github.com/Gessiux/neatchain/log"
	"github.com/Gessiux/neatchain/params"
	"gopkg.in/urfave/cli.v1"
	"os"
	"path/filepath"
)

type PrivValidatorForConsole struct {
	// IntChain Account Address
	Address string `json:"address"`
	// IntChain Consensus Public Key, in BLS format
	PubKey crypto.PubKey `json:"consensus_pub_key"`
	// IntChain Consensus Private Key, in BLS format
	// PrivKey should be empty if a Signer other than the default is being used.
	PrivKey crypto.PrivKey `json:"consensus_priv_key"`
}

func CreatePrivateValidatorCmd(ctx *cli.Context) error {
	var consolePrivVal *PrivValidatorForConsole
	address := ctx.Args().First()

	if address == "" {
		log.Info("address is empty, need an address")
		return nil
	}

	datadir := ctx.GlobalString(utils.DataDirFlag.Name)
	if err := os.MkdirAll(datadir, 0700); err != nil {
		return err
	}

	chainId := params.MainnetChainConfig.IntChainId

	if ctx.GlobalIsSet(utils.TestnetFlag.Name) {
		chainId = params.TestnetChainConfig.IntChainId
	}

	privValFilePath := filepath.Join(ctx.GlobalString(utils.DataDirFlag.Name), chainId)
	privValFile := filepath.Join(ctx.GlobalString(utils.DataDirFlag.Name), chainId, "priv_validator.json")

	err := os.MkdirAll(privValFilePath, os.ModePerm)
	if err != nil {
		panic(err)
	}

	validator := types.GenPrivValidatorKey(common.StringToAddress(address))

	consolePrivVal = &PrivValidatorForConsole{
		Address: validator.Address.String(),
		PubKey:  validator.PubKey,
		PrivKey: validator.PrivKey,
	}

	fmt.Printf(string(wire.JSONBytesPretty(consolePrivVal)))
	validator.SetFile(privValFile)
	validator.Save()

	return nil
}
