# Watchmarket

[![Build Status](https://dev.azure.com/TrustWallet/WatchMarket/_apis/build/status/trustwallet.watchmarket?branchName=master)](https://dev.azure.com/TrustWallet/WatchMarket/_build/latest?definitionId=45&branchName=master)
[![codecov](https://codecov.io/gh/trustwallet/watchmarket/branch/master/graph/badge.svg)](https://codecov.io/gh/trustwallet/watchmarket)
[![Go Report Card](https://goreportcard.com/badge/github.com/TrustWallet/watchmarket)](https://goreportcard.com/report/github.com/TrustWallet/watchmarket)

> Watchmarket is a Blockchain explorer API aggregator and caching layer. It's your one-stop-shop to get information for (almost) any coin in a common format

Watchmarket comes with three apps:
* API: RESTful API to retrieve coin info, charts, and tickers
* Worker: fetch and parse data from market APIs and store it at DB
* Swagger: API explorer

#### Supported Market APIs

<a href="https://coinmarketcap.com" target="_blank"><img src="https://coinmarketcap.com/apple-touch-icon.png" width="32" /></a>
<a href="https://www.binance.org/" target="_blank"><img src="https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/binance/info/logo.png" width="32" /></a>
<a href="https://fixer.io/" target="_blank"><img src="https://fixer.io/fixer_images/fixer_money.png" width="32" /></a>
<a href="https://www.coingecko.com/" target="_blank"><img src="https://static.coingecko.com/s/thumbnail-007177f3eca19695592f0b8b0eabbdae282b54154e1be912285c9034ea6cbaf2.png" width="32" /></a>

## Getting started

1. Spin up a Redis instance: ` docker run -d -p 5432:5432 -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -e POSTGRES_DB=my_db postgres`
2. Start the app: `make start`
3. Use the API (examples below):

  A. Get coin details about Ehtereum (coin 60 according to SLIPs) 
  
  ```curl -v "http://localhost:8421/v1/market/info?coin=60" | jq .``` 
  
  B. Get current ticker price of Ethereum in USD:
   
   ```curl -v -X POST 'http://localhost:8421/v1/market/ticker' -H 'Content-Type: application/json' -d '{"currency":"USD","assets":[{"type":"coin","coin":60}]}'```
   
  C. Get price interval of Ethereum to build chart starting from 1574483028 (UNIX time)
  
   ```curl -v "http://localhost:8421/v1/market/charts?coin=60&time_start=1574483028" | jq .```
 
Use `make stop` to stop the services

Run `make` to see a list of all available build directives.
