package nameservice

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	WhoisRecords   []Whois   `json:"whois_records"`
	WhichisRecords []Whichis `json:"whichis_records"`
}

func NewGenesisState(whoIsRecords []Whois) GenesisState {
	return GenesisState{WhoisRecords: nil, WhichisRecords: nil}
}

func ValidateGenesis(data GenesisState) error {
	for _, record := range data.WhoisRecords {
		if record.Owner == nil {
			return fmt.Errorf("Invalid WhoisRecord: Owner: %s. Error: Missing Owner", record.Owner)
		}
		if record.Value == "" {
			return fmt.Errorf("Invalid WhoisRecord: Value: %s. Error: Missing Value", record.Value)
		}
		if record.Price == nil {
			return fmt.Errorf("Invalid WhoisRecord: Price: %s. Error: Missing Price", record.Price)
		}
	}
	for _, record := range data.WhichisRecords {
		if record.Owner == nil {
			return fmt.Errorf("Invalid WhichisRecords: Owner: %s. Error: Missing Owner", record.Owner)
		}
		if record.Code == "" {
			return fmt.Errorf("Invalid WhichisRecords: Code: %s. Error: Missing Value", record.Code)
		}
		if record.Price == nil {
			return fmt.Errorf("Invalid WhichisRecords: Price: %s. Error: Missing Price", record.Price)
		}
	}
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		WhoisRecords:   []Whois{},
		WhichisRecords: []Whichis{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.WhoisRecords {
		keeper.SetWhois(ctx, record.Value, record)
	}
	for _, record := range data.WhichisRecords {
		keeper.SetWhichis(ctx, record.Code, record)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []Whois
	var records2 []Whichis
	iterator := k.GetNamesIterator(ctx)
	iterator2 := k.GetCodesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		name := string(iterator.Key())
		var whois Whois
		whois = k.GetWhois(ctx, name)
		records = append(records, whois)
	}
	for ; iterator2.Valid(); iterator2.Next() {
		code := string(iterator2.Key())
		var whichis Whichis
		whichis = k.GetWhichis(ctx, code)
		records2 = append(records2, whichis)
	}

	return GenesisState{WhoisRecords: records, WhichisRecords: records2}
}
