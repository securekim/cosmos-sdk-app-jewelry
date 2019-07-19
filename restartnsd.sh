#/bin/bash
_PATH=`pwd`
#./nsd unsafe-reset-all
rm -rf ~/.ns*
cd /home/securekim/go/src/github.com/cosmos/cosmos-sdk-app-jewelry
go build cmd/nscli/main.go
mv main nscli
go build cmd/nsd/main.go
mv main nsd
./nsd init securekim --chain-id diachain
./nscli keys add jack
./nsd add-genesis-account $(nscli keys show jack -a) 1000nametoken,100000000stake
./nscli config chain-id diachain
./nscli config output json
./nscli config indent true
./nscli config trust-node true
./nsd gentx --name jack
./nsd collect-gentxs
./nsd validate-genesis
./nsd start
