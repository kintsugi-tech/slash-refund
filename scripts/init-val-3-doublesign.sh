#/bin/bash

## "--------------------------------------------------------------------------------"
## "                                    WARNING                                     "
## "          check that slash-refundd home is set to \$HOME/.slash-refund          "
## "--------------------------------------------------------------------------------"

SR_ROOT=$HOME
MAINFOLDER=$SR_ROOT/.slash-refund
FLDVAL2=$SR_ROOT/.sr-node2
FLDVAL3=$SR_ROOT/.sr-node3

slash-refundd init val-3 --chain-id slashrefund --home $FLDVAL3

cp $MAINFOLDER/config/genesis.json $FLDVAL3/config/genesis.json
cp $MAINFOLDER/config/client.toml $FLDVAL3/config/client.toml
cp -r $MAINFOLDER/keyring-test $FLDVAL3/keyring-test

# copy private key 🥸 let's do a nice double sign
cp $FLDVAL2/config/priv_validator_key.json $FLDVAL3/config/priv_validator_key.json

sed -i 's/0.0.0.0:9090/0.0.0.0:9290/' $FLDVAL3/config/app.toml
sed -i 's/0.0.0.0:9091/0.0.0.0:9291/' $FLDVAL3/config/app.toml

sed -i 's/127.0.0.1:26658/127.0.0.1:46658/' $FLDVAL3/config/config.toml
sed -i 's/127.0.0.1:26657/127.0.0.1:46657/' $FLDVAL3/config/config.toml
sed -i 's/0.0.0.0:26656/0.0.0.0:46656/' $FLDVAL3/config/config.toml

NODE_ID=$(slash-refundd tendermint show-node-id --home $FLDVAL2)
NODE_ID2=$(slash-refundd tendermint show-node-id --home $MAINFOLDER)

PEERS=$NODE_ID@127.0.0.1:36656,$NODE_ID2@127.0.0.1:26656
sed -i.bak -e "s/^persistent_peers *=.*/persistent_peers = \"$PEERS\"/" $FLDVAL3/config/config.toml
sed -i -e "s/^allow_duplicate_ip *=.*/allow_duplicate_ip = true/" $FLDVAL3/config/config.toml

