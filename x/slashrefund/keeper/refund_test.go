package keeper_test

/*
import (
	//"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"

	"fmt"
	"time"

	"github.com/made-in-block/slash-refund/testutil/testsuite"

	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// Default initial state for all test. It creates two validators with a specified power.
func bootstrapRefundTest(t *testing.T, power int64) *KeeperTestSuite {
	return SetupTestSuite(t, power)
}

func defaultTestValues() (power int64, slashFactor sdk.Dec, slashAmt sdk.Int, depositAmt sdk.Int, infractionHeight int64, slashTime time.Time) {

	units := sdk.DefaultPowerReduction.Int64()
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

// Test the proper processing of a downtime slashing event.
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
	srApp, ctx, testAddrs, selfDelegation, valAddr, valAddrs, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.valAddrs, s.units

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
	srApp, ctx, testAddrs, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.units

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

func TestRefundFromSlash_MultipleEligibleUnbondingDeposit(t *testing.T) {

	power, slashFactor, slashAmt, _, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.units
	entryAmt := sdk.NewInt(2 * units)

	// Eligible Unbonding Deposits
	ctx = ctx.WithBlockHeight(infractionHeight + 1)

	// first unbonding deposit has two entries with initial balance of 2 units each
	ubd := types.NewUnbondingDeposit(testAddrs[0], valAddr, ctx.BlockHeight(), slashTime.Add(100), entryAmt)
	ubd.Entries = append(ubd.Entries, types.NewUnbondingDepositEntry(ctx.BlockHeight(), slashTime.Add(100), entryAmt))
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubd)

	// second unbonding deposit has three entries with initial balance of 2 units each
	ubd = types.NewUnbondingDeposit(testAddrs[1], valAddr, ctx.BlockHeight(), slashTime.Add(100), entryAmt)
	ubd.Entries = append(ubd.Entries, types.NewUnbondingDepositEntry(ctx.BlockHeight(), slashTime.Add(100), entryAmt))
	ubd.Entries = append(ubd.Entries, types.NewUnbondingDepositEntry(ctx.BlockHeight(), slashTime.Add(100), entryAmt))
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

	//  check unbonding deposits
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[0], valAddr)
	require.True(t, found)
	require.Equal(t, 2, len(ubd.Entries))
	require.Equal(t, entryAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, entryAmt, ubd.Entries[1].InitialBalance)
	require.Equal(t, sdk.NewInt(1*units), ubd.Entries[0].Balance)
	require.Equal(t, sdk.NewInt(1*units), ubd.Entries[1].Balance)

	ubd, found = srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[1], valAddr)
	require.True(t, found)
	require.Equal(t, 3, len(ubd.Entries))
	require.Equal(t, entryAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, entryAmt, ubd.Entries[1].InitialBalance)
	require.Equal(t, entryAmt, ubd.Entries[2].InitialBalance)
	require.Equal(t, sdk.NewInt(1*units), ubd.Entries[0].Balance)
	require.Equal(t, sdk.NewInt(1*units), ubd.Entries[1].Balance)
	require.Equal(t, sdk.NewInt(1*units), ubd.Entries[2].Balance)
}

func TestRefundFromSlash_NotEligibleUnbondingDelegations(t *testing.T) {

	power, slashFactor, valSlashAmt, depAmt, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, power)
	srApp, ctx, testAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0], s.units

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
	srApp, ctx, testAddrs, valAddrs, selfDelegation, valAddr, units := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.selfDelegation, s.valAddrs[0], s.units

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

func TestHandleRefundsFromSlash_DoubleSign(t *testing.T) {

	_, _, _, _, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, 10)
	srApp, ctx, testAddrs, selfDelegation, valAddr, valAddrs, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.valAddrs, s.units

	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)
	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)

	depAmt := sdk.NewInt(150 * units)
	ubdepAmt := sdk.NewInt(50 * units)
	delAmt := sdk.NewInt(90 * units)
	ubdelAmt := sdk.NewInt(5 * units)
	redelAmt := sdk.NewInt(15 * units)

	srh.Deposit(testAddrs[0], valAddr, depAmt)
	sth.Delegate(testAddrs[2], valAddr, delAmt)

	// eligible unbonding deposit
	// entry 1: eligible
	// entry 2: not eligible
	ubdep := types.NewUnbondingDeposit(testAddrs[1], valAddr, infractionHeight+1, slashTime.Add(100), ubdepAmt)
	ubdep.Entries = append(ubdep.Entries,
		types.NewUnbondingDepositEntry(
			infractionHeight-1,
			slashTime.Add(100),
			ubdepAmt,
		))
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep)

	// eligible unbonding delegation
	// entry 1: eligible
	// entry 2: not eligible
	ubdel := stakingtypes.NewUnbondingDelegation(testAddrs[3], valAddr, infractionHeight+1, slashTime.Add(100), ubdelAmt)
	ubdel.Entries = append(ubdel.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			ubdelAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubdel)

	// redelegation:
	// entry 1: eligible
	// entry 2: not eligible
	red := stakingtypes.NewRedelegation(testAddrs[4], valAddr, valAddrs[1], infractionHeight+1, slashTime.Add(100), redelAmt, sdk.NewDec(1))
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			redelAmt, sdk.NewDec(1)),
	)
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// validator update
	ctx = sth.TurnBlock(slashTime)

	// compute burned tokens
	consAddr, err := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsAddr()
	require.NoError(t, err)

	consPower := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsensusPower(srApp.StakingKeeper.PowerReduction(ctx))

	slashFactor := srApp.SlashingKeeper.SlashFractionDoubleSign(ctx)

	valBurnedTokens := slashFactor.MulInt(srApp.StakingKeeper.TokensFromConsensusPower(ctx, consPower)).TruncateInt()

	slashEventDS := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueDoubleSign),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, valBurnedTokens.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
	)

	burnedUbdel := slashFactor.MulInt(ubdelAmt).TruncateInt()
	burnedRedel := slashFactor.MulInt(redelAmt).TruncateInt()
	burnedTokens := valBurnedTokens.Add(burnedUbdel).Add(burnedRedel)

	// call slashrefund HandleRefundsFromSlash to trigger the refund process
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	srApp.SlashrefundKeeper.HandleRefundsFromSlash(ctx, slashEventDS)

	// process events to get the refund event
	valAddrEvent, refAmt := processEvents(t, ctx)
	require.Equal(t, valAddr, valAddrEvent)
	require.Equal(t, burnedTokens, refAmt)

	// check refunds
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))
	refund2 := s.RequireRefund(testAddrs[2], valAddr, slashFactor.MulInt(delAmt))
	refund3 := s.RequireRefund(testAddrs[3], valAddr, slashFactor.MulInt(ubdelAmt))
	refund4 := s.RequireRefund(testAddrs[4], valAddr, slashFactor.MulInt(redelAmt))
	s.RequireNoRefund(testAddrs[1], valAddr)

	// check refund pool
	s.RequireRefundPool(valAddr, burnedTokens, sdk.NewDecFromInt(burnedTokens), []types.Refund{refund0, refund2, refund3, refund4})

	// check deposits
	deposit0 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))
	s.RequireNoDeposit(testAddrs[1], valAddr)

	depTotal := depAmt.Add(ubdepAmt)
	drawFactor := sdk.NewDecFromInt(burnedTokens).QuoInt(depTotal)
	drawnFromUbdep := drawFactor.MulInt(ubdepAmt).TruncateInt()
	drawnFromPool := drawFactor.MulInt(depAmt).TruncateInt()
	require.Equal(t, drawnFromPool.Add(drawnFromUbdep), refAmt)

	// check unbonding deposits
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[1], valAddr)
	require.True(t, found)
	require.Equal(t, ubdepAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, ubdepAmt.Sub(drawnFromUbdep), ubd.Entries[0].Balance)
	require.Equal(t, ubdepAmt, ubd.Entries[1].InitialBalance)
	require.Equal(t, ubdepAmt, ubd.Entries[1].Balance)

	// check deposit pool
	s.RequireDepositPool(valAddr, depAmt.Sub(drawnFromPool), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit0})
}

func TestHandleRefundsFromSlash_DownTime(t *testing.T) {
	_, _, _, _, infractionHeight, slashTime := defaultTestValues()
	s := bootstrapRefundTest(t, 10)
	srApp, ctx, testAddrs, selfDelegation, valAddr, valAddrs, units := s.srApp, s.ctx, s.testAddrs, s.selfDelegation, s.valAddrs[0], s.valAddrs, s.units

	sth := teststaking.NewHelper(t, ctx, srApp.StakingKeeper)
	srh := testslashrefund.NewHelper(t, ctx, srApp.SlashrefundKeeper)

	depAmt := sdk.NewInt(150 * units)
	ubdepAmt := sdk.NewInt(50 * units)
	delAmt := sdk.NewInt(90 * units)
	ubdelAmt := sdk.NewInt(5 * units)
	redelAmt := sdk.NewInt(15 * units)

	srh.Deposit(testAddrs[0], valAddr, depAmt)
	sth.Delegate(testAddrs[2], valAddr, delAmt)

	// eligible unbonding deposit
	// entry 1: eligible
	// entry 2: not eligible
	ubdep := types.NewUnbondingDeposit(testAddrs[1], valAddr, infractionHeight+1, slashTime.Add(100), ubdepAmt)
	ubdep.Entries = append(ubdep.Entries,
		types.NewUnbondingDepositEntry(
			infractionHeight-1,
			slashTime.Add(100),
			ubdepAmt,
		))
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep)

	// eligible unbonding delegation
	// entry 1: eligible
	// entry 2: not eligible
	ubdel := stakingtypes.NewUnbondingDelegation(testAddrs[3], valAddr, infractionHeight+1, slashTime.Add(100), ubdelAmt)
	ubdel.Entries = append(ubdel.Entries,
		stakingtypes.NewUnbondingDelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			ubdelAmt,
		))
	srApp.StakingKeeper.SetUnbondingDelegation(ctx, ubdel)

	// redelegation:
	// entry 1: eligible
	// entry 2: not eligible
	red := stakingtypes.NewRedelegation(testAddrs[4], valAddr, valAddrs[1], infractionHeight+1, slashTime.Add(100), redelAmt, sdk.NewDec(1))
	red.Entries = append(red.Entries,
		stakingtypes.NewRedelegationEntry(
			infractionHeight-1,
			slashTime.Add(100),
			redelAmt, sdk.NewDec(1)),
	)
	srApp.StakingKeeper.SetRedelegation(ctx, red)

	// validator update
	ctx = sth.TurnBlock(slashTime)

	// compute burned tokens
	consAddr, err := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsAddr()
	require.NoError(t, err)

	consPower := srApp.StakingKeeper.Validator(ctx, valAddr).GetConsensusPower(srApp.StakingKeeper.PowerReduction(ctx))

	slashFactor := srApp.SlashingKeeper.SlashFractionDowntime(ctx)

	valBurnedTokens := slashFactor.MulInt(srApp.StakingKeeper.TokensFromConsensusPower(ctx, consPower)).TruncateInt()

	slashEventDT := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueMissingSignature),
		sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consAddr.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, valBurnedTokens.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
	)

	burnedUbdel := slashFactor.MulInt(ubdelAmt).TruncateInt()
	burnedRedel := slashFactor.MulInt(redelAmt).TruncateInt()
	burnedTokens := valBurnedTokens.Add(burnedUbdel).Add(burnedRedel)

	// call slashrefund HandleRefundsFromSlash to trigger the refund process
	ctx = ctx.WithBlockHeader(tmproto.Header{Height: infractionHeight + 2, Time: slashTime})
	srApp.SlashrefundKeeper.HandleRefundsFromSlash(ctx, slashEventDT)

	// process events to get the refund event
	valAddrEvent, refAmt := processEvents(t, ctx)
	require.Equal(t, valAddr, valAddrEvent)
	require.Equal(t, burnedTokens, refAmt)

	// check refunds
	refund0 := s.RequireRefund(testAddrs[0], valAddr, slashFactor.MulInt(selfDelegation))
	refund2 := s.RequireRefund(testAddrs[2], valAddr, slashFactor.MulInt(delAmt))
	refund3 := s.RequireRefund(testAddrs[3], valAddr, slashFactor.MulInt(ubdelAmt))
	refund4 := s.RequireRefund(testAddrs[4], valAddr, slashFactor.MulInt(redelAmt))
	s.RequireNoRefund(testAddrs[1], valAddr)

	// check refund pool
	s.RequireRefundPool(valAddr, burnedTokens, sdk.NewDecFromInt(burnedTokens), []types.Refund{refund0, refund2, refund3, refund4})

	// check deposits
	deposit0 := s.RequireDeposit(testAddrs[0], valAddr, sdk.NewDecFromInt(depAmt))
	s.RequireNoDeposit(testAddrs[1], valAddr)

	depTotal := depAmt.Add(ubdepAmt)
	drawFactor := sdk.NewDecFromInt(burnedTokens).QuoInt(depTotal)
	drawnFromUbdep := drawFactor.MulInt(ubdepAmt).TruncateInt()
	drawnFromPool := drawFactor.MulInt(depAmt).TruncateInt()
	require.Equal(t, drawnFromPool.Add(drawnFromUbdep), refAmt)

	// check unbonding deposits
	ubd, found := srApp.SlashrefundKeeper.GetUnbondingDeposit(ctx, testAddrs[1], valAddr)
	require.True(t, found)
	require.Equal(t, ubdepAmt, ubd.Entries[0].InitialBalance)
	require.Equal(t, ubdepAmt.Sub(drawnFromUbdep), ubd.Entries[0].Balance)
	require.Equal(t, ubdepAmt, ubd.Entries[1].InitialBalance)
	require.Equal(t, ubdepAmt, ubd.Entries[1].Balance)

	// check deposit pool
	s.RequireDepositPool(valAddr, depAmt.Sub(drawnFromPool), sdk.NewDecFromInt(depAmt), []types.Deposit{deposit0})
}

func createNRefund(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.Refund {

	items := make([]types.Refund, n)
	for i := range items {
		delPubk := secp256k1.GenPrivKey().PubKey()
		delAddr := sdk.AccAddress(delPubk.Address())
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].DelegatorAddress = delAddr.String()
		items[i].ValidatorAddress = valAddr.String()
		items[i].Shares = sdk.NewDec(int64(1000 * i))
		keeper.SetRefund(ctx, items[i])
	}
	return items
}

func TestRefundGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	for _, item := range items {
		delAddr, _ := sdk.AccAddressFromBech32(item.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		rst, found := keeper.GetRefund(ctx, delAddr, valAddr)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}
func TestRefundRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	for _, item := range items {
		delAddr, _ := sdk.AccAddressFromBech32(item.DelegatorAddress)
		valAddr, _ := sdk.ValAddressFromBech32(item.ValidatorAddress)
		refund, found := keeper.GetRefund(ctx, delAddr, valAddr)
		keeper.RemoveRefund(ctx, refund)
		_, found = keeper.GetRefund(ctx, delAddr, valAddr)
		require.False(t, found)
	}
}

func TestRefundGetAll(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefund(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllRefund(ctx))
}

// -------------------------------------------------------------------------------------------------
// Test refund pool
// -------------------------------------------------------------------------------------------------


func createNRefundPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.RefundPool {
	items := make([]types.RefundPool, n)
	for i := range items {
		valPubk := secp256k1.GenPrivKey().PubKey()
		valAddr := sdk.ValAddress(valPubk.Address())
		items[i].OperatorAddress = valAddr.String()
		items[i].Shares = sdk.NewDec(int64(1000 * i))
		items[i].Tokens = sdk.NewInt64Coin("stake", int64(1000*i))
		keeper.SetRefundPool(ctx, items[i])
	}
	return items
}

func TestRefundPoolGet(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		rst, found := keeper.GetRefundPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t, item, rst)
	}
}

func TestUpdateRefundPool(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	refPools := createNRefundPool(keeper, ctx, 10)
	for i, refPool := range refPools {
		valAddr, _ := sdk.ValAddressFromBech32(refPool.OperatorAddress)

		refPool.Tokens.Amount = refPool.Tokens.Amount.Add(sdk.NewInt(int64(i * 1000)))
		refPool.Shares = refPool.Shares.Add(sdk.NewDec(int64(i * 2000)))
		keeper.SetRefundPool(ctx, refPool)

		rst, found := keeper.GetRefundPool(ctx, valAddr)
		require.True(t, found)
		require.Equal(t, refPool, rst)
	}
}

func TestRefundPoolRemove(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	for _, item := range items {
		valAddr, _ := sdk.ValAddressFromBech32(item.OperatorAddress)
		keeper.RemoveRefundPool(ctx, valAddr)
		_, found := keeper.GetRefundPool(ctx, valAddr)
		require.False(t, found)
	}
}

func TestRefundPoolGetAll(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	items := createNRefundPool(keeper, ctx, 10)
	require.ElementsMatch(t, items, keeper.GetAllRefundPool(ctx))
}
*/
