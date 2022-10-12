package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) UnbondingDepositAll(c context.Context, req *types.QueryAllUnbondingDepositRequest) (*types.QueryAllUnbondingDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var unbondingDeposits []types.UnbondingDeposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	unbondingDepositStore := prefix.NewStore(store, types.KeyPrefix(types.UnbondingDepositKey))

	pageRes, err := query.Paginate(unbondingDepositStore, req.Pagination, func(key []byte, value []byte) error {
		var unbondingDeposit types.UnbondingDeposit
		if err := k.cdc.Unmarshal(value, &unbondingDeposit); err != nil {
			return err
		}

		unbondingDeposits = append(unbondingDeposits, unbondingDeposit)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllUnbondingDepositResponse{UnbondingDeposit: unbondingDeposits, Pagination: pageRes}, nil
}

func (k Keeper) UnbondingDeposit(c context.Context, req *types.QueryGetUnbondingDepositRequest) (*types.QueryGetUnbondingDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	unbondingDeposit, found := k.GetUnbondingDeposit(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: unbondingDeposit}, nil
}