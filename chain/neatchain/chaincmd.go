// Copyright 2015 The go-ethereum Authors
// This file is part of go-neatchain.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-neatchain. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/Gessiux/neatchain/neatdb"
	neatptc "github.com/Gessiux/neatchain/neatptc"
	"github.com/Gessiux/neatchain/params"

	"github.com/Gessiux/neatchain/chain/core"
	"github.com/Gessiux/neatchain/chain/core/rawdb"
	"github.com/Gessiux/neatchain/chain/core/state"
	"github.com/Gessiux/neatchain/chain/core/types"
	"github.com/Gessiux/neatchain/chain/log"
	"github.com/Gessiux/neatchain/neatptc/downloader"
	"github.com/Gessiux/neatchain/utilities/common"
	"github.com/Gessiux/neatchain/utilities/console"
	"github.com/Gessiux/neatchain/utilities/event"
	"github.com/Gessiux/neatchain/utilities/rlp"
	"github.com/Gessiux/neatchain/utilities/utils"
	"gopkg.in/urfave/cli.v1"
)

var (
	initNEATGenesisCmd = cli.Command{
		Action:    utils.MigrateFlags(initNeatGenesis),
		Name:      "init-neatchain",
		Usage:     "Initialize NEAT genesis.json file. init-neatchain {\"1000000000000000000000000000\",\"100000000000000000000000\"}",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Category:    "BLOCKCHAIN COMMANDS",
		Description: "The init-neatchain initializes a new NEAT genesis.json file for the network.",
	}

	initCommand = cli.Command{
		Action:    utils.MigrateFlags(initCmd),
		Name:      "init",
		Usage:     "Bootstrap and initialize a new genesis block",
		ArgsUsage: "<genesisPath>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The init command initializes a new genesis block and definition for the network.
This is a destructive action and changes the network in which you will be
participating.

It expects the genesis file as argument.`,
	}
	initSideChainCmd = cli.Command{
		Action:      utils.MigrateFlags(InitSideChainCmd),
		Name:        "init-side-chain",
		Usage:       "neatchain --sideChain=side_0,side_1,side_2 init-side-chain",
		Description: "Initialize side chain genesis from chain info db",
	}
	//	initCommand = cli.Command{
	//		Action:    utils.MigrateFlags(initGenesis),
	//		Name:      "init",
	//		Usage:     "Bootstrap and initialize a new genesis block",
	//		ArgsUsage: "<genesisPath>",
	//		Flags: []cli.Flag{
	//			utils.DataDirFlag,
	//		},
	//		Category: "BLOCKCHAIN COMMANDS",
	//		Description: `
	//The init command initializes a new genesis block and definition for the network.
	//This is a destructive action and changes the network in which you will be
	//participating.

	//It expects the genesis file as argument.`,
	//	}

	createValidatorCmd = cli.Command{
		//Action: GeneratePrivateValidatorCmd,
		Action: utils.MigrateFlags(CreatePrivateValidatorCmd),
		Name:   "create-validator",
		Usage:  "create-validator address", //create priv_validator.json for address
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Description: "Create priv_validator.json for address",
	}

	importCommand = cli.Command{
		Action:    utils.MigrateFlags(importChain),
		Name:      "import",
		Usage:     "Import a blockchain file",
		ArgsUsage: "<chainname> <filename> (<filename 2> ... <filename N>) ",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
			utils.GCModeFlag,
			utils.CacheDatabaseFlag,
			utils.CacheGCFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The import command imports blocks from an RLP-encoded form. The form can be one file
with several RLP-encoded blocks, or several files can be used.

If only one file is used, import error will result in failure. If several files are used,
processing will proceed even if an individual RLP-file import failure occurs.`,
	}
	exportCommand = cli.Command{
		Action:    utils.MigrateFlags(exportChain),
		Name:      "export",
		Usage:     "Export blockchain into file",
		ArgsUsage: "<chainname> <filename> [<blockNumFirst> <blockNumLast>]",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last block to write. In this mode, the file will be appended
if already existing.`,
	}
	importPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(importPreimages),
		Name:      "import-preimages",
		Usage:     "Import the preimage database from an RLP stream",
		ArgsUsage: "<datafile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
	The import-preimages command imports hash preimages from an RLP encoded stream.`,
	}
	exportPreimagesCommand = cli.Command{
		Action:    utils.MigrateFlags(exportPreimages),
		Name:      "export-preimages",
		Usage:     "Export the preimage database into an RLP stream",
		ArgsUsage: "<dumpfile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The export-preimages command export hash preimages to an RLP encoded stream`,
	}
	copydbCommand = cli.Command{
		Action:    utils.MigrateFlags(copyDb),
		Name:      "copydb",
		Usage:     "Create a local chain from a target chaindata folder",
		ArgsUsage: "<sourceChaindataDir>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
			utils.TestnetFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The first argument must be the directory containing the blockchain to download from`,
	}
	removedbCommand = cli.Command{
		Action:    utils.MigrateFlags(removeDB),
		Name:      "removedb",
		Usage:     "Remove blockchain and state databases",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			utils.DataDirFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
Remove blockchain and state databases`,
	}
	dumpCommand = cli.Command{
		Action:    utils.MigrateFlags(dump),
		Name:      "dump",
		Usage:     "Dump a specific block from storage",
		ArgsUsage: "[<blockHash> | <blockNum>]...",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
The arguments are interpreted as block numbers or hashes.
Use "ethereum dump 0" to dump the genesis block.`,
	}
	countBlockStateCommand = cli.Command{
		Action:    utils.MigrateFlags(countBlockState),
		Name:      "count-blockstate",
		Usage:     "Count the block state",
		ArgsUsage: "<datafile>",
		Flags: []cli.Flag{
			utils.DataDirFlag,
			utils.CacheFlag,
			utils.SyncModeFlag,
		},
		Category: "BLOCKCHAIN COMMANDS",
		Description: `
	The count-blockstate command count the block state from a given height.`,
	}

	versionCommand = cli.Command{
		Action:    utils.MigrateFlags(version),
		Name:      "version",
		Usage:     "Print version numbers",
		ArgsUsage: " ",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
The output of this command is supposed to be machine-readable.
`,
	}
)

// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
//func initGenesis(ctx *cli.Context) error {
//	// Make sure we have a valid genesis JSON
//	genesisPath := ctx.Args().First()
//	if len(genesisPath) == 0 {
//		utils.Fatalf("Must supply path to genesis JSON file")
//	}
//	file, err := os.Open(genesisPath)
//	if err != nil {
//		utils.Fatalf("Failed to read genesis file: %v", err)
//	}
//	defer file.Close()
//
//	genesis := new(core.Genesis)
//	if err := json.NewDecoder(file).Decode(genesis); err != nil {
//		utils.Fatalf("invalid genesis file: %v", err)
//	}
//	// Open an initialise both full and light databases
//	stack := makeFullNode(ctx)
//	for _, name := range []string{"chaindata", "lightchaindata"} {
//		chaindb, err := stack.OpenDatabase(name, 0, 0, "")
//		if err != nil {
//			utils.Fatalf("Failed to open database: %v", err)
//		}
//		_, hash, err := core.SetupGenesisBlock(chaindb, genesis)
//		if err != nil {
//			utils.Fatalf("Failed to write genesis block: %v", err)
//		}
//		log.Info("Successfully wrote genesis state", "database", name, "hash", hash)
//	}
//	return nil
//}

func importChain(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	chainName := ctx.Args().First()
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, cfg := makeConfigNode(ctx, chainName)
	cch := GetCMInstance(ctx).cch
	utils.RegisterIntService(stack, &cfg.Eth, ctx, cch)
	//stack := makeFullNode(ctx)
	defer stack.Close()

	chain, db := utils.MakeChain(ctx, stack)
	defer db.Close()

	// Start periodically gathering memory profiles
	var peakMemAlloc, peakMemSys uint64
	go func() {
		stats := new(runtime.MemStats)
		for {
			runtime.ReadMemStats(stats)
			if atomic.LoadUint64(&peakMemAlloc) < stats.Alloc {
				atomic.StoreUint64(&peakMemAlloc, stats.Alloc)
			}
			if atomic.LoadUint64(&peakMemSys) < stats.Sys {
				atomic.StoreUint64(&peakMemSys, stats.Sys)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	// Import the chain
	start := time.Now()

	if len(ctx.Args()) == 2 {
		if err := utils.ImportChain(chain, ctx.Args().Get(1)); err != nil {
			log.Error("Import error", "err", err)
		}
	} else {
		for i, arg := range ctx.Args() {
			if i == 0 {
				continue // skip the chain name
			}
			if err := utils.ImportChain(chain, arg); err != nil {
				log.Error("Import error", "file", arg, "err", err)
			}
		}
	}
	chain.Stop()
	fmt.Printf("Import done in %v.\n\n", time.Since(start))

	// Output pre-compaction stats mostly to see the import trashing
	stats, err := db.Stat("leveldb.stats")
	if err != nil {
		utils.Fatalf("Failed to read database stats: %v", err)
	}
	fmt.Println(stats)

	ioStats, err := db.Stat("leveldb.iostats")
	if err != nil {
		utils.Fatalf("Failed to read database iostats: %v", err)
	}
	fmt.Println(ioStats)

	// Print the memory statistics used by the importing
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	fmt.Printf("Object memory: %.3f MB current, %.3f MB peak\n", float64(mem.Alloc)/1024/1024, float64(atomic.LoadUint64(&peakMemAlloc))/1024/1024)
	fmt.Printf("System memory: %.3f MB current, %.3f MB peak\n", float64(mem.Sys)/1024/1024, float64(atomic.LoadUint64(&peakMemSys))/1024/1024)
	fmt.Printf("Allocations:   %.3f million\n", float64(mem.Mallocs)/1000000)
	fmt.Printf("GC pause:      %v\n\n", time.Duration(mem.PauseTotalNs))

	if ctx.GlobalIsSet(utils.NoCompactionFlag.Name) {
		return nil
	}

	// Compact the entire database to more accurately measure disk io and print the stats
	start = time.Now()
	fmt.Println("Compacting entire database...")
	if err = db.Compact(nil, nil); err != nil {
		utils.Fatalf("Compaction failed: %v", err)
	}
	fmt.Printf("Compaction done in %v.\n\n", time.Since(start))

	stats, err = db.Stat("leveldb.stats")
	if err != nil {
		utils.Fatalf("Failed to read database stats: %v", err)
	}
	fmt.Println(stats)

	ioStats, err = db.Stat("leveldb.iostats")
	if err != nil {
		utils.Fatalf("Failed to read database iostats: %v", err)
	}
	fmt.Println(ioStats)

	return nil
}

func exportChain(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	chainName := ctx.Args().First()
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, cfg := makeConfigNode(ctx, chainName)
	utils.RegisterIntService(stack, &cfg.Eth, ctx, GetCMInstance(ctx).cch)
	//stack := makeFullNode(ctx)
	defer stack.Close()

	chain, _ := utils.MakeChain(ctx, stack)
	start := time.Now()

	var err error
	fp := ctx.Args().Get(1)
	if len(ctx.Args()) < 4 {
		err = utils.ExportChain(chain, fp)
	} else {
		// This can be improved to allow for numbers larger than 9223372036854775807
		first, ferr := strconv.ParseInt(ctx.Args().Get(2), 10, 64)
		last, lerr := strconv.ParseInt(ctx.Args().Get(3), 10, 64)
		if ferr != nil || lerr != nil {
			utils.Fatalf("Export error in parsing parameters: block number not an integer\n")
		}
		if first < 0 || last < 0 {
			utils.Fatalf("Export error: block number must be greater than 0\n")
		}
		err = utils.ExportAppendChain(chain, fp, uint64(first), uint64(last))
	}

	if err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v", time.Since(start))
	return nil
}

// importPreimages imports preimage data from the specified file.
func importPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, cfg := makeConfigNode(ctx, chainName)
	utils.RegisterIntService(stack, &cfg.Eth, ctx, GetCMInstance(ctx).cch)
	defer stack.Close()

	db := utils.MakeChainDatabase(ctx, stack)
	start := time.Now()

	if err := utils.ImportPreimages(db, ctx.Args().First()); err != nil {
		utils.Fatalf("Import error: %v\n", err)
	}
	fmt.Printf("Import done in %v\n", time.Since(start))
	return nil
}

// exportPreimages dumps the preimage data to specified json file in streaming way.
func exportPreimages(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, cfg := makeConfigNode(ctx, chainName)
	utils.RegisterIntService(stack, &cfg.Eth, ctx, GetCMInstance(ctx).cch)
	defer stack.Close()

	db := utils.MakeChainDatabase(ctx, stack)
	start := time.Now()

	if err := utils.ExportPreimages(db, ctx.Args().First()); err != nil {
		utils.Fatalf("Export error: %v\n", err)
	}
	fmt.Printf("Export done in %v\n", time.Since(start))
	return nil
}
func copyDb(ctx *cli.Context) error {
	// Ensure we have a source chain directory to copy
	if len(ctx.Args()) != 1 {
		utils.Fatalf("Source chaindata directory path argument missing")
	}

	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	// Initialize a new chain for the running node to sync into
	stack, _ := makeConfigNode(ctx, chainName)
	chain, chainDb := utils.MakeChain(ctx, stack)

	syncmode := *utils.GlobalTextMarshaler(ctx, utils.SyncModeFlag.Name).(*downloader.SyncMode)
	dl := downloader.New(syncmode, chainDb, new(event.TypeMux), chain, nil, nil, nil)

	// Create a source peer to satisfy downloader requests from
	db, err := rawdb.NewLevelDBDatabase(ctx.Args().First(), ctx.GlobalInt(utils.CacheFlag.Name), 256, "")
	if err != nil {
		return err
	}
	hc, err := core.NewHeaderChain(db, chain.Config(), chain.Engine(), func() bool { return false })
	if err != nil {
		return err
	}
	peer := downloader.NewFakePeer("local", db, hc, dl)
	if err = dl.RegisterPeer("local", 63, peer); err != nil {
		return err
	}
	// Synchronise with the simulated peer
	start := time.Now()

	currentHeader := hc.CurrentHeader()
	if err = dl.Synchronise("local", currentHeader.Hash(), hc.GetTd(currentHeader.Hash(), currentHeader.Number.Uint64()), syncmode); err != nil {
		return err
	}
	for dl.Synchronising() {
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("Database copy done in %v\n", time.Since(start))

	// Compact the entire database to remove any sync overhead
	start = time.Now()
	fmt.Println("Compacting entire database...")
	if err = db.Compact(nil, nil); err != nil {
		utils.Fatalf("Compaction failed: %v", err)
	}
	fmt.Printf("Compaction done in %v.\n\n", time.Since(start))

	return nil
}

func removeDB(ctx *cli.Context) error {
	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, _ := makeConfigNode(ctx, chainName)

	for _, name := range []string{"chaindata", "lightchaindata"} {
		// Ensure the database exists in the first place
		logger := log.New("database", name)

		dbdir := stack.ResolvePath(name)
		if !common.FileExist(dbdir) {
			logger.Info("Database doesn't exist, skipping", "path", dbdir)
			continue
		}
		// Confirm removal and execute
		fmt.Println(dbdir)
		confirm, err := console.Stdin.PromptConfirm("Remove this database?")
		switch {
		case err != nil:
			utils.Fatalf("%v", err)
		case !confirm:
			logger.Warn("Database deletion aborted")
		default:
			start := time.Now()
			os.RemoveAll(dbdir)
			logger.Info("Database successfully deleted", "elapsed", common.PrettyDuration(time.Since(start)))
		}
	}
	return nil
}

func dump(ctx *cli.Context) error {
	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, _ := makeConfigNode(ctx, chainName)
	chain, chainDb := utils.MakeChain(ctx, stack)
	for _, arg := range ctx.Args() {
		var block *types.Block
		if hashish(arg) {
			block = chain.GetBlockByHash(common.HexToHash(arg))
		} else {
			num, _ := strconv.Atoi(arg)
			block = chain.GetBlockByNumber(uint64(num))
		}
		if block == nil {
			fmt.Println("{}")
			utils.Fatalf("block not found")
		} else {
			state, err := state.New(block.Root(), state.NewDatabase(chainDb))
			if err != nil {
				utils.Fatalf("could not create new state: %v", err)
			}
			fmt.Printf("%s\n", state.Dump())
		}
	}
	chainDb.Close()
	return nil
}

// hashish returns true for strings that look like hashes.
func hashish(x string) bool {
	_, err := strconv.Atoi(x)
	return err != nil
}

func countBlockState(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		utils.Fatalf("This command requires an argument.")
	}

	chainName := ctx.Args().Get(1)
	if chainName == "" {
		utils.Fatalf("This command requires chain name specified.")
	}

	stack, cfg := makeConfigNode(ctx, chainName)
	utils.RegisterIntService(stack, &cfg.Eth, ctx, GetCMInstance(ctx).cch)
	defer stack.Close()

	chainDb := utils.MakeChainDatabase(ctx, stack)

	height, _ := strconv.ParseUint(ctx.Args().First(), 10, 64)

	blockhash := rawdb.ReadCanonicalHash(chainDb, height)
	block := rawdb.ReadBlock(chainDb, blockhash, height)
	bsize := block.Size()

	root := block.Header().Root
	statedb, _ := state.New(block.Root(), state.NewDatabase(chainDb))
	accountTrie, _ := statedb.Database().OpenTrie(root)

	count := CountSize{}
	countTrie(chainDb, accountTrie, &count, func(addr common.Address, account state.Account) {
		if account.Root != emptyRoot {
			storageTrie, _ := statedb.Database().OpenStorageTrie(common.Hash{}, account.Root)
			countTrie(chainDb, storageTrie, &count, nil)
		}

		if account.TX1Root != emptyRoot {
			tx1Trie, _ := statedb.Database().OpenTX1Trie(common.Hash{}, account.TX1Root)
			countTrie(chainDb, tx1Trie, &count, nil)
		}

		if account.TX3Root != emptyRoot {
			tx3Trie, _ := statedb.Database().OpenTX3Trie(common.Hash{}, account.TX3Root)
			countTrie(chainDb, tx3Trie, &count, nil)
		}

		if account.ProxiedRoot != emptyRoot {
			proxiedTrie, _ := statedb.Database().OpenProxiedTrie(common.Hash{}, account.ProxiedRoot)
			countTrie(chainDb, proxiedTrie, &count, nil)
		}

		if account.RewardRoot != emptyRoot {
			rewardTrie, _ := statedb.Database().OpenRewardTrie(common.Hash{}, account.RewardRoot)
			countTrie(chainDb, rewardTrie, &count, nil)
		}
	})

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile("blockstate_nodedump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	// Write Node Data into file
	for _, data := range count.Data {
		fh.WriteString(data.key + " " + data.value + "\n")
	}

	fmt.Printf("Block %d, block size %v, state node %v, state size %v\n", height, bsize, count.Totalnode, count.Totalnodevaluesize)
	return nil
}

type CountSize struct {
	Totalnodevaluesize, Totalnode int
	Data                          []nodeData
}

type nodeData struct {
	key, value string
}

type processLeafTrie func(addr common.Address, account state.Account)

var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

func countTrie(db neatdb.Database, t state.Trie, count *CountSize, processLeaf processLeafTrie) {
	for it := t.NodeIterator(nil); it.Next(true); {
		if !it.Leaf() {
			// non leaf node -> count += value
			node, _ := db.Get(it.Hash().Bytes())
			count.Totalnodevaluesize += len(node)
			count.Totalnode++
			count.Data = append(count.Data, nodeData{it.Hash().String(), common.Bytes2Hex(node)})
		} else {
			// Process the Account -> Inner Trie
			if processLeaf != nil {
				addr := t.GetKey(it.LeafKey())
				if len(addr) == 20 {
					var data state.Account
					rlp.DecodeBytes(it.LeafBlob(), &data)

					processLeaf(common.BytesToAddress(addr), data)
				}
			}
		}
	}
}

func version(ctx *cli.Context) error {
	fmt.Println("Chain:", clientIdentifier)
	fmt.Println("Version:", params.VersionWithMeta)
	if gitCommit != "" {
		fmt.Println("Git Commit:", gitCommit)
	}
	if gitDate != "" {
		fmt.Println("Git Commit Date:", gitDate)
	}
	fmt.Println("Architecture:", runtime.GOARCH)
	fmt.Println("Protocol Versions:", neatptc.ProtocolVersions)
	fmt.Println("Go Version:", runtime.Version())
	fmt.Println("Operating System:", runtime.GOOS)
	fmt.Printf("GOPATH=%s\n", os.Getenv("GOPATH"))
	fmt.Printf("GOROOT=%s\n", runtime.GOROOT())
	return nil
}
