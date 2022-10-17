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

func (k Querier) DepositAll(c context.Context, req *types.QueryAllDepositRequest) (*types.QueryAllDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var deposits []types.Deposit
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	depositStore := prefix.NewStore(store, types.KeyPrefix(types.DepositKeyPrefix))

	pageRes, err := query.Paginate(depositStore, req.Pagination, func(key []byte, value []byte) error {
		var deposit types.Deposit
		if err := k.cdc.Unmarshal(value, &deposit); err != nil {
			return err
		}

		deposits = append(deposits, deposit)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllDepositResponse{Deposit: deposits, Pagination: pageRes}, nil
}

// TODO: aggiunta q a caso di query
func (k Querier) Deposit(c context.Context, req *types.QueryGetDepositRequest) (*types.QueryGetDepositResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetDeposit(
		ctx,
		sdk.AccAddress(req.DepositorAddress),
		sdk.ValAddress(req.ValidatorAddress),
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetDepositResponse{Deposit: val}, nil
}
