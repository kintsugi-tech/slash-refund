package testsuite

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/made-in-block/slash-refund/app"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// EmptyAppOptions is a stub implementing AppOptions
type emptyAppOptions struct{}

// Get implements AppOptions
func (ao emptyAppOptions) Get(o string) interface{} {
	return nil
}

type TestInputs struct {
	DelAddrs []sdk.AccAddress
	ValAddrs []sdk.ValAddress
	Balances []banktypes.Balance

	RedelAddrs           []sdk.AccAddress
	RedelValidatorsIndex []int64
	RedelHeightAndTimes  []HeigthAndTime

	UbdelAddrs           []sdk.AccAddress
	UbdelValidatorsIndex []int64
	UbdelHeightAndTimes  []HeigthAndTime
}

// CreateTestApp returns the context and the app just for testing purpose.
func CreateTestApp(
	testInputs TestInputs,
	isCheckTx bool,
) (*app.App, sdk.Context) {

	srApp := MakeTestApp(testInputs)
	ctx := srApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	return srApp, ctx
}

// MakeTestApp returns a new test app with a genesis state.
func MakeTestApp(
	testInputs TestInputs,
) *app.App {
	encodingConfig := app.MakeEncodingConfig()

	var invCheckPeriod uint
	testApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true, // load latest version
		map[int64]bool{},
		app.DefaultNodeHome,
		invCheckPeriod,
		encodingConfig,
		emptyAppOptions{},
	)

	delAddrs := testInputs.DelAddrs
	valAddrs := testInputs.ValAddrs
	balances := testInputs.Balances

	// Create validators from validator addresses and consensus keys.
	pks := GenerateNConsensusPubKeys(len(valAddrs))
	var validators []stakingtypes.Validator
	var valDelegations []stakingtypes.Delegation
	for i, valAddr := range valAddrs {
		val, valDel := GenerateValidator(valAddr, pks[i])
		valDelegations = append(valDelegations, valDel)
		validators = append(validators, val)
	}

	extractValidators := func(validators []stakingtypes.Validator, indexes []int64) []stakingtypes.Validator {
		var extracted []stakingtypes.Validator
		for _, idx := range indexes {
			if int(idx) < len(validators) {
				extracted = append(extracted, validators[idx])
			}
		}
		return extracted
	}

	updateValidators := func(original []stakingtypes.Validator, updates []stakingtypes.Validator, indexes []int64) []stakingtypes.Validator {
		for i, idx := range indexes {
			if int(idx) < len(validators) {
				original[idx] = updates[i]
			}
		}
		return original
	}

	// Create delegations from delegators.
	delDelegations, validators := GenerateRandomDelegations(delAddrs, validators)

	// Create redelegations from redelegators.
	var redelegations []stakingtypes.Redelegation
	var redDelegations []stakingtypes.Delegation
	if len(testInputs.RedelAddrs) > 0 &&
		len(testInputs.RedelValidatorsIndex) > 0 &&
		len(testInputs.RedelHeightAndTimes) > 0 {

		redValidators := extractValidators(validators, testInputs.RedelValidatorsIndex)
		if len(redValidators) > 1 {
			redelegations, redDelegations, redValidators = GenerateRandomRedelegations(
				testInputs.RedelAddrs,
				redValidators,
				testInputs.RedelHeightAndTimes,
			)
			validators = updateValidators(validators, redValidators, testInputs.RedelValidatorsIndex)
		}
	}

	// Create unbonding delegations from unbonding delegators.
	var ubdelegations []stakingtypes.UnbondingDelegation
	if len(testInputs.UbdelAddrs) > 0 &&
		len(testInputs.UbdelValidatorsIndex) > 0 &&
		len(testInputs.UbdelHeightAndTimes) > 0 {

		ubdelValidators := extractValidators(validators, testInputs.UbdelValidatorsIndex)
		if len(ubdelValidators) > 0 {
			ubdelegations, ubdelValidators = GenerateRandomUnbondingDelegations(
				testInputs.UbdelAddrs,
				ubdelValidators,
				testInputs.UbdelHeightAndTimes,
			)
			validators = updateValidators(validators, ubdelValidators, testInputs.UbdelValidatorsIndex)
		}
	}

	// Create custom genesis or empty.
	genesisState := CustomGenesisState(testApp, balances, validators, delDelegations, valDelegations, redelegations, redDelegations, ubdelegations)
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain.
	testApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	return testApp
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) app.GenesisState {
	return app.ModuleBasics.DefaultGenesis(cdc)
}

// CustomGenesisState generate a genesis state with an account and a single validator
func CustomGenesisState(
	srApp *app.App,
	balances []banktypes.Balance,
	validators []stakingtypes.Validator,
	delDelegations []stakingtypes.Delegation,
	valDelegations []stakingtypes.Delegation,
	redelegations []stakingtypes.Redelegation,
	redDelegations []stakingtypes.Delegation,
	ubdelegations []stakingtypes.UnbondingDelegation,
) (genesisState app.GenesisState) {

	// Generate default genesis
	genesisState = NewDefaultGenesisState(srApp.AppCodec())

	// ==== STAKING GENESIS STATE ====
	// Set validators and delegations.
	delegations := append(valDelegations, append(delDelegations, redDelegations...)...)
	stakingGenesis := stakingtypes.NewGenesisState(
		stakingtypes.DefaultParams(),
		validators,
		delegations,
	)
	stakingGenesis.Redelegations = redelegations
	stakingGenesis.UnbondingDelegations = ubdelegations
	genesisState[stakingtypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(stakingGenesis)

	// ==== BANK GENESIS STATE ====
	totalSupply := sdk.NewCoins()
	balancesSupply := sdk.NewCoins()
	bondedSupply := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))
	notbondedSupply := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))

	// Add genesis accounts tokens to total supply.
	for _, b := range balances {
		balancesSupply = balancesSupply.Add(b.Coins...)
	}

	// Add delegated tokens to total supply and bonded supply.
	// These contain also re-delegated tokens.
	for _, del := range delegations {
		// At this stage TruncateInt is safe since shares are always integer.
		// Since no slash occurred, shares and tokens are in 1:1 ratio.
		bondedSupply = bondedSupply.AddAmount(del.Shares.TruncateInt())
	}

	// Add unbonding tokens to total supply and not-bonded supply.
	for _, ubdel := range ubdelegations {
		for _, e := range ubdel.Entries {
			notbondedSupply = notbondedSupply.AddAmount(e.InitialBalance)
		}
	}

	// Add bonded amount to bonded pool module account.
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondedSupply.Amount)},
	})

	// Add unbonded amount to not-bonded pool module account.
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.NotBondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, notbondedSupply.Amount)},
	})

	// Update total supply.
	totalSupply = totalSupply.Add(balancesSupply...).Add(bondedSupply).Add(notbondedSupply)
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params,
		balances,
		totalSupply,
		[]banktypes.Metadata{},
	)
	genesisState[banktypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}
