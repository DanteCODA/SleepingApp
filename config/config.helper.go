
package config

import (
	"fmt"
)

// AllowDomain const
const AllowDomain = "ca.finance.yahoo.com"

// DomainGlob const
const DomainGlob = "*yahoo.*"

// GetPriceByTickerURL get price url
func GetPriceByTickerURL(ticker string) string {
	return fmt.Sprintf("https://ca.finance.yahoo.com/quote/%s", ticker)
}