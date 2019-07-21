#/bin/bash
_PATH=`pwd`
#./nsd unsafe-reset-all
killall nsd
killall nscli
rm -rf ~/.ns*
rm bin/*
cd /home/securekim/go/src/github.com/cosmos/cosmos-sdk-app-jewelry
go build cmd/nscli/main.go
mv main bin/nscli
go build cmd/nsd/main.go
mv main bin/nsd
cd bin
./nsd unsafe-reset-all
./nsd init securekim --chain-id diachain
./nscli keys add gemologist
./nscli keys add wholesaler
./nscli keys add retailer
./nsd add-genesis-account $($_PATH/bin/nscli keys show gemologist -a) 100nametoken,100000000stake
./nsd add-genesis-account $($_PATH/bin/nscli keys show wholesaler -a) 1500nametoken,100000000stake
./nsd add-genesis-account $($_PATH/bin/nscli keys show retailer -a) 1000nametoken,100000000stake
./nscli config chain-id diachain
./nscli config output json
./nscli config indent true
./nscli config trust-node true
./nsd gentx --name gemologist
./nsd collect-gentxs
./nsd validate-genesis
./nsd start
