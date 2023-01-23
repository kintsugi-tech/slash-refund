package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgDeposit_ValidateBasic(t *testing.T) {
	depAddress := sample.AccAddress()
	valAddress := sdk.ValAddress(sample.AccAddress()).String()
	coin := sdk.NewCoin(DefaultAllowedTokens[0], sdk.NewInt(500))
	tests := []struct {
		name string
		msg  MsgDeposit
		err  error
	}{
		{
			name: "empty depositor address",
			msg: MsgDeposit{
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty validator address",
			msg: MsgDeposit{
				DepositorAddress: depAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "empty amount",
			msg: MsgDeposit{
				ValidatorAddress: valAddress,
				DepositorAddress: depAddress,
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "invalid depositor address",
			msg: MsgDeposit{
				DepositorAddress: "invalid_address",
				ValidatorAddress: valAddress,
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid validator address",
			msg: MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: "invalid_address",
				Amount:           coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "negative amount",
			msg: MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.Coin{Denom: DefaultAllowedTokens[0], Amount: sdk.ZeroInt().SubRaw(1)},
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "zero amount",
			msg: MsgDeposit{
				DepositorAddress: depAddress,
				ValidatorAddress: valAddress,
				Amount:           sdk.NewCoin(DefaultAllowedTokens[0], sdk.NewInt(0)),
			},
			err: sdkerrors.ErrInvalidRequest,
		}, {
			name: "valid message",
			msg: MsgDeposit{
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
