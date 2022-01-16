
package entities

// Asset struct
type Asset struct {
	Ticker           string  `json:"ticker,omitempty"`
	Name             string  `json:"name,omitempty"`
	Type             string  `json:"type,omitempty"`
	AssetClass       string  `json:"assetClass,omitempty"`
	Currency         string  `json:"currency,omitempty"`
	AllocationStock  float64 `json:"allocationStock,omitempty"`
	AllocationBond   float64 `json:"allocationBond,omitempty"`
	AllocationCash   float64 `json:"allocationCash,omitempty"`
	DividendSchedule string  `json:"dividendSchedule,omitempty"`
	Yield12Month     float64 `json:"yield12Month,omitempty"`
	DistYield        float64 `json:"distYield,omitempty"`
	DistAmount       float64 `json:"distAmount,omitempty"`
}