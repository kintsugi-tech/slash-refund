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

func (k Keeper) DepositPoolAll(c context.Context, req *types.QueryAllDepositPoolRequest) (*types.QueryAllDepositPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var depositPools []types.DepositPool
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	depositPoolStore := prefix.NewStore(store, types.KeyPrefix(types.DepositPoolKeyPrefix))

	pageRes, err := query.Paginate(depositPoolStore, req.Pagination, func(key []byte, value []byte) error {
		var depositPool types.DepositPool
		if err := k.cdc.Unmarshal(value, &depositPool); err != nil {
			return err
		}

		depositPools = append(depositPools, depositPool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDepositPoolResponse{DepositPool: depositPools, Pagination: pageRes}, nil
}

func (k Keeper) DepositPool(c context.Context, req *types.QueryGetDepositPoolRequest) (*types.QueryGetDepositPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	valOperAddr, err := sdk.ValAddressFromBech32(req.OperatorAddress)
	if err != nil {
		panic(err)
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetDepositPool(
		ctx,
		valOperAddr,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDepositPoolResponse{DepositPool: val}, nil
}
