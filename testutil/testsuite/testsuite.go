package testsuite

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// returns context and app
func CreateTestApp(isCheckTx bool) (*app.App, sdk.Context) {

	srApp := setup(isCheckTx)
	ctx := srApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	return srApp, ctx
}

func setup(isCheckTx bool) *app.App {
	srApp, genesisState := genApp(!isCheckTx, 0)
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		srApp.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: simapp.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return srApp
}

// EmptyAppOptions is a stub implementing AppOptions
type emptyAppOptions struct{}

// Get implements AppOptions
func (ao emptyAppOptions) Get(o string) interface{} {
	return nil
}

func genApp(withGenesis bool, invCheckPeriod uint) (*app.App, app.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	cosmoscmdApp := app.New(
		log.NewNopLogger(),
		db,
		nil,
		true, // load latest version
		map[int64]bool{},
		app.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		emptyAppOptions{},
	)

	srApp := cosmoscmdApp.(*app.App)

	if withGenesis {
		return srApp, CustomGenesisState(srApp)
	}
	return srApp, app.GenesisState{}
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) app.GenesisState {
	return app.ModuleBasics.DefaultGenesis(cdc)
}

// NewCustomGenesisState generate a genesis state with an account and a single validator
func CustomGenesisState(srApp *app.App) (genesisState app.GenesisState) {

	var balances []banktypes.Balance
	var validators []stakingtypes.Validator
	var delegations []stakingtypes.Delegation
	// one unit in micro-units
	var units int64 = 1000000

	// generate genesis account
	privk := secp256k1.GenPrivKey()
	pk := privk.PubKey()
	acc := authtypes.NewBaseAccount(pk.Address().Bytes(), pk, 0, 0)
	addr := acc.GetAddress()
	balance := banktypes.Balance{
		Address: addr.String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1e9*units))),
	}

	// create validator set with single validator
	bondAmt := sdk.NewInt(10 * units)
	bondAmtDec := sdk.NewDecFromInt(bondAmt)
	pkAny, _ := codectypes.NewAnyWithValue(pk)
	valAddr := sdk.ValAddress(addr)
	validator := stakingtypes.Validator{
		OperatorAddress:   valAddr.String(),
		ConsensusPubkey:   pkAny,
		Jailed:            false,
		Status:            stakingtypes.Bonded,
		Tokens:            bondAmt,
		DelegatorShares:   bondAmtDec,
		Description:       stakingtypes.Description{},
		UnbondingHeight:   int64(0),
		UnbondingTime:     time.Unix(0, 0).UTC(),
		Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
		MinSelfDelegation: sdk.OneInt(),
	}

	delegation := stakingtypes.NewDelegation(addr, valAddr, bondAmtDec)

	balances = append(balances, balance)
	validators = append(validators, validator)
	delegations = append(delegations, delegation)

	// generate default genesis
	genesisState = NewDefaultGenesisState(srApp.AppCodec())

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	for range delegations {
		// add delegated tokens to total supply
		totalSupply = totalSupply.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))
	}

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt)},
	})

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = srApp.AppCodec().MustMarshalJSON(bankGenesis)

	return genesisState
}

type Helper struct {
	t       *testing.T
	k       keeper.Keeper
	msgSrvr types.MsgServer
	ctx     sdk.Context
}

func NewHelper(t *testing.T, k keeper.Keeper, ctx sdk.Context) *Helper {
	return &Helper{t, k, keeper.NewMsgServerImpl(k), ctx}
}

func (srh *Helper) Deposit(depAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Int) {
	types.DefaultGenesis().Params.GetAllowedTokens()
	coin := sdk.NewCoin(srh.k.AllowedTokensList(srh.ctx)[0], amount)
	msg := types.NewMsgDeposit(depAddr.String(), valAddr.String(), coin)
	res, err := srh.msgSrvr.Deposit(sdk.WrapSDKContext(srh.ctx), msg)
	require.NoError(srh.t, err)
	require.NotNil(srh.t, res)
}
