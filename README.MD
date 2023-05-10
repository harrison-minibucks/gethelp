# GETHELP
This is a simple project to demo some capabilities of Go Etherium and provide helpful insights.

*Note*: This project is named as _geth-help_.

## Tools Installation
1. Install [golang](https://go.dev/dl/)
1. Install [geth](https://geth.ethereum.org/downloads) cli

## Initial Setup
1. Create a datadir folder for the eth node, e.g. `blockchain-data`, and `cd` into the folder
1. Run `geth --datadir . --keystore %LocalAppData%\Ethereum\Keystore init ..\gethelp\geth\genesis.json`
1. Create a new account using `geth account new`, then enter your password
1. Retrieve your account address using `geth account list`, example "929548598a3b93362c5aa2a24de190d18e657ae0"
1. Start local eth node using `geth --datadir . --keystore %LocalAppData%\Ethereum\Keystore --mine --miner.etherbase "0x929548598a3b93362c5aa2a24de190d18e657ae0" --miner.threads 3`
   - **NOTE**: Your node will start mining immediately, use `geth attach` and `miner.stop()` to stop mining while keeping the node alive. 

### Reproduce Steps
1. To reset the blockchain, remove the contents in datadir (`geth` folder).
   - You may also remove the contents in `%LocalAppData%\Ethereum\Keystore` for accounts removal.
2. Redo [Initial Setup](#initial-setup)

## Running the Application
1. Start by running `go run main.go`

## TODO List
- [x] Setup local eth node
- [x] Perform transactions on the node
  - [x] Send eth
  - [ ] Validate recipient (Stretch goal)
- [ ] Validate transactions
- [ ] Withdraw money into cold wallet
  - [ ] Make a HD Wallet (Stretch goal)
- [ ] Simple REST API
  - [ ] GET Account balance (+Pending Balance)
    - [x] Format output
  - [ ] GET Gas price
  - [ ] GET Transaction verification
  - [ ] Implement challenge for client (Stretch goal)
- [ ] Create & Deploy Smart Contract (Stretch goal)

### Tools Versions
Name | Version
-----|--------
Go | 1.20.4
Geth | 1.11.6-stable-ea9e62ca