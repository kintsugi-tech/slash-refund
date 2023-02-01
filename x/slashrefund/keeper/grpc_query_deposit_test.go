package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDepositQuerySingle(t *testing.T) {
	s := SetupTestSuite(t, 100)
	srApp, ctx, testAddrs, valAddrs, querier := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	dep1 := types.NewDeposit(testAddrs[0], valAddrs[0], sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetDeposit(ctx, dep1)

	dep2 := types.NewDeposit(testAddrs[1], valAddrs[0], sdk.NewDec(100))
	srApp.SlashrefundKeeper.SetDeposit(ctx, dep2)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDepositRequest
		response *types.QueryGetDepositResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: testAddrs[0].String(),
				ValidatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetDepositResponse{Deposit: dep1},
		},
		{
			desc: "Second",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: testAddrs[1].String(),
				ValidatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetDepositResponse{Deposit: dep2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: testAddrs[1].String(),
				ValidatorAddress: valAddrs[1].String(),
			},
			err: status.Errorf(codes.NotFound, "deposit with depositor %s not found for validator %s", testAddrs[1].String(), valAddrs[1].String()),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {

			response, err := querier.Deposit(wctx, tc.request)

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
func TestDepositQueryPaginated(t *testing.T) {
	s := SetupTestSuite(t, 100)
	srApp, ctx, querier := s.srApp, s.ctx, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	deposits := createNDeposit(&srApp.SlashrefundKeeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllDepositRequest {
		return &types.QueryAllDepositRequest{
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
		for i := 0; i < len(deposits); i += step {
			resp, err := querier.DepositAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t,
				nullify.Fill(deposits),
				nullify.Fill(resp.Deposit),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(deposits); i += step {
			resp, err := querier.DepositAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t,
				nullify.Fill(deposits),
				nullify.Fill(resp.Deposit),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.DepositAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(deposits), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(deposits),
			nullify.Fill(resp.Deposit),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.DepositAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
