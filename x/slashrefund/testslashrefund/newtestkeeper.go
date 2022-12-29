package testslashrefund

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v5/modules/apps/transfer/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func NewTestKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	memStoreKey_auth := storetypes.NewMemoryStoreKey("mem_acc")
	memStoreKey_bank := storetypes.NewMemoryStoreKey("mem_bank")
	memStoreKey_slashing := storetypes.NewMemoryStoreKey("mem_slashing")
	memStoreKey_staking := storetypes.NewMemoryStoreKey("mem_staking")

	storeKey_auth := sdk.NewKVStoreKey(authtypes.StoreKey)
	storeKey_bank := sdk.NewKVStoreKey(banktypes.StoreKey)
	storeKey_staking := sdk.NewKVStoreKey(stakingtypes.StoreKey)
	storeKey_slashing := sdk.NewKVStoreKey(slashingtypes.StoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"SlashrefundParams",
	)
	paramsSubspace_auth := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey_auth,
		memStoreKey_auth,
		"AuthParams",
	)
	paramsSubspace_bank := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey_bank,
		memStoreKey_bank,
		"BankParams",
	)
	paramsSubspace_slashing := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey_slashing,
		memStoreKey_slashing,
		"SlashingParams",
	)
	paramsSubspace_staking := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey_staking,
		memStoreKey_staking,
		"StakingParams",
	)

	keys := sdk.NewKVStoreKeys(banktypes.StoreKey, stakingtypes.StoreKey, slashingtypes.StoreKey, authtypes.StoreKey)

	// module account permissions
	maccPerms := map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		icatypes.ModuleName:            nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		types.ModuleName:               nil,
		// this line is used by starport scaffolding # stargate/app/maccPerms
	}

	authKeeper := authkeeper.NewAccountKeeper(
		cdc,
		keys[authtypes.StoreKey],
		paramsSubspace_auth,
		authtypes.ProtoBaseAccount,
		maccPerms,
		sdk.Bech32PrefixAccAddr,
	)

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		keys[banktypes.StoreKey],
		authKeeper,
		paramsSubspace_bank,
		BlockedModuleAccountAddrs(maccPerms),
	)

	stakingKeeper := stakingkeeper.NewKeeper(
		cdc,
		keys[stakingtypes.StoreKey],
		authKeeper,
		bankKeeper,
		paramsSubspace_staking,
	)

	slashingKeeper := slashingkeeper.NewKeeper(
		cdc,
		keys[slashingtypes.StoreKey],
		stakingKeeper,
		paramsSubspace_slashing,
	)

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bankKeeper,
		stakingKeeper,
		slashingKeeper,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func BlockedModuleAccountAddrs(maccPerms map[string][]string) map[string]bool {
	blockedModAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedModAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	delete(blockedModAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return blockedModAccAddrs
}
