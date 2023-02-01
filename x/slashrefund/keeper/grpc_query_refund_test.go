package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestRefundQuerySingle(t *testing.T) {
	s := SetupTestSuite(t, 100)
	srApp, ctx, testAddrs, valAddrs, querier := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.querier

	wctx := sdk.WrapSDKContext(ctx)

	ref1 := types.NewRefund(testAddrs[0], valAddrs[0], sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetRefund(ctx, ref1)

	ref2 := types.NewRefund(testAddrs[1], valAddrs[0], sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetRefund(ctx, ref2)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetRefundRequest
		response *types.QueryGetRefundResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRefundRequest{
				Delegator: testAddrs[0].String(),
				Validator: valAddrs[0].String(),
			},
			response: &types.QueryGetRefundResponse{Refund: ref1},
		},
		{
			desc: "Second",
			request: &types.QueryGetRefundRequest{
				Delegator: testAddrs[1].String(),
				Validator: valAddrs[0].String(),
			},
			response: &types.QueryGetRefundResponse{Refund: ref2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRefundRequest{
				Delegator: testAddrs[1].String(),
				Validator: valAddrs[1].String(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := querier.Refund(wctx, tc.request)
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
	s := SetupTestSuite(t, 100)
	srApp, ctx, querier := s.srApp, s.ctx, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	refunds := createNRefund(&srApp.SlashrefundKeeper, ctx, 5)

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
		for i := 0; i < len(refunds); i += step {
			resp, err := querier.RefundAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Refund), step)
			require.Subset(t,
				nullify.Fill(refunds),
				nullify.Fill(resp.Refund),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(refunds); i += step {
			resp, err := querier.RefundAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Refund), step)
			require.Subset(t,
				nullify.Fill(refunds),
				nullify.Fill(resp.Refund),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.RefundAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(refunds), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(refunds),
			nullify.Fill(resp.Refund),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.RefundAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
