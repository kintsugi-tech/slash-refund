package keeper

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Manages the request of a user to withdraw previously deposited token from the module. The amount received will be
// based on the amount of shares the user holds and the amount of tokens associated to a validator. The ration between
// the tokens associated to a validator and the shares can change due to slashing events.
func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// === VALIDATION CHECKS ===
	// Check if valid validator address.
	validatorAddress, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	// Check if valid depositor address.
	depositorAddress, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, err
	}

	// Check if requested amount is valid and returns associated shares.
	shares, err := k.ComputeAssociatedShares(ctx, depositorAddress, validatorAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	// === STATE TRANSITION ===
	_, completionTime, err := k.Keeper.Withdraw(ctx, depositorAddress, validatorAddress, msg.Amount.Denom, shares)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyToken, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress),
		),
	})

	return &types.MsgWithdrawResponse{CompletionTime: completionTime}, nil
}
