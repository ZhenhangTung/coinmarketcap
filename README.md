# Coinmarketcap
Go wrapper for [Coinmarketcap's API V2](https://coinmarketcap.com/api/).

## How to use it?
### [Listings](https://coinmarketcap.com/api/#endpoint_listings)

```
import "github.com/ZhenhangTung/coinmarketcap"

listings, err := coinmarketcap.GetListings()
if err != nil {
    log.Error(err)
    return
}
fmt.Println(listings)
```

### [Ticker](https://coinmarketcap.com/api/#endpoint_ticker)

### [Ticker (Specific Currency)](https://coinmarketcap.com/api/#endpoint_ticker_specific_cryptocurrency)

### [Global Data](https://coinmarketcap.com/api/#endpoint_global_data)

## TODO: Unit tests