package keeper

import (
	//"fmt"
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) Withdraw(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	witShares sdk.Dec,
) (time.Time, error) {
	//logger := k.Logger(ctx)

	returnAmount, err := k.Unbond(ctx, depAddr, valAddr, witShares)
	if err != nil {
		return time.Time{}, err
	}

	return sdk.NewDec(depCoin.Amount.Int64()), nil
}

func (k Keeper) ValidateWithdrawdAmount(
	ctx sdk.Context, 
	depAddr sdk.AccAddress, 
	valAddr sdk.ValAddress, 
	tokens sdk.Coin,
) (shares sdk.Dec, err error) {

	dep, found := k.GetDeposit(ctx, depAddr, valAddr)
	if !found {
		return shares, types.ErrNoDepositForAddress
	}

	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return shares, types.ErrNoDepositPoolForValidator
	}

	shares, err = depPool.SharesFromTokens(tokens)
	if err != nil {
		return shares, err
	}

	sharesTruncated, err := depPool.SharesFromTokensTruncated(tokens)
	if err != nil {
		return shares, err
	}

	depShares := dep.GetShares()
	if sharesTruncated.GT(depShares) {
		return shares, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token amount")
	}

	if shares.GT(depShares) {
		shares = depShares
	}

	return shares, nil
}

func (k Keeper) Unbond(
	ctx sdk.Context, 
	delAddr sdk.AccAddress, 
	valAddr sdk.ValAddress, 
	shares sdk.Dec,
) (amount sdk.Int, err error) {
	// check if a delegation object exists in the store
	deposit, found := k.GetDeposit(ctx, delAddr, valAddr)
	if !found {
		return amount, types.ErrNoDepositForAddress
	}

	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return amount, types.ErrNoDepositPoolForValidator
	}

	// ensure that we have enough shares to remove
	if deposit.Shares.LT(shares) {
		return amount, sdkerrors.Wrap(types.ErrNotEnoughDepositShares, deposit.Shares.String())
	}

	// get validator
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return amount, stakingtypes.ErrNoValidatorFound
	}

	// subtract shares from delegation
	deposit.Shares = deposit.Shares.Sub(shares)

	// remove the delegation
	if deposit.Shares.IsZero() {
		k.RemoveDeposit(ctx, deposit)
	} else {
		k.SetDeposit(ctx, deposit)
	}

	_ = k.RemovePoolTokensAndShares(ctx, depPool, shares)

	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveValidator(ctx, validator.GetOperator())
	}

	return amount, nil
}