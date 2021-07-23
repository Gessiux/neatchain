package types

import (
	"errors"
	"math/big"

	"github.com/Gessiux/neatchain/common"
	"github.com/Gessiux/neatchain/common/hexutil"
)

var (
	// TendermintDigest represents a hash of "Tendermint practical byzantine fault tolerance"
	// to identify whether the block is from Tendermint consensus engine
	TendermintDigest = common.HexToHash("0x63746963616c2062797a616e74696e65206661756c7420746f6c6572616e6365")

	TendermintExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
	TendermintExtraSeal   = 65 // Fixed number of extra-data bytes reserved for validator seal

	TendermintDefaultDifficulty = big.NewInt(1)
	TendermintNilUncleHash      = CalcUncleHash(nil) // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	TendermintEmptyNonce        = BlockNonce{}
	TendermintNonce             = hexutil.MustDecode("0x88ff88ff88ff88ff") // Magic nonce number to vote on adding a new validator

	MagicExtra = []byte("neatchain_tmp_extra")

	// ErrInvalidTendermintHeaderExtra is returned if the length of extra-data is less than 32 bytes
	ErrInvalidTendermintHeaderExtra = errors.New("invalid ipbft header extra-data")
)

// TendermintFilteredHeader returns a filtered header which some information (like seal, committed seals)
// are clean to fulfill the Tendermint hash rules. It returns nil if the extra-data cannot be
// decoded/encoded by rlp.
func TendermintFilteredHeader(h *Header, keepSeal bool) *Header {
	newHeader := CopyHeader(h)
	payload := MagicExtra
	newHeader.Extra = payload
	return newHeader
}
