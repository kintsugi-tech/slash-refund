package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) Withdraw(
	ctx sdk.Context,
	depAddr sdk.AccAddress,
	valAddr sdk.ValAddress,
	witShares sdk.Dec,
) (sdk.Coin, time.Time, error) {
	//logger := k.Logger(ctx)

	witAmt, err := k.Unbond(ctx, depAddr, valAddr, witShares)
	if err != nil {
		return sdk.NewCoin("", sdk.NewInt(0)), time.Time{}, err
	}

	completionTime := ctx.BlockHeader().Time.Add(k.stakingKeeper.UnbondingTime(ctx))

	ubd := k.SetUnbondingDepositEntry(ctx, depAddr, valAddr, ctx.BlockHeight(), completionTime, witAmt)

	k.InsertUBDQueue(ctx, ubd, completionTime)

	// TODO: change "stake"
	return sdk.NewCoin("stake", witAmt), completionTime, nil
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
		return shares, types.ErrZeroWithdraw
	}

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
) (issuedTokensAmt sdk.Int, err error) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("entered: Unbond"))

	// check if a deposit object exists in the store
	deposit, found := k.GetDeposit(ctx, delAddr, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositForAddress
	}

	depPool, found := k.GetDepositPool(ctx, valAddr)
	if !found {
		return issuedTokensAmt, types.ErrNoDepositPoolForValidator
	}

	logger.Error(fmt.Sprintf("    depPool.Shares: %s", depPool.Shares.String()))

	// ensure that we have enough shares to remove
	if deposit.Shares.LT(shares) {
		return issuedTokensAmt, sdkerrors.Wrap(types.ErrNotEnoughDepositShares, deposit.Shares.String())
	}

	// get validator
	// TODO: if a validator is no more active we have to send back tokens
	// TODO: remove validator if not used
	_, found = k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return issuedTokensAmt, stakingtypes.ErrNoValidatorFound
	}

	// subtract shares from deposit
	deposit.Shares = deposit.Shares.Sub(shares)

	// remove the deposit if zero or set a new doposit
	if deposit.Shares.IsZero() {
		k.RemoveDeposit(ctx, deposit)
	} else {
		k.SetDeposit(ctx, deposit)
	}

	logger.Error(fmt.Sprintf("    call k.RemovePoolTokensAndShares:"))

	depPool, issuedTokensAmt = k.RemovePoolTokensAndShares(ctx, depPool, shares)

	logger.Error(fmt.Sprintf("returned to: Unbond"))
	logger.Error(fmt.Sprintf("    depPool.Shares: %s", depPool.Shares.String()))

	if depPool.Shares.IsZero() {
		//TODO: if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveDepositPool(ctx, valAddr)
	}

	return issuedTokensAmt, nil
}
