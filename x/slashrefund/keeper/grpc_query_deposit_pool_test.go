package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/nullify"
	//"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestDepositPoolQuerySingle(t *testing.T) {
	s := SetupTestSuite(t, 100)
	srApp, ctx, testAddrs, valAddrs, querier := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	depPool1 := types.NewDepositPool(valAddrs[0], sdk.NewCoin("stake", sdk.NewInt(100)), sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetDepositPool(ctx, depPool1)

	depPool2 := types.NewDepositPool(valAddrs[1], sdk.NewCoin("stake", sdk.NewInt(200)), sdk.NewDec(200))
	srApp.SlashrefundKeeper.SetDepositPool(ctx, depPool2)

	val, err := sdk.ValAddressFromBech32(sdk.ValAddress(testAddrs[2]).String())
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDepositPoolRequest
		response *types.QueryGetDepositPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetDepositPoolResponse{DepositPool: depPool1},
		},
		{
			desc: "Second",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: valAddrs[1].String(),
			},
			response: &types.QueryGetDepositPoolResponse{DepositPool: depPool2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: val.String(),
			},
			err: status.Errorf(codes.NotFound, "deposit pool not found for operator %s", val.String()),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := querier.DepositPool(wctx, tc.request)
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
	s := SetupTestSuite(t, 100)
	srApp, ctx, querier := s.srApp, s.ctx, s.querier

	wctx := sdk.WrapSDKContext(ctx)
	depPools := createNDepositPool(&srApp.SlashrefundKeeper, ctx, 5)

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
			resp, err := querier.DepositPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
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
			resp, err := querier.DepositPoolAll(wctx, request(next, 0, uint64(step), false))
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
		resp, err := querier.DepositPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(depPools), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(depPools),
			nullify.Fill(resp.DepositPool),
		)
	})
	t.Run("Invalid Request", func(t *testing.T) {
		_, err := querier.DepositPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
