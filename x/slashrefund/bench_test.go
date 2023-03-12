package slashrefund_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	"github.com/stretchr/testify/require"
)

func BenchmarkHandleRefundFromSlash(b *testing.B) {
	benchmarkHandleRefundFromSlash(b, 1, 1, 1)
}

type BenchmarkTestSetup struct {
	app      *app.App
	ctx      sdk.Context
	delAddrs []sdk.AccAddress
	valAddrs []sdk.ValAddress
	depAddrs []sdk.AccAddress

	deposits          []types.Deposit
	unbondingDeposits []types.UnbondingDeposit

	delegations          []stakingtypes.Delegation
	UnbondingDelegations []stakingtypes.UnbondingDelegation
	redelegations        []stakingtypes.Redelegation

	infractionHeight int64
	slashingTime     time.Time

	slashEvent sdk.Event
}

func benchmarkHandleRefundFromSlash(b *testing.B, numDelAddrs, numRedelAddrs, numUbdelAddrs int) {
	//b.ReportAllocs()
	for i := 0; i < b.N; i++ {
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

		valAddrs := s.valAddrs
		depAddrs := s.depAddrs
		app := s.app
		ctx := s.ctx
		valAddr := valAddrs[0]
		validator, found := app.StakingKeeper.GetValidator(ctx, valAddr)
		require.True(b, found)
		consAddr, err := validator.GetConsAddr()
		require.NoError(b, err)
		consPower := validator.ConsensusPower(sdk.DefaultPowerReduction)
		slashFactor := app.SlashingKeeper.SlashFractionDowntime(ctx)
		valBurnedTokens := slashFactor.MulInt(app.StakingKeeper.TokensFromConsensusPower(ctx, consPower)).TruncateInt()
		slashEventDT := sdk.NewEvent(
			slashingtypes.EventTypeSlash,
			sdk.NewAttribute(slashingtypes.AttributeKeyAddress, consAddr.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyPower, fmt.Sprintf("%d", consPower)),
			sdk.NewAttribute(slashingtypes.AttributeKeyReason, slashingtypes.AttributeValueMissingSignature),
			sdk.NewAttribute(slashingtypes.AttributeKeyJailed, consAddr.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyBurnedCoins, valBurnedTokens.String()),
			sdk.NewAttribute(slashingtypes.AttributeKeyInfractionHeight, fmt.Sprintf("%d", infractionHeight)),
		)

		// Generate deposit.
		depCoin := sdk.NewCoin(app.SlashrefundKeeper.AllowedTokens(ctx)[0], valBurnedTokens)
		res, err := app.SlashrefundKeeper.Deposit(ctx, depAddrs[0], depCoin, validator)
		require.NoError(b, err)
		require.NotNil(b, res)

		// Set block height.
		ctx = ctx.WithBlockHeight(int64(slashingHeight))
		ctx = ctx.WithBlockTime(slashingTime)
		if i == 0 {
			b.ResetTimer()
		}

		// Run HandleRefundsFromSlash.
		b.StartTimer()
		app.SlashrefundKeeper.HandleRefundsFromSlash(ctx, slashEventDT)
		b.StopTimer()

		// Check HandleRefundsFromSlash execution.
		valAddrE, refAmount := processEvents(b, ctx)
		require.Equal(b, valAddr, valAddrE)
		require.Equal(b, valBurnedTokens, refAmount)
	}

}

func makeRandomAddressesAndPublicKeys(n int) (accL []sdk.ValAddress, pkL []*ed25519.PubKey) {
	for i := 0; i < n; i++ {
		pk := ed25519.GenPrivKey().PubKey().(*ed25519.PubKey)
		pkL = append(pkL, pk)
		accL = append(accL, sdk.ValAddress(pk.Address()))
	}
	return accL, pkL
}

func SetupBenchmarkTest(numDelAddrs, numRedelAddrs, numUbdelAddrs int, creationHeight int64, slashingTime time.Time) (
	s BenchmarkTestSetup,
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

	newHeightAndTime := func(creationHeight int64, completionTime time.Time) testsuite.HeigthAndTime {
		return testsuite.HeigthAndTime{
			CreationHeight: creationHeight,
			CompletionTime: completionTime,
		}
	}

	// Same height, different completion time in order to make entries unique.
	heightAndTimes := []testsuite.HeigthAndTime{
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(1))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(2))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(3))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(4))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(5))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(6))),
		newHeightAndTime(creationHeight, slashingTime.Add(time.Hour*time.Duration(7))),
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
	valAddr := valAddrs[0]
	validator, found := app.StakingKeeper.GetValidator(ctx, valAddr)
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

	s.app = app
	s.ctx = ctx
	s.delAddrs = delAddrs
	s.valAddrs = valAddrs
	s.depAddrs = depAddrs
	s.slashEvent = slashEventDT

	return s, nil
}

func processEvents(tb testing.TB, ctx sdk.Context) (valAddr sdk.ValAddress, amount sdk.Int) {

	var slashEvents []sdk.Event

	events := ctx.EventManager().Events()

	for _, event := range events {
		if event.Type == types.EventTypeRefund {
			slashEvents = append(slashEvents, event)
		}
	}
	require.Equal(tb, 1, len(slashEvents), fmt.Sprintf("expected one refund event, got: %d", len(slashEvents)))

	valAddr = sdk.ValAddress{}
	amount = sdk.ZeroInt()

	for _, attr := range slashEvents[0].Attributes {
		switch string(attr.GetKey()) {
		case "validator":
			valAddrE, err := sdk.ValAddressFromBech32(string(attr.GetValue()))
			require.NoError(tb, err)
			valAddr = valAddrE
		case "amount":
			amountE, ok := sdk.NewIntFromString(string(attr.GetValue()))
			require.True(tb, ok)
			amount = amountE
		}
	}
	return valAddr, amount
}
