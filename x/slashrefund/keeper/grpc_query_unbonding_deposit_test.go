package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/made-in-block/slash-refund/testutil/keeper"
	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestUnbondingDepositQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNUnbondingDeposit(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetUnbondingDepositRequest
		response *types.QueryGetUnbondingDepositResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetUnbondingDepositRequest{Id: msgs[0].Id},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetUnbondingDepositRequest{Id: msgs[1].Id},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetUnbondingDepositRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.UnbondingDeposit(wctx, tc.request)
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
	keeper, ctx := keepertest.SlashrefundKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNUnbondingDeposit(keeper, ctx, 5)

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
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.UnbondingDepositAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.UnbondingDeposit),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.UnbondingDepositAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.UnbondingDeposit),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.UnbondingDepositAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.UnbondingDeposit),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.UnbondingDepositAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
