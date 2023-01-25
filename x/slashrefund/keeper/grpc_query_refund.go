package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Querier) RefundAll(c context.Context, req *types.QueryAllRefundRequest) (*types.QueryAllRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var refunds []types.Refund
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	refundStore := prefix.NewStore(store, types.KeyPrefix(types.RefundKeyPrefix))

	pageRes, err := query.Paginate(refundStore, req.Pagination, func(key []byte, value []byte) error {
		var refund types.Refund
		if err := k.cdc.Unmarshal(value, &refund); err != nil {
			return err
		}

		refunds = append(refunds, refund)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRefundResponse{Refund: refunds, Pagination: pageRes}, nil
}

func (k Querier) Refund(c context.Context, req *types.QueryGetRefundRequest) (*types.QueryGetRefundResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Delegator == "" {
		return nil, status.Error(codes.InvalidArgument, "depositor address cannot be empty")
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

	val, found := k.GetRefund(ctx, delAddr, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRefundResponse{Refund: val}, nil
}
