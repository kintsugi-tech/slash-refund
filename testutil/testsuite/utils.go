package testsuite

import (
	"math/rand"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bytes"
)

var (
	bondAmt = sdk.DefaultPowerReduction
	denom = sdk.DefaultBondDenom
	delegationMultiplier = sdk.NewDec(5)
	balanceMultiplier = math.NewInt(3)
)

// GenerateBalances returns a bank type balance structure list with a fixed amount of tokens
// assigned to addresses.
func GenerateBalances(addresses []sdk.AccAddress) []banktypes.Balance {
	balances := make([]banktypes.Balance, len(addresses))
	for i, addr := range(addresses) {
		balances[i] = banktypes.Balance{
			Address: addr.String(),
			Coins: sdk.NewCoins(sdk.NewCoin(denom, bondAmt.Mul(balanceMultiplier))),
		}
	}

	return balances
}

// GenerateNConsensusPubKeys returns a specified number of public keys using ed25519. These are the
// keys used by Tendermint for for the consensus key.
func GenerateNConsensusPubKeys(number int) []*codectypes.Any {
	pks:= simapp.CreateTestPubKeys(number)
	pksAny := make([]*codectypes.Any, number)
	for i := 0; i < number; i++ {
		pkAny, err := codectypes.NewAnyWithValue(pks[i])
		if err != nil {
			panic(err)
		}
		pksAny[i] = pkAny
	}
	return pksAny
}

// GenerateNAddresses returns a specified number of a addresses using secp256k1. 
func GenerateNAddresses(number int) []bytes.HexBytes {
	addresses := make([]bytes.HexBytes, number)
	for i := 0; i < number; i++ {
		pk := secp256k1.GenPrivKey().PubKey()
		addresses[i] = pk.Address()
	}

	return addresses 
}

// ConvertAddressesToValAddr convert a list of addresses generated with secp256k1 into validator
// addresses.
func ConvertAddressesToValAddr(addresses []bytes.HexBytes) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addresses))
	for i, addr := range(addresses) {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs 
}

// ConvertAddressesToAccAddr convert a list of addresses generated with secp256k1 into generic users
// addresses.
func ConvertAddressesToAccAddr(addresses []bytes.HexBytes) []sdk.AccAddress {
	accAddrs := make([]sdk.AccAddress, len(addresses))
	for i, addr := range(addresses) {
		accAddrs[i] = sdk.AccAddress(addr)
	}

	return accAddrs 
}

// GenerateRandomDelegations generates random delegations given a set of delegators and validators. 
// All delegations are equal to 3 times the tokens required for a unit of voting power.
func GenerateRandomDelegations(
	delegators []sdk.AccAddress, 
	validators []stakingtypes.Validator,
) ([]stakingtypes.Delegation, []stakingtypes.Validator) {

	delegations := make([]stakingtypes.Delegation, len(delegators))
	for i, del := range(delegators) {
		valIndex := rand.Intn(len(validators))
		delegations[i] = stakingtypes.Delegation{
			DelegatorAddress: del.String(),
			ValidatorAddress: validators[valIndex].OperatorAddress,
			Shares: sdk.NewDecFromInt(bondAmt).Mul(delegationMultiplier),
		}
		validators[valIndex].DelegatorShares = validators[valIndex].
			DelegatorShares.Add(sdk.NewDecFromInt(bondAmt).Mul(delegationMultiplier))
	}

	return delegations, validators
}

// GenerateValidator returns a basic validator with 4 
func GenerateValidator(
	valAddr sdk.ValAddress, 
	consKey *codectypes.Any,
) (stakingtypes.Validator, stakingtypes.Delegation) {

	zero := sdk.ZeroDec()
	shares := sdk.NewDecFromInt(bondAmt)

	// Each validator has tokens corresponding to 1 point of consensus power. Since the validator
	// operator will be the first delegator the shares are equal to the bondAmt.
	val := stakingtypes.Validator{
		OperatorAddress:   valAddr.String(),
		ConsensusPubkey:   consKey,
		Jailed:            false,
		Status:            stakingtypes.Bonded, // important
		Tokens:            bondAmt,
		DelegatorShares:   shares,
		Description:       stakingtypes.Description{},
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        stakingtypes.NewCommission(zero, zero, zero),
		MinSelfDelegation: sdk.ZeroInt(), // so it can be slashed without problem
	}

	del := stakingtypes.Delegation{
		DelegatorAddress: sdk.AccAddress(valAddr.Bytes()).String(),
		ValidatorAddress: valAddr.String(),
		Shares: shares,
	}

	return val, del
}