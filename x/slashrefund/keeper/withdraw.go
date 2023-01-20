package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Withdraw implements the state transition logic associated to a valid amount of tokens that a user
// wants to withdraw from the module.
func (k Keeper) Withdraw(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	denom string,
	witShares sdk.Dec,
) (sdk.Coin, time.Time, error) {

	witAmt, err := k.Unbond(ctx, depAddr, valAddr, witShares)
	if err != nil {
		return sdk.NewCoin(denom, sdk.NewInt(0)), time.Time{}, err
	}

	// Time at which the withdrawn tokens become available.
	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))

	ubd := k.SetUnbondingDepositEntry(ctx, depAddr, valAddr, ctx.BlockHeight(), completionTime, witAmt)

	k.InsertUBDQueue(ctx, ubd, completionTime)

	return sdk.NewCoin(k.AllowedTokens(ctx)[0], witAmt), completionTime, nil
}

// Returns user's shares associated with desired withdrawal tokens if available,
// or an error.
func (k Keeper) ComputeAssociatedShares(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	tokens sdk.Coin,
) (shares sdk.Dec, err error) {



	deposit, found := k.GetDeposit(ctx, depAddr, valAddr)
	if !found {
		return sdk.NewDec(0), types.ErrNoDepositForAddress
	}

	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return sdk.NewDec(0), types.ErrNoDepositPoolForValidator
	}

	// Compute shares from desired withdrawal tokens.
	shares, err = depPool.SharesFromTokens(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// Compute rounded down shares from desired withdrawal tokens.
	sharesTruncated, err := depPool.SharesFromTokensTruncated(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// Check if desired withdrawal tokens converted to truncated shares are greater than actual
	// total of depositor shares.
	depositorShares := deposit.GetShares()
	if sharesTruncated.GT(depositorShares) {
		return sdk.NewDec(0), sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token amount")
	}

	// Cap shares (not-truncated) at total depositor shares.
	if shares.GT(depositorShares) {
		shares = depositorShares
	}

	return shares, nil
}

func (k Keeper) Unbond(
	ctx sdk.Context,
	delAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	shares sdk.Dec,
) (issuedTokensAmt sdk.Int, err error) {

	// Check if the deposit exists in the store.
	deposit, found := k.GetDeposit(ctx, delAddr, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositForAddress
	}

	// Check if deposit pool exists in the store.
	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositPoolForValidator
	}

	// Ensure that we have enough shares to remove.
	if deposit.Shares.LT(shares) {
		return issuedTokensAmt, sdkerrors.Wrap(types.ErrNotEnoughDepositShares, deposit.Shares.String())
	}

	// Subtract shares from deposit.
	deposit.Shares = deposit.Shares.Sub(shares)

	// Remove the deposit if zero or set a new deposit.
	if deposit.Shares.IsZero() {
		k.RemoveDeposit(ctx, deposit)
	} else {
		k.SetDeposit(ctx, deposit)
	}

	depPool, issuedTokensAmt = k.RemoveDepPoolTokensAndShares(ctx, depPool, shares)

	if depPool.Shares.IsZero() {
		k.RemoveDepositPool(ctx, valAddr)
	}

	return issuedTokensAmt, nil
}
