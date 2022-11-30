package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestRefundQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRefund(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetRefundRequest
		response *types.QueryGetRefundResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRefundRequest{
				Delegator: msgs[0].Delegator,
				Validator: msgs[0].Validator,
			},
			response: &types.QueryGetRefundResponse{Refund: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetRefundRequest{
				Delegator: msgs[1].Delegator,
				Validator: msgs[1].Validator,
			},
			response: &types.QueryGetRefundResponse{Refund: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRefundRequest{
				Delegator: strconv.Itoa(100000),
				Validator: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Refund(wctx, tc.request)
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

func TestRefundQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNRefund(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllRefundRequest {
		return &types.QueryAllRefundRequest{
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
			resp, err := keeper.RefundAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Refund), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Refund),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.RefundAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Refund), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Refund),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.RefundAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Refund),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.RefundAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
