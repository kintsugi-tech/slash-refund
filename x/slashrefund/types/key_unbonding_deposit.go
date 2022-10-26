package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ binary.ByteOrder

// UnbondingQueueKey defined as staking module does: empty byte slice
var UnbondingQueueKey = []byte{0x41} // prefix for the timestamps in unbonding queue

const (
	// UnbondingDepositKeyPrefix is the prefix to retrieve all UnbondingDeposit
	UnbondingDepositKeyPrefix = "UnbondingDeposit/value/"
)

// UnbondingDepositKey returns the store key to retrieve a UnbondingDeposit from the index fields
func UnbondingDepositKey(
	depositorAddress string,
	validatorAddress string,
) []byte {
	var key []byte

	depositorAddressBytes := []byte(depositorAddress)
	key = append(key, depositorAddressBytes...)
	key = append(key, []byte("/")...)

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}

// GetUnbondingDepositTimeKey creates the prefix for all unbonding deposits from a delegator
func GetUnbondingDepositTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}
