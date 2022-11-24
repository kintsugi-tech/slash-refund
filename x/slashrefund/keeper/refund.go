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
	// Iterate attributes to find which validators has been slashed
	logger := k.Logger(ctx)
	logger.Error("        |_ Entered HandleRefundsFromSlash")

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
				logger.Error("        ERROR cant find validator.")
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
				//TODO Handle errorsa
				logger.Error("        ERROR: Unknown slashing reason.")
			}
		case "burned_coins":
			burnedTokens, isFound = sdk.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				//TODO Handle errors
				logger.Error("        ERROR in converting burnedTokens into int.")
			}
		case "infraction_height":
			infractionHeight, isFound = math.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				//TODO Handle errors
				logger.Error("        ERROR in converting infractionHeight into int.")
			}
		}
	}

	logger.Error(fmt.Sprintf("            BurnedTokens amt is %s", burnedTokens.String()))

	valAddr, err := sdk.ValAddressFromBech32(string(validator.OperatorAddress))
	if err != nil {
		//TODO return error
		logger.Error("		ERROR: Can't transform OperatorAddress into sdk.valAddress")
	}

	//This is not an error because we can still have an unbonding deposit queue we can access.
	depPool, isFoundDepositPool := k.GetDepositPool(ctx, valAddr)
	if !isFoundDepositPool {
		logger.Error(fmt.Sprintf("No pool for the slashed validator: %s", valAddr.String()))
	}

	// Compute how much to refund
	switch {

	// impossible case
	case infractionHeight.Int64() > ctx.BlockHeight():
		panic(fmt.Sprintf(
			"impossible attempt to handle a slash: future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	// special case: unbonding delegations and redelegations were not slashed
	case infractionHeight.Int64() == ctx.BlockHeight():
		if !isFoundDepositPool {
			return sdk.NewInt(0)
		}

		// draw from pool
		//TODO: depPool Tokens has also a denom, should be managed
		refundAmount, err = k.UpdateValidatorDepositPool(ctx, burnedTokens, depPool, validator)
		if err != nil {
			//TODO Handle error
			logger.Error("        |_ ERROR RefundFromValidatorPool")
			return sdk.NewInt(0)
		}
		////////////////////////////////////////////////////////////////////////////
		// k.stakingKeeper.AddValidatorTokens(ctx, validator, refundAmount)
		////////////////////////////////////////////////////////////////////////////
		k.AddValidatorTokens_SR(ctx, validator, refundAmount)
		logger.Error(fmt.Sprintf("        |_ Refunded %s to validator %s", refundAmount.String(), valAddr.String()))

	// must check for unbondings between slash and evidence
	case infractionHeight.Int64() < ctx.BlockHeight():
		// Iterate through unbonding deposits from slashed validator
		unbondingDeposits := k.GetUnbondingDepositsFromValidator(ctx, validator.OperatorAddress)

		// pool+ubds tokens
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
			logger.Error("        |_ Exited: No refunds available.")
			return sdk.NewInt(0)
		}

		// drawFactor is not capped at 1 because deposit and unbonding deposit update methods
		// handles the cap on the maximum available amount to draw.
		drawFactor = sdk.NewDecFromInt(burnedTokens).QuoInt(availableRefundTokens)

		// log
		logger.Error(fmt.Sprintf("              |_ unbondingRefunds %s", unbondingRefunds.String()))
		logger.Error(fmt.Sprintf("                 depPool %s", depPool.Tokens.Amount.String()))
		logger.Error(fmt.Sprintf("                 drawFactor %s", drawFactor.String()))

		// Update pool:
		drawnFromPool := sdk.NewInt(0)
		if isFoundDepositPool {
			amtToDrawFromPoolDec := drawFactor.MulInt(depPool.Tokens.Amount)
			amtToDrawFromPool := amtToDrawFromPoolDec.TruncateInt()
			drawnFromPool, err = k.UpdateValidatorDepositPool(ctx, amtToDrawFromPool, depPool, validator)
			logger.Error(fmt.Sprintf("                     |_ DepPool updated: drawnFromPool %s", drawnFromPool.String()))
		}
		if err != nil {
			logger.Error("        |_ ERROR RefundFromValidatorPool")
		}

		//  Update ubdeps
		drawnFromUBDs, err := k.UpdateValidatorUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64(), drawFactor)
		logger.Error(fmt.Sprintf("                     |_ Ubds updated: drawnFromUBDs %s", drawnFromUBDs.String()))
		if err != nil {
			logger.Error("        |_ ERROR RefundFromUnbondingDeposits")
		}

		// Compute total refunds
		refundAmount = drawnFromPool.Add(drawnFromUBDs)

		// compute refund factor
		refFactor := sdk.NewDecFromInt(refundAmount).QuoInt(burnedTokens)

		amtRefundedUBDs, err := k.RefundSlashedUnbondingDelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor)
		if err != nil {
			//TODO Handle errors
			logger.Error("        |_ ERROR RefundSlashedUnbondingDelegations")
		}
		amtRefundedRedel, err := k.RefundSlashedRedelegations(ctx, valAddr, infractionHeight.Int64(), slashFactor, refFactor)
		if err != nil {
			//TODO Handle errors
			logger.Error("        |_ ERROR RefundSlashedUnbondingDelegations")
		}

		refundForValidator := refundAmount.Sub(amtRefundedUBDs).Sub(amtRefundedRedel)
		k.stakingKeeper.AddValidatorTokens(ctx, validator, refundForValidator)
		logger.Error(fmt.Sprintf("        |_ Refunded %s to validator %s", refundForValidator.String(), valAddr.String()))
		logger.Error(fmt.Sprintf("        |_ Total Refunded is %s", refundAmount.String()))

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

	logger := k.Logger(ctx)
	logger.Error("        |_ Entered ComputeEligibleRefundFromUnbondingDeposits")
	logger.Error(fmt.Sprintf("          |_ num of ubds: %d", len(unbondingDeposits)))

	for _, unbondingDeposit := range unbondingDeposits {
		for _, entry := range unbondingDeposit.Entries {

			// If unbonding deposit entry started before infractionHeight, this entry is not eligible for refund
			if entry.CreationHeight < infractionHeight {
				logger.Error("            entry discarded due to CreationHeight")
				continue
			}

			if entry.IsMature(now) {
				// Unbonding deposit entry no longer eligible for refund, skip it
				logger.Error("            entry discarded due IsMature")
				continue
			}
			logger.Error(fmt.Sprintf("            entry eligible and adding %s", entry.Balance.String()))
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
	logger := k.Logger(ctx)
	logger.Error("                     |_ Entered UpdateValidatorUnbondingDeposits")
	totalRefundAmount = sdk.NewInt(0)
	for _, unbondingDeposit := range unbondingDeposits {
		refundAmount, err := k.UpdateUnbondingDepositEntries(ctx, unbondingDeposit, infractionHeight, drawFactor)
		_ = err
		totalRefundAmount = totalRefundAmount.Add(refundAmount)
	}
	return totalRefundAmount, nil
}

func (k Keeper) UpdateUnbondingDepositEntries(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (refundAmount sdk.Int, err error) {

	logger := k.Logger(ctx)

	now := ctx.BlockHeader().Time
	refundAmount = sdk.ZeroInt()

	// look at all entries within the unbonding deposit
	for i, entry := range unbondingDeposit.Entries {
		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			logger.Error("                     entry skip because CreationHeight")
			continue
		}

		if entry.IsMature(now) {
			logger.Error("                     entry skip because mature")
			// Unbonding deposit no longer eligible for withdraw, skip it
			continue
		}

		// Calculate refund amount proportional to deposit contributing to cover the infraction
		entryRefundAmountDec := drawFactor.MulInt(entry.InitialBalance)
		entryRefundAmount := entryRefundAmountDec.TruncateInt()

		// Don't refund more tokens than held.
		// Possible since the unbonding deposit may already have been drawn
		entryRefundAmount = sdk.MinInt(entryRefundAmount, entry.Balance)
		logger.Error(fmt.Sprintf("                     entry: balance:     %s", entry.Balance.String()))
		logger.Error(fmt.Sprintf("                            amt to draw: %s", entryRefundAmount.String()))

		// Update unbonding deposit entry only if necessary
		if !entryRefundAmount.IsZero() {
			entry.Balance = entry.Balance.Sub(entryRefundAmount)
			unbondingDeposit.Entries[i] = entry
			k.SetUnbondingDeposit(ctx, unbondingDeposit)
			logger.Error(fmt.Sprintf("                            new balance: %s", entry.Balance.String()))
		}
		//TODO Check if entry balance is zero, and in this case remove the entry
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
) (totalRefundedAmt sdk.Int, err error) {

	logger := k.Logger(ctx)

	now := ctx.BlockHeader().Time

	totalRefundedAmt = sdk.ZeroInt()

	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegationsFromValidator(ctx, validator)

	for _, ubd := range unbondingDelegations {

		// perform refund on all slashed entries within the unbonding delegation
		for i, entry := range ubd.Entries {

			// get delegator address (used only in logger)
			delegator, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
			if err != nil {
				return sdk.ZeroInt(), err
			}

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

			entry.Balance = entry.Balance.Add(refundAmt)
			ubd.Entries[i] = entry
			k.stakingKeeper.SetUnbondingDelegation(ctx, ubd)

			/*	OLD VERSION: send coins directly to depositor account:

				//TODO: generalize refundDenom with all the AllowedTokens
				refundDenom := k.AllowedTokensList(ctx)[0]
				coins := sdk.NewCoins(sdk.NewCoin(refundDenom, refundAmt))

				k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, coins)
			*/

			totalRefundedAmt = totalRefundedAmt.Add(refundAmt)

			logger.Error(fmt.Sprintf("        |_ Undelegation: Refunded %s to delegator %s, entry.CreationHeight=%d", refundAmt.String(), delegator.String(), entry.CreationHeight))
		}
	}
	return totalRefundedAmt, err
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedRedelegations(
	ctx sdk.Context,
	validator sdk.ValAddress,
	infractionHeight int64,
	slashFactor sdk.Dec,
	refFactor sdk.Dec,
) (totalRefundedAmt sdk.Int, err error) {

	logger := k.Logger(ctx)

	now := ctx.BlockHeader().Time

	totalRefundedAmt = sdk.ZeroInt()

	redelegations := k.stakingKeeper.GetRedelegationsFromSrcValidator(ctx, validator)

	for _, red := range redelegations {

		// perform refund on all slashed entries within the unbonding delegation
		for _, entry := range red.Entries {
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
			//TODO: generalize refundDenom with all the AllowedTokens
			refundDenom := k.AllowedTokensList(ctx)[0]

			coins := sdk.NewCoins(sdk.NewCoin(refundDenom, refundAmt))
			delegator, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
			if err != nil {
				return sdk.ZeroInt(), err
			}

			k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, coins)
			totalRefundedAmt = totalRefundedAmt.Add(refundAmt)

			logger.Error(fmt.Sprintf("        |_ Redelegation: Refunded %s to delegator %s", refundAmt.String(), delegator.String()))
		}
	}
	return totalRefundedAmt, err
}
