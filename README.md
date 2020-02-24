# Watchmarket Blockatlas version
Is under development at master branch. Was migrated from [Blockatlas](https://github.com/trustwallet/blockatlas) at version/blockatlas branch

[![Build Status](https://dev.azure.com/TrustWallet/WatchMarket/_apis/build/status/trustwallet.watchmarket?branchName=master)](https://dev.azure.com/TrustWallet/WatchMarket/_build/latest?definitionId=45&branchName=master)

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

### More details

Current data providers: Coinmarketcap, BinanceDex, Compound, Fixer, Coingecko




