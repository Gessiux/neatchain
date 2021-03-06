package state

import (
	"errors"
	"fmt"

	"github.com/Gessiux/neatchain/chain/consensus"
	ep "github.com/Gessiux/neatchain/chain/consensus/neatcon/epoch"
	"github.com/Gessiux/neatchain/chain/consensus/neatcon/types"
	"github.com/Gessiux/neatchain/chain/core"
	neatTypes "github.com/Gessiux/neatchain/chain/core/types"
)

//--------------------------------------------------

// return a bit array of validators that signed the last commit
// NOTE: assumes commits have already been authenticated
/*
func commitBitArrayFromBlock(block *types.NCBlock) *BitArray {

	signed := NewBitArray(uint64(len(block.NTCExtra.SeenCommit.Precommits)))
	for i, precommit := range block.NTCExtra.SeenCommit.Precommits {
		if precommit != nil {
			signed.SetIndex(uint64(i), true) // val_.LastCommitHeight = block.Height - 1
		}
	}
	return signed
}*/

//-----------------------------------------------------
// Validate block

func (s *State) ValidateBlock(block *types.NCBlock) error {
	return s.validateBlock(block)
}

//Very current block
func (s *State) validateBlock(block *types.NCBlock) error {
	// Basic block validation.
	err := block.ValidateBasic(s.NTCExtra)
	if err != nil {
		return err
	}

	// Validate block SeenCommit.
	epoch := s.Epoch.GetEpochByBlockNumber(block.NTCExtra.Height)
	if epoch == nil || epoch.Validators == nil {
		return errors.New("no epoch for current block height")
	}

	valSet := epoch.Validators
	err = valSet.VerifyCommit(block.NTCExtra.ChainID, block.NTCExtra.Height,
		block.NTCExtra.SeenCommit)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------

func init() {
	core.RegisterInsertBlockCb("UpdateLocalEpoch", updateLocalEpoch)
	core.RegisterInsertBlockCb("AutoStartMining", autoStartMining)
}

func updateLocalEpoch(bc *core.BlockChain, block *neatTypes.Block) {
	if block.NumberU64() == 0 {
		return
	}

	ncExtra, _ := types.ExtractNeatConExtra(block.Header())
	//here handles the proposed next epoch
	epochInBlock := ep.FromBytes(ncExtra.EpochBytes)

	eng := bc.Engine().(consensus.NeatCon)
	currentEpoch := eng.GetEpoch()

	if epochInBlock != nil {
		if epochInBlock.Number == currentEpoch.Number+1 {
			//fmt.Printf("update local epoch 1\n")
			// Save the next epoch
			if block.NumberU64() == currentEpoch.StartBlock+1 || block.NumberU64() == 2 {
				//fmt.Printf("update local epoch block number %v, current epoch start block %v\n", block.NumberU64(), currentEpoch.StartBlock)
				// Propose next epoch
				//epochInBlock.SetEpochValidatorVoteSet(ep.NewEpochValidatorVoteSet())
				epochInBlock.Status = ep.EPOCH_VOTED_NOT_SAVED
				epochInBlock.SetRewardScheme(currentEpoch.GetRewardScheme())
				currentEpoch.SetNextEpoch(epochInBlock)
			} else if block.NumberU64() == currentEpoch.EndBlock {
				//fmt.Printf("update local epoch 2\n")
				// Finalize next epoch
				// Validator set in next epoch will not finalize and send to mainchain
				nextEp := currentEpoch.GetNextEpoch()
				nextEp.Validators = epochInBlock.Validators
				nextEp.Status = ep.EPOCH_VOTED_NOT_SAVED
			}
			currentEpoch.Save()
		} else if epochInBlock.Number == currentEpoch.Number {
			//fmt.Printf("update local epoch 3\n")
			// Update the current epoch Start Time from proposer
			currentEpoch.StartTime = epochInBlock.StartTime
			currentEpoch.Save()

			// Update the previous epoch End Time
			if currentEpoch.Number > 0 {
				currentEpoch.GetPreviousEpoch().EndTime = epochInBlock.StartTime // update previous epoch end time to the state to avoid different node epoch end time mismatch
				ep.UpdateEpochEndTime(currentEpoch.GetDB(), currentEpoch.Number-1, epochInBlock.StartTime)
			}
		}
	}
}

func autoStartMining(bc *core.BlockChain, block *neatTypes.Block) {
	eng := bc.Engine().(consensus.NeatCon)
	currentEpoch := eng.GetEpoch()

	// At one block before epoch end block, we should able to calculate the new validator
	if block.NumberU64() == currentEpoch.EndBlock-1 {
		fmt.Printf("auto start mining first %v\n", block.Number())
		// Re-Calculate the next epoch validators
		nextEp := currentEpoch.GetNextEpoch()
		state, _ := bc.State()
		nextValidators := currentEpoch.Validators.Copy()
		dryrunErr := ep.DryRunUpdateEpochValidatorSet(state, nextValidators, nextEp.GetEpochValidatorVoteSet())
		if dryrunErr != nil {
			panic("can not update the validator set base on the vote, error: " + dryrunErr.Error())
		}
		nextEp.Validators = nextValidators

		if nextValidators.HasAddress(eng.PrivateValidator().Bytes()) && !eng.IsStarted() {
			fmt.Printf("auto start mining first, post start mining event")
			bc.PostChainEvents([]interface{}{core.StartMiningEvent{}}, nil)
		}
	}
}
