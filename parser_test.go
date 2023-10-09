package parser

import (
	"testing"

	"github.com/HariniThangavel-Stack/market-data-parser-golang/model"
)

// incomplete test
func TestMarketDataParser(t *testing.T) {
	t.Run("Parser main function", func(t *testing.T) {
		options := model.Options{OutputFormat: "json", FilePath: "./examples/MW-NIFTY-BANK-05-Aug-2021.csv"}
		MarketDataParser(options)
	})
}
