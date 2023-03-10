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

// CreateTestApp returns the context and the app just for testing purpose.
func CreateTestApp(
	delAddrs []sdk.AccAddress,
	valAddrs []sdk.ValAddress,
	balances []banktypes.Balance,
	isCheckTx bool,
) (*app.App, sdk.Context) {

	srApp := MakeTestApp(delAddrs, valAddrs, balances)
	ctx := srApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	return srApp, ctx
}

// MakeTestApp returns a new test app with a genesis state.
func MakeTestApp(
	delAddrs []sdk.AccAddress,
	valAddrs []sdk.ValAddress,
	balances []banktypes.Balance,
) *app.App {
	encodingConfig := app.MakeEncodingConfig()

	var invCheckPeriod uint = 0
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

	// Create validators from validator addresses and consensus keys.
	pks := GenerateNConsensusPubKeys(len(valAddrs))
	var validators []stakingtypes.Validator
	var valDelegations []stakingtypes.Delegation
	for i, valAddr := range(valAddrs) {
		val, valDel := GenerateValidator(valAddr, pks[i])
		valDelegations = append(valDelegations, valDel)
		validators = append(validators, val)
	}

	// Create delegations from delegators
	delDelegations, validators := GenerateRandomDelegations(delAddrs, validators)

	// Create custom genesis or empty
	genesisState := CustomGenesisState(testApp, balances, validators, delDelegations, valDelegations)
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
) (genesisState app.GenesisState) {

	// Generate default genesis
	genesisState = NewDefaultGenesisState(srApp.AppCodec())

	// Set validators and delegations.
	stakingGenesis := stakingtypes.NewGenesisState(
		stakingtypes.DefaultParams(), 
		validators, 
		append(delDelegations, valDelegations...),
	)
	genesisState[stakingtypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// Add genesis accounts tokens to total supply.
		totalSupply = totalSupply.Add(b.Coins...)
	}

	bondedSupply := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(0))
	for _, valDel := range valDelegations {
		// At this stage TruncateInt is safe since shares are always integer.
		bondedSupply = bondedSupply.AddAmount(valDel.Shares.TruncateInt())
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, valDel.Shares.TruncateInt()))
	}

	/*
	for _, del := range delegations {
		// Add delegated tokens to total supply and bonded supply. 
	}
	*/

	// Add bonded amount to bonded pool module account.
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondedSupply.Amount)},
	})

	// Update total supply.
	bankGenesis := banktypes.NewGenesisState(
		banktypes.DefaultGenesisState().Params, 
		balances, 
		totalSupply, 
		[]banktypes.Metadata{},
	)
	genesisState[banktypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}