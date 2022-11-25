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

func (k Keeper) RefundPoolAll(c context.Context, req *types.QueryAllRefundPoolRequest) (*types.QueryAllRefundPoolResponse, error) {
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

func (k Keeper) RefundPool(c context.Context, req *types.QueryGetRefundPoolRequest) (*types.QueryGetRefundPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	valOperAddr, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		//TODO: change "panic(err) with "return nil, err" where needed in all files
		panic(err)
	}

	val, found := k.GetRefundPool(
		ctx,
		valOperAddr,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetRefundPoolResponse{RefundPool: val}, nil
}
