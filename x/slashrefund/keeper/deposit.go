package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// SetDeposit set a specific deposit in the store from its index
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	b := k.cdc.MustMarshal(&deposit)
	store.Set(types.DepositKey(
		deposit.Address,
		deposit.ValidatorAddress,
	), b)
}

// GetDeposit returns a deposit from its index
func (k Keeper) GetDeposit(
	ctx sdk.Context,
	address string,
	validatorAddress string,

) (val types.Deposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))

	b := store.Get(types.DepositKey(
		address,
		validatorAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveDeposit removes a deposit from the store
func (k Keeper) RemoveDeposit(
	ctx sdk.Context,
	address string,
	validatorAddress string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	store.Delete(types.DepositKey(
		address,
		validatorAddress,
	))
}

// GetAllDeposit returns all deposit
func (k Keeper) GetAllDeposit(ctx sdk.Context) (list []types.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Deposit
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetDeposits of specific validator (TODO: Handle secondary tokens ie. stATOM)
func (k Keeper) GetDepositOfValidator(ctx sdk.Context, valAddr sdk.ValAddress) (list []types.Deposit, total sdk.Coin) {

	deposits := k.GetAllDeposit(ctx)

	var valDeposits []types.Deposit

	totalDeposit := sdk.NewInt(0)

	for _, deposit := range deposits {

		if deposit.ValidatorAddress == valAddr.String() {

			valDeposits = append(valDeposits, deposit)
			totalDeposit = totalDeposit.Add(deposit.Balance.Amount)
		}
	}

	// TODO: Use staking module denom from param
	return valDeposits, sdk.NewCoin("stake", totalDeposit)
}

// Slash balance in proportional way to all the elements of a list
func (k Keeper) SlashDepositsProportionally(ctx sdk.Context, deposits []types.Deposit, totalRefund sdk.Int, totalDeposit sdk.Int) {

	// Examples
	// Validator total insurance deposit: 1000
	// Validator slash 500
	// Deposit 1 900: 900 * (slash/deposit) = 450 slash
	// Deposit 2 100: 100 * (slash/deposit) = 50 slash
	// --
	// Validator total insurance deposit: 500
	// Validator slash 500
	// Deposit 1 400: 400 * (slash/deposit) = 400 slash
	// Deposit 2 100: 100 * (slash/deposit) = 100 slash

	// Calculate ratio
	ratio := sdk.NewDecFromInt(totalRefund)
	ratio = ratio.QuoInt(totalDeposit)

	for _, deposit := range deposits {

		// if ratio = 1, simply delete all deposits 🔥
		if ratio.Equal(sdk.NewDec(1)) {
			k.RemoveDeposit(ctx, deposit.Address, deposit.ValidatorAddress)
			continue
		}

		k.Logger(ctx).Error("lol", deposit)
	}
}
