package types

import (
	"time"
	
	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/types/address"
	kv "github.com/cosmos/cosmos-sdk/types/kv"
)

const (
	// ModuleName defines the module name
	ModuleName = "slashrefund"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_slashrefund"

	// DepositKeyPrefix is the prefix to retrieve all Deposit
	DepositKeyPrefix = "Deposit/value/"
	
	// RefundKeyPrefix is the prefix to retrieve all Refund
	RefundKeyPrefix = "Refund/value/"

	// DepositPoolKeyPrefix is the prefix to retrieve all DepositPool
	DepositPoolKeyPrefix = "DepositPool/value/"

	// RefundPoolKeyPrefix is the prefix to retrieve all RefundPool
	RefundPoolKeyPrefix = "RefundPool/value/"

)

var _ binary.ByteOrder

var UnbondingDepositsKeyPrefix = []byte{0x32}          // prefix for each key for an unbonding-deposit
var UnbondingDepositByValIndexKeyPrefix = []byte{0x33} // prefix for each key for an unbonding-deposit, by validator operator
var UnbondingQueueKey = []byte{0x41}                   // key for the timestamps in unbonding queue


func KeyPrefix(p string) []byte {
	return []byte(p)
}

// DepositKey returns the store key to retrieve a Deposit from the index fields
func DepositKey(
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
) []byte {
	var key []byte

	depAddrBytes := []byte(depAddr)
	key = append(key, depAddrBytes...)
	key = append(key, []byte("/")...)

	valAddrBytes := []byte(valAddr)
	key = append(key, valAddrBytes...)
	key = append(key, []byte("/")...)

	return key
}

// RefundKey returns the store key to retrieve a Refund from the index fields
func RefundKey(
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
) []byte {
	var key []byte

	delegatorBytes := []byte(delegator)
	key = append(key, delegatorBytes...)
	key = append(key, []byte("/")...)

	validatorBytes := []byte(validator)
	key = append(key, validatorBytes...)
	key = append(key, []byte("/")...)

	return key
}

// DepositPoolKey returns the store key to retrieve a DepositPool from the index fields
func DepositPoolKey(
	validatorAddress sdk.ValAddress,
) []byte {
	var key []byte

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}

// RefundPoolKey returns the store key to retrieve a RefundPool from the index fields
func RefundPoolKey(
	validatorAddress sdk.ValAddress,
) []byte {
	var key []byte

	validatorAddressBytes := []byte(validatorAddress)
	key = append(key, validatorAddressBytes...)
	key = append(key, []byte("/")...)

	return key
}

// GetUBDsKey creates the prefix for all unbonding deposits
func GetUBDsKeyPrefix() []byte {
	return UnbondingDepositsKeyPrefix
}

// GetUBDsKey creates the prefix for all unbonding deposits from a depositor
func GetUBDsKey(depAddr sdk.AccAddress) []byte {
	return append(GetUBDsKeyPrefix(), address.MustLengthPrefix(depAddr)...)
}

// GetUBDKey creates the key for an unbonding deposit by depositor and validator addr
func GetUBDKey(depAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetUBDsKey(depAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
}

// GetUBDByValIndexKey creates the index-key for an unbonding deposits, stored by validator-index.
// This will return empty bytes, because the key (validator-depositor) is used only to return to
// the actual key (depositor-validator) and get data from that.
func GetUBDByValIndexKey(valAddr sdk.ValAddress, depAddr sdk.AccAddress) []byte {
	return append(GetUBDsByValIndexKey(valAddr), address.MustLengthPrefix(depAddr)...)
}

// GetUBDsByValIndexKey creates the prefix keyspace for the indexes of unbonding deposits for a validator
func GetUBDsByValIndexKey(valAddr sdk.ValAddress) []byte {
	return append(UnbondingDepositByValIndexKeyPrefix, address.MustLengthPrefix(valAddr)...)
}

func GetUBDKeyFromValIndexKey(indexKey []byte) []byte {
	kv.AssertKeyAtLeastLength(indexKey, 2)
	addrs := indexKey[1:] // remove prefix bytes

	valAddrLen := addrs[0]
	kv.AssertKeyAtLeastLength(addrs, 2+int(valAddrLen))
	valAddr := addrs[1 : 1+valAddrLen]
	kv.AssertKeyAtLeastLength(addrs, 3+int(valAddrLen))
	depAddr := addrs[valAddrLen+2:]

	return GetUBDKey(depAddr, valAddr)
}

// GetUnbondingDepositTimeKey creates the prefix for all unbonding deposits from a delegator
func GetUnbondingDepositTimeKey(timestamp time.Time) []byte {
	bz := sdk.FormatTimeBytes(timestamp)
	return append(UnbondingQueueKey, bz...)
}
