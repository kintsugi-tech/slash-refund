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
	bondAmt              = sdk.DefaultPowerReduction
	denom                = sdk.DefaultBondDenom
	delegationMultiplier = sdk.NewInt(5)
	balanceMultiplier    = math.NewInt(100)
)

type HeigthAndTime struct {
	CreationHeight int64
	CompletionTime time.Time
}

// GenerateBalances returns a bank type balance structure list with a fixed amount of tokens
// assigned to addresses.
func GenerateBalances(addresses []sdk.AccAddress) []banktypes.Balance {
	balances := make([]banktypes.Balance, len(addresses))
	for i, addr := range addresses {
		balances[i] = banktypes.Balance{
			Address: addr.String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(denom, bondAmt.Mul(balanceMultiplier))),
		}
	}

	return balances
}

// GenerateNConsensusPubKeys returns a specified number of public keys using ed25519. These are the
// keys used by Tendermint for for the consensus key.
func GenerateNConsensusPubKeys(number int) []*codectypes.Any {
	pks := simapp.CreateTestPubKeys(number)
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
	for i, addr := range addresses {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

// ConvertAddressesToAccAddr convert a list of addresses generated with secp256k1 into generic users
// addresses.
func ConvertAddressesToAccAddr(addresses []bytes.HexBytes) []sdk.AccAddress {
	accAddrs := make([]sdk.AccAddress, len(addresses))
	for i, addr := range addresses {
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
	for i, del := range delegators {
		valIndex := rand.Intn(len(validators))
		delegations[i] = stakingtypes.Delegation{
			DelegatorAddress: del.String(),
			ValidatorAddress: validators[valIndex].OperatorAddress,
			Shares:           sdk.NewDecFromInt(bondAmt).Mul(delegationMultiplier),
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
		MinSelfDelegation: sdk.ZeroInt(), // so it can be slashed without problems
	}

	del := stakingtypes.Delegation{
		DelegatorAddress: sdk.AccAddress(valAddr.Bytes()).String(),
		ValidatorAddress: valAddr.String(),
		Shares:           shares,
	}

	return val, del
}

// GenerateRandomDelegations generates random delegations given a set of delegators and validators.
// All delegations are equal to 3 times the tokens required for a unit of voting power.
func GenerateRandomDelegations(
	delegators []sdk.AccAddress,
	validators []stakingtypes.Validator,
) ([]stakingtypes.Delegation, []stakingtypes.Validator) {

	delAmt := bondAmt.Mul(delegationMultiplier)

	var delegations []stakingtypes.Delegation
	for _, del := range delegators {
		valIndex := rand.Intn(len(validators))
		delegations = append(delegations, stakingtypes.Delegation{
			DelegatorAddress: del.String(),
			ValidatorAddress: validators[valIndex].OperatorAddress,
			Shares:           sdk.NewDecFromInt(delAmt),
		})
		validators[valIndex].DelegatorShares = validators[valIndex].DelegatorShares.Add(sdk.NewDecFromInt(delAmt))
		validators[valIndex].Tokens = validators[valIndex].Tokens.Add(delAmt)
	}

	return delegations, validators
}

// Generates random delegations given a set of delegators and validators.
// All delegations are equal to 3 times the tokens required for a unit of voting power.
func GenerateRandomUnbondingDelegations(
	ubdelegators []sdk.AccAddress,
	validators []stakingtypes.Validator,
	heightAndTimes []HeigthAndTime,
) ([]stakingtypes.UnbondingDelegation, []stakingtypes.Validator) {

	delAmt := bondAmt.Mul(delegationMultiplier)
	nentries := len(heightAndTimes)
	if nentries > int(stakingtypes.DefaultMaxEntries) {
		nentries = int(stakingtypes.DefaultMaxEntries)
	}

	var ubdelegations []stakingtypes.UnbondingDelegation
	for _, ubdel := range ubdelegators {
		// Generate entries.
		var entries []stakingtypes.UnbondingDelegationEntry
		for i := 0; i < nentries; i++ {
			entries = append(entries,
				stakingtypes.NewUnbondingDelegationEntry(
					heightAndTimes[i].CreationHeight,
					heightAndTimes[i].CompletionTime,
					delAmt,
				),
			)
		}

		valIndex := rand.Intn(len(validators))
		ubdelegations = append(ubdelegations, stakingtypes.UnbondingDelegation{
			DelegatorAddress: ubdel.String(),
			ValidatorAddress: validators[valIndex].OperatorAddress,
			Entries:          entries,
		})
	}

	return ubdelegations, validators
}

// Generates random redelegations given a set of delegators and validators.
// All delegations are equal to 3 times the tokens required for a unit of voting power.
// For each redelegator, a new redelegation will be created from source validator
// to a random destination validator among the dstValidators set.
// No duplicated redelegations can be created, since each (delegator,dstVal) couple
// is unique.
// For each redelegation, a delegation linked to this will be created.
// This function check the delegationsToUpdate given in input, and if it finds an
// already existant delegation among these, then it updates this delegation. If
// no delegation has to be updated then it creates a new delegation.
func GenerateRandomRedelegationsFromValidator(
	redelegators []sdk.AccAddress,
	srcValidator stakingtypes.Validator,
	dstValidators []stakingtypes.Validator,
	heightAndTimes []HeigthAndTime,
) ([]stakingtypes.Redelegation, []stakingtypes.Delegation, []stakingtypes.Validator) {

	redelAmt := bondAmt.Mul(delegationMultiplier)
	nentries := len(heightAndTimes)
	if nentries > int(stakingtypes.DefaultMaxEntries) {
		nentries = int(stakingtypes.DefaultMaxEntries)
	}

	// Generate redelegations.
	var redelegations []stakingtypes.Redelegation
	var delegations []stakingtypes.Delegation
	dstVals := dstValidators
	for _, del := range redelegators {
		// Generate entries.
		var entries []stakingtypes.RedelegationEntry
		tokensToAdd := sdk.ZeroInt()
		for i := 0; i < nentries; i++ {
			entries = append(entries,
				stakingtypes.NewRedelegationEntry(
					heightAndTimes[i].CreationHeight,
					heightAndTimes[i].CompletionTime,
					redelAmt,
					sdk.NewDecFromInt(redelAmt),
				),
			)
			tokensToAdd = tokensToAdd.Add(redelAmt)
		}
		// Generate redelegation.
		valIndex := rand.Intn(len(dstVals))
		redelegations = append(redelegations, stakingtypes.Redelegation{
			DelegatorAddress:    del.String(),
			ValidatorSrcAddress: srcValidator.OperatorAddress,
			ValidatorDstAddress: dstVals[valIndex].OperatorAddress,
			Entries:             entries,
		},
		)
		delegations = append(delegations, stakingtypes.Delegation{
			DelegatorAddress: del.String(),
			ValidatorAddress: dstVals[valIndex].OperatorAddress,
			Shares:           sdk.NewDecFromInt(tokensToAdd),
		})

		// Update destination validator.
		dstVals[valIndex].DelegatorShares = dstVals[valIndex].DelegatorShares.Add(sdk.NewDecFromInt(tokensToAdd))
		dstVals[valIndex].Tokens = dstVals[valIndex].Tokens.Add(tokensToAdd)
	}

	return redelegations, delegations, dstVals
}

// Generates random redelegations given a set of delegators and validators.
// All delegations are equal to 3 times the tokens required for a unit of voting power.
func GenerateRandomRedelegations(
	redelegators []sdk.AccAddress,
	validators []stakingtypes.Validator,
	heightAndTimes []HeigthAndTime,
) ([]stakingtypes.Redelegation, []stakingtypes.Delegation, []stakingtypes.Validator) {

	var redelegations []stakingtypes.Redelegation
	var delegations []stakingtypes.Delegation
	updated := validators

	for i := 0; i < len(validators); i++ {

		// Get source validator.
		srcVal := updated[i]

		// Remove srcVal from validators to obtain dstValidators.
		var dstValidators []stakingtypes.Validator
		dstValidators = append(dstValidators, validators[:i]...)
		dstValidators = append(dstValidators, validators[i+1:]...)

		// Generate redelegations from srcVal.
		redels, newDelegations, dstValidators := GenerateRandomRedelegationsFromValidator(redelegators, srcVal, dstValidators, heightAndTimes)
		redelegations = append(redelegations, redels...)
		delegations = append(delegations, newDelegations...)

		// Update validators.
		var vals []stakingtypes.Validator
		vals = append(vals, dstValidators[:i]...)
		vals = append(vals, srcVal)
		vals = append(vals, dstValidators[i:]...)
		updated = vals
	}

	return redelegations, delegations, updated
}

/*
	// Check if srcValidator is contained in dstValidators and if it is then remove it.
	var index int
	var found bool
	var dstVals []stakingtypes.Validator
	for i, dstVal := range dstValidators {
		if dstVal == srcValidator {
			index = i
			found = true
		}
	}
	if found {
		dstVals = append(dstValidators[:index], dstValidators[index+1:]...)

	} else {
		dstVals = dstValidators
	}

	// Re-insert srcVal if it was found among dstVals
	if found {
		dstVals = append(dstVals[:index], srcValidator)
		dstVals = append(dstVals, dstVals[index+1:]...)
	}
*/
