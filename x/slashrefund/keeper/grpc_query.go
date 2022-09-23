package keeper

import (
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

var _ types.QueryServer = Keeper{}
