#/bin/bash

# Init new chain
cd ..
MAINFOLDER=$(pwd)/.dimi-sr-test
VAL2FOL=$(pwd)/.dimi-sr-test-2
FOLDER=$(pwd)/.dimi-sr-test-3

rm -rf $FOLDER
slash-refundd init val-3 --chain-id slashrefund --home $FOLDER
cp $MAINFOLDER/config/genesis.json $FOLDER/config/genesis.json
cp $MAINFOLDER/config/client.toml $FOLDER/config/client.toml
cp -r $MAINFOLDER/keyring-test $FOLDER/keyring-test

# copy private key 🥸 let's do a nice double sign
cp $VAL2FOL/config/priv_validator_key.json $FOLDER/config/priv_validator_key.json

sed -i "" 's/0.0.0.0:9090/0.0.0.0:9290/' $FOLDER/config/app.toml
sed -i "" 's/0.0.0.0:9091/0.0.0.0:9291/' $FOLDER/config/app.toml

sed -i "" 's/127.0.0.1:26658/127.0.0.1:46658/' $FOLDER/config/config.toml
sed -i "" 's/127.0.0.1:26657/127.0.0.1:46657/' $FOLDER/config/config.toml
sed -i "" 's/0.0.0.0:26656/0.0.0.0:46656/' $FOLDER/config/config.toml

NODE_ID=$(slash-refundd tendermint show-node-id --home $VAL2FOL)
NODE_ID2=$(slash-refundd tendermint show-node-id --home $MAINFOLDER)

PEERS=$NODE_ID@127.0.0.1:36656,$NODE_ID2@127.0.0.1:26656
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $FOLDER/config/config.toml
sed -i -e "s/^allow_duplicate_ip *=.*/allow_duplicate_ip = true/" $FOLDER/config/config.toml

# Run node
screen -S "val3" -dm slash-refundd start --home $FOLDER

