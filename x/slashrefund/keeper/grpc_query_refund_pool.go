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

func (k Querier) RefundPoolAll(c context.Context, req *types.QueryAllRefundPoolRequest) (*types.QueryAllRefundPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var refundPools []types.RefundPool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	refundPoolStore := prefix.NewStore(store, types.KeyPrefix(types.RefundPoolKeyPrefix))

	pageRes, err := query.Paginate(refundPoolStore, req.Pagination, func(key []byte, value []byte) error {
		var refundPool types.RefundPool
		if err := k.cdc.Unmarshal(value, &refundPool); err != nil {
			return err
		}

		refundPools = append(refundPools, refundPool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllRefundPoolResponse{RefundPool: refundPools, Pagination: pageRes}, nil
}

func (k Querier) RefundPool(c context.Context, req *types.QueryGetRefundPoolRequest) (*types.QueryGetRefundPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valAddr, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		return nil, err
	}

	refPool, found := k.GetRefundPool(ctx, valAddr)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRefundPoolResponse{RefundPool: refPool}, nil
}
