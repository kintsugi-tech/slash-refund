package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Manages the deposit of funds from a user to a particular validator into a module KVStore.
// TODO: add param for allowed tokens.
// TODO: check if allowed token.
func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	logger := k.Logger(ctx)
	logger.Error("Entrati nel Msg Server Deposit")

	// Validation checks
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	depositorAddress, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, err
	}

	// TODO: implementare -> bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != "stake" {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s, expected %s", msg.Amount.Denom, "stake",
		)
	}
	
	// Keeper method for state transition
	shares, err := k.Keeper.Deposit(ctx, depositorAddress, msg.Amount.Amount, validator)
	if err != nil {
		return nil, err
	}

	_ = shares

	return &types.MsgDepositResponse{}, nil
}
