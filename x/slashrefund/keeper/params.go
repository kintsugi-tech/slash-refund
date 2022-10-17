package keeper

import (
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
