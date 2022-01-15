package entities

// AssetPrice struct
type AssetPrice struct {
	Ticker   string  `json:"ticker,omitempty"`
	Price    float64 `json:"price,omitempty"`
	Currency string  `json:"currency,omitempty"`
}
