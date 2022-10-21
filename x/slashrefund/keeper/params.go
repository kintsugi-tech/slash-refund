package keeper

import (
	"strings"
	
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.AllowedTokens(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// AllowedTokens returns the AllowedTokens param
func (k Keeper) AllowedTokens(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyAllowedTokens, &res)
	return
}

func (k Keeper) AllowedTokensList(ctx sdk.Context) (re []string) {
	return strings.Split(k.AllowedTokens(ctx), ",")
}

func (k Keeper) CheckAllowedTokens(ctx sdk.Context, denom string) (bool, error) {
	var isAcceptable bool // default is false
	for _, validToken := range k.AllowedTokensList(ctx) {
		if denom == validToken {
			isAcceptable = true
			break
		}
	}
	if !isAcceptable {
		return false, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s. Allowed tokens are %s", denom, k.AllowedTokens(ctx),
		)
	}
	return true, nil
}