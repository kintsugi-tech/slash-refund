package types

import (
	"encoding/binary"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	kv "github.com/cosmos/cosmos-sdk/types/kv"
)

var _ binary.ByteOrder

var UnbondingDepositsKeyPrefix = []byte{0x32}          // prefix for each key for an unbonding-deposit
var UnbondingDepositByValIndexKeyPrefix = []byte{0x33} // prefix for each key for an unbonding-deposit, by validator operator
var UnbondingQueueKey = []byte{0x41}                   // key for the timestamps in unbonding queue

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
