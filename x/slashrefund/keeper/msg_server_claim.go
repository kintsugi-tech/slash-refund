package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k msgServer) Claim(goCtx context.Context, msg *types.MsgClaim) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// === VALIDATION CHECKS ===
	//Check if valid validator address
	validatorAddress, err := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return nil, err
	}

	// Check if valid depositor address
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// === STATE TRANSITION ===
	coins, err := k.Keeper.Claim(ctx, delegatorAddress, validatorAddress)
	if err != nil {
		return nil, err
	}

	for _, coin := range coins {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeWithdraw,
				sdk.NewAttribute(types.AttributeKeyDelegator, msg.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyToken, coin.Denom),
				sdk.NewAttribute(sdk.AttributeKeyAmount, coin.Amount.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			),
		})
	}

	return &types.MsgClaimResponse{}, nil
}
