#/bin/bash

# Init new chain
MAINFOLDER=$(pwd)/.dimi-sr-test
FOLDER=$(pwd)/.dimi-sr-test-2

rm -rf $FOLDER
slash-refundd init val-2 --chain-id slashrefund --home $FOLDER
cp $MAINFOLDER/config/genesis.json $FOLDER/config/genesis.json
cp $MAINFOLDER/config/client.toml $FOLDER/config/client.toml
cp -r $MAINFOLDER/keyring-test $FOLDER/keyring-test

sed -i "" 's/0.0.0.0:9090/0.0.0.0:9190/' $FOLDER/config/app.toml
sed -i "" 's/0.0.0.0:9091/0.0.0.0:9191/' $FOLDER/config/app.toml

sed -i "" 's/127.0.0.1:26658/127.0.0.1:36658/' $FOLDER/config/config.toml
sed -i "" 's/127.0.0.1:26657/127.0.0.1:36657/' $FOLDER/config/config.toml
sed -i "" 's/0.0.0.0:26656/0.0.0.0:36656/' $FOLDER/config/config.toml

NODE_ID=$(slash-refundd tendermint show-node-id --home $MAINFOLDER)

PEERS=$NODE_ID@127.0.0.1:26656
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $FOLDER/config/config.toml
sed -i -e "s/^allow_duplicate_ip *=.*/allow_duplicate_ip = true/" $FOLDER/config/config.toml

# Run node
screen -S "val2" -dm slash-refundd start --home $FOLDER

# Create validator tx
PUBKEY=$(slash-refundd tendermint show-validator --home $FOLDER)
slash-refundd tx staking create-validator -y --from pippo --amount 10000000stake --commission-max-change-rate 1 --commission-max-rate 1 --commission-rate 1 --moniker "dimi2" --home $FOLDER --pubkey ''"$PUBKEY"'' --min-self-delegation 1
