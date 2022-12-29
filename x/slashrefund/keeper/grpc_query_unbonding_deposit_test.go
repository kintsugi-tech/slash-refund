package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/nullify"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func TestUnbondingDepositQuerySingle(t *testing.T) {
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	type testcase struct {
		desc     string
		request  *types.QueryGetUnbondingDepositRequest
		response *types.QueryGetUnbondingDepositResponse
		err      error
	}

	// test not found
	var tc testcase

	depPubk := secp256k1.GenPrivKey().PubKey()
	depAddr := sdk.AccAddress(depPubk.Address())
	valPubk := secp256k1.GenPrivKey().PubKey()
	valAddr := sdk.ValAddress(valPubk.Address())

	tc.desc = "KeyNotFound"
	tc.request = &types.QueryGetUnbondingDepositRequest{
		DepositorAddress: depAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
	tc.err = status.Error(codes.NotFound, "unbonding deposit not found")

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

	// test others
	msgs := createNUnbondingDeposit(keeper, ctx, 2, 2)
	for _, tc := range []testcase{
		{
			desc: "First",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: msgs[0].DepositorAddress,
				ValidatorAddress: msgs[0].ValidatorAddress,
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: msgs[1].DepositorAddress,
				ValidatorAddress: msgs[1].ValidatorAddress,
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: msgs[1]},
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
	keeper, ctx := testslashrefund.NewTestKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNUnbondingDeposit(keeper, ctx, 5, 2)

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
