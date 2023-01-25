package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestDepositPoolQuerySingle(t *testing.T) {
	k, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	depPools := createNDepositPool(k, ctx, 3)

	// remove last deposit pool to test "KeyNotFound"
	valAddr, err := sdk.ValAddressFromBech32(depPools[2].OperatorAddress)
	require.NoError(t, err)
	k.RemoveDepositPool(ctx, valAddr)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDepositPoolRequest
		response *types.QueryGetDepositPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: depPools[0].OperatorAddress,
			},
			response: &types.QueryGetDepositPoolResponse{DepositPool: depPools[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: depPools[1].OperatorAddress,
			},
			response: &types.QueryGetDepositPoolResponse{DepositPool: depPools[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: depPools[2].OperatorAddress,
			},
			err: status.Errorf(codes.NotFound, "deposit pool not found for operator %s", depPools[2].OperatorAddress),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := k.DepositPool(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestDepositPoolQueryPaginated(t *testing.T) {
	k, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	depPools := createNDepositPool(k, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDepositPoolRequest {
		return &types.QueryAllDepositPoolRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(depPools); i += step {
			resp, err := k.DepositPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DepositPool), step)
			require.Subset(t,
				nullify.Fill(depPools),
				nullify.Fill(resp.DepositPool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(depPools); i += step {
			resp, err := k.DepositPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DepositPool), step)
			require.Subset(t,
				nullify.Fill(depPools),
				nullify.Fill(resp.DepositPool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := k.DepositPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(depPools), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(depPools),
			nullify.Fill(resp.DepositPool),
		)
	})
	t.Run("Invalid Request", func(t *testing.T) {
		_, err := k.DepositPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
