package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k Keeper) HandleRefundsFromSlash(ctx sdk.Context, slashEvent sdk.Event) (refundAmount sdk.Int, err error) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("Handling refunds from slash:"))

	// Iterate attributes to find which validators has been slashed
	var validator stakingtypes.Validator
	var isFound bool
	var validatorBurnedTokens sdk.Int
	var infractionHeight sdk.Int
	var slashFactor sdk.Dec

	for _, attr := range slashEvent.Attributes {

		switch string(attr.GetKey()) {
		case "address":
			validator, isFound, err = k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
			if !isFound {
				logger.Error("ERROR:  ", err.Error())
				return sdk.NewInt(0), err
			}
		case "reason":
			switch string(attr.GetValue()) {
			case slashingtypes.AttributeValueDoubleSign:
				slashFactor = k.slashingKeeper.SlashFractionDoubleSign(ctx)
			case slashingtypes.AttributeValueMissingSignature:
				slashFactor = k.slashingKeeper.SlashFractionDowntime(ctx)
			default:
				slashFactor = sdk.ZeroDec()
				logger.Error("ERROR:  cant reason")
				return sdk.NewInt(0), types.ErrUnknownSlashingReasonFromSlashEvent
			}
		case "burned_coins":
			validatorBurnedTokens, isFound = sdk.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				logger.Error("ERROR:  cant validator burned tokens")
				return sdk.NewInt(0), types.ErrCantGetValidatorBurnedTokensFromSlashEvent
			}
		case "infraction_height":
			infractionHeight, isFound = math.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				logger.Error("ERROR:  cant infraction height")
				return sdk.NewInt(0), types.ErrCantGetInfractionHeightFromSlashEvent
			}
		}
	}

	valAddr, err := sdk.ValAddressFromBech32(string(validator.OperatorAddress))
	if err != nil {
		logger.Error("ERROR:  ", err.Error())
		return sdk.NewInt(0), err
	}

	//No error if not found the deposit pool because we can still have an unbonding deposit queue we can access.
	depPool, isFoundDepositPool := k.GetDepositPool(ctx, valAddr)

	// Check if the deposit pool exists or create it
	refPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		// TODO: should be initialized with actual Coins allowed. Now the hp is of just one allowed token.
		refPool = types.NewRefundPool(
			valAddr,
			sdk.NewCoin(k.AllowedTokens(ctx)[0], sdk.ZeroInt()),
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
			logger.Error("ERROR:  zero deposit available")
			return sdk.NewInt(0), types.ErrZeroDepositAvailable
		}

		// draw from pool
		//TODO: depPool Tokens has also a denom, should be managed
		refundAmount = k.UpdateValidatorDepositPool(ctx, validatorBurnedTokens, depPool, validator)

		// compute refund factor
		refFactor := sdk.NewDecFromInt(refundAmount).QuoInt(validatorBurnedTokens)

		// get refund pool shares-token ratio
		poolShTkRatio, err := refPool.GetSharesOverTokensRatio()
		if err != nil {
			// zero tokens in pool means issued shares are 1 to 1 with added tokens
			poolShTkRatio = sdk.NewDec(1)
		}

		// refund delegations
		amtRefundedDel, sharesRefundDel := k.RefundSlashedDelegations(ctx, valAddr, infractionHeight.Int64(), refundAmount, poolShTkRatio)

		// compute total refund shares issued
		refundDiff := refundAmount.Sub(amtRefundedDel)
		_ = refundDiff
		// TODO: check refundDiff

		// update refund pool
		if !refundAmount.IsZero() && !sharesRefundDel.IsZero() {
			refundTokens := sdk.NewCoin(k.AllowedTokens(ctx)[0], refundAmount)
			refPool.Tokens = refPool.Tokens.Add(refundTokens)
			refPool.Shares = refPool.Shares.Add(sharesRefundDel)
			k.SetRefundPool(ctx, refPool)
			logger.Error("  refund pool updated.")
		}

		// LOGGER
		refPoolLOG, foundRefPoolLOG := k.GetRefundPool(ctx, valAddr)
		logger.Error(fmt.Sprintf("--------------- refund recap ---------------"))
		logger.Error(fmt.Sprintf("Slash info:"))
		logger.Error(fmt.Sprintf("  validator=%s", validator.OperatorAddress))
		logger.Error(fmt.Sprintf("  slash factor=%s", slashFactor.String()))
		logger.Error(fmt.Sprintf("  infraction height=%s", infractionHeight.String()))
		logger.Error(fmt.Sprintf("  burned from validator=%s", validatorBurnedTokens.String()))
		logger.Error(fmt.Sprintf("infractionHeight == currentHeight :"))
		logger.Error(fmt.Sprintf("  depPool found=%t", isFoundDepositPool))
		logger.Error(fmt.Sprintf("  refPool found=%t", found))
		logger.Error(fmt.Sprintf("  avlbl  in  depPool=%s", depPool.Tokens.Amount.String()))
		logger.Error(fmt.Sprintf("  drawn from depPool=%s", refundAmount.String()))
		logger.Error(fmt.Sprintf("  refFactor=%s", refFactor.String()))
		logger.Error(fmt.Sprintf("  poolShTkRatio=%s", poolShTkRatio.String()))
		logger.Error(fmt.Sprintf("  amount refunded to delegators=%s", amtRefundedDel.String()))
		logger.Error(fmt.Sprintf("  shares issued   to delegators=%s", sharesRefundDel.String()))
		logger.Error(fmt.Sprintf("  refund pool: added tokens=%s , added shares=%s", refundAmount.String(), sharesRefundDel.String()))
		logger.Error(fmt.Sprintf("  refund pool now: found=%t , tokens=%s , shares=%s", foundRefPoolLOG, refPoolLOG.Tokens.Amount.String(), refPoolLOG.Shares.String()))
		logger.Error(fmt.Sprintf("---------------  done refund ---------------"))

		// ---- standard case: ----
	// must check for unbondings between slash and evidence
	case infractionHeight.Int64() < ctx.BlockHeight():
		// Iterate through unbonding deposits from slashed validator
		unbondingDeposits := k.GetUnbondingDepositsFromValidator(ctx, valAddr)

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
			logger.Error("ERROR:  zero deposit available")
			return sdk.NewInt(0), types.ErrZeroDepositAvailable
		}

		// = COMPUTE BURNED =
		ubdelBurnedTokens := k.ComputeSlashedUnbondingDelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor)
		redelBurnedTokens := k.ComputeSlashedRedelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor)
		burnedTokens := validatorBurnedTokens.Add(ubdelBurnedTokens).Add(redelBurnedTokens)

		// ====== DRAW ======
		// drawFactor is not capped at 1 because deposit and unbonding deposit update methods
		// handles the cap on the maximum available amount to draw.
		drawFactor = sdk.NewDecFromInt(burnedTokens).QuoInt(availableRefundTokens)

		drawnFromPool := sdk.NewInt(0)
		if isFoundDepositPool {
			amtToDrawFromPoolDec := drawFactor.MulInt(depPool.Tokens.Amount)
			amtToDrawFromPool := amtToDrawFromPoolDec.TruncateInt()
			drawnFromPool = k.UpdateValidatorDepositPool(ctx, amtToDrawFromPool, depPool, validator)
		}

		drawnFromUBDs := k.UpdateValidatorUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64(), drawFactor)

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
		amtRefundedUBDs, sharesRefundUBDS := k.RefundSlashedUnbondingDelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor, poolShTkRatio)

		// refund redelegations
		amtRefundedRedel, sharesRefundRedel := k.RefundSlashedRedelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor, poolShTkRatio)

		// refund delegations
		refundForDelegators := refundAmount.Sub(amtRefundedUBDs).Sub(amtRefundedRedel)
		amtRefundedDel, sharesRefundDel := k.RefundSlashedDelegations(ctx, valAddr, infractionHeight.Int64(), refundForDelegators, poolShTkRatio)

		// compute total refund shares issued
		totalRefundShares := sharesRefundUBDS.Add(sharesRefundRedel).Add(sharesRefundDel)
		refundDiff := refundAmount.Sub(amtRefundedUBDs).Sub(amtRefundedRedel).Sub(amtRefundedDel)
		_ = refundDiff
		// TODO: check refundDiff

		// update refund pool
		if !refundAmount.IsZero() && !totalRefundShares.IsZero() {
			refundTokens := sdk.NewCoin(k.AllowedTokens(ctx)[0], refundAmount)
			refPool.Tokens = refPool.Tokens.Add(refundTokens)
			refPool.Shares = refPool.Shares.Add(totalRefundShares)
			k.SetRefundPool(ctx, refPool)
			logger.Error("  refund pool updated.")
		}

		// LOGGER
		refPoolLOG, foundRefPoolLOG := k.GetRefundPool(ctx, valAddr)
		logger.Error(fmt.Sprintf("--------------- refund recap ---------------"))
		logger.Error(fmt.Sprintf("Slash info:"))
		logger.Error(fmt.Sprintf("  validator=%s", validator.OperatorAddress))
		logger.Error(fmt.Sprintf("  slash factor=%s", slashFactor.String()))
		logger.Error(fmt.Sprintf("  infraction height=%s", infractionHeight.String()))
		logger.Error(fmt.Sprintf("  burned from validator=%s", validatorBurnedTokens.String()))
		logger.Error(fmt.Sprintf("infractionHeight < currentHeight :"))
		logger.Error(fmt.Sprintf("  burned from   undel  =%s", ubdelBurnedTokens.String()))
		logger.Error(fmt.Sprintf("  burned from   redel  =%s", redelBurnedTokens.String()))
		logger.Error(fmt.Sprintf("  burned     total     =%s", burnedTokens.String()))
		logger.Error(fmt.Sprintf("  depPool found=%t", isFoundDepositPool))
		logger.Error(fmt.Sprintf("  refPool found=%t", found))
		logger.Error(fmt.Sprintf("  avlbl  in  depPool=%s", depPool.Tokens.Amount.String()))
		logger.Error(fmt.Sprintf("  avlbl  in  unbndng=%s", unbondingRefunds.String()))
		logger.Error(fmt.Sprintf("  total ref availabl=%s", availableRefundTokens.String()))
		logger.Error(fmt.Sprintf("  drawFactor=%s", drawFactor.String()))
		logger.Error(fmt.Sprintf("  drawn from deposit pool=%s", drawnFromPool.String()))
		logger.Error(fmt.Sprintf("  drawn from unbnd depost=%s", drawnFromUBDs.String()))
		logger.Error(fmt.Sprintf("  refund amount=%s", refundAmount.String()))
		logger.Error(fmt.Sprintf("  refFactor=%s", refFactor.String()))
		logger.Error(fmt.Sprintf("  poolShTkRatio=%s", poolShTkRatio.String()))
		logger.Error(fmt.Sprintf("  amount refunded to unbnd dels=%s", amtRefundedUBDs.String()))
		logger.Error(fmt.Sprintf("  shares issued   to unbnd dels=%s", sharesRefundUBDS.String()))
		logger.Error(fmt.Sprintf("  amount refunded to redelegats=%s", amtRefundedRedel.String()))
		logger.Error(fmt.Sprintf("  shares issued   to redelegats=%s", sharesRefundRedel.String()))
		logger.Error(fmt.Sprintf("  amount refunded to delegators=%s", amtRefundedDel.String()))
		logger.Error(fmt.Sprintf("  shares issued   to delegators=%s", sharesRefundDel.String()))
		logger.Error(fmt.Sprintf("  total shares issued=%s", totalRefundShares.String()))
		logger.Error(fmt.Sprintf("  refund pool: added tokens=%s , added shares=%s", refundAmount.String(), totalRefundShares.String()))
		logger.Error(fmt.Sprintf("  refund pool now: found=%t , tokens=%s , shares=%s", foundRefPoolLOG, refPoolLOG.Tokens.Amount.String(), refPoolLOG.Shares.String()))
		logger.Error(fmt.Sprintf("---------------  done refund ---------------"))
	}

	return refundAmount, err
}

func (keeper Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddrByte []byte) (validator stakingtypes.Validator, found bool, err error) {
	// Decode address
	consAddr, err := sdk.ConsAddressFromBech32(string(consAddrByte))
	if err != nil {
		return validator, false, err
	}
	validator, found = keeper.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	return validator, found, types.ErrCantGetValidatorFromSlashEvent
}

func (k Keeper) ComputeEligibleRefundFromUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64,
) (totalUBDSAmount sdk.Int) {

	logger := k.Logger(ctx)
	logger.Error("  computing eligible unbonding deposits:")

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
			logger.Error(fmt.Sprintf("    found eligible entry: entry.balance=%s", entry.Balance.String()))
		}
	}
	//TODO make a list of indexes to ease the refund procedure
	return totalUBDSAmount
}

// set the deposit pool and returns the amount that will be refunded from the pool.
func (k Keeper) UpdateValidatorDepositPool(ctx sdk.Context, amt sdk.Int, depPool types.DepositPool, validator stakingtypes.Validator,
) (refundTokens sdk.Int) {
	if amt.GT(depPool.Tokens.Amount) {
		refundTokens = depPool.Tokens.Amount
		depPool.Tokens.Amount = sdk.ZeroInt()
		// TODO: remove pool
	} else {
		refundTokens = amt
		depPool.Tokens.Amount = depPool.Tokens.Amount.Sub(amt)
	}
	k.SetDepositPool(ctx, depPool)
	return refundTokens
}

func (k Keeper) UpdateValidatorUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (totalRefundAmount sdk.Int) {

	totalRefundAmount = sdk.NewInt(0)
	for _, unbondingDeposit := range unbondingDeposits {
		refundAmount := k.UpdateUnbondingDepositEntries(ctx, unbondingDeposit, infractionHeight, drawFactor)
		totalRefundAmount = totalRefundAmount.Add(refundAmount)
	}
	return totalRefundAmount
}

func (k Keeper) UpdateUnbondingDepositEntries(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (refundAmount sdk.Int) {

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

	return refundAmount

}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) ComputeSlashedUnbondingDelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
) (totalSlashedAmt sdk.Int) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("  computing slashed unbonding delegations :"))

	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegationsFromValidator(ctx, valAddr)
	now := ctx.BlockHeader().Time
	totalSlashedAmt = sdk.NewInt(0)

	for _, ubd := range unbondingDelegations {

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

			slashedAmtDec := slashFactor.MulInt(entry.InitialBalance)
			slashedAmt := slashedAmtDec.TruncateInt()

			totalSlashedAmt = totalSlashedAmt.Add(slashedAmt)
			logger.Error(fmt.Sprintf("    undel entry: salshed=%s , initialBalance=%s , balance=%s , delegator=%s", slashedAmt.String(), entry.InitialBalance.String(), entry.Balance.String(), ubd.DelegatorAddress))
		}
	}

	return totalSlashedAmt
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) ComputeSlashedRedelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
) (totalSlashedAmt sdk.Int) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("  computing slashed redelegations:"))

	redelegations := k.stakingKeeper.GetRedelegationsFromSrcValidator(ctx, valAddr)
	now := ctx.BlockHeader().Time
	totalSlashedAmt = sdk.NewInt(0)

	for _, red := range redelegations {

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

			slashedAmtDec := slashFactor.MulInt(entry.InitialBalance)
			slashedAmt := slashedAmtDec.TruncateInt()

			totalSlashedAmt = totalSlashedAmt.Add(slashedAmt)
			logger.Error(fmt.Sprintf("    redel entry: salshed=%s , initialBalance=%s , delegator=%s", slashedAmt.String(), entry.InitialBalance.String(), red.DelegatorAddress))
		}
	}

	return totalSlashedAmt
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedUnbondingDelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("  refunding unbonding delegations :"))

	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegationsFromValidator(ctx, valAddr)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, ubd := range unbondingDelegations {

		delAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
		if err != nil {
			logger.Error(fmt.Sprintf("PANIC: converting delAddr: %s", err.Error()))
			panic(err)
		}
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
			logger.Error(fmt.Sprintf("    undel entry: refund=%s , shares=%s , delegator=%s", refundAmt.String(), entryShares.String(), ubd.DelegatorAddress))
		}

		// issue shares
		refund, found := k.GetRefund(ctx, delAddr, valAddr)
		if !found {
			refund = types.NewRefund(delAddr, valAddr, sdk.ZeroDec())
		}
		refund.Shares = refund.Shares.Add(delegatorShares)
		k.SetRefund(ctx, refund)
		logger.Error(fmt.Sprintf("    undel: total delegator shares=%s , delegator=%s", delegatorShares.String(), ubd.DelegatorAddress))

		totalRefundShares = totalRefundShares.Add(delegatorShares)
	}

	return totalRefundedAmt, totalRefundShares
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedRedelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("  refunding redelegations :"))

	redelegations := k.stakingKeeper.GetRedelegationsFromSrcValidator(ctx, valAddr)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, red := range redelegations {

		delAddr, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
		if err != nil {
			logger.Error(fmt.Sprintf("PANIC: converting delAddr: %s", err.Error()))
			panic(err)
		}

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
			logger.Error(fmt.Sprintf("    redel entry: refund=%s , shares=%s , delegator=%s", refundAmt.String(), entryShares.String(), red.DelegatorAddress))
		}

		// issue shares
		refund, found := k.GetRefund(ctx, delAddr, valAddr)
		if !found {
			refund = types.NewRefund(delAddr, valAddr, sdk.ZeroDec())
		}
		refund.Shares = refund.Shares.Add(delegatorShares)
		k.SetRefund(ctx, refund)
		logger.Error(fmt.Sprintf("    redel: total delegator shares=%s , delegator=%s", delegatorShares.String(), red.DelegatorAddress))

		totalRefundShares = totalRefundShares.Add(delegatorShares)
	}

	return totalRefundedAmt, totalRefundShares
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedDelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	refund sdk.Int,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt sdk.Int, totalRefundShares sdk.Dec) {

	logger := k.Logger(ctx)
	logger.Error(fmt.Sprintf("  refunding delegations :"))

	delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		panic(fmt.Sprintf("validator record not found for address: %X\n", valAddr))
	}

	refundPerShare := sdk.NewDecFromInt(refund).Quo(validator.GetDelegatorShares())
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, del := range delegations {

		delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
			logger.Error(fmt.Sprintf("PANIC: converting delAddr: %s", err.Error()))
			panic(err)
		}

		delRefundAmtDec := refundPerShare.Mul(del.Shares)
		delRefundAmt := delRefundAmtDec.TruncateInt()

		// compute shares to issue
		delRefundShares := poolShTkRatio.MulInt(delRefundAmt)

		// issue shares
		refund, found := k.GetRefund(ctx, delAddr, valAddr)
		if !found {
			refund = types.NewRefund(delAddr, valAddr, sdk.ZeroDec())
		}
		refund.Shares = refund.Shares.Add(delRefundShares)
		k.SetRefund(ctx, refund)
		logger.Error(fmt.Sprintf("    deleg: refund=%s , shares=%s , delegator=%s", delRefundAmt.String(), delRefundShares.String(), del.DelegatorAddress))

		// update totals
		totalRefundedAmt = totalRefundedAmt.Add(delRefundAmt)
		totalRefundShares = totalRefundShares.Add(delRefundShares)
	}

	return totalRefundedAmt, totalRefundShares
}
