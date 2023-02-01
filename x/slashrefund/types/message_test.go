package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func TestMsgDeposit_ValidateBasic(t *testing.T) {
	depAddress := sample.AccAddress()
	valAddress := sdk.ValAddress(sample.AccAddress()).String()
	coin := sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(500))
	tests := []struct {
		name string
		msg  types.MsgDeposit
		err  error
	}{
		{
			name: "empty depositor address",
			msg: types.MsgDeposit{
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty validator address",
			msg: types.MsgDeposit{
				DepositorAddress: depAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty amount",
			msg: types.MsgDeposit{
				ValidatorAddress: valAddress,
				DepositorAddress: depAddress,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid depositor address",
			msg: types.MsgDeposit{
				DepositorAddress: "invalid_address",
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid validator address",
			msg: types.MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: "invalid_address",
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "negative amount",
			msg: types.MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.Coin{Denom: types.DefaultAllowedTokens[0], Amount: sdk.ZeroInt().SubRaw(1)},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "zero amount",
			msg: types.MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "valid message",
			msg: types.MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgWithdraw_ValidateBasic(t *testing.T) {
	depAddress := sample.AccAddress()
	valAddress := sdk.ValAddress(sample.AccAddress()).String()
	coin := sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(500))

	tests := []struct {
		name string
		msg  types.MsgWithdraw
		err  error
	}{
		{
			name: "empty depositor address",
			msg: types.MsgWithdraw{
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty validator address",
			msg: types.MsgWithdraw{
				DepositorAddress: depAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty amount",
			msg: types.MsgWithdraw{
				ValidatorAddress: valAddress,
				DepositorAddress: depAddress,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid depositor address",
			msg: types.MsgWithdraw{
				DepositorAddress: "invalid_address",
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid validator address",
			msg: types.MsgWithdraw{
				DepositorAddress: depAddress,
				ValidatorAddress: "invalid_address",
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "negative amount",
			msg: types.MsgWithdraw{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.Coin{Denom: types.DefaultAllowedTokens[0], Amount: sdk.ZeroInt().SubRaw(1)},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "zero amount",
			msg: types.MsgWithdraw{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "valid message",
			msg: types.MsgWithdraw{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestMsgClaim_ValidateBasic(t *testing.T) {

	delAddress := sample.AccAddress()
	valAddress := sdk.ValAddress(sample.AccAddress()).String()

	tests := []struct {
		name string
		msg  types.MsgClaim
		err  error
	}{{
		name: "empty delegator address",
		msg: types.MsgClaim{
			ValidatorAddress: valAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "empty validator address",
		msg: types.MsgClaim{
			DelegatorAddress: delAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "invalid delegator address",
		msg: types.MsgClaim{
			DelegatorAddress: "invalid_address",
			ValidatorAddress: valAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "invalid validator address",
		msg: types.MsgClaim{
			DelegatorAddress: delAddress,
			ValidatorAddress: "invalid_address",
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "valid message",
		msg: types.MsgClaim{
			DelegatorAddress: delAddress,
			ValidatorAddress: valAddress,
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
