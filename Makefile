run:
	go run main.go

geth-setup:
	geth init --datadir ..\geth-blockchain .\geth\genesis.json
