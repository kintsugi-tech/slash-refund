package types

import (
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Append an entry to the unbonding deposit entry list.
func (ubd *UnbondingDeposit) AddEntry(creationHeight int64, minTime time.Time, balance math.Int) {
	// Check if another entry with equal time and creation height exists.
	entryIndex := -1
	for index, ubdEntry := range ubd.Entries {
		if ubdEntry.CreationHeight == creationHeight && ubdEntry.CompletionTime.Equal(minTime) {
			entryIndex = index
			break
		}
	}
	// If the eentry exists, update the balance
	if entryIndex != -1 {
		ubdEntry := ubd.Entries[entryIndex]
		ubdEntry.Balance = ubdEntry.Balance.Add(balance)
		ubdEntry.InitialBalance = ubdEntry.InitialBalance.Add(balance)

		ubd.Entries[entryIndex] = ubdEntry
	} else {
		// Append the new unbonding deposit entry
		entry := NewUnbondingDepositEntry(creationHeight, minTime, balance)
		ubd.Entries = append(ubd.Entries, entry)
	}
}

func NewUnbondingDepositEntry(creationHeight int64, completionTime time.Time, balance math.Int) UnbondingDepositEntry {
	return UnbondingDepositEntry{
		CreationHeight: creationHeight,
		CompletionTime: completionTime,
		InitialBalance: balance,
		Balance:        balance,
	}
}

func NewUnbondingDeposit(
	depositorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	creationHeight int64, minTime time.Time, balance math.Int,
) UnbondingDeposit {
	return UnbondingDeposit{
		DepositorAddress: depositorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
		Entries: []UnbondingDepositEntry{
			NewUnbondingDepositEntry(creationHeight, minTime, balance),
		},
	}
}

// IsMature - is the current entry mature
func (e UnbondingDepositEntry) IsMature(currentTime time.Time) bool {
	return !e.CompletionTime.After(currentTime)
}

// RemoveEntry - remove entry at index i to the unbonding delegation
func (ubd *UnbondingDeposit) RemoveEntry(i int64) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

// unmarshal a unbonding delegation from a store value
func UnmarshalUBD(cdc codec.BinaryCodec, value []byte) (ubd UnbondingDeposit, err error) {
	err = cdc.Unmarshal(value, &ubd)
	return ubd, err
}

// unmarshal a unbonding delegation from a store value
func MustUnmarshalUBD(cdc codec.BinaryCodec, value []byte) UnbondingDeposit {
	ubd, err := UnmarshalUBD(cdc, value)
	if err != nil {
		panic(err)
	}

	return ubd
}
