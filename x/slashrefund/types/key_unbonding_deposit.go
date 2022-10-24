package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// UnbondingDepositKeyPrefix is the prefix to retrieve all UnbondingDeposit
	UnbondingDepositKeyPrefix = "UnbondingDeposit/value/"
)

// UnbondingDepositKey returns the store key to retrieve a UnbondingDeposit from the index fields
func UnbondingDepositKey(
	delegatorAddress string,
	validatorAddress string,
) []byte {
	var key []byte

	delegatorAddressBytes := []byte(delegatorAddress)
	key = append(key, delegatorAddressBytes...)
	key = append(key, []byte("/")...)

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}
