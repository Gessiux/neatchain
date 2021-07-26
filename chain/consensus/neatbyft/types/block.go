package types

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	. "github.com/Gessiux/go-common"
	"github.com/Gessiux/go-crypto"
	"github.com/Gessiux/go-merkle"
	"github.com/Gessiux/go-wire"
	"github.com/Gessiux/neatchain/chain/core/state"
	"github.com/Gessiux/neatchain/chain/core/types"
	"github.com/Gessiux/neatchain/chain/log"
	"github.com/Gessiux/neatchain/utilities/rlp"
)

const MaxBlockSize = 22020096 // 21MB TODO make it configurable

// IntermediateBlockResult represents intermediate block execute result.
type IntermediateBlockResult struct {
	Block *types.Block
	// followed by block execute result
	State    *state.StateDB
	Receipts types.Receipts
	Ops      *types.PendingOps
}

type NCBlock struct {
	Block              *types.Block             `json:"block"`
	NCExtra            *NeatConExtra            `json:"tdmexdata"`
	TX3ProofData       []*types.TX3ProofData    `json:"tx3proofdata"`
	IntermediateResult *IntermediateBlockResult `json:"-"`
}

func MakeBlock(height uint64, chainID string, commit *Commit,
	block *types.Block, valHash []byte, epochNumber uint64, epochBytes []byte, tx3ProofData []*types.TX3ProofData, partSize int) (*NCBlock, *PartSet) {
	NCExtra := &NeatConExtra{
		ChainID:        chainID,
		Height:         uint64(height),
		Time:           time.Now(),
		EpochNumber:    epochNumber,
		ValidatorsHash: valHash,
		SeenCommit:     commit,
		EpochBytes:     epochBytes,
	}

	ncBlock := &NCBlock{
		Block:        block,
		NCExtra:      NCExtra,
		TX3ProofData: tx3ProofData,
	}
	return ncBlock, ncBlock.MakePartSet(partSize)
}

// Basic validation that doesn't involve state data.
func (b *NCBlock) ValidateBasic(ncExtra *NeatConExtra) error {

	if b.NCExtra.ChainID != ncExtra.ChainID {
		return errors.New(Fmt("Wrong Block.Header.ChainID. Expected %v, got %v", ncExtra.ChainID, b.NCExtra.ChainID))
	}
	if b.NCExtra.Height != ncExtra.Height+1 {
		return errors.New(Fmt("Wrong Block.Header.Height. Expected %v, got %v", ncExtra.Height+1, b.NCExtra.Height))
	}

	/*
		if !b.NCExtra.BlockID.Equals(blockID) {
			return errors.New(Fmt("Wrong Block.Header.LastBlockID.  Expected %v, got %v", blockID, b.NCExtra.BlockID))
		}
		if !bytes.Equal(b.NCExtra.SeenCommitHash, b.NCExtra.SeenCommit.Hash()) {
			return errors.New(Fmt("Wrong Block.Header.LastCommitHash.  Expected %X, got %X", b.NCExtra.SeenCommitHash, b.NCExtra.SeenCommit.Hash()))
		}
		if b.NCExtra.Height != 1 {
			if err := b.NCExtra.SeenCommit.ValidateBasic(); err != nil {
				return err
			}
		}
	*/
	return nil
}

func (b *NCBlock) FillSeenCommitHash() {
	if b.NCExtra.SeenCommitHash == nil {
		b.NCExtra.SeenCommitHash = b.NCExtra.SeenCommit.Hash()
	}
}

// Computes and returns the block hash.
// If the block is incomplete, block hash is nil for safety.
func (b *NCBlock) Hash() []byte {
	// fmt.Println(">>", b.Data)
	if b == nil || b.NCExtra.SeenCommit == nil {
		return nil
	}
	b.FillSeenCommitHash()
	return b.NCExtra.Hash()
}

func (b *NCBlock) MakePartSet(partSize int) *PartSet {

	return NewPartSetFromData(b.ToBytes(), partSize)
}

func (b *NCBlock) ToBytes() []byte {

	type TmpBlock struct {
		BlockData    []byte
		NCExtra      *NeatConExtra
		TX3ProofData []*types.TX3ProofData
	}
	//fmt.Printf("NCBlock.toBytes 0 with block: %v\n", b)

	bs, err := rlp.EncodeToBytes(b.Block)
	if err != nil {
		log.Warnf("NCBlock.toBytes error\n")
	}
	bb := &TmpBlock{
		BlockData:    bs,
		NCExtra:      b.NCExtra,
		TX3ProofData: b.TX3ProofData,
	}

	ret := wire.BinaryBytes(bb)
	return ret
}

func (b *NCBlock) FromBytes(reader io.Reader) (*NCBlock, error) {

	type TmpBlock struct {
		BlockData    []byte
		NCExtra      *NeatConExtra
		TX3ProofData []*types.TX3ProofData
	}

	//fmt.Printf("NCBlock.FromBytes \n")

	var n int
	var err error
	bb := wire.ReadBinary(&TmpBlock{}, reader, MaxBlockSize, &n, &err).(*TmpBlock)
	if err != nil {
		log.Warnf("NCBlock.FromBytes 0 error: %v\n", err)
		return nil, err
	}

	var block types.Block
	err = rlp.DecodeBytes(bb.BlockData, &block)
	if err != nil {
		log.Warnf("NCBlock.FromBytes 1 error: %v\n", err)
		return nil, err
	}

	ncBlock := &NCBlock{
		Block:        &block,
		NCExtra:      bb.NCExtra,
		TX3ProofData: bb.TX3ProofData,
	}

	log.Debugf("NCBlock.FromBytes 2 with: %v\n", ncBlock)
	return ncBlock, nil
}

// Convenience.
// A nil block never hashes to anything.
// Nothing hashes to a nil hash.
func (b *NCBlock) HashesTo(hash []byte) bool {
	if len(hash) == 0 {
		return false
	}
	if b == nil {
		return false
	}
	return bytes.Equal(b.Hash(), hash)
}

func (b *NCBlock) String() string {
	return b.StringIndented("")
}

func (b *NCBlock) StringIndented(indent string) string {
	if b == nil {
		return "nil-Block"
	}

	return fmt.Sprintf(`Block{
%s  %v
%s  %v
%s  %v
%s}#%X`,
		indent, b.Block.String(),
		indent, b.NCExtra,
		indent, b.NCExtra.SeenCommit.StringIndented(indent+""),
		indent, b.Hash())
}

func (b *NCBlock) StringShort() string {
	if b == nil {
		return "nil-Block"
	} else {
		return fmt.Sprintf("Block#%X", b.Hash())
	}
}

//-------------------------------------

// NOTE: Commit is empty for height 1, but never nil.
type Commit struct {
	// NOTE: The Precommits are in order of address to preserve the bonded ValidatorSet order.
	// Any peer with a block can gossip precommits by index with a peer without recalculating the
	// active ValidatorSet.
	BlockID BlockID `json:"blockID"`
	Height  uint64  `json:"height"`
	Round   int     `json:"round"`

	// BLS signature aggregation to be added here
	SignAggr crypto.BLSSignature `json:"SignAggr"`
	BitArray *BitArray

	// Volatile
	hash []byte
}

func (commit *Commit) Type() byte {
	return VoteTypePrecommit
}

func (commit *Commit) Size() int {
	return (int)(commit.BitArray.Size())
}

func (commit *Commit) NumCommits() int {
	return (int)(commit.BitArray.NumBitsSet())
}

func (commit *Commit) ValidateBasic() error {
	if commit.BlockID.IsZero() {
		return errors.New("Commit cannot be for nil block")
	}
	/*
		if commit.Type() != VoteTypePrecommit {
			return fmt.Errorf("Invalid commit type. Expected VoteTypePrecommit, got %v",
				precommit.Type)
		}

		// shall we validate the signature aggregation?
	*/

	return nil
}

func (commit *Commit) Hash() []byte {
	if commit.hash == nil {
		hash := merkle.SimpleHashFromBinary(*commit)
		commit.hash = hash
	}
	return commit.hash
}

func (commit *Commit) StringIndented(indent string) string {
	if commit == nil {
		return "nil-Commit"
	}
	return fmt.Sprintf(`Commit{
%s  BlockID:    %v
%s  Height:     %v
%s  Round:      %v
%s  Type:       %v
%s  BitArray:   %v
%s}#%X`,
		indent, commit.BlockID,
		indent, commit.Height,
		indent, commit.Round,
		indent, commit.Type(),
		indent, commit.BitArray.String(),
		indent, commit.hash)
}

//--------------------------------------------------------------------------------

type BlockID struct {
	Hash        []byte        `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}

func (blockID BlockID) IsZero() bool {
	return len(blockID.Hash) == 0 && blockID.PartsHeader.IsZero()
}

func (blockID BlockID) Equals(other BlockID) bool {
	return bytes.Equal(blockID.Hash, other.Hash) &&
		blockID.PartsHeader.Equals(other.PartsHeader)
}

func (blockID BlockID) Key() string {
	return string(blockID.Hash) + string(wire.BinaryBytes(blockID.PartsHeader))
}

func (blockID BlockID) WriteSignBytes(w io.Writer, n *int, err *error) {
	if blockID.IsZero() {
		wire.WriteTo([]byte("null"), w, n, err)
	} else {
		wire.WriteJSON(CanonicalBlockID(blockID), w, n, err)
	}

}

func (blockID BlockID) String() string {
	return fmt.Sprintf(`%X:%v`, blockID.Hash, blockID.PartsHeader)
}
