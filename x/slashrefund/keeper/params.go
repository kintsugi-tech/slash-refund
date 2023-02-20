package keeper

import (

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.AllowedTokens(ctx),
		k.MaxEntries(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// AllowedTokens returns the AllowedTokens param
func (k Keeper) AllowedTokens(ctx sdk.Context) (res []string) {
	k.paramstore.Get(ctx, types.KeyAllowedTokens, &res)
	return
}

// Returns the maxmimum number of unbnding deposit entries.
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return 
}

// Checks if a specific denom is among module's allowed tokens. Returns true or false with an error.
func (k Keeper) CheckAllowedTokens(ctx sdk.Context, denom string) (bool, error) {
	for _, validToken := range k.AllowedTokens(ctx) {
		if denom == validToken {
			return true, nil
		}
	}

	return false, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s. Allowed tokens are %s", denom, k.AllowedTokens(ctx),
	)
}
