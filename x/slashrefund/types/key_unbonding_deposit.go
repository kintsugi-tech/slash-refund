package types

import "encoding/binary"

var _ binary.ByteOrder

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
