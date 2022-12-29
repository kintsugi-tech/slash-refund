package keeper_test

import (
	"fmt"
	"testing"

	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/stretchr/testify/require"
)

func CreateNTestAccounts(srApp *app.App, ctx sdk.Context, N int, initAmt sdk.Int) ([]sdk.AccAddress, []cryptotypes.PubKey) {
	initCoins := sdk.NewCoins(sdk.NewCoin(srApp.StakingKeeper.BondDenom(ctx), initAmt))
	pks := simapp.CreateTestPubKeys(N)
	addrs := make([]sdk.AccAddress, 0, N)
	for _, pk := range pks {
		addr := sdk.AccAddress(pk.Address())
		addrs = append(addrs, addr)
		err := srApp.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins)
		if err != nil {
			panic(err)
		}
		err = srApp.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, initCoins)
		if err != nil {
			panic(err)
		}
	}

	return addrs, pks
}
func TestProcessSlashEventDoubleSign(t *testing.T) {

	// init state
	srApp, ctx := testsuite.CreateTestApp(false)
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)

	var Nacc int = 1
	var units int64 = 1e6

	initAmt := sdk.NewInt(int64(1000 * units))
	addrs, pks := CreateNTestAccounts(srApp, ctx, Nacc, initAmt)

	operator := addrs[0]
	selfDelegation := sdk.NewInt(100 * units)

	burnedAmtExpectedDS := srApp.SlashingKeeper.SlashFractionDoubleSign(ctx).MulInt(selfDelegation).TruncateInt()

	// create validator
	powerReduction := srApp.StakingKeeper.PowerReduction(ctx)
	sth.CreateValidatorWithValPower(sdk.ValAddress(operator), pks[0], selfDelegation.Quo(powerReduction).Int64(), true)
	validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(operator))
	require.True(t, found)
	valAddr := validator.GetOperator()

	consPow := sdk.TokensToConsensusPower(sdk.NewInt(100*units), sdk.DefaultPowerReduction)

	slashEventDS := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, sdk.ConsAddress(operator).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPow)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueDoubleSign),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, burnedAmtExpectedDS.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", 0)),
	)

	// Double sign
	gotValAddr, valBurnedAmt, infractionHeight, gotSlashFactor, err := srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEventDS)
	require.NoError(t, err)
	require.Equal(t, valAddr, gotValAddr)
	require.Equal(t, srApp.SlashingKeeper.SlashFractionDoubleSign(ctx), gotSlashFactor)
	require.Equal(t, sdk.NewInt(0), infractionHeight)
	require.Equal(t, burnedAmtExpectedDS.String(), valBurnedAmt.String())
}
func TestProcessSlashEventDownTime(t *testing.T) {

	// init state
	srApp, ctx := testsuite.CreateTestApp(false)
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)

	var Nacc int = 1
	var units int64 = 1e6

	initAmt := sdk.NewInt(int64(1000 * units))
	addrs, pks := CreateNTestAccounts(srApp, ctx, Nacc, initAmt)

	operator := addrs[0]
	selfDelegation := sdk.NewInt(100 * units)

	burnedAmtExpectedDT := srApp.SlashingKeeper.SlashFractionDowntime(ctx).MulInt(selfDelegation).TruncateInt()

	// create validator
	powerReduction := srApp.StakingKeeper.PowerReduction(ctx)
	sth.CreateValidatorWithValPower(sdk.ValAddress(operator), pks[0], selfDelegation.Quo(powerReduction).Int64(), true)
	validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(operator))
	require.True(t, found)
	valAddr := validator.GetOperator()

	consPow := sdk.TokensToConsensusPower(sdk.NewInt(100*units), sdk.DefaultPowerReduction)

	slashEventDT := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, sdk.ConsAddress(operator).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPow)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueMissingSignature),
		sdk.NewAttribute(slashingtypes.AttributeKeyJailed, sdk.ConsAddress(operator).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, burnedAmtExpectedDT.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", 0)),
	)

	// Downtime
	gotValAddr, valBurnedAmt, infractionHeight, gotSlashFactor, err := srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEventDT)
	require.NoError(t, err)
	require.Equal(t, valAddr, gotValAddr)
	require.Equal(t, srApp.SlashingKeeper.SlashFractionDowntime(ctx), gotSlashFactor)
	require.Equal(t, sdk.NewInt(0), infractionHeight)
	require.Equal(t, burnedAmtExpectedDT.String(), valBurnedAmt.String())
}

func TestProcessSlashEventErrors(t *testing.T) {

	// init state
	srApp, ctx := testsuite.CreateTestApp(false)

	//Error wrong validator address
	slashEvent := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, "not a validator consensous address"),
	)
	_, _, _, _, err := srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEvent)
	require.ErrorIs(t, err, types.ErrCantGetValidatorFromSlashEvent)

	//Error unknown slashing reason
	slashEvent = sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, "not a double sign"),
	)
	_, _, _, _, err = srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEvent)
	require.ErrorIs(t, err, types.ErrUnknownSlashingReasonFromSlashEvent)

	//Error in converting burned tokens into a number
	slashEvent = sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, "not a number"),
	)
	_, _, _, _, err = srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEvent)
	require.ErrorIs(t, err, types.ErrCantGetValidatorBurnedTokensFromSlashEvent)

	//Error in converting infraction height into a number
	slashEvent = sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, "not a number"),
	)
	_, _, _, _, err = srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEvent)
	require.ErrorIs(t, err, types.ErrCantGetInfractionHeightFromSlashEvent)
}

func TestRefundFromSlash(t *testing.T) {
	// init state
	srApp, ctx := testsuite.CreateTestApp(false)
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)
	srh := testslashrefund.NewHelper(t, srApp.SlashrefundKeeper, ctx)

	var Nacc int = 5
	var units int64 = 1e6

	initAmt := sdk.NewInt(int64(1000 * units))
	addrs, pks := CreateNTestAccounts(srApp, ctx, Nacc, initAmt)

	operator := addrs[0]
	depositor1 := addrs[0]
	depositor2 := addrs[1]
	delegator1 := addrs[2]
	delegator2 := addrs[3]
	stranger := addrs[4]

	selfDelegation := sdk.NewInt(10 * units)
	delAmt1 := sdk.NewInt(50 * units)
	delAmt2 := sdk.NewInt(40 * units)

	depAmt1 := sdk.NewInt(8 * units)
	depAmt2 := sdk.NewInt(2 * units)

	slashFactor := sdk.NewDec(5).QuoInt(sdk.NewInt(100))
	slashAmt := sdk.NewInt(5 * units)

	refAmtExpected := sdk.NewInt(5 * units)
	refShrExpected := sdk.NewDec(5 * units)

	// create validator
	powerReduction := srApp.StakingKeeper.PowerReduction(ctx)
	sth.CreateValidatorWithValPower(sdk.ValAddress(operator), pks[0], selfDelegation.Quo(powerReduction).Int64(), true)
	validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(operator))
	require.True(t, found)
	valAddr := validator.GetOperator()

	// ==== new block ====
	sth.TurnBlock(ctx.BlockTime().Add(1))

	validator, found = srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(valAddr))
	require.True(t, found)

	//check status
	require.Equal(t, stakingtypes.BondStatusBonded, validator.GetStatus().String())

	//check stacked: selfdelegation = 10 -> bondedTokens = 10
	require.Equal(t,
		selfDelegation.String(),
		validator.GetBondedTokens().String())

	// check operator balance
	require.Equal(t,
		sdk.NewCoins(sdk.NewCoin(srApp.StakingKeeper.BondDenom(ctx), initAmt.Sub(selfDelegation))).String(),
		srApp.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(validator.GetOperator())).String(),
	)

	// delegate
	sth.Delegate(delegator1, valAddr, delAmt1)
	sth.Delegate(delegator2, valAddr, delAmt2)

	// deposit
	srh.Deposit(depositor1, valAddr, depAmt1)
	srh.Deposit(depositor2, valAddr, depAmt2)

	// ==== new block ====
	sth.TurnBlock(ctx.BlockTime().Add(1))

	// check validator
	consAddr, err := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsAddr()
	require.NoError(t, err)
	validator, found = srApp.StakingKeeper.GetValidatorByConsAddr(ctx, consAddr)
	require.True(t, found)
	require.Equal(t, selfDelegation.Add(delAmt1).Add(delAmt2).String(), validator.GetBondedTokens().String())

	// check deposit pool
	_, found = srApp.SlashrefundKeeper.GetDepositPool(ctx, valAddr)
	require.True(t, found)

	// check deposit 1
	deposit1, found := srApp.SlashrefundKeeper.GetDeposit(ctx, depositor1, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt1), deposit1.Shares)

	// check deposit 2
	deposit2, found := srApp.SlashrefundKeeper.GetDeposit(ctx, depositor2, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt2), deposit2.Shares)

	// Check refund pool
	_, found = srApp.SlashrefundKeeper.GetRefundPool(ctx, valAddr)
	require.False(t, found)

	// Sub tokens to validator as it was slashed
	validator.Tokens = validator.Tokens.Sub(slashAmt)
	srApp.StakingKeeper.SetValidator(ctx, validator)
	srApp.StakingKeeper.SetValidatorByPowerIndex(ctx, validator)

	// Refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, 0, slashFactor)
	require.NoError(t, err)
	require.Equal(t, refAmtExpected.String(), refAmt.String())

	refPool, found := srApp.SlashrefundKeeper.GetRefundPool(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, refAmt, refPool.Tokens.Amount)

	//check issued shares to stranger
	_, found = srApp.SlashrefundKeeper.GetRefund(ctx, stranger, valAddr)
	require.False(t, found)

	// check self-delegation issued shares
	ref, found := srApp.SlashrefundKeeper.GetRefund(ctx, operator, valAddr)
	require.True(t, found)
	require.Equal(t, ref.Shares, sdk.NewDecFromInt(selfDelegation).Mul(slashFactor))
	refSharesTotal := sdk.NewDec(0)
	refSharesTotal = refSharesTotal.Add(ref.Shares)

	//check delegation 1 issued shares
	ref, found = srApp.SlashrefundKeeper.GetRefund(ctx, delegator1, valAddr)
	require.True(t, found)
	require.Equal(t, ref.Shares, sdk.NewDecFromInt(delAmt1).Mul(slashFactor))
	refSharesTotal = refSharesTotal.Add(ref.Shares)

	// check delegation 2 issued shares
	ref, found = srApp.SlashrefundKeeper.GetRefund(ctx, delegator2, valAddr)
	require.True(t, found)
	require.Equal(t, ref.Shares, sdk.NewDecFromInt(delAmt2).Mul(slashFactor))
	refSharesTotal = refSharesTotal.Add(ref.Shares)

	// check total issued shares
	require.Equal(t, refPool.Shares, refSharesTotal)
	require.Equal(t, refPool.Shares, refShrExpected)

	// check deposit pool
	depPool, found := srApp.SlashrefundKeeper.GetDepositPool(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewInt(5*units), depPool.Tokens.Amount)

	// check deposit 1
	deposit1, found = srApp.SlashrefundKeeper.GetDeposit(ctx, depositor1, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt1), deposit1.Shares)

	// check deposit 2
	deposit2, found = srApp.SlashrefundKeeper.GetDeposit(ctx, depositor2, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt2), deposit2.Shares)
}
