package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	abci "github.com/tendermint/tendermint/abci/types"
	//"github.com/tendermint/tendermint/libs/bytes"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/stretchr/testify/require"
)

var units int64 = sdk.DefaultPowerReduction.Int64()

// bootstrapRefundTest creates a validator with given power and bootstrap the app
func bootstrapRefundTest(t *testing.T, power int64) *KeeperTestSuite {

	srApp, ctx := testsuite.CreateTestApp(false)

	initAmt := sdk.NewInt(int64(1000 * units))
	testAddrs, pubks := CreateNTestAccounts(srApp, ctx, 5, initAmt)

	selfDelegation := sdk.NewInt(power * units)

	// create 2 validators with consensous power equal to input power
	powerReduction := srApp.StakingKeeper.PowerReduction(ctx)
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)

	valAddrs := make([]sdk.ValAddress, 0, 2)

	for i := 0; i < 2; i++ {
		sth.CreateValidatorWithValPower(sdk.ValAddress(testAddrs[i]), pubks[i], selfDelegation.Quo(powerReduction).Int64(), true)
		validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(testAddrs[i]))
		require.True(t, found)
		valAddr := validator.GetOperator()
		sd, found := srApp.StakingKeeper.GetDelegation(ctx, testAddrs[i], valAddr)
		require.True(t, found)
		require.Equal(t, selfDelegation, sd.Shares.TruncateInt())
		valAddrs = append(valAddrs, valAddr)
	}

	s := KeeperTestSuite{}
	s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.t = srApp, ctx, testAddrs, valAddrs, selfDelegation, t

	return &s
}

func defaultTestValues() (power int64, slashFactor sdk.Dec, slashAmt sdk.Int, depositAmt sdk.Int, infractionHeight int64, slashTime time.Time) {

	power = int64(100)
	slashFactor = sdk.NewDecWithPrec(5, 2)
	slashAmt = sdk.NewInt(5 * units)
	depositAmt = sdk.NewInt(10 * units)
	infractionHeight = int64(10)
	slashTime = time.Unix(100, 0)

	return power, slashFactor, slashAmt, depositAmt, infractionHeight, slashTime

}

func TestProcessSlashEvent_DoubleSign(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs, selfDelegation := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation

	validator, found := srApp.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)
	consPower := validator.ConsensusPower(srApp.StakingKeeper.PowerReduction(ctx))

	expectedBurned := srApp.SlashingKeeper.SlashFractionDoubleSign(ctx).MulInt(selfDelegation).TruncateInt()

	slashEventDS := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, sdk.ConsAddress(testAddrs[0]).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueDoubleSign),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, expectedBurned.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", 12345)),
	)

	// Double sign
	gotValAddr, valBurnedAmt, infractionHeight, gotSlashFactor, err := srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEventDS)
	require.NoError(t, err)
	require.Equal(t, valAddrs[0], gotValAddr)
	require.Equal(t, srApp.SlashingKeeper.SlashFractionDoubleSign(ctx), gotSlashFactor)
	require.Equal(t, sdk.NewInt(12345), infractionHeight)
	require.Equal(t, expectedBurned.String(), valBurnedAmt.String())
}

func TestProcessSlashEvent_DownTime(t *testing.T) {

	s := bootstrapRefundTest(t, 100)
	srApp, ctx, testAddrs, valAddrs, selfDelegation := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation

	validator, found := srApp.StakingKeeper.GetValidator(ctx, valAddrs[0])
	require.True(t, found)
	consPower := validator.ConsensusPower(srApp.StakingKeeper.PowerReduction(ctx))

	expectedBurned := srApp.SlashingKeeper.SlashFractionDoubleSign(ctx).MulInt(selfDelegation).TruncateInt()

	slashEventDT := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, sdk.ConsAddress(testAddrs[0]).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueMissingSignature),
		sdk.NewAttribute(slashingtypes.AttributeKeyJailed, sdk.ConsAddress(testAddrs[0]).String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, expectedBurned.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", 12345)),
	)

	// Downtime
	gotValAddr, valBurnedAmt, infractionHeight, gotSlashFactor, err := srApp.SlashrefundKeeper.ProcessSlashEvent(ctx, slashEventDT)
	require.NoError(t, err)
	require.Equal(t, valAddrs[0], gotValAddr)
	require.Equal(t, srApp.SlashingKeeper.SlashFractionDowntime(ctx), gotSlashFactor)
	require.Equal(t, sdk.NewInt(12345), infractionHeight)
	require.Equal(t, expectedBurned.String(), valBurnedAmt.String())
}

func TestProcessSlashEvent_Errors(t *testing.T) {

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

func TestRefundFromSlash_CurrentHeight(t *testing.T) {

	power, slashFactor, slashAmt, depAmt, infractionHeight, _ := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	// deposit
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// set ctx at infraction height
	ctx = ctx.WithBlockHeight(infractionHeight)

	// set unbonding deposit (should be ignored because infractionHeight = currentHeight)
	ubd := types.NewUnbondingDeposit(testAddrs[1], valAddr, infractionHeight, ctx.BlockTime().Add(1), depAmt)
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubd)

	// set unbonding delegation (should be ignored because infractionHeight = currentHeight)
	ubdel := stakingtypes.NewUnbondingDelegation(testAddrs[1], valAddr, infractionHeight, ctx.BlockTime().Add(1), depAmt)
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubdel)

	// set redelegation (should be ignored because infractionHeight = currentHeight)
	redel := stakingtypes.NewRedelegation(testAddrs[1], valAddr, valAddrs[1], infractionHeight, ctx.BlockTime().Add(1), depAmt, sdk.NewDecFromInt(depAmt))
	srApp.StakingKeeper.SetRedelegation(ctx, redel)

	// call refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, slashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	s.RequireNoRefund(testAddrs[2], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, sdk.NewDecFromInt(selfDelegation).Mul(slashFactor))

	// check refund pool
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(slashAmt), []types.Refund{refund0})

	// check deposit
	deposit := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pool
	s.RequireDepositPool(valAddr, refAmt, sdk.NewDecFromInt(depAmt), []types.Deposit{deposit})

	//  check unbonding deposit
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[1], valAddr)
	require.True(t, found)
	require.Equal(t, 1, len(ubd.Entries))
	require.Equal(t, depAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, depAmt, ubd.Entries[0].Balance)

	//  check unbonding delegation
	ubdel, found = srApp.StakingKeeper.GetUnbondingDelegation(ctx, testAddrs[1], valAddr)
	require.True(t, found)
	require.Equal(t, 1, len(ubdel.Entries))
	require.Equal(t, depAmt, ubdel.Entries[0].InitialBalance)
	require.Equal(t, depAmt, ubdel.Entries[0].Balance)

	//  check redelegation
	redel, found = srApp.StakingKeeper.GetRedelegation(ctx, testAddrs[1], valAddr, valAddrs[1])
	require.True(t, found)
	require.Equal(t, 1, len(redel.Entries))
	require.Equal(t, depAmt, ubdel.Entries[0].InitialBalance)
	require.Equal(t, depAmt, ubdel.Entries[0].Balance)
}

func TestRefundFromSlash_EqualToDepositPool(t *testing.T) {

	power, slashFactor, slashAmt, _, _, _ := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, selfDelegation, valAddr, valAddrs := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.valAddrs

	// deposit
	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)
	srh.Deposit(testAddrs[0], valAddr, sdk.NewInt(2*units))
	srh.Deposit(testAddrs[1], valAddr, sdk.NewInt(2*units))
	srh.Deposit(testAddrs[2], valAddr, sdk.NewInt(1*units))

	// depposits to other validator
	srh.Deposit(testAddrs[0], valAddrs[1], sdk.NewInt(2*units))
	srh.Deposit(testAddrs[1], valAddrs[1], sdk.NewInt(2*units))

	// call refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, 0, slashFactor)
	require.NoError(t, err)
	require.Equal(t, slashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, sdk.NewDecFromInt(selfDelegation).Mul(slashFactor))

	// check refund pool
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(slashAmt), []types.Refund{refund0})

	// check deposit
	s.RequireNoDeposit(testAddrs[0], valAddr)
	s.RequireNoDeposit(testAddrs[1], valAddr)
	s.RequireNoDeposit(testAddrs[2], valAddr)
	deposit1 := s.RequireDeposit(testAddrs[0], valAddrs[1], sdk.NewDec(2*units))
	deposit2 := s.RequireDeposit(testAddrs[1], valAddrs[1], sdk.NewDec(2*units))

	// check deposit pool
	s.RequireNoDepositPool(valAddr)
	s.RequireDepositPool(valAddrs[1], sdk.NewInt(4*units), sdk.NewDec(4*units), []types.Deposit{deposit1, deposit2})

}

func TestRefundFromSlash_GTDepositPool(t *testing.T) {

	power, slashFactor, slashAmt, _, _, _ := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs[0]

	// deposit
	depAmt := sdk.NewInt(1 * units)
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// call refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, 0, slashFactor)
	require.NoError(t, err)
	require.Equal(t, depAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check refund pool
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(depAmt), []types.Refund{refund0})

	// check deposit
	s.RequireNoDeposit(testAddrs[0], valAddr)

	// check deposit pool
	s.RequireNoDepositPool(valAddr)

}

func TestRefundFromSlash_MultipleDelegations(t *testing.T) {

	_, slashFactor, slashAmt, depAmt, _, _ := defaultTestValues()
	s := bootstrapRefundTest(t, 10)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	delAmt1 := sdk.NewInt(50 * units)
	delAmt2 := sdk.NewInt(40 * units)

	operator, delegator1, delegator2, stranger := testAddrs[0], testAddrs[1], testAddrs[2], testAddrs[3]

	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)
	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)

	srh.Deposit(operator, valAddr, depAmt)

	sth.Delegate(delegator1, valAddr, delAmt1)
	sth.Delegate(delegator2, valAddr, delAmt2)

	// delegation to not slashed validator, must not be refunded
	sth.Delegate(delegator1, valAddrs[1], delAmt1)

	// call refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, 0, slashFactor)
	require.NoError(t, err)
	require.Equal(t, slashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(stranger, valAddr)
	s.RequireNoRefund(delegator1, valAddrs[1])
	refund0 := s.RequireRefund(operator, valAddr, slashFactor.MulInt(selfDelegation))
	refund1 := s.RequireRefund(delegator1, valAddr, slashFactor.MulInt(delAmt1))
	refund2 := s.RequireRefund(delegator2, valAddr, slashFactor.MulInt(delAmt2))

	// check refund pools
	s.RequireNoRefundPool(valAddrs[1])
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(slashAmt), []types.Refund{refund0, refund1, refund2})

	// check deposit
	deposit := s.RequireDeposit(operator, valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pool
	s.RequireDepositPool(valAddr, sdk.NewInt(5*units), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit})

}

func TestRefundFromSlash_MultipleDeposits(t *testing.T) {

	power, slashFactor, slashAmt, _, _, _ := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	depAmt1 := sdk.NewInt(8 * units)
	depAmt2 := sdk.NewInt(2 * units)

	operator, depositor2, stranger := testAddrs[0], testAddrs[1], testAddrs[2]

	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)

	srh.Deposit(operator, valAddr, depAmt1)
	srh.Deposit(depositor2, valAddr, depAmt2)

	// deposit for not slashed validator: must not be accounted for
	srh.Deposit(operator, valAddrs[1], depAmt1)

	// call refund from slash
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, 0, slashFactor)
	require.NoError(t, err)
	require.Equal(t, slashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(stranger, valAddr)
	s.RequireNoRefund(operator, valAddrs[1])
	refund0 := s.RequireRefund(operator, valAddr, slashFactor.MulInt(selfDelegation))

	// check refund pools
	s.RequireNoRefundPool(valAddrs[1])
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(slashAmt), []types.Refund{refund0})

	// check deposits
	deposit1 := s.RequireDeposit(operator, valAddr, sdk.NewDecFromInt(depAmt1))
	deposit2 := s.RequireDeposit(depositor2, valAddr, sdk.NewDecFromInt(depAmt2))
	deposit3 := s.RequireDeposit(operator, valAddrs[1], sdk.NewDecFromInt(depAmt1))

	// check deposit pools
	s.RequireDepositPool(valAddr, sdk.NewInt(5*units), sdk.NewDecFromInt(depAmt1.Add(depAmt2)), []types.Deposit{deposit1, deposit2})
	s.RequireDepositPool(valAddrs[1], depAmt1, sdk.NewDecFromInt(depAmt1), []types.Deposit{deposit3})
}

func TestRefundFromSlash_NotEligibleUnbondingDeposits(t *testing.T) {

	power, slashFactor, slashAmt, depAmt, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs[0]

	// first entry uneligible due to creation height
	ctx = ctx.WithBlockHeight(infractionHeight - 1)
	ubd := types.NewUnbondingDeposit(testAddrs[0], valAddr, ctx.BlockHeight(), slashTime.Add(100), depAmt)
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubd)

	// second entry uneligible due to completion time
	ctx = ctx.WithBlockHeight(infractionHeight + 1)
	ubd.Entries = append(ubd.Entries, types.NewUnbondingDepositEntry(ctx.BlockHeight(), slashTime, depAmt))
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubd)

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, sdk.NewInt(0), refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	s.RequireNoRefund(testAddrs[0], valAddr)

	// check refund pool
	s.RequireNoRefundPool(valAddr)

	//  check unbonding deposit
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[0], valAddr)
	require.True(t, found)
	require.Equal(t, 2, len(ubd.Entries))
	require.Equal(t, depAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, depAmt, ubd.Entries[0].Balance)
	require.Equal(t, depAmt, ubd.Entries[1].InitialBalance)
	require.Equal(t, depAmt, ubd.Entries[1].Balance)
}

func TestRefundFromSlash_EligibleUnbondingDeposit(t *testing.T) {

	power, slashFactor, slashAmt, depAmt, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0]

	// Eligible Unbonding Deposit
	ctx = ctx.WithBlockHeight(infractionHeight + 1)
	ubd := types.NewUnbondingDeposit(testAddrs[0], valAddr, ctx.BlockHeight(), slashTime.Add(100), depAmt)
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubd)

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, slashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, slashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))

	// check refund pool
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(slashAmt), []types.Refund{refund0})

	//  check unbonding deposit
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[0], valAddr)
	require.True(t, found)
	require.Equal(t, depAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, sdk.NewInt(5*units), ubd.Entries[0].Balance)
}

func TestRefundFromSlash_NotEligibleUnbondingDelegations(t *testing.T) {

	power, slashFactor, valSlashAmt, depAmt, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0]

	// deposit
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// unbonding delegation 1: not eligible
	//  entry 1 n.e. due to creation height
	//  entry 2 n.e. due to completion time
	ubd := stakingtypes.NewUnbondingDelegation(testAddrs[1], valAddr, infractionHeight-1, slashTime.Add(100), depAmt)
	ubd.Entries = append(ubd.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			infractionHeight+1,
			slashTime,
			depAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	// unbonding delegation 2:
	//  entry 1 n.e. due to completion time
	//  entry 2 n.e. due to creation height
	ubd = stakingtypes.NewUnbondingDelegation(testAddrs[2], valAddr, infractionHeight+1, slashTime, depAmt)
	ubd.Entries = append(ubd.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			depAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	require.Equal(t, 2, len(srApp.StakingKeeper.GetUnbondingDelegationsFromValidator(ctx, valAddr)))

	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, valSlashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, valSlashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	s.RequireNoRefund(testAddrs[2], valAddr)
	s.RequireNoRefund(testAddrs[3], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))

	// check refund pool
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(valSlashAmt), []types.Refund{refund0})

	// check deposits
	deposit1 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pools
	s.RequireDepositPool(valAddr, sdk.NewInt(5*units), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit1})

}

func TestRefundFromSlash_EligibleUnbondingDelegations(t *testing.T) {

	power, slashFactor, valSlashAmt, _, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	// will set 2 unbonding delegations: ubd1 and ubd2, each with 2 eligible entries of 100 units each
	ubdelAmt := sdk.NewInt(100 * units)

	depAmt := sdk.NewInt(50 * units)
	refAmtExpected := sdk.NewInt(25 * units)

	// deposit
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// eligible unbonding delegation 1: two eligible entries with initial balance 100 units
	ctx = ctx.WithBlockHeight(infractionHeight + 1)

	ubd := stakingtypes.NewUnbondingDelegation(testAddrs[1], valAddr, ctx.BlockHeight(), slashTime.Add(100), ubdelAmt)
	ubd.Entries = append(ubd.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			ctx.BlockHeight(), slashTime.Add(100), ubdelAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	// eligible unbonding delegation 2: two eligible entries with initial balance 100 units
	ubd = stakingtypes.NewUnbondingDelegation(testAddrs[2], valAddr, ctx.BlockHeight(), slashTime.Add(100), ubdelAmt)
	ubd.Entries = append(ubd.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			ctx.BlockHeight(), slashTime.Add(100), ubdelAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubd)

	// unbonding delegation 3: not from the slashed validator, must not be refunded
	ubd = stakingtypes.NewUnbondingDelegation(testAddrs[3], valAddrs[1], ctx.BlockHeight(), slashTime.Add(100), ubdelAmt)

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, valSlashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, refAmtExpected, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[3], valAddr)
	s.RequireNoRefund(testAddrs[3], valAddrs[1])
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))
	refund1 := s.RequireRefund(testAddrs[1], valAddr, slashFactor.MulInt(ubdelAmt.MulRaw(2)))
	refund2 := s.RequireRefund(testAddrs[2], valAddr, slashFactor.MulInt(ubdelAmt.MulRaw(2)))

	// check refund pool
	s.RequireNoRefundPool(valAddrs[1])
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(refAmtExpected), []types.Refund{refund0, refund1, refund2})

	// check deposits
	deposit1 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pools
	s.RequireDepositPool(valAddr, sdk.NewInt(25*units), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit1})
}

func TestRefundFromSlash_EligibleRedelegations(t *testing.T) {

	power, slashFactor, valSlashAmt, _, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	// will set 2 redelegations: ubd1 and ubd2, each with 2 eligible entries of 100 units each
	amt := sdk.NewInt(100 * units)

	depAmt := sdk.NewInt(50 * units)
	refAmtExpected := sdk.NewInt(25 * units)

	// deposit
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// not eligible redelegation 1: two eligible entries with initial balance 100 units
	ctx = ctx.WithBlockHeight(infractionHeight + 1)

	red := stakingtypes.NewRedelegation(testAddrs[1], valAddr, valAddrs[1], ctx.BlockHeight(), slashTime.Add(100), amt, sdk.NewDecFromInt(amt))
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			ctx.BlockHeight(), slashTime.Add(100), amt, sdk.NewDecFromInt(amt),
		))
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// eligible redelegation 2: two eligible entries with initial balance 100 units
	red = stakingtypes.NewRedelegation(testAddrs[2], valAddr, valAddrs[1], ctx.BlockHeight(), slashTime.Add(100), amt, sdk.NewDecFromInt(amt))
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			ctx.BlockHeight(), slashTime.Add(100), amt, sdk.NewDecFromInt(amt),
		))
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// redelegation 3: not from the slashed validator, must not be refunded
	red = stakingtypes.NewRedelegation(testAddrs[3], valAddrs[1], valAddr, ctx.BlockHeight(), slashTime.Add(100), amt, sdk.NewDecFromInt(amt))
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, valSlashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, refAmtExpected, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[3], valAddr)
	s.RequireNoRefund(testAddrs[3], valAddrs[1])
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))
	refund1 := s.RequireRefund(testAddrs[1], valAddr, slashFactor.MulInt(amt.MulRaw(2)))
	refund2 := s.RequireRefund(testAddrs[2], valAddr, slashFactor.MulInt(amt.MulRaw(2)))

	// check refund pool
	s.RequireNoRefundPool(valAddrs[1])
	s.RequireRefundPool(valAddr, refAmt, sdk.NewDecFromInt(refAmtExpected), []types.Refund{refund0, refund1, refund2})

	// check deposits
	deposit1 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pools
	s.RequireDepositPool(valAddr, sdk.NewInt(25*units), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit1})
}

func TestRefundFromSlash_NotEligibleRedelegations(t *testing.T) {

	power, slashFactor, valSlashAmt, depAmt, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0]

	// deposit
	testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper).Deposit(testAddrs[0], valAddr, depAmt)

	// redelegation 1:
	//  entry 1 n.e. due to creation height
	//  entry 2 n.e. due to completion time
	red := stakingtypes.NewRedelegation(testAddrs[1], valAddr, valAddrs[1],
		infractionHeight-1,
		slashTime.Add(100),
		depAmt, sdk.NewDecFromInt(depAmt),
	)
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			infractionHeight+1,
			slashTime,
			depAmt, sdk.NewDecFromInt(depAmt)),
	)
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// redelegation 2:
	//  entry 1 n.e. due to completion time
	//  entry 2 n.e. due to creation height
	red = stakingtypes.NewRedelegation(testAddrs[2], valAddr, valAddrs[1],
		infractionHeight+1,
		slashTime,
		depAmt, sdk.NewDecFromInt(depAmt),
	)
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			depAmt, sdk.NewDecFromInt(depAmt)),
	)
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	require.Equal(t, 2, len(srApp.StakingKeeper.GetRedelegationsFromSrcValidator(ctx, valAddr)))

	// call refund from slash
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	refAmt, err := srApp.SlashrefundKeeper.RefundFromSlash(ctx, valAddr, valSlashAmt, infractionHeight, slashFactor)
	require.NoError(t, err)
	require.Equal(t, valSlashAmt, refAmt)

	// check refunds
	s.RequireNoRefund(testAddrs[1], valAddr)
	s.RequireNoRefund(testAddrs[2], valAddr)
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))

	// check refund pool
	s.RequireNoRefundPool(valAddrs[1])
	s.RequireRefundPool(valAddr, valSlashAmt, sdk.NewDecFromInt(valSlashAmt), []types.Refund{refund0})

	// check deposit
	deposit1 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))

	// check deposit pool
	s.RequireDepositPool(valAddr, sdk.NewInt(5*units), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit1})
}

func TestHandleRefundsFromSlashDoubleSign(t *testing.T) {
	//TODO
}

func TestHandleRefundsFromSlashDownTime(t *testing.T) {
	//TODO
}

func TestSlashRefundDoubleSign(t *testing.T) {
	// init state
	srApp, ctx := testsuite.CreateTestApp(false)
	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)
	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)

	initAmt := sdk.NewInt(int64(1000 * units))
	addrs, pks := CreateNTestAccounts(srApp, ctx, 5, initAmt)

	operator := addrs[0]
	depositor1 := addrs[0]
	depositor2 := addrs[1]
	delegator1 := addrs[2]
	delegator2 := addrs[3]
	stranger := addrs[4]

	selfDelegation := sdk.NewInt(10 * units)
	delAmt1 := sdk.NewInt(50 * units)
	delAmt2 := sdk.NewInt(40 * units)

	depAmt1 := sdk.NewInt(800 * units)
	depAmt2 := sdk.NewInt(200 * units)

	// create validator
	powerReduction := srApp.StakingKeeper.PowerReduction(ctx)
	sth.CreateValidatorWithValPower(sdk.ValAddress(operator), pks[0], selfDelegation.Quo(powerReduction).Int64(), true)
	validator, found := srApp.StakingKeeper.GetValidatorByConsAddr(ctx, sdk.ConsAddress(operator))
	require.True(t, found)
	valAddr := validator.GetOperator()

	// ==== new block ====
	sth.TurnBlock(ctx.BlockTime().Add(time.Duration(1) * time.Second))

	sth.Delegate(delegator1, valAddr, delAmt1)
	sth.Delegate(delegator2, valAddr, delAmt2)

	srh.Deposit(depositor1, valAddr, depAmt1)
	srh.Deposit(depositor2, valAddr, depAmt2)

	// ==== new block ====
	sth.TurnBlock(ctx.BlockTime().Add(time.Duration(1) * time.Second))

	// slash for double sign
	consAddr, err := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsAddr()
	require.NoError(t, err)
	initialTokAmt := srApp.StakingKeeper.Validator(ctx, valAddr).GetTokens()
	consPower := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsensusPower(srApp.StakingKeeper.PowerReduction(ctx))
	//srApp.SlashingKeeper.HandleValidatorSignature(ctx, bytes.HexBytes(valAddr.Bytes()), selfDelegation.Int64(), true)
	srApp.EvidenceKeeper.HandleEquivocationEvidence(ctx,
		&evidencetypes.Equivocation{
			Height:           0,
			Time:             time.Unix(0, 0),
			Power:            consPower,
			ConsensusAddress: consAddr.String(),
		},
	)

	// compute burned tokens
	sf := srApp.SlashingKeeper.SlashFractionDoubleSign(ctx)
	valBurnedTokens := sf.MulInt(srApp.StakingKeeper.TokensFromConsensusPower(ctx, consPower)).TruncateInt()
	require.Equal(t, initialTokAmt.Sub(valBurnedTokens), srApp.StakingKeeper.Validator(ctx, valAddr).GetTokens())

	// check no refund pool
	_, found = srApp.SlashrefundKeeper.GetRefundPool(ctx, valAddr)
	require.False(t, found)

	// call slash refund
	slashrefund.BeginBlocker(ctx, abci.RequestBeginBlock{}, srApp.SlashrefundKeeper)

	// check refunds
	_, found = srApp.SlashrefundKeeper.GetRefund(ctx, stranger, valAddr)
	require.False(t, found)

	refund0, found := srApp.SlashrefundKeeper.GetRefund(ctx, operator, valAddr)
	require.True(t, found)
	require.Equal(t, refund0.Shares, sdk.NewDecFromInt(selfDelegation).Mul(sf))

	refund1, found := srApp.SlashrefundKeeper.GetRefund(ctx, delegator1, valAddr)
	require.True(t, found)
	require.Equal(t, refund1.Shares, sdk.NewDecFromInt(delAmt1).Mul(sf))

	refund2, found := srApp.SlashrefundKeeper.GetRefund(ctx, delegator2, valAddr)
	require.True(t, found)
	require.Equal(t, refund2.Shares, sdk.NewDecFromInt(delAmt2).Mul(sf))

	// check refund pool
	refPool, found := srApp.SlashrefundKeeper.GetRefundPool(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, valBurnedTokens, refPool.Tokens.Amount)
	require.Equal(t, refund0.Shares.Add(refund1.Shares).Add(refund2.Shares), refPool.Shares)
	require.Equal(t, sdk.NewDecFromInt(valBurnedTokens), refPool.Shares)

	// check deposits
	deposit1, found := srApp.SlashrefundKeeper.GetDeposit(ctx, depositor1, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt1), deposit1.Shares)
	deposit2, found := srApp.SlashrefundKeeper.GetDeposit(ctx, depositor2, valAddr)
	require.True(t, found)
	require.Equal(t, sdk.NewDecFromInt(depAmt2), deposit2.Shares)

	// check deposit pool
	depPool, found := srApp.SlashrefundKeeper.GetDepositPool(ctx, valAddr)
	require.True(t, found)
	require.Equal(t, depPool.Shares, sdk.NewDecFromInt(depAmt1.Add(depAmt2)))
	require.Equal(t, depAmt1.Add(depAmt2), depPool.Tokens.Amount.Add(valBurnedTokens))
}

func TestSlashRefundDownTime(t *testing.T) {
	//TODO
}
