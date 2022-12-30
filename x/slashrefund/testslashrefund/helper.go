package testslashrefund

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund/keeper"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

// Helper is a structure which wraps the slashrefund message server
// and provides methods useful in tests
type Helper struct {
	t       *testing.T
	k       keeper.Keeper
	msgSrvr types.MsgServer
	ctx     sdk.Context
}

func NewHelper(t *testing.T, ctx sdk.Context, k keeper.Keeper) *Helper {
	return &Helper{t, k, keeper.NewMsgServerImpl(k), ctx}
}

func (srh *Helper) Deposit(depAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Int) {
	coin := sdk.NewCoin(srh.k.AllowedTokens(srh.ctx)[0], amount)
	msg := types.NewMsgDeposit(depAddr.String(), valAddr.String(), coin)
	res, err := srh.msgSrvr.Deposit(sdk.WrapSDKContext(srh.ctx), msg)
	require.NoError(srh.t, err)
	require.NotNil(srh.t, res)
}

func (srh *Helper) Withdraw(depAddr sdk.AccAddress, valAddr sdk.ValAddress, amount sdk.Int) {
	coin := sdk.NewCoin(srh.k.AllowedTokens(srh.ctx)[0], amount)
	msg := types.NewMsgWithdraw(depAddr.String(), valAddr.String(), coin)
	res, err := srh.msgSrvr.Withdraw(sdk.WrapSDKContext(srh.ctx), msg)
	require.NoError(srh.t, err)
	require.NotNil(srh.t, res)
}
