package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgCreateValSetPreference = "create_validator_set_preference"
)

var _ sdk.Msg = &MsgCreateValidatorSetPreference{}

// NewMsgCreateValidatorSetPreference creates a msg to create a new denom
func NewMsgCreateValidatorSetPreference(owner sdk.AccAddress, preferences []ValidatorPreference) *MsgCreateValidatorSetPreference {
	return &MsgCreateValidatorSetPreference{
		Owner:       owner.String(),
		Preferences: preferences,
	}
}

func (m MsgCreateValidatorSetPreference) Route() string { return RouterKey }
func (m MsgCreateValidatorSetPreference) Type() string  { return TypeMsgCreateValSetPreference }
func (m MsgCreateValidatorSetPreference) ValidateBasic() error {
	for _, validator := range m.Preferences {
		_, err := sdk.AccAddressFromBech32(validator.ValOperAddress)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid validator address (%s)", err)
		}
	}

	return nil
}

func (m MsgCreateValidatorSetPreference) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgCreateValidatorSetPreference) GetSigners() []sdk.AccAddress {
	var validators []sdk.AccAddress
	for _, validator := range m.Preferences {
		valAddr, _ := sdk.AccAddressFromBech32(validator.ValOperAddress)
		validators = append(validators, valAddr)
	}

	return validators
}
