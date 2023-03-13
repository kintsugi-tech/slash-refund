package slashrefund_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	"github.com/stretchr/testify/require"
)

func BenchmarkHandleRefundFromSlash100(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 100, 100, 100)
}

func BenchmarkHandleRefundFromSlash200(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 200, 200, 200)
}

func BenchmarkHandleRefundFromSlash500(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 500, 500, 500)
}

func BenchmarkHandleRefundFromSlash1000(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 1000, 1000, 1000)
}

func BenchmarkHandleRefundFromSlash2000(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 2000, 2000, 2000)
}

func BenchmarkHandleRefundFromSlash4000(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 4000, 4000, 4000)
}

type BenchmarkTestInputs struct {
	app        *app.App
	ctx        sdk.Context
	delAddrs   []sdk.AccAddress
	valAddrs   []sdk.ValAddress
	redelAddrs []sdk.AccAddress
	ubdelAddrs []sdk.AccAddress
	depAddrs   []sdk.AccAddress

	deposits          []types.Deposit
	unbondingDeposits []types.UnbondingDeposit
	depositPool       types.DepositPool

	delegations          []stakingtypes.Delegation
	unbondingDelegations []stakingtypes.UnbondingDelegation
	redelegations        []stakingtypes.Redelegation

	infractionHeight int64
	slashingTime     time.Time

	slashEvent sdk.Event
}

func benchmarkHandleRefundFromSlash(b *testing.B, numDelAddrs, numRedelAddrs, numUbdelAddrs int) {
	// Generate test inputs.
	// Account for all redelegations' and unbonding delegations' entries:
	//
	//   creationHeight > infractionHeight && slashTime < completionTime
	//
	// creationHeight is set to 15 in redelegations' and unbonding delegations'
	// entries.
	// completionTime is set to time.Now() plus 1 to 7 hours.
	infractionHeight := 10
	creationHeight := infractionHeight + 2
	slashingTime := time.Now()
	slashingHeight := creationHeight + 2
	s, err := SetupBenchmarkTest(numDelAddrs, numRedelAddrs, numUbdelAddrs, int64(creationHeight), slashingTime)
	require.NoError(b, err)
	k := s.app.SlashrefundKeeper
	ctx := s.ctx.WithBlockHeight(int64(slashingHeight)).WithBlockTime(slashingTime)

	for i := 0; i < b.N; i++ {

		k.SetDeposit(ctx, s.deposits[0])
		k.SetDepositPool(ctx, s.depositPool)

		// Run HandleRefundsFromSlash.
		if i == 0 {
			b.ResetTimer()
		}
		b.StartTimer()
		s.app.SlashrefundKeeper.HandleRefundsFromSlash(ctx, s.slashEvent)
		b.StopTimer()

		// Remove refunds and refunds pool
		refunds := s.app.SlashrefundKeeper.GetValidatorRefunds(ctx, s.valAddrs[0])
		require.Greater(b, len(refunds), 0)
		for _, ref := range refunds {
			k.RemoveRefund(ctx, ref)
		}
		k.RemoveRefundPool(ctx, s.valAddrs[0])

		// Remove deposits and deposit pool
		//numDepAddrs := 1
		deposits := k.GetValidatorDeposits(ctx, s.valAddrs[0])
		for _, dep := range deposits {
			k.RemoveDeposit(ctx, dep)
		}
		k.RemoveDepositPool(ctx, s.valAddrs[0])
	}
}

func SetupBenchmarkTest(numDelAddrs, numRedelAddrs, numUbdelAddrs int, creationHeight int64, slashingTime time.Time) (
	s BenchmarkTestInputs,
	err error,
) {

	numValAddrs := 2
	numDepAddrs := 1

	// Setup delegators
	delAddrsHex := testsuite.GenerateNAddresses(numDelAddrs)
	delAddrs := testsuite.ConvertAddressesToAccAddr(delAddrsHex)
	balances := testsuite.GenerateBalances(delAddrs)

	// Setup validators
	valAddrsHex := testsuite.GenerateNAddresses(numValAddrs)
	valAddrs := testsuite.ConvertAddressesToValAddr(valAddrsHex)

	// Setup depositors
	depAddrsHex := testsuite.GenerateNAddresses(numDepAddrs)
	depAddrs := testsuite.ConvertAddressesToAccAddr(depAddrsHex)
	depBalances := testsuite.GenerateBalances(depAddrs)

	// Setup redelegators
	redAddrsHex := testsuite.GenerateNAddresses(numRedelAddrs)
	redAddrs := testsuite.ConvertAddressesToAccAddr(redAddrsHex)
	redBalances := testsuite.GenerateBalances(redAddrs)

	// Setup unbonding delegators
	ubdAddrsHex := testsuite.GenerateNAddresses(numUbdelAddrs)
	ubdAddrs := testsuite.ConvertAddressesToAccAddr(ubdAddrsHex)
	ubdBalances := testsuite.GenerateBalances(ubdAddrs)

	balances = append(append(append(balances, depBalances...), redBalances...), ubdBalances...)

	// Same height, different completion time in order to make entries unique.
	heightAndTimes := []testsuite.HeigthAndTime{
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(1))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(2))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(3))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(4))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(5))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(6))),
		testsuite.NewHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(7))),
	}

	testInputs := testsuite.TestInputs{}
	testInputs.DelAddrs = delAddrs
	testInputs.ValAddrs = valAddrs
	testInputs.Balances = balances
	testInputs.RedelAddrs = redAddrs
	testInputs.RedelValidatorsIndex = []int64{0, 1}
	testInputs.RedelHeightAndTimes = heightAndTimes
	testInputs.UbdelAddrs = ubdAddrs
	testInputs.UbdelValidatorsIndex = []int64{0}
	testInputs.UbdelHeightAndTimes = heightAndTimes

	app, ctx := testsuite.CreateTestApp(testInputs, false)

	// Get validator and validator consensus address.
	validator, found := app.StakingKeeper.GetValidator(ctx, valAddrs[0])
	if !found {
		return s, stakingtypes.ErrNoValidatorFound
	}
	consAddr, err := validator.GetConsAddr()
	if err != nil {
		return s, err
	}

	// Generate slash event.
	consPower := validator.ConsensusPower(sdk.DefaultPowerReduction)
	slashFactor := app.SlashingKeeper.SlashFractionDowntime(ctx)
	valBurnedTokens := slashFactor.MulInt(app.StakingKeeper.TokensFromConsensusPower(ctx, consPower)).TruncateInt()
	infractionHeight := 10
	slashEventDT := sdk.NewEvent(
		slashingtypes.EventTypeSlash,
		sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
		sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueMissingSignature),
		sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consAddr.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, valBurnedTokens.String()),
		sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
	)

	depAmt := valBurnedTokens.MulRaw(1)
	depCoin := sdk.NewCoin(app.SlashrefundKeeper.AllowedTokens(ctx)[0], depAmt)

	s.app = app
	s.ctx = ctx
	s.delAddrs = delAddrs
	s.valAddrs = valAddrs
	s.redelAddrs = redAddrs
	s.ubdelAddrs = ubdAddrs
	s.depAddrs = depAddrs
	s.slashEvent = slashEventDT
	s.deposits = []types.Deposit{
		types.NewDeposit(depAddrs[0], valAddrs[0], sdk.NewDecFromInt(depAmt)),
	}
	s.depositPool = types.NewDepositPool(valAddrs[0], depCoin, sdk.NewDecFromInt(depAmt))

	return s, nil
}
