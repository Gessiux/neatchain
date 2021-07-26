package neatbyft

import (
	"crypto/ecdsa"
	"sync"

	"github.com/Gessiux/neatchain/chain/consensus"
	"github.com/Gessiux/neatchain/chain/consensus/neatbyft/types"
	"github.com/Gessiux/neatchain/chain/core"
	ethTypes "github.com/Gessiux/neatchain/chain/core/types"
	"github.com/Gessiux/neatchain/chain/log"
	"github.com/Gessiux/neatchain/params"
	"github.com/Gessiux/neatchain/utilities/common"
	"github.com/Gessiux/neatchain/utilities/event"
	"gopkg.in/urfave/cli.v1"
)

// New creates an Ethereum backend for NeatCon core engine.
func New(chainConfig *params.ChainConfig, cliCtx *cli.Context,
	privateKey *ecdsa.PrivateKey, cch core.CrossChainHelper) consensus.NeatByFT {
	// Allocate the snapshot caches and create the engine
	//recents, _ := lru.NewARC(inmemorySnapshots)
	//recentMessages, _ := lru.NewARC(inmemoryPeers)
	//knownMessages, _ := lru.NewARC(inmemoryMessages)

	config := GetNeatConConfig(chainConfig.NeatChainId, cliCtx)

	backend := &backend{
		//config:             config,
		chainConfig:     chainConfig,
		neatconEventMux: new(event.TypeMux),
		privateKey:      privateKey,
		//address:          crypto.PubkeyToAddress(privateKey.PublicKey),
		//core:             node,
		//chain:     chain,
		logger:    chainConfig.ChainLogger,
		commitCh:  make(chan *ethTypes.Block, 1),
		vcommitCh: make(chan *types.IntermediateBlockResult, 1),
		//recents:          recents,
		//candidates:  make(map[common.Address]bool),
		coreStarted: false,
		//recentMessages:   recentMessages,
		//knownMessages:    knownMessages,
	}
	backend.core = MakeNeatConNode(backend, config, chainConfig, cch)
	return backend
}

type backend struct {
	//config             cfg.Config
	chainConfig     *params.ChainConfig
	neatconEventMux *event.TypeMux
	privateKey      *ecdsa.PrivateKey
	address         common.Address
	core            *Node
	logger          log.Logger
	chain           consensus.ChainReader
	currentBlock    func() *ethTypes.Block
	hasBadBlock     func(hash common.Hash) bool

	// the channels for istanbul engine notifications
	commitCh          chan *ethTypes.Block
	vcommitCh         chan *types.IntermediateBlockResult
	proposedBlockHash common.Hash
	sealMu            sync.Mutex
	shouldStart       bool
	coreStarted       bool
	coreMu            sync.RWMutex

	// Current list of candidates we are pushing
	//candidates map[common.Address]bool
	// Protects the signer fields
	//candidatesLock sync.RWMutex
	// Snapshots for recent block to speed up reorgs
	//recents *lru.ARCCache

	// event subscription for ChainHeadEvent event
	broadcaster consensus.Broadcaster

	//recentMessages *lru.ARCCache // the cache of peer's messages
	//knownMessages  *lru.ARCCache // the cache of self messages
}

// WaitForTxs returns true if the consensus should wait for transactions before entering the propose step
//func (b *backend) WaitForTxs() bool {
//
//	return !b.config.GetBool("create_empty_blocks") || b.config.GetInt("create_empty_blocks_interval") > 0
//}
//
//func (b *backend) GetCreateEmptyBlocks() bool {
//	return b.config.GetBool("create_empty_blocks")
//}
//
//func (b *backend) GetCreateEmptyBlocksInterval() int {
//	return b.config.GetInt("create_empty_blocks_interval")
//}

func GetBackend() backend {
	return backend{}
}
