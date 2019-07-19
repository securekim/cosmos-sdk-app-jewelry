package nameservice

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

// Gets the entire Whois metadata struct for a name
func (k Keeper) GetWhois(ctx sdk.Context, name string) Whois {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(name)) {
		return NewWhois()
	}
	bz := store.Get([]byte(name))
	var whois Whois
	k.cdc.MustUnmarshalBinaryBare(bz, &whois)
	return whois
}

//GetWhichis Gets the entire Whichis metadata struct for a code
func (k Keeper) GetWhichis(ctx sdk.Context, code string) Whichis {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(code)) {
		return NewWhichis()
	}
	bz := store.Get([]byte(code))
	var whichis Whichis
	k.cdc.MustUnmarshalBinaryBare(bz, &whichis)
	return whichis
}

//SetWhois Sets the entire Whois metadata struct for a name
func (k Keeper) SetWhois(ctx sdk.Context, name string, whois Whois) {
	if whois.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(whois))
}

//SetWhichis Sets the entire Whichis metadata struct for a code
func (k Keeper) SetWhichis(ctx sdk.Context, code string, whichis Whichis) {
	if whichis.Owner.Empty() {
		return
	}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(code), k.cdc.MustMarshalBinaryBare(whichis))
}

// // ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	return k.GetWhois(ctx, name).Value
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveCode(ctx sdk.Context, code string) string {
	return k.GetWhichis(ctx, code).Code
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, name, whois)
}

// SetCode - sets the value string that a name resolves to
func (k Keeper) SetCode(ctx sdk.Context, code string, carat string, cut string, clarity string, color string, fluorescence string) {
	//[code] [carat] [cut] [clarity] [color] [fluorescence]
	whichis := k.GetWhichis(ctx, code)
	whichis.Code = code
	whichis.Carat = carat
	whichis.Cut = cut
	whichis.Clarity = clarity
	whichis.Color = color
	whichis.Fluorescence = fluorescence
	k.SetWhichis(ctx, code, whichis)
}

// HasOwner - returns whether or not the code already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, code string) bool {
	return !k.GetWhichis(ctx, code).Owner.Empty()
}

// GetOwner - get the current owner of a code
func (k Keeper) GetOwner(ctx sdk.Context, code string) sdk.AccAddress {
	return k.GetWhichis(ctx, code).Owner
}

// SetOwner - sets the current owner of a code
func (k Keeper) SetOwner(ctx sdk.Context, code string, owner sdk.AccAddress) {
	whichis := k.GetWhichis(ctx, code)
	whichis.Owner = owner
	k.SetWhichis(ctx, code, whichis)
}

// GetPrice - gets the current price of a code
func (k Keeper) GetPrice(ctx sdk.Context, code string) sdk.Coins {
	return k.GetWhichis(ctx, code).Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, code string, price sdk.Coins) {
	whichis := k.GetWhichis(ctx, code)
	whichis.Price = price
	k.SetWhichis(ctx, code, whichis)
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetCodesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
