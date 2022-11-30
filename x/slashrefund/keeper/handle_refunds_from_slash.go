package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) HandleRefundsFromSlash(ctx sdk.Context, slashEvent sdk.Event) (refundAmount sdk.Int) {

	//TODO Handle errors

	// Iterate attributes to find which validators has been slashed
	var validator stakingtypes.Validator
	var isFound bool
	var burnedTokens sdk.Int
	var infractionHeight sdk.Int
	var slashFactor sdk.Dec

	for _, attr := range slashEvent.Attributes {

		switch string(attr.GetKey()) {
		case "address":
			validator, isFound = k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
			if !isFound {
				return sdk.NewInt(0)
			}
		case "reason":
			switch string(attr.GetValue()) {
			case slashingtypes.AttributeValueDoubleSign:
				slashFactor = k.slashingKeeper.SlashFractionDoubleSign(ctx)
			case slashingtypes.AttributeValueMissingSignature:
				slashFactor = k.slashingKeeper.SlashFractionDowntime(ctx)
			default:
				slashFactor = sdk.ZeroDec()
				//ERROR: Unknown slashing reason.")
				return sdk.NewInt(0)
			}
		case "burned_coins":
			burnedTokens, isFound = sdk.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				//ERROR in converting burnedTokens into int.
				return sdk.NewInt(0)
			}
		case "infraction_height":
			infractionHeight, isFound = math.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				//ERROR in converting infractionHeight into int.
				return sdk.NewInt(0)
			}
		}
	}

	valAddr, err := sdk.ValAddressFromBech32(string(validator.OperatorAddress))
	if err != nil {
		//ERROR: Can't transform OperatorAddress into sdk.valAddress
		return sdk.NewInt(0)
	}

	//No error if not found the deposit pool because we can still have an unbonding deposit queue we can access.
	depPool, isFoundDepositPool := k.GetDepositPool(ctx, valAddr)

	// Check if the deposit pool exists or create it
	refPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		// TODO: should be initialized with actual Coins allowed. Now the hp is of just one allowed token.
		refPool = types.NewRefundPool(
			valAddr,
			sdk.NewCoin(k.AllowedTokensList(ctx)[0], sdk.ZeroInt()),
			sdk.ZeroDec(),
		)
	}

	// Compute how much to refund and refund
	switch {

	// ---- impossible case ----
	case infractionHeight.Int64() > ctx.BlockHeight():
		panic(fmt.Sprintf(
			"impossible attempt to handle a slash: future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	// ---- special case: ----
	// unbonding delegations and redelegations were not slashed
	case infractionHeight.Int64() == ctx.BlockHeight():
		if !isFoundDepositPool {
			return sdk.NewInt(0)
		}

		// draw from pool
		//TODO: depPool Tokens has also a denom, should be managed
		refundAmount, err = k.UpdateValidatorDepositPool(ctx, burnedTokens, depPool, validator)
		if err != nil {
			//ERROR in RefundFromValidatorPool
			return sdk.NewInt(0)
		}

		// refund
		refundTokens := sdk.NewCoin(k.AllowedTokensList(ctx)[0], refundAmount)
		k.AddRefPoolTokensAndShares(ctx, refPool, refundTokens)

		// distribute shares
		delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)
		for _, del := range delegations {
			delAddr := del.GetDelegatorAddr()
			_ = delAddr
			// TODO: check if refund exists for (delegator,validator), if not create it
			// TODO: update shares of refund
		}

	// ---- standard case: ----
	// must check for unbondings between slash and evidence
	case infractionHeight.Int64() < ctx.BlockHeight():
		// Iterate through unbonding deposits from slashed validator
		unbondingDeposits := k.GetUnbondingDepositsFromValidator(ctx, validator.OperatorAddress)

		// compute pool+ubds amount
		var availableRefundTokens sdk.Int
		availableRefundTokens = sdk.ZeroInt()

		unbondingRefunds := sdk.ZeroInt()
		if len(unbondingDeposits) > 0 {
			unbondingRefunds = k.ComputeEligibleRefundFromUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64())
		}

		if !isFoundDepositPool {
			availableRefundTokens = unbondingRefunds
		} else {
			availableRefundTokens = depPool.Tokens.Amount.Add(unbondingRefunds)
		}

		// compute percentage to draw from pool and ubdeps
		drawFactor := sdk.NewDec(0)
		if availableRefundTokens.IsZero() {
			return sdk.NewInt(0)
		}

		// ====== DRAW ======
		// drawFactor is not capped at 1 because deposit and unbonding deposit update methods
		// handles the cap on the maximum available amount to draw.
		drawFactor = sdk.NewDecFromInt(burnedTokens).QuoInt(availableRefundTokens)

		drawnFromPool := sdk.NewInt(0)
		if isFoundDepositPool {
			amtToDrawFromPoolDec := drawFactor.MulInt(depPool.Tokens.Amount)
			amtToDrawFromPool := amtToDrawFromPoolDec.TruncateInt()
			drawnFromPool, err = k.UpdateValidatorDepositPool(ctx, amtToDrawFromPool, depPool, validator)
		}
		if err != nil {
			//ERROR in RefundFromValidatorPool
			return sdk.NewInt(0)
		}

		drawnFromUBDs, err := k.UpdateValidatorUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64(), drawFactor)
		if err != nil {
			//ERROR in RefundFromUnbondingDeposits
			return sdk.NewInt(0)
		}

		// ====== REFUND ======
		// Compute total refunds
		refundAmount = drawnFromPool.Add(drawnFromUBDs)

		// compute refund factor
		refFactor := sdk.NewDecFromInt(refundAmount).QuoInt(burnedTokens)

		// get refund pool shares-token ratio
		poolShTkRatio, err := refPool.GetSharesOverTokensRatio()
		if err != nil {
			// zero tokens in pool means issued shares are 1 to 1 with added tokens
			poolShTkRatio = sdk.NewDec(1)
		}

		// refund undelegations
		amtRefundedUBDs, sharesRefundUBDS, err := k.RefundSlashedUnbondingDelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor, poolShTkRatio)
		if err != nil {
			//ERROR in RefundSlashedUnbondingDelegations
		}

		// refund redelegations
		amtRefundedRedel, sharesRefundRedel, err := k.RefundSlashedRedelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor, poolShTkRatio)
		if err != nil {
			//ERROR in RefundSlashedUnbondingDelegations
		}

		// refund delegations
		amtRefundedDel, sharesRefundDel, err := k.RefundSlashedDelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor, poolShTkRatio)
		if err != nil {
			//ERROR in RefundSlashedDelegations
		}

		// compute total refund shares issued
		totalRefundShares := sharesRefundUBDS.Add(sharesRefundRedel).Add(sharesRefundDel)
		refundDiff := refundAmount.Sub(amtRefundedUBDs).Sub(amtRefundedRedel).Sub(amtRefundedDel)
		_ = refundDiff
		// TODO: check refundDiff

		// update refund pool
		if !refundAmount.IsZero() && !totalRefundShares.IsZero() {
			refundTokens := sdk.NewCoin(k.AllowedTokensList(ctx)[0], refundAmount)
			refPool.Tokens.Add(refundTokens)
			refPool.Shares.Add(totalRefundShares)
			k.SetRefundPool(ctx, refPool)
		}

	}

	return refundAmount
}

func (keeper Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddr []byte) (validator stakingtypes.Validator, found bool) {
	// Decode address
	addr, _ := sdk.ConsAddressFromBech32(string(consAddr))
	// TODO: Handle error
	return keeper.stakingKeeper.GetValidatorByConsAddr(ctx, addr)
}

func (k Keeper) ComputeEligibleRefundFromUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64,
) (totalUBDSAmount sdk.Int) {
	now := ctx.BlockHeader().Time
	totalUBDSAmount = sdk.NewInt(0)

	for _, unbondingDeposit := range unbondingDeposits {
		for _, entry := range unbondingDeposit.Entries {

			// If unbonding deposit entry started before infractionHeight, this entry is not eligible for refund
			if entry.CreationHeight < infractionHeight {
				continue
			}

			// If mature the unbonding deposit entry is no longer eligible for refund
			if entry.IsMature(now) {
				continue
			}

			totalUBDSAmount = totalUBDSAmount.Add(entry.Balance)
		}
	}
	//TODO make a list of indexes to ease the refund procedure
	return totalUBDSAmount
}

// set the deposit pool and returns the amount that will be refunded from the pool.
func (k Keeper) UpdateValidatorDepositPool(ctx sdk.Context, amt sdk.Int, depPool types.DepositPool, validator stakingtypes.Validator,
) (refundTokens sdk.Int, err error) {
	if amt.GT(depPool.Tokens.Amount) {
		refundTokens = depPool.Tokens.Amount
		depPool.Tokens.Amount = sdk.ZeroInt()
		// TODO: remove pool
	} else {
		refundTokens = amt
		depPool.Tokens.Amount = depPool.Tokens.Amount.Sub(amt)
	}
	k.SetDepositPool(ctx, depPool)
	return refundTokens, nil
}

func (k Keeper) UpdateValidatorUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (totalRefundAmount sdk.Int, err error) {

	//TODO: Handle error

	totalRefundAmount = sdk.NewInt(0)
	for _, unbondingDeposit := range unbondingDeposits {
		refundAmount, err := k.UpdateUnbondingDepositEntries(ctx, unbondingDeposit, infractionHeight, drawFactor)
		if err != nil {
			//ERROR in UpdateUnbondingDepositEntries
			refundAmount = sdk.NewInt(0)
		}
		totalRefundAmount = totalRefundAmount.Add(refundAmount)
	}
	return totalRefundAmount, nil
}

func (k Keeper) UpdateUnbondingDepositEntries(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (refundAmount sdk.Int, err error) {

	now := ctx.BlockHeader().Time
	refundAmount = sdk.ZeroInt()

	// look at all entries within the unbonding deposit
	for i, entry := range unbondingDeposit.Entries {
		// If unbonding entry started before this height, entry were not eligible for refunding, so skip it
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding deposit were no longer eligible for refunding, so skip it
			continue
		}

		// Calculate refund amount proportional to deposit contributing to cover the infraction
		entryRefundAmountDec := drawFactor.MulInt(entry.InitialBalance)
		entryRefundAmount := entryRefundAmountDec.TruncateInt()

		// Don't refund more tokens than held.
		// Possible since the unbonding deposit may already have been drawn
		entryRefundAmount = sdk.MinInt(entryRefundAmount, entry.Balance)

		// Update unbonding deposit entry only if necessary
		if !entryRefundAmount.IsZero() {
			entry.Balance = entry.Balance.Sub(entryRefundAmount)
			unbondingDeposit.Entries[i] = entry
			k.SetUnbondingDeposit(ctx, unbondingDeposit)
		}
		//TODO remove the entry if entry balance is zero
		refundAmount = refundAmount.Add(entryRefundAmount)
	}

	return refundAmount, nil

}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedUnbondingDelegations(
	ctx sdk.Context,
	validator sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec, err error) {

	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegationsFromValidator(ctx, validator)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, ubd := range unbondingDelegations {

		delegator := ubd.DelegatorAddress
		_ = delegator
		delegatorShares := sdk.ZeroDec()

		// process slashed entries
		for _, entry := range ubd.Entries {

			// If unbonding started before this height, stake didn't contribute to infraction
			if entry.CreationHeight < infractionHeight {
				continue
			}
			if entry.IsMature(now) {
				// Unbonding delegation were no longer eligible for slashing, skip it
				continue
			}

			refundAmtDec := refFactor.Mul(slashFactor).MulInt(entry.InitialBalance)
			refundAmt := refundAmtDec.TruncateInt()

			// compute un-delegator refund shares
			entryShares := poolShTkRatio.MulInt(refundAmt)
			delegatorShares = delegatorShares.Add(entryShares)
			totalRefundedAmt = totalRefundedAmt.Add(refundAmt)
		}

		// distribute shares
		// TODO: check if refund exists for (delegator,validator), if not create it
		// TODO: update refund shares of delegator
		totalRefundShares = totalRefundShares.Add(delegatorShares)
	}

	return totalRefundedAmt, totalRefundShares, err
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedRedelegations(
	ctx sdk.Context,
	validator sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec, err error) {

	redelegations := k.stakingKeeper.GetRedelegationsFromSrcValidator(ctx, validator)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, red := range redelegations {

		delegator := red.DelegatorAddress
		_ = delegator
		delegatorShares := sdk.NewDec(0)

		// process slashed entries
		for _, entry := range red.Entries {
			// If unbonding started before this height, stake didn't contribute to infraction
			if entry.CreationHeight < infractionHeight {
				continue
			}
			if entry.IsMature(now) {
				// Unbonding delegation were no longer eligible for slashing, skip it
				continue
			}

			// compute refund amount for this entry
			refundAmtDec := refFactor.Mul(slashFactor).MulInt(entry.InitialBalance)
			refundAmt := refundAmtDec.TruncateInt()

			// compute redelegator refund shares
			entryShares := poolShTkRatio.MulInt(refundAmt)
			delegatorShares = delegatorShares.Add(entryShares)
			totalRefundedAmt = totalRefundedAmt.Add(refundAmt)
		}

		// distribute shares
		// TODO: check if refund exists for (delegator,validator), if not create it
		// TODO: update refund shares of delegator
		totalRefundShares = totalRefundShares.Add(delegatorShares)
	}

	return totalRefundedAmt, totalRefundShares, err
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedDelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec, err error) {

	delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		//TODO Handle error
		return sdk.NewInt(0), sdk.NewDec(0), sdk.ErrEmptyHexAddress
	}

	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, del := range delegations {

		delegator := del.DelegatorAddress
		_ = delegator

		balanceDec := validator.TokensFromShares(del.Shares)
		balance := balanceDec.TruncateInt()

		//TODO check if zero-balance-check is needed to speed up the iteration

		initialBalanceDec := balanceDec.Quo(sdk.NewDec(1).Sub(slashFactor))
		delInitialBalance := initialBalanceDec.TruncateInt()

		burnedAmt := delInitialBalance.Sub(balance)

		refundAmtDec := refFactor.MulInt(burnedAmt)
		refundAmt := refundAmtDec.TruncateInt()

		delegatorShares := poolShTkRatio.MulInt(refundAmt)
		if err != nil {
			//ERROR in refundPool.SharesFromTokens
			return sdk.NewInt(0), sdk.NewDec(0), err
		}

		totalRefundedAmt = totalRefundedAmt.Add(refundAmt)

		// distribute shares
		// TODO: check if refund exists for (delegator,validator), if not create it
		// TODO: update refund shares of delegator
		totalRefundShares = totalRefundShares.Add(delegatorShares)
	}

	return totalRefundedAmt, totalRefundShares, err
}
