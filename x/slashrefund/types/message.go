package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeposit = "deposit"
const TypeMsgWithdraw = "withdraw"
const TypeMsgClaim = "claim"

var _ sdk.Msg = &MsgDeposit{}
var _ sdk.Msg = &MsgWithdraw{}
var _ sdk.Msg = &MsgClaim{}

func NewMsgDeposit(depositorAddress string, validatorAddress string, amount sdk.Coin) *MsgDeposit {
	return &MsgDeposit{
		DepositorAddress: depositorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg *MsgDeposit) Route() string {
	return RouterKey
}

func (msg *MsgDeposit) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDeposit) GetSigners() []sdk.AccAddress {
	depositor, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depositor}
}

func (msg *MsgDeposit) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeposit) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if !msg.Amount.IsValid() || msg.Amount.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid denom or non-positive amount: (%s)", msg.Amount.String())
	}

	return nil
}

func NewMsgWithdraw(depositorAddress string, validatorAddress string, amount sdk.Coin) *MsgWithdraw {
	return &MsgWithdraw{
		DepositorAddress: depositorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg *MsgWithdraw) Route() string {
	return RouterKey
}

func (msg *MsgWithdraw) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgWithdraw) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdraw) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid depositor address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	if !msg.Amount.IsValid() || msg.Amount.Amount.IsZero() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid denom or non-positive amount: (%s)", msg.Amount.String())
	}

	return nil
}

func NewMsgClaim(delegatorAddress string, validatorAddress string, amount sdk.Coin) *MsgClaim {
	return &MsgClaim{
		DelegatorAddress: delegatorAddress,
		ValidatorAddress: validatorAddress,
		Amount:           amount,
	}
}

func (msg *MsgClaim) Route() string {
	return RouterKey
}

func (msg *MsgClaim) Type() string {
	return TypeMsgClaim
}

func (msg *MsgClaim) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaim) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaim) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid account address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address (%s)", err)
	}

	return nil
}
