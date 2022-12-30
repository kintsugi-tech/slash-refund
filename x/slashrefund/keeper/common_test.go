package keeper_test

import (
	"github.com/made-in-block/slash-refund/app"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
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
