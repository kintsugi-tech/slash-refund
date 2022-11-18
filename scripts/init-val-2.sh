#/bin/bash

## "--------------------------------------------------------------------------------"
## "                                    WARNING                                     "
## "          check that slash-refundd home is set to \$HOME/.slash-refund          "
## "--------------------------------------------------------------------------------"

SR_ROOT=$HOME
MAINFOLDER=$SR_ROOT/.slash-refund
FLDVAL2=$SR_ROOT/.sr-node2

slash-refundd init val-2 --chain-id slashrefund --home $FLDVAL2

cp $MAINFOLDER/config/genesis.json $FLDVAL2/config/genesis.json
cp $MAINFOLDER/config/client.toml $FLDVAL2/config/client.toml
cp -r $MAINFOLDER/keyring-test $FLDVAL2/keyring-test

sed -i 's/0.0.0.0:9090/0.0.0.0:9190/' $FLDVAL2/config/app.toml
sed -i 's/0.0.0.0:9091/0.0.0.0:9191/' $FLDVAL2/config/app.toml

sed -i 's/127.0.0.1:26658/127.0.0.1:36658/' $FLDVAL2/config/config.toml
sed -i 's/127.0.0.1:26657/127.0.0.1:36657/' $FLDVAL2/config/config.toml
sed -i 's/0.0.0.0:26656/0.0.0.0:36656/' $FLDVAL2/config/config.toml

NODE_ID=$(slash-refundd tendermint show-node-id --home $MAINFOLDER)

PEERS=$NODE_ID@127.0.0.1:26656
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $FLDVAL2/config/config.toml
sed -i -e "s/^allow_duplicate_ip *=.*/allow_duplicate_ip = true/" $FLDVAL2/config/config.toml