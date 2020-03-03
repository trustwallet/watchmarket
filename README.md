# Watchmarket Blockatlas version
Is under development at master branch. Was migrated from [Blockatlas](https://github.com/trustwallet/blockatlas) at version/blockatlas branch

[![Build Status](https://dev.azure.com/TrustWallet/WatchMarket/_apis/build/status/trustwallet.watchmarket?branchName=master)](https://dev.azure.com/TrustWallet/WatchMarket/_build/latest?definitionId=45&branchName=master)
[![codecov](https://codecov.io/gh/trustwallet/watchmarket/branch/develop/graph/badge.svg)](https://codecov.io/gh/trustwallet/watchmarket)

**IMPORTANT!**
Is under development at master branch. Was migrated from [Blockatlas](https://github.com/trustwallet/blockatlas) at version/blockatlas branch

### What is “Watchmarket”?
Watchmarket is an aggregation and caching service for blockchain market information. 
The main features of it are:
1. **Aggregation**: it standardizes information from different blockchain data providers e.g. Coinmarketcap into a unified format
2. **Caching**: it acts as a caching layer for this information mainly catering to the goal of cost savings as API calls to data providers are expensive

### Architecture

This project consists of 2 main parts: Worker and REST API 

Worker - service that periodically fetch latest data from **data providers** (like coinmarketcap), parse it to the common data structure, set the parsed data to the cache (Redis)

REST API - allows to get cached data through REST HTTP API

### Getting started

1. Spin up a Redis instance: `docker run -it -p 6379:6379 redis`
1. Start the app: `make start`
   1. Explore the API: [http://localhost:8423/swagger/index.html](http://localhost:8423/swagger/index.html)
   1. Use the API:
      * `curl -v "http://localhost:8421/v1/market/info?coin=60" | jq .`
      * `curl -v -X POST 'http://localhost:8421/v1/market/ticker' -H 'Content-Type: application/json' -d '{"currency":"ETH","assets":[{"type":"coin","coin":60}]}'`
      * `curl -v "http://localhost:8421/v1/market/charts?coin=60&time_start=1574483028" | jq .`
1. When done: `make stop`

Run `make` to see a list of all available build directives.

### More details

Current data providers: Coinmarketcap, BinanceDex, Compound, Fixer, Coingecko




