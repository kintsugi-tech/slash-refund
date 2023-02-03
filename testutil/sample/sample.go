package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AccAddress returns a sample account address
func AccAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address()).String()
}

// AccAddress returns a sample validator address
func ValAddress() string {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.ValAddress(sdk.AccAddress(pk.Address())).String()
}

// MockAddress returns "test________________", a string that can be used as account or validator test address.
func MockAddress() string {
	return "test________________"
}

// MockAddress returns "test_______________1", a string that can be used as account or validator test address.
func MockAddress1() string {
	return "test_______________1"
}

// MockAddress returns "test_______________2", a string that can be used as account or validator test address.
func MockAddress2() string {
	return "test_______________2"
}
