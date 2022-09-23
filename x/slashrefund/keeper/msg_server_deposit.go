package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var deposit = types.Deposit{
		Address:          msg.Creator,
		ValidatorAddress: msg.ValidatorAddress,
		Balance:          msg.Amount.String(),
	}

	k.SetDeposit(ctx, deposit)
	// TODO: Handling the message
	_ = ctx

	return &types.MsgDepositResponse{}, nil
}
