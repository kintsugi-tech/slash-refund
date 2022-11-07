package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

	for _, attr := range slashEvent.Attributes {

		switch string(attr.GetKey()) {
		case "address":
			validator, isFound = k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
			if !isFound {
				logger.Error("        ERROR cant find validator.")
				return sdk.NewInt(0)
			}
		case "burned_coins":
			burnedTokens, isFound = sdk.NewIntFromString(string(attr.GetValue()))
			if !isFound {
				logger.Error("        ERROR in converting burnedTokens into int.")
			}
		case "infraction_height":
			infractionHeight, isFound = math.NewIntFromString(string(attr.GetValue()))
			if !isFound {
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

	//This is not an error because we can still have an unbonding deposit queue
	//we can access.
	depPool, isFoundDepositPool := k.GetDepositPool(ctx, valAddr)
	if !isFoundDepositPool {
		logger.Error(fmt.Sprintf("No pool for the slashed validator: %s", valAddr.String()))
	}

	switch {
	case infractionHeight.Int64() > ctx.BlockHeight():
		panic(fmt.Sprintf(
			"impossible attempt to handle a slash: future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight.Int64() == ctx.BlockHeight():
		if !isFoundDepositPool {
			return sdk.NewInt(0)
		}
		//TODO: depPool Tokens has also a denom, should be managed
		// draw from pool
		refundAmount, err = k.UpdateValidatorDepositPool(ctx, burnedTokens, depPool, validator)
		if err != nil {
			//TODO Handle error
			logger.Error("        |_ ERROR RefundFromValidatorPool")
		}

	case infractionHeight.Int64() < ctx.BlockHeight():
		// Iterate through unbonding deposits from slashed validator
		unbondingDeposits := k.GetUnbondingDepositsFromValidator(ctx, valAddr)
		// pool+ubds tokens
		var availableRefundTokens sdk.Int

		unbondingRefunds := k.ComputeEligibleRefundFromUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64())
		logger.Error(fmt.Sprintf("          |_ unbondingRefunds %s", unbondingRefunds.String()))
		if !isFoundDepositPool {
			availableRefundTokens = unbondingRefunds
		} else {
			availableRefundTokens = depPool.Tokens.Amount.Add(unbondingRefunds)
		}
		// compute percentage to draw
		drawFactor := sdk.NewDec(0)
		//TODO CODE STABILITY: check if there is a possibility that depPool.Tokens.Amount and unbondingRefunds are both zero
		drawFactor = sdk.NewDecFromInt(burnedTokens).QuoInt(availableRefundTokens)
		// Update:
		//  pool
		refundFromPool := sdk.NewInt(0)
		if isFoundDepositPool {
			amtToDrawFromPoolDec := drawFactor.MulInt(depPool.Tokens.Amount)
			amtToDrawFromPool := amtToDrawFromPoolDec.TruncateInt()
			refundFromPool, err = k.UpdateValidatorDepositPool(ctx, amtToDrawFromPool, depPool, validator)
		}
		if err != nil {
			logger.Error("        |_ ERROR RefundFromValidatorPool")
		}
		//  ubds
		refundFromUbds, err := k.UpdateValidatorUnbondingDeposits(ctx, unbondingDeposits, infractionHeight.Int64(), drawFactor)
		if err != nil {
			logger.Error("        |_ ERROR RefundFromUnbondingDeposits")
		}
		// Compute total to refund
		refundAmount = refundFromPool.Add(refundFromUbds)
	}

	k.stakingKeeper.AddValidatorTokens(ctx, validator, refundAmount)

	logger.Error(fmt.Sprintf("        |_ Refunded %s to validator %s", refundAmount.String(), valAddr.String()))

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
	logger.Error("          |_ num of ubds:", len(unbondingDeposits))

	for _, unbondingDeposit := range unbondingDeposits {
		for _, entry := range unbondingDeposit.Entries {

			// If unbonding deposit entry started before infractionHeight, this entry is not eligible for refund
			if entry.CreationHeight < infractionHeight {
				logger.Error("        Discarded due to CreationHeight")
				continue
			}

			if entry.IsMature(now) {
				// Unbonding deposit entry no longer eligible for refund, skip it
				logger.Error("        Discarded due IsMature")
				continue
			}
			logger.Error(fmt.Sprintf("        Found and adding %s", entry.Balance.String()))
			totalUBDSAmount = totalUBDSAmount.Add(entry.Balance)
		}
	}
	//TODO make a list of indexes to ease the refund procedure
	return totalUBDSAmount
}

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

	now := ctx.BlockHeader().Time
	refundAmount = sdk.ZeroInt()

	// look at all entries within the unbonding deposit
	for i, entry := range unbondingDeposit.Entries {
		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding deposit no longer eligible for withdraw, skip it
			continue
		}

		// Calculate refund amount proportional to deposit contributing to cover the infraction
		refundAmountDec := drawFactor.MulInt(entry.InitialBalance)
		refundAmount := refundAmountDec.TruncateInt()

		// Don't refund more tokens than held.
		// Possible since the unbonding deposit may already
		// have been used
		entryRefundAmount := sdk.MinInt(refundAmount, entry.Balance)

		// Update unbonding deposit entry only if necessary
		if !entryRefundAmount.IsZero() {
			entry.Balance = entry.Balance.Sub(entryRefundAmount)
			unbondingDeposit.Entries[i] = entry
			k.SetUnbondingDeposit(ctx, unbondingDeposit)
		}
		//TODO Check if entry balance is zero, and in this case remove the entry
		refundAmount = refundAmount.Add(entryRefundAmount)
	}

	return refundAmount, nil

}
