package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// GetDeposit returns a deposit from its index: depAddr & valAddr
func (k Keeper) GetDeposit(ctx sdk.Context, depAddr sdk.AccAddress, valAddr sdk.ValAddress) (deposit types.Deposit, found bool) {
	moduleStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(moduleStore, types.KeyPrefix(types.DepositKeyPrefix))
	key := types.DepositKey(depAddr, valAddr)
	b := store.Get(key)
	if b == nil {
		return deposit, false
	}

	k.cdc.MustUnmarshal(b, &deposit)
	return deposit, true
}

// SetDeposit set a specific deposit in the store from its index
func (k Keeper) SetDeposit(ctx sdk.Context, deposit types.Deposit) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	b := k.cdc.MustMarshal(&deposit)
	store.Set(types.DepositKey(
		deposit.MustGetDepositorAddr(),
		deposit.MustGetValidatorAddr(),
	), b)
}

// RemoveDeposit removes a deposit from the store
func (k Keeper) RemoveDeposit(
	ctx sdk.Context,
	deposit types.Deposit,
) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DepositKeyPrefix))
	store.Delete(types.DepositKey(
		deposit.MustGetDepositorAddr(),
		deposit.MustGetValidatorAddr(),
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

// GetDeposits of specific validator
func (k Keeper) GetValidatorDeposits(ctx sdk.Context, valAddr sdk.ValAddress) (deposits []types.Deposit) {

	for _, deposit := range k.GetAllDeposit(ctx) {
		if deposit.ValidatorAddress == valAddr.String() {
			deposits = append(deposits, deposit)
		}
	}
	return deposits
}

// Deposit implements the state transition logic for a deposit. It checks if a pool for the
// validator exists or create it. Sended tokens are deposited into the module and shares are
// created for the depositor to take into account its balance inside the module.
func (k Keeper) Deposit(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	depCoin sdk.Coin,
	validator stakingtypes.Validator,
) (newShares sdk.Dec, err error) {

	// Check if a validator has zero token but shares. This situation can arise due to slashing
	// of the considered validator. Deposit for these validators are allowed.
	if validator.InvalidExRate() {
		// Return zero shares and an error
		return sdk.ZeroDec(), types.ErrDepositorShareExRateInvalid
	}

	// Operator address of the validator
	valOperAddr := validator.GetOperator()

	// Check if the deposit pool exists or create it
	var deposit types.Deposit
	depPool, found := k.GetDepositPool(ctx, valOperAddr)
	if !found {
		depPool = types.NewDepositPool(
			valOperAddr,
			sdk.NewCoin(depCoin.Denom, sdk.ZeroInt()),
			sdk.ZeroDec(),
		)

		// If the pool does not exists no deposit can exists.
		deposit = types.NewDeposit(depAddr, valOperAddr, sdk.ZeroDec())
	} else {
		// Check if the deposit exists or create it
		deposit, found = k.GetDeposit(ctx, depAddr, valOperAddr)
		if !found {
			// If a previous deposit does not exist initialize one with zero shares
			deposit = types.NewDeposit(depAddr, valOperAddr, sdk.ZeroDec())
		}
	}

	// Send deposited tokens to the slashrefund module
	coins := sdk.NewCoins(sdk.NewCoin(depCoin.Denom, depCoin.Amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depAddr, types.ModuleName, coins); err != nil {
		return sdk.Dec{}, err
	}

	// Deposited tokens are treated as pool shares, similarly to the staking module.
	newShares = k.AddDepPoolTokensAndShares(ctx, depPool, depCoin)

	deposit.Shares = deposit.Shares.Add(newShares)
	k.SetDeposit(ctx, deposit)

	return sdk.NewDec(depCoin.Amount.Int64()), nil
}
