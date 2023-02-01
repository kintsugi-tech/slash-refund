package keeper_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestUnbondingDepositQuerySingle(t *testing.T) {

	s := SetupTestSuite(t, 100)
	srApp, ctx, testAddrs, valAddrs, querier := s.srApp, s.ctx, s.testAddrs, s.valAddrs, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	ubdep1 := types.NewUnbondingDeposit(testAddrs[0], valAddrs[0], 10, time.Unix(10, 0), sdk.NewInt(100))
	entry2 := types.NewUnbondingDepositEntry(20, time.Unix(20, 0), sdk.NewInt(200))
	ubdep1.Entries = append(ubdep1.Entries, entry2)
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep1)

	ubdep2 := types.NewUnbondingDeposit(testAddrs[1], valAddrs[0], 0, time.Unix(30, 0), sdk.NewInt(300))
	entry2 = types.NewUnbondingDepositEntry(40, time.Unix(40, 0), sdk.NewInt(400))
	ubdep2.Entries = append(ubdep2.Entries, entry2)
	srApp.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep2)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetUnbondingDepositRequest
		response *types.QueryGetUnbondingDepositResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: testAddrs[0].String(),
				ValidatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: ubdep1},
		},
		{
			desc: "Second",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: testAddrs[1].String(),
				ValidatorAddress: valAddrs[0].String(),
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: ubdep2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: testAddrs[1].String(),
				ValidatorAddress: valAddrs[1].String(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := querier.UnbondingDeposit(wctx, tc.request)
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

func TestUnbondingDepositQueryPaginated(t *testing.T) {

	s := SetupTestSuite(t, 100)
	srApp, ctx, querier := s.srApp, s.ctx, s.querier
	wctx := sdk.WrapSDKContext(ctx)

	ubds := createNUnbondingDeposit(&srApp.SlashrefundKeeper, ctx, 5, 2)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllUnbondingDepositRequest {
		return &types.QueryAllUnbondingDepositRequest{
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
		for i := 0; i < len(ubds); i += step {
			resp, err := querier.UnbondingDepositAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t,
				nullify.Fill(ubds),
				nullify.Fill(resp.UnbondingDeposit),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(ubds); i += step {
			resp, err := querier.UnbondingDepositAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t,
				nullify.Fill(ubds),
				nullify.Fill(resp.UnbondingDeposit),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.UnbondingDepositAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(ubds), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(ubds),
			nullify.Fill(resp.UnbondingDeposit),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.UnbondingDepositAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
