package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgClaim_ValidateBasic(t *testing.T) {

	delAddress := sample.AccAddress()
	valAddress := sdk.ValAddress(sample.AccAddress()).String()

	tests := []struct {
		name string
		msg  MsgClaim
		err  error
	}{{
		name: "empty delegator address",
		msg: MsgClaim{
			ValidatorAddress: valAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "empty validator address",
		msg: MsgClaim{
			DelegatorAddress: delAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "invalid delegator address",
		msg: MsgClaim{
			DelegatorAddress: "invalid_address",
			ValidatorAddress: valAddress,
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "invalid validator address",
		msg: MsgClaim{
			DelegatorAddress: delAddress,
			ValidatorAddress: "invalid_address",
		},
		err: sdkerrors.ErrInvalidAddress,
	}, {
		name: "valid message",
		msg: MsgClaim{
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
