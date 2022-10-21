package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// logger := k.Logger(ctx)
	// logger.Error("Entrati nel Msg Server Withdraw")

	// === VALIDATION CHECKS ===
	//Check if valid validator address
	validatorAddress, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	// Check if valid depositor address
	depositorAddress, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, err
	}

	// Check if requested amount is valid
	shares, err := k.ValidateWithdrawdAmount(ctx, depositorAddress, validatorAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Check if allowed token for deposit
	isValid, err := k.CheckAllowedTokens(ctx, msg.Amount.Denom)
	if !isValid {
		return nil, err
	}

	// === STATE TRANSITION ===
	newShares, err := k.Keeper.Withdraw(ctx, depositorAddress, validatorAddress, shares)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDeposit,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyToken, msg.Amount.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyNewShares, newShares.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DepositorAddress),
		),
	})

	return &types.MsgWithdrawResponse{}, nil
	// sender, _ := sdk.AccAddressFromBech32(msg.Creator)

	/*
		deposit, isFound := k.GetDeposit(ctx, sdk.AccAddress(msg.DepositorAddress), sdk.ValAddress(msg.ValidatorAddress))

		if !isFound {
			return nil, errors.New("Don't fuck with mib")
		} else if deposit.Shares.Amount.LT(msg.Amount.Amount) {
			return nil, errors.New("Too much zio")
		}

		updated_balance := deposit.Balance.SubAmount(msg.Amount.Amount)

		deposit = types.Deposit{
			Address:          msg.Creator,
			ValidatorAddress: msg.ValidatorAddress,
			Balance:          updated_balance,
		}

		k.SetDeposit(ctx, deposit)

		unbonding_deposit := types.UnbondingDeposit{
			Id: k.GetUnbondingDepositCount(ctx),
			UnbondingStart: ctx.BlockTime(),
			Address: msg.Creator,
			ValidatorAddress: msg.ValidatorAddress,
			Balance: msg.Amount,
		}

		k.AppendUnbondingDeposit(ctx, unbonding_deposit)

	*/

	return &types.MsgWithdrawResponse{}, nil
}
