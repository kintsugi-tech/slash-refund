package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

var _ binary.ByteOrder

// UnbondingQueueKey defined as staking module does: empty byte slice
var UnbondingDepositsKey = []byte{0x32}
var UnbondingDepositByValIndexKey = []byte{0x33} // prefix for each key for an unbonding-deposit, by validator operator
var UnbondingQueueKey = []byte{0x41}             // prefix for the timestamps in unbonding queue

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

func UnbondingDepositKeyByValIndex(
	depositorAddress string,
	validatorAddress string,
) []byte {
	var key []byte

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	depositorAddressBytes := []byte(depositorAddress)
	key = append(key, depositorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}

// GetUnbondingDepositTimeKey creates the prefix for all unbonding deposits from a delegator
func GetUnbondingDepositTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}

// GetUBDsByValIndexKey creates the prefix keyspace for the indexes of unbonding deposits for a validator
func GetUBDsByValIndexKey(validatorAddress string) []byte {

	var key []byte

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}

// GetUBDsKey creates the prefix for all unbonding deposits from a delegator
func GetUBDsKey(depAddr sdk.AccAddress) []byte {
	return append(UnbondingDepositsKey, address.MustLengthPrefix(depAddr)...)
}

// GetUBDKey creates the key for an unbonding deposit by delegator and validator addr
// VALUE: staking/UnbondingDeposit
func GetUBDKey(depAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsKey(depAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}
