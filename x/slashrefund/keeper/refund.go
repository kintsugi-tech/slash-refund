package keeper

import (
	"fmt"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// TODO: update undelying logic in order to handle different tokens to be refunded.
func (k Keeper) HandleRefundsFromSlash(ctx sdk.Context, slashEvent sdk.Event) {

	// Iterate attributes to find which validator has been slashed
	valAddr, valBurnedTokens, infractionHeight, slashFactor, err := k.ProcessSlashEvent(ctx, slashEvent)
	if err != nil {
		return
	}

	refundAmount, err := k.RefundFromSlash(ctx, valAddr, valBurnedTokens, infractionHeight.Int64(), slashFactor)

	if err == nil {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeRefund,
				sdk.NewAttribute(types.AttributeKeyValidator, valAddr.String()),
				sdk.NewAttribute(types.AttributeKeyToken, k.AllowedTokens(ctx)[0]),
				sdk.NewAttribute(sdk.AttributeKeyAmount, refundAmount.String()),
			),
		})
	}
}

func (k Keeper) ProcessSlashEvent(ctx sdk.Context, event sdk.Event) (
	valAddr sdk.ValAddress,
	valBurnedAmt math.Int,
	infractionHeight math.Int,
	slashFactor sdk.Dec,
	err error,
) {

	var found bool

	for _, attr := range event.Attributes {

		switch string(attr.GetKey()) {

		case "address":
			validator, found := k.GetValidatorByConsAddrBytes(ctx, attr.GetValue())
			if !found {
				err = types.ErrCantGetValidatorFromSlashEvent
				break
			}
			valAddr, err = sdk.ValAddressFromBech32(validator.OperatorAddress)
			if err != nil {
				break
			}

		case "reason":
			switch string(attr.GetValue()) {

			case slashingtypes.AttributeValueDoubleSign:
				slashFactor = k.slashingKeeper.SlashFractionDoubleSign(ctx)

			case slashingtypes.AttributeValueMissingSignature:
				slashFactor = k.slashingKeeper.SlashFractionDowntime(ctx)

			default:
				err = types.ErrUnknownSlashingReasonFromSlashEvent
				break
			}

		case "burned_coins":
			valBurnedAmt, found = sdk.NewIntFromString(string(attr.GetValue()))
			if !found {
				err = types.ErrCantGetValidatorBurnedTokensFromSlashEvent
				break
			}

		case "infraction_height":
			infractionHeight, found = math.NewIntFromString(string(attr.GetValue()))
			if !found {
				err = types.ErrCantGetInfractionHeightFromSlashEvent
				break
			}
		}
	}

	return valAddr, valBurnedAmt, infractionHeight, slashFactor, err
}

func (k Keeper) RefundFromSlash(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	valBurnedTokens math.Int,
	infractionHeight int64,
	slashFactor sdk.Dec,
) (refundAmount math.Int, err error) {

	// If the deposit pool is not found it is not an error because there could be eligible unbonding deposits.
	depPool, isFoundDepositPool := k.GetDepositPool(ctx, valAddr)

	// Check if the refund pool exists or create it
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
	case infractionHeight > ctx.BlockHeight():
		panic(fmt.Sprintf(
			"impossible attempt to handle a slash: future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	// ---- special case: ----
	// unbonding delegations and redelegations were not slashed
	case infractionHeight == ctx.BlockHeight():
		if !isFoundDepositPool {
			return sdk.NewInt(0), nil
		}

		// draw from pool
		//TODO: depPool Tokens has also a denom, should be managed
		refundAmount = k.UpdateValidatorDepositPool(ctx, valBurnedTokens, depPool)

		// get refund pool shares-token ratio
		poolShTkRatio, err := refPool.GetSharesOverTokensRatio()
		if err != nil {
			// zero tokens in pool means issued shares are 1 to 1 with added tokens
			poolShTkRatio = sdk.NewDec(1)
		}

		// refund delegations
		amtRefundedDel, sharesRefundDel := k.RefundSlashedDelegations(ctx, valAddr, infractionHeight, refundAmount, poolShTkRatio)

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
		}

	// ---- standard case: ----
	// must check for unbondings between slash and evidence
	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding deposits from slashed validator
		unbondingDeposits := k.GetUnbondingDepositsFromValidator(ctx, valAddr)

		// compute pool+ubds amount
		var availableRefundTokens = sdk.ZeroInt()

		unbondingRefunds := sdk.ZeroInt()
		if len(unbondingDeposits) > 0 {
			unbondingRefunds = k.ComputeEligibleRefundFromUnbondingDeposits(ctx, unbondingDeposits, infractionHeight)
		}

		if !isFoundDepositPool {
			availableRefundTokens = unbondingRefunds
		} else {
			availableRefundTokens = depPool.Tokens.Amount.Add(unbondingRefunds)
		}

		// compute percentage to draw from pool and ubdeps
		if availableRefundTokens.IsZero() {
			return sdk.NewInt(0), nil
		}

		// = COMPUTE BURNED =
		ubdelBurnedTokens := k.ComputeSlashedUnbondingDelegations(ctx, valAddr, infractionHeight, slashFactor)
		redelBurnedTokens := k.ComputeSlashedRedelegations(ctx, valAddr, infractionHeight, slashFactor)
		burnedTokens := valBurnedTokens.Add(ubdelBurnedTokens).Add(redelBurnedTokens)

		// ====== DRAW ======
		// drawFactor is not capped at 1 because deposit and unbonding deposit update methods
		// handles the cap on the maximum available amount to draw.
		drawFactor := sdk.NewDecFromInt(burnedTokens).QuoInt(availableRefundTokens)

		drawnFromPool := sdk.NewInt(0)
		if isFoundDepositPool {
			amtToDrawFromPoolDec := drawFactor.MulInt(depPool.Tokens.Amount)
			amtToDrawFromPool := amtToDrawFromPoolDec.TruncateInt()
			drawnFromPool = k.UpdateValidatorDepositPool(ctx, amtToDrawFromPool, depPool)
		}

		drawnFromUBDs := k.UpdateValidatorUnbondingDeposits(ctx, unbondingDeposits, infractionHeight, drawFactor)

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
		amtRefundedUBDs, sharesRefundUBDS := k.RefundSlashedUnbondingDelegations(ctx, valAddr, infractionHeight, slashFactor, refFactor, poolShTkRatio)

		// refund redelegations
		amtRefundedRedel, sharesRefundRedel := k.RefundSlashedRedelegations(ctx, valAddr, infractionHeight, slashFactor, refFactor, poolShTkRatio)

		// refund delegations
		refundForDelegators := refundAmount.Sub(amtRefundedUBDs).Sub(amtRefundedRedel)
		amtRefundedDel, sharesRefundDel := k.RefundSlashedDelegations(ctx, valAddr, infractionHeight, refundForDelegators, poolShTkRatio)

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
		}
	}

	return refundAmount, err
}

func (k Keeper) GetValidatorByConsAddrBytes(ctx sdk.Context, consAddrByte []byte) (validator stakingtypes.Validator, found bool) {
	// Decode address
	consAddr, err := sdk.ConsAddressFromBech32(string(consAddrByte))
	if err != nil {
		return validator, false
	}
	validator, found = k.stakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	return validator, found
}

func (k Keeper) ComputeEligibleRefundFromUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64,
) (totalUBDSAmount math.Int) {

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
	// TODO make a list of indexes to ease the refund procedure
	return totalUBDSAmount
}

// set the deposit pool and returns the amount that will be refunded from the pool.
func (k Keeper) UpdateValidatorDepositPool(ctx sdk.Context, amt math.Int, depPool types.DepositPool,
) (refundTokens math.Int) {

	if amt.GTE(depPool.Tokens.Amount) {
		refundTokens = depPool.Tokens.Amount

		// Remove validator deposit pool and associated deposits
		valAddr, _ := sdk.ValAddressFromBech32(depPool.OperatorAddress)
		k.RemoveDepositPool(ctx, valAddr)
		for _, d := range k.GetValidatorDeposits(ctx, valAddr) {
			k.RemoveDeposit(ctx, d)
		}

	} else {
		refundTokens = amt
		depPool.Tokens.Amount = depPool.Tokens.Amount.Sub(amt)
		k.SetDepositPool(ctx, depPool)
	}

	return refundTokens
}

func (k Keeper) UpdateValidatorUnbondingDeposits(ctx sdk.Context, unbondingDeposits []types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (totalRefundAmount math.Int) {

	totalRefundAmount = sdk.NewInt(0)
	for _, unbondingDeposit := range unbondingDeposits {
		refundAmount := k.UpdateUnbondingDepositEntries(ctx, unbondingDeposit, infractionHeight, drawFactor)
		totalRefundAmount = totalRefundAmount.Add(refundAmount)
	}
	return totalRefundAmount
}

func (k Keeper) UpdateUnbondingDepositEntries(ctx sdk.Context, unbondingDeposit types.UnbondingDeposit, infractionHeight int64, drawFactor sdk.Dec,
) (refundAmount math.Int) {

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
		// TODO remove the entry if entry balance is zero
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
) (totalSlashedAmt math.Int) {

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
) (totalSlashedAmt math.Int) {

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
) (totalRefundedAmt math.Int, totalRefundShares sdk.Dec) {

	unbondingDelegations := k.stakingKeeper.GetUnbondingDelegationsFromValidator(ctx, valAddr)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, ubd := range unbondingDelegations {

		delAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
		if err != nil {
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
		}

		// issue shares
		if !delegatorShares.IsZero() {
			refund, found := k.GetRefund(ctx, delAddr, valAddr)
			if !found {
				refund = types.NewRefund(delAddr, valAddr, sdk.ZeroDec())
			}
			refund.Shares = refund.Shares.Add(delegatorShares)
			k.SetRefund(ctx, refund)

			totalRefundShares = totalRefundShares.Add(delegatorShares)
		}
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
) (totalRefundedAmt math.Int, totalRefundShares sdk.Dec) {

	redelegations := k.stakingKeeper.GetRedelegationsFromSrcValidator(ctx, valAddr)

	now := ctx.BlockHeader().Time
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, red := range redelegations {

		delAddr, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
		if err != nil {
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
		}

		// issue shares
		if !delegatorShares.IsZero() {
			refund, found := k.GetRefund(ctx, delAddr, valAddr)
			if !found {
				refund = types.NewRefund(delAddr, valAddr, sdk.ZeroDec())
			}
			refund.Shares = refund.Shares.Add(delegatorShares)
			k.SetRefund(ctx, refund)

			totalRefundShares = totalRefundShares.Add(delegatorShares)
		}
	}

	return totalRefundedAmt, totalRefundShares
}

// TODO: handle output of different denoms (return skd.Coins)
func (k Keeper) RefundSlashedDelegations(
	ctx sdk.Context,
	valAddr sdk.ValAddress,
	infractionHeight int64,
	refund math.Int,
	poolShTkRatio sdk.Dec,
) (totalRefundedAmt math.Int, totalRefundShares sdk.Dec) {

	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		panic(fmt.Sprintf("validator record not found for address: %X\n", valAddr))
	}
	delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)

	refundPerShare := sdk.NewDecFromInt(refund).Quo(validator.GetDelegatorShares())
	totalRefundedAmt = sdk.ZeroInt()
	totalRefundShares = sdk.ZeroDec()

	for _, del := range delegations {

		delAddr, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
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

		// update totals
		totalRefundedAmt = totalRefundedAmt.Add(delRefundAmt)
		totalRefundShares = totalRefundShares.Add(delRefundShares)
	}

	return totalRefundedAmt, totalRefundShares
}

// SetRefund set a specific refund in the store from its index
func (k Keeper) SetRefund(ctx sdk.Context, refund types.Refund) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	b := k.cdc.MustMarshal(&refund)
	store.Set(types.RefundKey(
		refund.MustGetDelegatorAddr(),
		refund.MustGetValidatorAddr(),
	), b)
}

// GetRefund returns a refund from its index
func (k Keeper) GetRefund(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (refund types.Refund, found bool) {
	moduleStore := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(moduleStore, types.KeyPrefix(types.RefundKeyPrefix))
	key := types.RefundKey(delAddr, valAddr)
	b := store.Get(key)
	if b == nil {
		return refund, false
	}

	k.cdc.MustUnmarshal(b, &refund)
	return refund, true
}

// RemoveRefund removes a refund from the store
func (k Keeper) RemoveRefund(
	ctx sdk.Context,
	refund types.Refund,
) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	store.Delete(types.RefundKey(
		refund.MustGetDelegatorAddr(),
		refund.MustGetValidatorAddr(),
	))
}

// GetAllRefund returns all refund
func (k Keeper) GetAllRefund(ctx sdk.Context) (list []types.Refund) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Refund
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// -------------------------------------------------------------------------------------------------
// Refund pool
// -------------------------------------------------------------------------------------------------

// SetRefundPool set a specific refundPool in the store from its index
func (k Keeper) SetRefundPool(ctx sdk.Context, refundPool types.RefundPool) {

	valAddr, err := sdk.ValAddressFromBech32(refundPool.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	b := k.cdc.MustMarshal(&refundPool)
	store.Set(types.RefundPoolKey(valAddr), b)
}

// GetRefundPool returns a refundPool from its index
func (k Keeper) GetRefundPool(
	ctx sdk.Context,
	valAddr sdk.ValAddress,

) (refPool types.RefundPool, found bool) {

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	b := store.Get(types.RefundPoolKey(valAddr))
	if b == nil {
		return refPool, false
	}
	k.cdc.MustUnmarshal(b, &refPool)
	return refPool, true
}

// RemoveRefundPool removes a refundPool from the store
func (k Keeper) RemoveRefundPool(
	ctx sdk.Context,
	valAddr sdk.ValAddress,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	store.Delete(types.RefundPoolKey(valAddr))
}

// GetAllRefundPool returns all refundPool
func (k Keeper) GetAllRefundPool(ctx sdk.Context) (list []types.RefundPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.RefundPoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var refPool types.RefundPool
		k.cdc.MustUnmarshal(iterator.Value(), &refPool)
		list = append(list, refPool)
	}

	return
}

func (k Keeper) AddRefPoolTokensAndShares(
	ctx sdk.Context,
	refundPool types.RefundPool,
	tokensToAdd sdk.Coin,
) (addedShares sdk.Dec) {

	var issuedShares sdk.Dec
	if refundPool.Shares.IsZero() {
		issuedShares = sdk.NewDecFromInt(tokensToAdd.Amount)
	} else {
		// TODO: we have to manage post slashing send of tokens. We have to put zero shares when  tokens -> 0
		shares, err := refundPool.SharesFromTokens(tokensToAdd)
		if err != nil {
			panic(err)
		}
		issuedShares = shares
	}

	refundPool.Tokens = refundPool.Tokens.Add(tokensToAdd)
	refundPool.Shares = refundPool.Shares.Add(issuedShares)

	k.SetRefundPool(ctx, refundPool)

	return issuedShares
}

// TODO: generalize it considering AllowedTokens param
func (k Keeper) RemoveRefPoolTokensAndShares(
	ctx sdk.Context,
	refundPool types.RefundPool,
	sharesToRemove sdk.Dec,
) (types.RefundPool, math.Int) {

	var issuedTokensAmt math.Int

	remainingShares := refundPool.Shares.Sub(sharesToRemove)

	if remainingShares.IsZero() {
		// last share gets any trimmings
		issuedTokensAmt = refundPool.Tokens.Amount
		refundPool.Tokens.Amount = sdk.ZeroInt()

	} else {
		// leave excess tokens in the refund pool
		issuedTokensAmt = refundPool.TokensFromShares(sharesToRemove).TruncateInt()
		refundPool.Tokens.Amount = refundPool.Tokens.Amount.Sub(issuedTokensAmt)
		if refundPool.Tokens.Amount.IsNegative() {
			panic("attempting to remove more tokens than available in refund pool")
		}
	}

	refundPool.Shares = remainingShares
	k.SetRefundPool(ctx, refundPool)

	return refundPool, issuedTokensAmt
}
