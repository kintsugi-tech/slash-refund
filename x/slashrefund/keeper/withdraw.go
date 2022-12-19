package keeper

import (
	"fmt"

	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) Withdraw(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	witShares sdk.Dec,
) (sdk.Coin, time.Time, error) {

	logger := k.Logger(ctx)
	logger.Error("creating unbonding deposit entry:")

	witAmt, err := k.Unbond(ctx, depAddr, valAddr, witShares)
	if err != nil {
		// TODO: change k.AllowedTokensList(ctx)[0] to handle different denoms
		return sdk.NewCoin(k.AllowedTokens(ctx)[0], sdk.NewInt(0)), time.Time{}, err
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))

	ubd := k.SetUnbondingDepositEntry(ctx, depAddr, valAddr, ctx.BlockHeight(), completionTime, witAmt)

	k.InsertUBDQueue(ctx, ubd, completionTime)

	// logger
	logger.Error(fmt.Sprintf("  new entry:"))
	logger.Error(fmt.Sprintf("    - initialBalance=%s", witAmt.String()))
	logger.Error(fmt.Sprintf("    - creationHeight=%s", fmt.Sprint(ctx.BlockHeight())))
	logger.Error(fmt.Sprintf("    - completionTime=%s", completionTime.String()))
	logger.Error(fmt.Sprintf("    - depositor=%s", ubd.DepositorAddress))
	logger.Error(fmt.Sprintf("    - validator=%s", ubd.ValidatorAddress))

	return sdk.NewCoin(k.AllowedTokens(ctx)[0], witAmt), completionTime, nil
}

func (k Keeper) ValidateWithdrawAmount(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	tokens sdk.Coin,
) (shares sdk.Dec, err error) {

	isValid, err := k.CheckAllowedTokens(ctx, tokens.Denom)
	if !isValid {
		return sdk.NewDec(0), err
	}

	if tokens.Amount.IsZero() {
		return sdk.NewDec(0), types.ErrZeroWithdraw
	}

	deposit, found := k.GetDeposit(ctx, depAddr, valAddr)
	if !found {
		return sdk.NewDec(0), types.ErrNoDepositForAddress
	}

	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return sdk.NewDec(0), types.ErrNoDepositPoolForValidator
	}

	// compute shares from wanted withdraw tokens
	shares, err = depPool.SharesFromTokens(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// compute shares from wanted withdraw tokens, rounded down
	sharesTruncated, err := depPool.SharesFromTokensTruncated(tokens)
	if err != nil {
		return sdk.NewDec(0), err
	}

	// check if wanted tokens converted to truncated shares are greater than actual total of depositor shares
	depositorShares := deposit.GetShares()
	if sharesTruncated.GT(depositorShares) {
		return sdk.NewDec(0), sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid token amount")
	}

	// cap shares (not-truncated) at total depositor shares
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

	// check if the deposit exists in the store
	deposit, found := k.GetDeposit(ctx, delAddr, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositForAddress
	}

	// check if deposit pool exists in the store
	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositPoolForValidator
	}

	// ensure that we have enough shares to remove
	if deposit.Shares.LT(shares) {
		return issuedTokensAmt, sdkerrors.Wrap(types.ErrNotEnoughDepositShares, deposit.Shares.String())
	}

	// subtract shares from deposit
	deposit.Shares = deposit.Shares.Sub(shares)

	// remove the deposit if zero or set a new deposit
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
