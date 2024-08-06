# Dex Price API
An API to retrieve prices directly from DEXs.

## What is it?
This project was initiated as an Open Source token price API. Its primary goal is to provide free access to token prices across various DEXs (Decentralized Exchanges). The API is designed to be particularly useful for companies whose products rely on token prices, such as wallets and other Web3 applications. All prices accessible through this API are publicly available and stored within DEX contracts on public blockchains. 

## Supported protocols
Now are supported only 2 protocols:
- [STonFi](https://ston.fi/)
- [UniSwap](https://uniswap.org)

## Supported chains
- Ethereum
- TON

## Run docker
To run docker compose just:

`docker compose up --build -d`

Before using docker compose be sure that you have your ip correct in your config file. Also check ports and config path.

## Config file
Config file now contains only 2 params:
- ip
- port


## Examples
On request:
```sh
curl --location 'http://localhost:15001/ton' \
--header 'Content-Type: application/json' \
--data '{
    "method": "TONDex.GetPoolPrice",
    "params": [
        {
            "pool":"EQD8TJ8xEWB1SpnRE4d89YO3jl0W0EiBnNS4IBaHaUmdfizE"
        }
    ],
    "id": "1"
}'
```
You will se the response:
```json
{
    "result": {
        "token0": "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe",
        "token1": "0:8cdc1d7640ad5ee326527fc1ad0514f468b30dc84b0173f0e155f451b4e11f7c",
        "price_of_token0": "149599453",
        "price_of_token1": "6684516415",
        "decimals": "9"
    },
    "error": null,
    "id": "1"
}
```
And the second type of request:
```sh
curl --location 'http://localhost:15001/ton' \
--header 'Content-Type: application/json' \
--data '{
    "method": "TONDex.GetTokensPrices",
    "params": [
        {
           "token0" : "EQAvlWFDxGF2lXm67y4yzC17wYKD9A0guwPkMs1gOsM__NOT",
           "token1" : "EQCxE6mUtQJKFnGfaROTKOt1lZbDiiX1kCixRv7Nw2Id_sDs"
        }
    ],
    "id": "1"
}'
```
Response:
```json
{
    "result": {
        "token0": "0:2f956143c461769579baef2e32cc2d7bc18283f40d20bb03e432cd603ac33ffc",
        "token1": "0:b113a994b5024a16719f69139328eb759596c38a25f59028b146fecdc3621dfe",
        "price_of_token0": "12590087",
        "price_of_token1": "79427567228",
        "decimals": "9"
    },
    "error": null,
    "id": "1"
}
```
Be aware that prices are changing all the time and this is only example. 