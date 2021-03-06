// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package neatapi implements the general NEAT Chain API functions.
package neatapi

import (
	"context"
	"math/big"

	"github.com/Gessiux/neatchain/chain/accounts"
	"github.com/Gessiux/neatchain/chain/core"
	"github.com/Gessiux/neatchain/chain/core/state"
	"github.com/Gessiux/neatchain/chain/core/types"
	"github.com/Gessiux/neatchain/chain/core/vm"
	"github.com/Gessiux/neatchain/neatdb"
	"github.com/Gessiux/neatchain/neatptc/downloader"
	"github.com/Gessiux/neatchain/network/rpc"
	"github.com/Gessiux/neatchain/params"
	"github.com/Gessiux/neatchain/utilities/common"
	"github.com/Gessiux/neatchain/utilities/event"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	// General Ethereum API
	Downloader() *downloader.Downloader
	ProtocolVersion() int
	SuggestPrice(ctx context.Context) (*big.Int, error)
	ChainDb() neatdb.Database
	EventMux() *event.TypeMux
	AccountManager() *accounts.Manager

	// BlockChain API
	SetHead(number uint64)
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error)
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
	GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error)
	GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error)
	GetTd(blockHash common.Hash) *big.Int
	GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error)
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
	SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription

	// TxPool API
	SendTx(ctx context.Context, signedTx *types.Transaction) error
	GetPoolTransactions() (types.Transactions, error)
	GetPoolTransaction(txHash common.Hash) *types.Transaction
	GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error)
	Stats() (pending int, queued int)
	TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions)
	SubscribeTxPreEvent(chan<- core.TxPreEvent) event.Subscription

	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block

	//SetInnerAPIBridge(inBridge InnerAPIBridge)
	//GetInnerAPIBridge() InnerAPIBridge
	GetCrossChainHelper() core.CrossChainHelper

	BroadcastTX3ProofData(proofData *types.TX3ProofData)
}

func GetAPIs(apiBackend Backend, solcPath string) []rpc.API {
	compiler := makeCompilerAPIs(solcPath)
	nonceLock := new(AddrLocker)
	txapi := NewPublicTransactionPoolAPI(apiBackend, nonceLock)
	//apiBackend.SetInnerAPIBridge(&APIBridge{txapi: txapi})

	all := []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicNEATChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   txapi,
			Public:    true,
		}, {
			Namespace: "int",
			Version:   "1.0",
			Service:   NewPublicNEATChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "int",
			Version:   "1.0",
			Service:   NewPublicBlockChainAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "int",
			Version:   "1.0",
			Service:   txapi,
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewPublicTxPoolAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(apiBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(apiBackend),
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "int",
			Version:   "1.0",
			Service:   NewPublicAccountAPI(apiBackend.AccountManager()),
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   NewPrivateAccountAPI(apiBackend, nonceLock),
			Public:    false,
		}, {
			Namespace: "int",
			Version:   "1.0",
			Service:   NewPublicNEATAPI(apiBackend, nonceLock),
			Public:    true,
		},
	}
	return append(compiler, all...)
}
