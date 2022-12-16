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

func (k Keeper) UnbondingDepositAll(c context.Context, req *types.QueryAllUnbondingDepositRequest) (*types.QueryAllUnbondingDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var unbondingDeposits []types.UnbondingDeposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	unbondingDepositStore := prefix.NewStore(store, types.KeyPrefix(string(types.GetUBDsKeyPrefix())))

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

	depAddr, err := sdk.AccAddressFromBech32(req.DepositorAddress)
	if err != nil {
		return nil, err
	}
	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	val, found := k.GetUnbondingDeposit(
		ctx,
		depAddr,
		valAddr,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "unbonding deposit not found")
	}

	return &types.QueryGetUnbondingDepositResponse{UnbondingDeposit: val}, nil
}
