package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// MsgSetName defines a SetName message
type MsgSetName struct {
	Name  string         `json:"name"`
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}

// MsgSetCode defines a SetName message
type MsgSetCode struct {
	Code         string         `json:"code"`
	Carat        string         `json:"carat"`
	Cut          string         `json:"cut"`
	Clarity      string         `json:"clarity"`
	Color        string         `json:"color"`
	Fluorescence string         `json:"fluorescence"`
	Owner        sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgSetName(name string, value string, owner sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Name:  name,
		Value: value,
		Owner: owner,
	}
}

// NewMsgSetCode is a constructor function for MsgSetCode
func NewMsgSetCode(code string, carat string, cut string, clarity string, color string, fluorescence string, owner sdk.AccAddress) MsgSetCode {
	return MsgSetCode{
		Code:         code,
		Carat:        carat,
		Cut:          cut,
		Clarity:      clarity,
		Color:        color,
		Fluorescence: fluorescence,
		Owner:        owner,
	}
}

// Route should return the name of the module
func (msg MsgSetName) Route() string { return RouterKey }

// Route should return the name of the module
func (msg MsgSetCode) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSetName) Type() string { return "set_name" }

// Type should return the action
func (msg MsgSetCode) Type() string { return "set_code" }

// ValidateBasic runs stateless checks on the message
func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

// ValidateBasic runs stateless checks on the message
func (msg MsgSetCode) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Code) == 0 || len(msg.Carat) == 0 || len(msg.Cut) == 0 || len(msg.Clarity) == 0 || len(msg.Color) == 0 || len(msg.Fluorescence) == 0 {
		return sdk.ErrUnknownRequest("Code and/or 4C cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgSetCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// GetSigners defines whose signature is required
func (msg MsgSetCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// MsgBuyName defines the BuyName message
type MsgBuyName struct {
	Name  string         `json:"name"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// MsgBuyCode defines the BuyCode message
type MsgBuyCode struct {
	Code  string         `json:"code"`
	Bid   sdk.Coins      `json:"bid"`
	Buyer sdk.AccAddress `json:"buyer"`
}

// NewMsgBuyName is the constructor function for MsgBuyName
func NewMsgBuyName(name string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyName {
	return MsgBuyName{
		Name:  name,
		Bid:   bid,
		Buyer: buyer,
	}
}

// NewMsgBuyCode is the constructor function for MsgBuyCode
func NewMsgBuyCode(code string, bid sdk.Coins, buyer sdk.AccAddress) MsgBuyCode {
	return MsgBuyCode{
		Code:  code,
		Bid:   bid,
		Buyer: buyer,
	}
}

// Route should return the name of the module
func (msg MsgBuyName) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyName) Type() string { return "buy_name" }

// Route should return the name of the module
func (msg MsgBuyCode) Route() string { return RouterKey }

// Type should return the action
func (msg MsgBuyCode) Type() string { return "buy_code" }

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyName) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Name) == 0 {
		return sdk.ErrUnknownRequest("Name cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// ValidateBasic runs stateless checks on the message
func (msg MsgBuyCode) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}
	if len(msg.Code) == 0 {
		return sdk.ErrUnknownRequest("Code cannot be empty")
	}
	if !msg.Bid.IsAllPositive() {
		return sdk.ErrInsufficientCoins("Bids must be positive")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyName) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSignBytes encodes the message for signing
func (msg MsgBuyCode) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners defines whose signature is required
func (msg MsgBuyName) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}

// GetSigners defines whose signature is required
func (msg MsgBuyCode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
