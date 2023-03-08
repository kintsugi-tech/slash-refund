package keeper_test

import (
	"testing"
	"time"

	"github.com/made-in-block/slash-refund/testutil/testsuite"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/made-in-block/slash-refund/app"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	queryDelAddrs = 2
	queryValAddrs = 2
	queryDepAddrs = 2
	depToken = types.DefaultAllowedTokens[0]
)

func SetupQueryServerTest() (
	*app.App, 
	sdk.Context,
	types.QueryServer,
	[]sdk.AccAddress, 
	[]sdk.ValAddress,
	[]sdk.AccAddress,
) {

	// Setup delegators
	delAddrs := testsuite.GenerateNAddresses(queryDelAddrs)
	delAccs := testsuite.ConvertAddressesToAccAddr(delAddrs)
	balances := testsuite.GenerateBalances(delAccs)

	// Setup validators
	valAddrs := testsuite.GenerateNAddresses(queryValAddrs)
	valAccs := testsuite.ConvertAddressesToValAddr(valAddrs)

	// Setup depositors
	depAddrs := testsuite.GenerateNAddresses(queryDepAddrs)
	depAccs := testsuite.ConvertAddressesToAccAddr(depAddrs)
	depBalances := testsuite.GenerateBalances(depAccs)

	balances = append(balances, depBalances...)

	app, ctx := testsuite.CreateTestApp(delAccs, valAccs, balances, false)
	qs := keeper.NewQueryServerImpl(app.SlashrefundKeeper)

	return app, ctx, qs, delAccs, valAccs, depAccs
}

// -------------------------------------------------------------------------------------------------
// Params
// -------------------------------------------------------------------------------------------------
func TestQueryParams(t *testing.T) {
	app, ctx, qs,  _, _, _ := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	params := types.DefaultParams()
	app.SlashrefundKeeper.SetParams(ctx, params)

	resp, err := qs.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, resp)
}

// -------------------------------------------------------------------------------------------------
// Test Deposit
// -------------------------------------------------------------------------------------------------
func TestDepositQuerySingle(t *testing.T) {
	app, ctx, qs,  _, validators, depositors := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	dep1 := types.NewDeposit(depositors[0], validators[0], sdk.NewDec(100))
	app.SlashrefundKeeper.SetDeposit(ctx, dep1)

	dep2 := types.NewDeposit(depositors[1], validators[0], sdk.NewDec(100))
	app.SlashrefundKeeper.SetDeposit(ctx, dep2)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDepositRequest
		response *types.QueryGetDepositResponse
		err      error
	}{
		{
			desc: "FirstValid",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: depositors[0].String(),
				ValidatorAddress: validators[0].String(),
			},
			response: &types.QueryGetDepositResponse{Deposit: dep1},
		},
		{
			desc: "SecondValid",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: depositors[1].String(),
				ValidatorAddress: validators[0].String(),
			},
			response: &types.QueryGetDepositResponse{Deposit: dep2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetDepositRequest{
				DepositorAddress: depositors[1].String(),
				ValidatorAddress: validators[1].String(),
			},
			err: status.Errorf(codes.NotFound, "deposit with depositor %s not found for validator %s", depositors[1].String(), validators[1].String()),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {

			response, err := qs.Deposit(wctx, tc.request)

			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestDepositQueryPaginated(t *testing.T) {
	app, ctx, qs,  _, _, _ := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	deposits := testslashrefund.CreateNDeposit(&app.SlashrefundKeeper, ctx, 5)

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
			resp, err := qs.DepositAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t, deposits, resp.Deposit)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(deposits); i += step {
			resp, err := qs.DepositAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t, deposits, resp.Deposit)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.DepositAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(deposits), int(resp.Pagination.Total))
		require.ElementsMatch(t, deposits, resp.Deposit)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.DepositAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

// -------------------------------------------------------------------------------------------------
// Test DepositPool
// -------------------------------------------------------------------------------------------------
func TestDepositPoolQuerySingle(t *testing.T) {
	app, ctx, qs,  _, validators, depositors := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	depPool1 := types.NewDepositPool(
		validators[0], 
		sdk.NewCoin(depToken, sdk.NewInt(100)), 
		sdk.NewDec(100),
	)
	app.SlashrefundKeeper.SetDepositPool(ctx, depPool1)

	depPool2 := types.NewDepositPool(validators[1], 
		sdk.NewCoin(depToken, sdk.NewInt(200)), 
		sdk.NewDec(200),
	)
	app.SlashrefundKeeper.SetDepositPool(ctx, depPool2)

	val, err := sdk.ValAddressFromBech32(sdk.ValAddress(depositors[0]).String())
	require.NoError(t, err)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetDepositPoolRequest
		response *types.QueryGetDepositPoolResponse
		err      error
	}{
		{
			desc: "FirstValid",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: validators[0].String(),
			},
			response: &types.QueryGetDepositPoolResponse{DepositPool: depPool1},
		},
		{
			desc: "SecondValid",
			request: &types.QueryGetDepositPoolRequest{
				OperatorAddress: validators[1].String(),
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
			response, err := qs.DepositPool(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TestDepositPoolQueryPaginated(t *testing.T) {
	app, ctx, qs,  _, _, _ := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	depPools := testslashrefund.CreateNDepositPool(&app.SlashrefundKeeper, ctx, 5)

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
			resp, err := qs.DepositPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DepositPool), step)
			require.Subset(t, depPools, resp.DepositPool)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(depPools); i += step {
			resp, err := qs.DepositPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.DepositPool), step)
			require.Subset(t, depPools, resp.DepositPool)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.DepositPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(depPools), int(resp.Pagination.Total))
		require.ElementsMatch(t, depPools, resp.DepositPool)
	})
	t.Run("Invalid Request", func(t *testing.T) {
		_, err := qs.DepositPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
// -------------------------------------------------------------------------------------------------
// Test UnbondingDeposit
// -------------------------------------------------------------------------------------------------
func TestUnbondingDepositQuerySingle(t *testing.T) {
	app, ctx, qs,  _, validators, depositors := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	ubdep1 := types.NewUnbondingDeposit(
		depositors[0], 
		validators[0], 
		10, 
		time.Unix(10, 0).UTC(), 
		sdk.NewInt(100),
	)
	entry2 := types.NewUnbondingDepositEntry(20, time.Unix(20, 0).UTC(), sdk.NewInt(200))
	ubdep1.Entries = append(ubdep1.Entries, entry2)
	app.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep1)

	ubdep2 := types.NewUnbondingDeposit(
		depositors[1], 
		validators[0], 
		0, 
		time.Unix(30, 0).UTC(), 
		sdk.NewInt(300),
	)
	entry2 = types.NewUnbondingDepositEntry(40, time.Unix(40, 0).UTC(), sdk.NewInt(400))
	ubdep2.Entries = append(ubdep2.Entries, entry2)
	app.SlashrefundKeeper.SetUnbondingDeposit(ctx, ubdep2)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetUnbondingDepositRequest
		response *types.QueryGetUnbondingDepositResponse
		err      error
	}{
		{
			desc: "FirstValid",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: depositors[0].String(),
				ValidatorAddress: validators[0].String(),
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: ubdep1},
		},
		{
			desc: "SecondValid",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: depositors[1].String(),
				ValidatorAddress: validators[0].String(),
			},
			response: &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: ubdep2},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: depositors[1].String(),
				ValidatorAddress: validators[1].String(),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := qs.UnbondingDeposit(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.response, response)
			}
		})
	}
}

/*
func TestUnbondingDepositQueryPaginated(t *testing.T) {
	app, ctx, qs,  _, validators, depositors := SetupQueryServerTest()
	wctx := sdk.WrapSDKContext(ctx)

	ubds := createNUnbondingDeposit(&app.SlashrefundKeeper, ctx, 5, 2)

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
			resp, err := qs.UnbondingDepositAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t, ubds, resp.UnbondingDeposit)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(ubds); i += step {
			resp, err := qs.UnbondingDepositAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t, ubds, resp.UnbondingDeposit)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := qs.UnbondingDepositAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(ubds), int(resp.Pagination.Total))
		require.ElementsMatch(t, ubds, resp.UnbondingDeposit)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := qs.UnbondingDepositAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

// -------------------------------------------------------------------------------------------------
// Test Refund
// -------------------------------------------------------------------------------------------------
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
				require.Equal(t, tc.response, response)
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
			require.Subset(t, refunds, resp.Refund)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(refunds); i += step {
			resp, err := querier.RefundAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Refund), step)
			require.Subset(t, refunds, resp.Refund)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.RefundAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(refunds), int(resp.Pagination.Total))
		require.ElementsMatch(t, refunds, resp.Refund)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.RefundAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}

// -------------------------------------------------------------------------------------------------
// Test RefundPool
// -------------------------------------------------------------------------------------------------
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
				require.Equal(t, tc.response, response)
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
			require.Subset(t, refPools, resp.RefundPool)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(refPools); i += step {
			resp, err := querier.RefundPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.RefundPool), step)
			require.Subset(t, refPools, resp.RefundPool)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := querier.RefundPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(refPools), int(resp.Pagination.Total))
		require.ElementsMatch(t, refPools, resp.RefundPool)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := querier.RefundPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
*/