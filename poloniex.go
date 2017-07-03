// Package Poloniex is an implementation of the Poloniex API in Golang.
package poloniex

import (
	"encoding/json"
	"fmt"
)

// New return a instantiate poloniex struct
func New(apiKey, apiSecret string) *Poloniex {
	client := NewClient(apiKey, apiSecret)
	return &Poloniex{client}
}

// poloniex represent a poloniex client
type Poloniex struct {
	client *client
}

// GetTickers is used to get the ticker for all markets
func (b *Poloniex) ReturnTicker() (tickers map[string]Ticker, err error) {
	r, err := b.client.do("GET", "returnTicker", nil)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &tickers); err != nil {
		return
	}
	return
}

// GetVolumes is used to get the volume for all markets
func (b *Poloniex) Return24Volume() (vc VolumeCollection, err error) {
	r, err := b.client.do("GET", "return24hVolume", nil)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &vc); err != nil {
		return
	}
	return
}

func (b *Poloniex) GetCurrencies() (currencies Currencies, err error) {
	r, err := b.client.do("GET", "returnCurrencies", nil)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &currencies.Pair); err != nil {
		return
	}
	return
}

// GetOrderBook is used to get retrieve the orderbook for a given market
// market: a string literal for the market (ex: BTC_NXT). 'all' not implemented.
// cat: bid, ask or both to identify the type of orderbook to return.
// depth: how deep of an order book to retrieve
func (b *Poloniex) ReturnOrderBook(currencyPair string, depth int) (orderBook OrderBook, err error) {
	if depth > 100 {
		depth = 100
	}
	if depth < 1 {
		depth = 1
	}

	args := map[string]string{
		"currencyPair": currencyPair,
		"depth":        fmt.Sprintf("%v", depth)}

	r, err := b.client.do("GET", "returnOrderBook", args)
	if err != nil {
		return
	}
	if err = json.Unmarshal(r, &orderBook); err != nil {
		return
	}

	return
}

// Returns candlestick chart data. Required GET parameters are "currencyPair",
// "period" (candlestick period in seconds; valid values are 300, 900, 1800,
// 7200, 14400, and 86400), "start", and "end". "Start" and "end" are given in
// UNIX timestamp format and used to specify the date range for the data
// returned.
func (b *Poloniex) ReturnChartData(currencyPair string, start, end, period int) (chart []Chart, err error) {
	args := map[string]string{
		"currencyPair": currencyPair,
		"start":        fmt.Sprintf("%v", start),
		"end":          fmt.Sprintf("%v", end),
		"period":       fmt.Sprintf("%v", period)}

	r, err := b.client.do("GET", "returnChartData", args)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &chart); err != nil {
		return
	}

	return
}

func (b *Poloniex) ReturnLoanOrders(currency string) (loans Loan, err error) {
	args := map[string]string{
		"currency": currency}

	r, err := b.client.do("GET", "returnLoanOrders", args)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &loans); err != nil {
		return
	}

	return
}

func (b *Poloniex) ReturnTradeHistory(currencyPair string, start, end int) (trades []*Trade, err error) {
	args := map[string]string{
		"currencyPair": currencyPair,
		"start":        fmt.Sprintf("%v", start),
		"end":          fmt.Sprintf("%v", end)}

	r, err := b.client.do("GET", "returnTradeHistory", args)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &trades); err != nil {
		return
	}

	return
}

func (b *Poloniex) ReturnBalances() (balances []Balance, err error) {
	r, err := b.client.do("POST", "returnBalances", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &balances); err != nil {
		return
	}

	return
}

func (b *Poloniex) ReturnCompleteBalances() (completeBalances []CompleteBalance, err error) {
	r, err := b.client.do("POST", "returnCompleteBalances", nil)
	if err != nil {
		return
	}

	if err = json.Unmarshal(r, &completeBalances); err != nil {
		return
	}

	return
}
