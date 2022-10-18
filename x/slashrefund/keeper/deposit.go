package keeper

import (
	//"fmt"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// GetDeposit returns a deposit from its index: depAddr & valAddr
func (k Keeper) GetDeposit(ctx sdk.Context, depAddr sdk.AccAddress, valAddr sdk.ValAddress) (deposit types.Deposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	keys := types.DepositKey(depAddr,valAddr)
	b := store.Get(keys)
	if b == nil {
		return deposit, false
	}

	k.cdc.MustUnmarshal(b, &deposit)
	return deposit, true
}

// SetDeposit set a specific deposit in the store from its index
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {

	depositorAddress := sdk.MustAccAddressFromBech32(deposit.DepositorAddress)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	b := k.cdc.MustMarshal(&deposit)
	store.Set(types.DepositKey(
		depositorAddress,
		deposit.GetValidatorAddr(),
	), b)
}

// RemoveDeposit removes a deposit from the store
func (k Keeper) RemoveDeposit(
	ctx sdk.Context,
	deposit types.Deposit,
) {

	depositorAddress := sdk.MustAccAddressFromBech32(deposit.DepositorAddress)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	store.Delete(types.DepositKey(
		depositorAddress,
		deposit.GetValidatorAddr(),
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

// CUSTOM IMPLEMENTATIONS
// GetDeposits of specific validator
func (k Keeper) GetDepositOfValidator(ctx sdk.Context, valAddr sdk.ValAddress) (list []types.Deposit, total sdk.Coin) {
	// TODO: Handle secondary tokens ie. stATOM
	deposits := k.GetAllDeposit(ctx)

	var valDeposits []types.Deposit

	totalDeposit := sdk.NewInt(0)

	for _, deposit := range deposits {

		if deposit.ValidatorAddress == valAddr.String() {
			valDeposits = append(valDeposits, deposit)
			// TODO: fix math.NewInt with proper logic
			totalDeposit = totalDeposit.Add(math.NewInt(1))
		}
	}

	return valDeposits, sdk.NewCoin("stake", totalDeposit)
}

// Deposit implements the state transition logic for a deposit
// TODO: controllare hook: logiche da eseguire se deposito viene creato o modificato.
// TODO: definire i diversi ModuleAccount account a cui mandare i token
func (k Keeper) Deposit(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	depCoin sdk.Coin,
	validator stakingtypes.Validator,
) (newShares sdk.Dec, err error) {
	//logger := k.Logger(ctx)

	// Check if a validator has zero token but shares.
	if validator.InvalidExRate() {
		return sdk.ZeroDec(), types.ErrDepositorShareExRateInvalid
	}

	// Get, if exists, a previous deposit
	deposit, found := k.GetDeposit(ctx, depAddr, validator.GetOperator())
	if !found {
		// If a previous deposit does not exist initialize one with zero shares
		deposit = types.NewDeposit(depAddr, validator.GetOperator(), sdk.ZeroDec())
	}

	coins := sdk.NewCoins(sdk.NewCoin(depCoin.Denom, depCoin.Amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depAddr, types.ModuleName, coins); err != nil {
		return sdk.Dec{}, err
	}

	_, newShares = k.AddValidatorTokensAndShares(ctx, validator, depCoin)
	/*
			balance := msg.Amount
		if isFound {
			balance = balance.AddAmount(deposit.Balance.Amount)
		}
	*/

	k.SetDeposit(ctx, deposit)

	return sdk.NewDec(depCoin.Amount.Int64()), nil

}
