package keeper

import (
	"context"
	"time"

	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type msgServer struct {
	k Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = msgServer{}

// Deposit manages the deposit of funds from a user to a particular validator into the module's 
// KVStore.
func (ms msgServer) Deposit(
	goCtx context.Context, 
	msg *types.MsgDeposit,
) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// === VALIDATION CHECKS ===
	// Check if valid validator address
	valAddr, valErr := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if valErr != nil {
		return nil, valErr
	}

	// Check if valAddr correspond to a validator
	validator, found := ms.k.stakingKeeper.GetValidator(ctx, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoValidatorFound
	}

	// Check if valid depositor address
	depositorAddress, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return nil, err
	}

	// Check if allowed token
	isValid, err := ms.k.CheckAllowedTokens(ctx, msg.Amount.Denom)
	if !isValid {
		return nil, err
	}

	// Check if is non-zero deposit
	if msg.Amount.Amount.IsZero() {
		return nil, types.ErrZeroDeposit
	}

	// === STATE TRANSITION ===
	newShares, err := ms.k.Deposit(ctx, depositorAddress, msg.Amount, validator)
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
	})

	return &types.MsgDepositResponse{}, nil
}

// Manages the request of a user to withdraw previously deposited tokens from the module. The amount
// received will be based on the amount of shares the user holds and the amount of tokens associated
// to a validator. The tokens associated to a validator and the shares ratio can change due to
// slashing events.
func (ms msgServer) Withdraw(
	goCtx context.Context, 
	msg *types.MsgWithdraw,
) (*types.MsgWithdrawResponse, error) {
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

	isValid, err := ms.k.CheckAllowedTokens(ctx, msg.Amount.Denom)
	if !isValid {
		return nil, err
	}

	if msg.Amount.Amount.IsZero() {
		return nil, types.ErrZeroWithdraw
	}

	// === STATE TRANSITION ===
	witTokens, completionTime, err := ms.k.Withdraw(
		ctx, 
		depositorAddress, 
		validatorAddress, 
		msg.Amount,
	)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeWithdraw,
			sdk.NewAttribute(types.AttributeKeyValidator, msg.ValidatorAddress),
			sdk.NewAttribute(types.AttributeKeyToken, witTokens.Denom),
			sdk.NewAttribute(sdk.AttributeKeyAmount, witTokens.Amount.String()),
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

func (ms msgServer) Claim(
	goCtx context.Context, 
	msg *types.MsgClaim,
) (*types.MsgClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// === VALIDATION CHECKS ===
	// Check if valid validator address
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
	coins, err := ms.k.Claim(ctx, delegatorAddress, validatorAddress)
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
