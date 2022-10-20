package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	_ = ctx
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
