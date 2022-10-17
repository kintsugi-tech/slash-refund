package keeper

import (
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}
