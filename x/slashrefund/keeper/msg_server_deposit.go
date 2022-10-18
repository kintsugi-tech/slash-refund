package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	// Check if allowed token
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
	
	// Keeper method for state transition
	shares, err := k.Keeper.Deposit(ctx, depositorAddress, msg.Amount, validator)
	if err != nil {
		return nil, err
	}

	_ = shares

	return &types.MsgDepositResponse{}, nil
}
