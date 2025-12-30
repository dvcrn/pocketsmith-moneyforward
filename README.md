# Pocketsmith MoneyForward Sync

A tool to sync transactions from MoneyForward (Japan) to Pocketsmith.

## Features

- Automatically syncs transactions from MoneyForward to Pocketsmith
- Normalizes payees for consistency
- Avoids duplicate transactions

## Setup

### Required Environment Variables


MONEYFORWARD_COOKIE=your_cookie
POCKETSMITH_TOKEN=your_token


### Command Line Flags

Alternatively, you can provide credentials via command line flags:


./pocketsmith-moneyforward -mf-cookie=xxx -pocketsmith-token=xxx

### Run with docker (recommended)

`docker run -e MONEYFORWARD_COOKIE=xxx -e POCKETSMITH_TOKEN=xxx ghcr.io/dvcrn/pocketsmith-moneyforward:latest`

## Building


go build ./...


## Running


./pocketsmith-moneyforward
