package nameservice

import (
	"github.com/cosmos/cosmos-sdk-app-jewelry/x/nameservice/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewMsgBuyName = types.NewMsgBuyName
	NewMsgBuyCode = types.NewMsgBuyCode
	NewMsgSetName = types.NewMsgSetName
	NewMsgSetCode = types.NewMsgSetCode
	NewWhois      = types.NewWhois
	NewWhichis    = types.NewWhichis
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgSetName      = types.MsgSetName
	MsgSetCode      = types.MsgSetCode
	MsgBuyName      = types.MsgBuyName
	MsgBuyCode      = types.MsgBuyCode
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	QueryResCodes   = types.QueryResCodes
	Whois           = types.Whois
	Whichis         = types.Whichis
)
