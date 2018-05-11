package coinmarketcap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const ApiURLV2 = "https://api.coinmarketcap.com/v2"
const ApiTicker = "/ticker/"
const ApiListings = "/listings/"
const ApiGlobal = "/global/"

type Listing struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	WebsiteSlug string `json:"website_slug"`
}

type ListingsResponse struct {
	Data     []Listing `json:"data"`
	MetaData MetaData  `json:"metadata"`
}

type Ticker struct {
	Id                int64            `json:"id"`
	Name              string           `json:"name"`
	Symbol            string           `json:"symbol"`
	WebsiteSlug       string           `json:"website_slug"`
	Rank              int64            `json:"rank"`
	CirculatingSupply float64          `json:"circulating_supply"`
	TotalSupply       float64          `json:"total_supply"`
	MaxSupply         float64          `json:"max_supply"`
	Quotes            map[string]Quote `json:"quotes"`
	LastUpdated       int64            `json:"last_updated"`
}

type Quote struct {
	Price          float64 `json:"price"`
	Volume24h      float64 `json:"volume_24h"`
	MarketCap      float64 `json:"market_cap"`
	PriceChange1h  float64 `json:"price_change_1h"`
	PriceChange24h float64 `json:"price_change_24h"`
	PriceChange7d  float64 `json:"price_change_7d"`
}

type TickersResponse struct {
	MetaData MetaData          `json:"metadata"`
	Data     map[string]Ticker `json:"data"`
}

type TickerResponse struct {
	MetaData MetaData `json:"metadata"`
	Data     Ticker   `json:"data"`
}

type MetaData struct {
	Timestamp           int64  `json:"timestamp"`
	NumCryptoCurrencies string `json:"num_cryptocurrencies,omitempty"`
	Error               string `json:"error"`
}

type GlobalData struct {
	ActiveCryptoCurrencies       int64                  `json:"active_cryptocurrencies"`
	ActiveMarkets                int64                  `json:"active_markets"`
	BitcoinPercentageOfMarketCap float64                `json:"bitcoin_percentage_of_market_cap"`
	Quotes                       map[string]GlobalQuote `json:"quotes"`
	LastUpdatedAt                int64                  `json:"last_updated_at"`
}

type GlobalQuote struct {
	TotalMarketCap float64 `json:"total_market_cap"`
	TotalVolume24h float64 `json:"total_volume_24h"`
}

type GlobalDataResponse struct {
	Data     GlobalData `json:"data"`
	MetaData MetaData   `json:"meta_data"`
}

func GetListings() ([]Listing, error) {
	url := fmt.Sprintf("%s%s", ApiURLV2, ApiListings)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("something went wrong with listings API and the response code is %v", resp.StatusCode))
	}
	respBody := ListingsResponse{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	defer resp.Body.Close()
	if respBody.MetaData.Error != "" {
		return nil, errors.New(respBody.MetaData.Error)
	}
	return respBody.Data, nil
}

func GetTicks(start, limit int, convert string) (map[string]Ticker, error) {
	url := fmt.Sprintf("%s%s", ApiURLV2, ApiTicker)
	if start != 0 {
		url += fmt.Sprintf("?start=%v", start)
	}
	if limit != 0 {
		url += fmt.Sprintf("?limit=%v", limit)
	}
	if convert != "" {
		url += fmt.Sprintf("?convert=%v", convert)
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("something went wrong with ticker API and the response code is %v", resp.StatusCode))
	}
	respBody := TickersResponse{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	defer resp.Body.Close()
	if respBody.MetaData.Error != "" {
		return nil, errors.New(respBody.MetaData.Error)
	}
	return respBody.Data, nil
}

func GetTick(id int) (Ticker, error) {
	if id <= 0 {
		return Ticker{}, errors.New("id is required")
	}
	url := fmt.Sprintf(ApiURLV2+ApiTicker+"%v/", id)
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Ticker{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return Ticker{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Ticker{}, errors.New(fmt.Sprintf("something went wrong with ticker/(id) API and the response code is %v", resp.StatusCode))
	}
	respBody := TickerResponse{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	defer resp.Body.Close()
	if respBody.MetaData.Error != "" {
		return Ticker{}, errors.New(respBody.MetaData.Error)
	}
	return respBody.Data, nil
}

func GetGlobalData(convert string) (GlobalData, error) {
	url := fmt.Sprintf("%s%s", ApiURLV2, ApiGlobal)
	if convert != "" {
		url = fmt.Sprintf("%s?convert=%s", url, convert)
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GlobalData{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return GlobalData{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GlobalData{}, errors.New(fmt.Sprintf("something went wrong with global data API and the response code is %v", resp.StatusCode))
	}
	respBody := GlobalDataResponse{}
	json.NewDecoder(resp.Body).Decode(&respBody)
	defer resp.Body.Close()
	if respBody.MetaData.Error != "" {
		return GlobalData{}, errors.New(respBody.MetaData.Error)
	}
	return respBody.Data, nil
}
