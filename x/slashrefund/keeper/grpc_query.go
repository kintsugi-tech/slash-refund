package keeper

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"google.golang.org/grpc/codes"
)

type queryServer struct {
	K Keeper
}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{k}
}

// -------------------------------------------------------------------------------------------------
// Params
// -------------------------------------------------------------------------------------------------
func (q queryServer) Params(
	c context.Context,
	req *types.QueryParamsRequest,
) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: q.K.GetParams(ctx)}, nil
}

// -------------------------------------------------------------------------------------------------
// Deposit
// -------------------------------------------------------------------------------------------------
// Query to get a single deposit associated to the tuple (depositor, validator).
func (q queryServer) Deposit(
	c context.Context,
	req *types.QueryGetDepositRequest,
) (*types.QueryGetDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.DepositorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "depositor address cannot be empty")
	}
	if req.ValidatorAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	depAddr, err := sdk.AccAddressFromBech32(req.DepositorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	val, found := q.K.GetDeposit(ctx, depAddr, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound, "deposit with depositor %s not found for validator %s",
			req.DepositorAddress, req.ValidatorAddress,
		)
	}

	return &types.QueryGetDepositResponse{Deposit: val}, nil
}

// Query to get all stored deposits.
func (q queryServer) DepositAll(
	c context.Context,
	req *types.QueryAllDepositRequest,
) (*types.QueryAllDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var deposits []types.Deposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.K.storeKey)
	depositStore := prefix.NewStore(store, types.KeyPrefix(types.DepositKeyPrefix))

	pageRes, err := query.Paginate(
		depositStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			var deposit types.Deposit
			if err := q.K.cdc.Unmarshal(value, &deposit); err != nil {
				return err
			}
			deposits = append(deposits, deposit)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDepositResponse{Deposit: deposits, Pagination: pageRes}, nil
}

// -------------------------------------------------------------------------------------------------
// DepositPool
// -------------------------------------------------------------------------------------------------
// Query to get, if exists, the deposit pool associated to a single validator.
func (q queryServer) DepositPool(
	c context.Context,
	req *types.QueryGetDepositPoolRequest,
) (*types.QueryGetDepositPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	valAddr, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := q.K.GetDepositPool(ctx, valAddr)
	if !found {
		return nil, status.Errorf(
			codes.NotFound,
			"deposit pool not found for operator %s", req.OperatorAddress,
		)
	}

	return &types.QueryGetDepositPoolResponse{DepositPool: val}, nil
}

// Query to get all stored deposit pools.
func (q queryServer) DepositPoolAll(
	c context.Context,
	req *types.QueryAllDepositPoolRequest,
) (*types.QueryAllDepositPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var depositPools []types.DepositPool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.K.storeKey)
	depositPoolStore := prefix.NewStore(store, types.KeyPrefix(types.DepositPoolKeyPrefix))

	pageRes, err := query.Paginate(
		depositPoolStore,
		req.Pagination, func(key []byte, value []byte) error {
			var depositPool types.DepositPool
			if err := q.K.cdc.Unmarshal(value, &depositPool); err != nil {
				return err
			}
			depositPools = append(depositPools, depositPool)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDepositPoolResponse{DepositPool: depositPools, Pagination: pageRes}, nil
}

// -------------------------------------------------------------------------------------------------
// UnbondingDeposit
// -------------------------------------------------------------------------------------------------
// Query to get a single unbonding deposit associated to the tuple (depositor, validator).
func (q queryServer) UnbondingDeposit(
	c context.Context,
	req *types.QueryGetUnbondingDepositRequest,
) (*types.QueryGetUnbondingDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	depAddr, err := sdk.AccAddressFromBech32(req.DepositorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	val, found := q.K.GetUnbondingDeposit(ctx, depAddr, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: val}, nil
}

// Query to get all stored unbonding deposits.
func (q queryServer) UnbondingDepositAll(
	c context.Context,
	req *types.QueryAllUnbondingDepositRequest,
) (*types.QueryAllUnbondingDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var unbondingDeposits []types.UnbondingDeposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.K.storeKey)
	unbondingDepositStore := prefix.NewStore(
		store,
		types.KeyPrefix(string(types.GetUBDsKeyPrefix())),
	)

	pageRes, err := query.Paginate(
		unbondingDepositStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			var unbondingDeposit types.UnbondingDeposit
			if err := q.K.cdc.Unmarshal(value, &unbondingDeposit); err != nil {
				return err
			}
			unbondingDeposits = append(unbondingDeposits, unbondingDeposit)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllUnbondingDepositResponse{
		UnbondingDeposit: unbondingDeposits,
		Pagination:       pageRes,
	}, nil
}

// -------------------------------------------------------------------------------------------------
// Refund
// -------------------------------------------------------------------------------------------------
// Query to get the refund associated to the touple (delegator, validator).
func (q queryServer) Refund(
	c context.Context,
	req *types.QueryGetRefundRequest,
) (*types.QueryGetRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "delegator address cannot be empty")
	}
	if req.Validator == "" {
		return nil, status.Error(codes.InvalidArgument, "validator address cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)

	delAddr, err := sdk.AccAddressFromBech32(req.Delegator)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(req.Validator)
	if err != nil {
		return nil, err
	}

	val, found := q.K.GetRefund(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRefundResponse{Refund: val}, nil
}

// Query to get all the stored refund.
func (q queryServer) RefundAll(
	c context.Context,
	req *types.QueryAllRefundRequest,
) (*types.QueryAllRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var refunds []types.Refund
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.K.storeKey)
	refundStore := prefix.NewStore(store, types.KeyPrefix(types.RefundKeyPrefix))

	pageRes, err := query.Paginate(
		refundStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			var refund types.Refund
			if err := q.K.cdc.Unmarshal(value, &refund); err != nil {
				return err
			}

			refunds = append(refunds, refund)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRefundResponse{Refund: refunds, Pagination: pageRes}, nil
}

// -------------------------------------------------------------------------------------------------
// RefundPool
// -------------------------------------------------------------------------------------------------
// Query to get the refund pool associated to a validator.
func (q queryServer) RefundPool(
	c context.Context,
	req *types.QueryGetRefundPoolRequest,
) (*types.QueryGetRefundPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		return nil, err
	}

	refPool, found := q.K.GetRefundPool(ctx, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRefundPoolResponse{RefundPool: refPool}, nil
}

// Query to get all stored refund pools.
func (q queryServer) RefundPoolAll(
	c context.Context,
	req *types.QueryAllRefundPoolRequest,
) (*types.QueryAllRefundPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var refundPools []types.RefundPool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(q.K.storeKey)
	refundPoolStore := prefix.NewStore(store, types.KeyPrefix(types.RefundPoolKeyPrefix))

	pageRes, err := query.Paginate(
		refundPoolStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			var refundPool types.RefundPool
			if err := q.K.cdc.Unmarshal(value, &refundPool); err != nil {
				return err
			}

			refundPools = append(refundPools, refundPool)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRefundPoolResponse{RefundPool: refundPools, Pagination: pageRes}, nil
}
