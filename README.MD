# GETHELP
This is a simple project (Built with Kratos) to demo some capabilities of Go Etherium and provide helpful insights.

**Note**: This project is named as _geth-help_, not Get Help.

## Tools Installation
1. Install [golang](https://go.dev/dl/)
2. Install [geth](https://geth.ethereum.org/downloads) cli
3. Install [Make](https://gnuwin32.sourceforge.net/packages/make.htm) cli
4. Run `make init`

## Initial Setup
1. Create a datadir folder for the eth node, e.g. `blockchain-data`, and `cd` into the folder
2. Run `geth --datadir ..\blockchain-data --keystore %LocalAppData%\Ethereum\Keystore init .\geth\genesis.json`
3. Create a new account using `geth account new`, then enter your password
4. Retrieve your account address using `geth account list`, example "929548598a3b93362c5aa2a24de190d18e657ae0"
5. Start local eth node using `geth --datadir ..\blockchain-data --keystore %LocalAppData%\Ethereum\Keystore --mine --miner.etherbase "0x929548598a3b93362c5aa2a24de190d18e657ae0" --miner.threads 3`
   - **NOTE**: Your node will start mining immediately, use `geth attach` and `miner.stop()` to stop mining while keeping the node alive. 

### Reproduce Steps
1. To reset the blockchain, remove the contents in datadir (`geth` folder).
   - You may also remove the contents in `%LocalAppData%\Ethereum\Keystore` for accounts removal.
2. Redo [Initial Setup](#initial-setup)

## Running the Application
1. Start by running `make run`

## Docker (TODO)
```bash
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

## TODO List
- [x] Setup local eth node
- [x] Perform transactions on the node
  - [x] Send eth
  - [ ] Validate recipient (Stretch goal)
- [ ] Withdraw money into cold wallet
  - [ ] Make a HD Wallet (Stretch goal)
- [x] Simple REST API
  - [x] GET Account balance (+Pending Balance)
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
Kratos | v2.0.0-20230515030202-6d741828c2d4
Protoc-gen-go | v1.30.0