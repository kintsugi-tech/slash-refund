package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// Manages the deposit of funds from a user to a particular validator into a module KVStore.

// TODO: check if  valid validator.
func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check if allowed token
	var isAcceptable bool // default is false
	for _, validToken := range strings.Split(k.AllowedTokens(ctx), ",") {
		if msg.Amount.Denom == validToken {
			isAcceptable = true
			break
		}
	}
	if !isAcceptable {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s. Allowed tokens are %s", msg.Amount.Denom, k.AllowedTokens(ctx),
		)
	}

	// deposit coins

	sender, _ := sdk.AccAddressFromBech32(msg.Creator)

	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{msg.Amount})
	if err != nil {
		return nil, err
	}

	deposit, isFound := k.GetDeposit(ctx, msg.Creator, msg.ValidatorAddress)

	balance := msg.Amount
	if isFound {
		balance = balance.AddAmount(deposit.Balance.Amount)
	}

	deposit = types.Deposit{
		Address:          msg.Creator,
		ValidatorAddress: msg.ValidatorAddress,
		Balance:          balance,
	}

	k.SetDeposit(ctx, deposit)

	return &types.MsgDepositResponse{}, nil
}
