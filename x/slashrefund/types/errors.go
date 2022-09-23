package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/slashrefund module sentinel errors
var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)
