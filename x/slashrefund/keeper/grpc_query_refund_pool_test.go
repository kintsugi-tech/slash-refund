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

func TestRefundPoolQuerySingle(t *testing.T) {
	s := SetupTestSuite(t, 100)
	srApp, ctx, testAddrs, valAddrs, querier := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	refPool1 := types.NewRefundPool(valAddrs[0], sdk.NewCoin("stake", sdk.NewInt(100)), sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetRefundPool(ctx, refPool1)

	refPool2 := types.NewRefundPool(valAddrs[1], sdk.NewCoin("stake", sdk.NewInt(200)), sdk.NewDec(200))
	srApp.SlashrefundKeeper.SetRefundPool(ctx, refPool2)

	val, err := sdk.ValAddressFromBech32(sdk.ValAddress(testAddrs[2]).String())
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetRefundPoolRequest
		response *types.QueryGetRefundPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetRefundPoolResponse{RefundPool: refPool1},
		},
		{
			desc: "Second",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: valAddrs[1].String(),
			},
			response: &types.QueryGetRefundPoolResponse{RefundPool: refPool2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetRefundPoolRequest{
				OperatorAddress: val.String(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := querier.RefundPool(wctx, tc.request)
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
	s := SetupTestSuite(t, 100)
	srApp, ctx, querier := s.srApp, s.ctx, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	refPools := createNRefundPool(&srApp.SlashrefundKeeper, ctx, 5)

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
		for i := 0; i < len(refPools); i += step {
			resp, err := querier.RefundPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RefundPool), step)
			require.Subset(t,
				nullify.Fill(refPools),
				nullify.Fill(resp.RefundPool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(refPools); i += step {
			resp, err := querier.RefundPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RefundPool), step)
			require.Subset(t,
				nullify.Fill(refPools),
				nullify.Fill(resp.RefundPool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.RefundPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(refPools), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(refPools),
			nullify.Fill(resp.RefundPool),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.RefundPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
