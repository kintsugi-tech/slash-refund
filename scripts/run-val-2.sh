#/bin/bash

# Init new chain
cd ..
echo $(pwd)
#MAINFOLDER=$(pwd)/.dimi-sr-test
MAINFOLDER="~/.slash-refund"
FOLDER=$(pwd)/.dimi-sr-test-2

rm -rf $FOLDER
echo INIT NODE VAL-2 IN $FOLDER: --------------
slash-refundd init val-2 --chain-id slashrefund --home $FOLDER
echo DONE. ----------------------------------------------------
cp $MAINFOLDER/config/genesis.json $FOLDER/config/genesis.json
cp $MAINFOLDER/config/client.toml $FOLDER/config/client.toml
cp -r $MAINFOLDER/keyring-test $FOLDER/keyring-test

sed -i 's/0.0.0.0:9090/0.0.0.0:9190/' $FOLDER/config/app.toml
sed -i 's/0.0.0.0:9091/0.0.0.0:9191/' $FOLDER/config/app.toml

sed -i 's/127.0.0.1:26658/127.0.0.1:36658/' $FOLDER/config/config.toml
sed -i 's/127.0.0.1:26657/127.0.0.1:36657/' $FOLDER/config/config.toml
sed -i 's/0.0.0.0:26656/0.0.0.0:36656/' $FOLDER/config/config.toml

NODE_ID=$(slash-refundd tendermint show-node-id --home $MAINFOLDER)

PEERS=$NODE_ID@127.0.0.1:26656
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $FOLDER/config/config.toml
sed -i -e "s/^allow_duplicate_ip *=.*/allow_duplicate_ip = true/" $FOLDER/config/config.toml

# Run node
screen -S "val2" -dm slash-refundd start --home $FOLDER

# Create validator tx
PUBKEY=$(slash-refundd tendermint show-validator --home $FOLDER)
echo CREATE VAL-2 : ----------------------------------------
slash-refundd tx staking create-validator -y --from pippo --amount 10000000stake --commission-max-change-rate 1 --commission-max-rate 1 --commission-rate 1 --moniker "dimi2" --home $FOLDER --pubkey ''"$PUBKEY"'' --min-self-delegation 1
echo DONE --------------------------------------------------
slash-refundd tx staking delegate cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 100000000stake --from bob --home $FOLDER -y
slash-refundd tx slashrefund deposit cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 69000000stake --from alice --home $FOLDER -y