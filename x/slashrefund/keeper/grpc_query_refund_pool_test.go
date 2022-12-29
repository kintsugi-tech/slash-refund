package keeper_test

import (
	"strconv"
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

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRefundPoolQuerySingle(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRefundPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetRefundPoolRequest
		response *types.QueryGetRefundPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: msgs[0].OperatorAddress,
			},
			response: &types.QueryGetRefundPoolResponse{RefundPool: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: msgs[1].OperatorAddress,
			},
			response: &types.QueryGetRefundPoolResponse{RefundPool: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.RefundPool(wctx, tc.request)
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

func TestRefundPoolQueryPaginated(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRefundPool(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRefundPoolRequest {
		return &types.QueryAllRefundPoolRequest{
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
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RefundPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RefundPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RefundPool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RefundPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RefundPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.RefundPool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RefundPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.RefundPool),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RefundPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
