package keeper

import (
	"fmt"

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

// Deposit implements the state transition logic for a deposit
func (k Keeper) Deposit(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	depCoin sdk.Coin,
	validator stakingtypes.Validator,
) (newShares sdk.Dec, err error) {

	logger := k.Logger(ctx)
	logger.Error("creating/updating deposit:")

	// Check if a validator has zero token but shares. This situation can arise due to slashing
	// of the considered validator.
	if validator.InvalidExRate() {
		// Return zero shares and an error
		return sdk.ZeroDec(), types.ErrDepositorShareExRateInvalid
	}

	// Operator address of the validator
	valOperAddr := validator.GetOperator()

	// Check if the deposit exists or create it
	deposit, found := k.GetDeposit(ctx, depAddr, valOperAddr)
	if !found {
		// If a previous deposit does not exist initialize one with zero shares
		deposit = types.NewDeposit(depAddr, valOperAddr, sdk.ZeroDec())
	}

	// Check if the deposit pool exists or create it
	depPool, found := k.GetDepositPool(ctx, valOperAddr)
	if !found {
		// TODO: should be initialized with actual Coins allowed. Now the hp is of just one allowed token.
		depPool = types.NewDepositPool(
			valOperAddr,
			sdk.NewCoin(k.AllowedTokensList(ctx)[0], sdk.ZeroInt()),
			sdk.ZeroDec(),
		)
	}

	// Send the deposited tokens to the slashrefund module
	coins := sdk.NewCoins(sdk.NewCoin(depCoin.Denom, depCoin.Amount))
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, depAddr, types.ModuleName, coins); err != nil {
		return sdk.Dec{}, err
	}

	// Deposited tokens are treated as pool shares, similarly to the staking module.
	newShares = k.AddDepPoolTokensAndShares(ctx, depPool, depCoin)

	deposit.Shares = deposit.Shares.Add(newShares)
	k.SetDeposit(ctx, deposit)

	// logger
	logger.Error(fmt.Sprintf("  deposit pool: added tokens=%s%s , added shares=%s", depCoin.Amount.String(), depCoin.Denom, newShares.String()))
	logger.Error(fmt.Sprintf("  deposit: added shares=%s , depositor=%s , validator=%s", newShares.String(), deposit.DepositorAddress, deposit.ValidatorAddress))

	return sdk.NewDec(depCoin.Amount.Int64()), nil
}
