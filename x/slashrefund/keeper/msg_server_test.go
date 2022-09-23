package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.SlashrefundKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
