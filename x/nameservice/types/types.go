package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Initial Starting Price for a name that was never previously owned
var MinNamePrice = sdk.Coins{sdk.NewInt64Coin("nametoken", 1)}

// Whois is a struct that contains all the metadata of a name
type Whois struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
	Price sdk.Coins      `json:"price"`
}

// Whichis is a struct that contains all the metadata of a name
type Whichis struct {
	Code         string         `json:"code"`
	Carat        string         `json:"carat"`
	Cut          string         `json:"cut"`
	Clarity      string         `json:"clarity"`
	Color        string         `json:"color"`
	Fluorescence string         `json:"fluorescence"`
	Owner        sdk.AccAddress `json:"owner"`
	Price        sdk.Coins      `json:"price"`
}

// Returns a new Whois with the minprice as the price
func NewWhois() Whois {
	return Whois{
		Price: MinNamePrice,
	}
}

// Returns a new Whois with the minprice as the price
func NewWhichis() Whichis {
	return Whichis{
		Price: MinNamePrice,
	}
}

// implement fmt.Stringer
func (w Whois) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Value: %s
Price: %s`, w.Owner, w.Value, w.Price))
}

// implement fmt.Stringer
func (w Whichis) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
Code: %s
Carat: %s
Cut: %s
Clarity: %s
Color: %s
Fluorescence: %s
Price: %s`, w.Owner, w.Code, w.Carat, w.Cut, w.Clarity, w.Color, w.Fluorescence, w.Price))
}
