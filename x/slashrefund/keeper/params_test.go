package keeper_test

import (
	"testing"

	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testslashrefund.NewTestKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.AllowedTokens, k.AllowedTokens(ctx))
}
