package keeper_test

import (
	"testing"

	testkeeper "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.SlashrefundKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.AllowedTokens, k.AllowedTokens(ctx))
}
