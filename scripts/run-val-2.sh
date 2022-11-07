#/bin/bash

# Init new chain
echo "--------------------------------------------------------------------------------"
echo "                                    WARNING                                     "
echo "       check that slash-refundd home must be set to \$HOME/.slash-refund        "
echo "--------------------------------------------------------------------------------"

SR_ROOT=$HOME
MAINFOLDER=$SR_ROOT/.slash-refund
FLDVAL2=$(pwd)/../.sr-node2
FLDVAL3=$(pwd)/../.sr-node3
rm -rf $FLDVAL2
rm -rf $FLDVAL3

echo INIT NODE VAL-2 IN $FLDVAL2: --------------
slash-refundd init val-2 --chain-id slashrefund --home $FLDVAL2
echo DONE. ----------------------------------------------------

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

# Run node
####echo RUN NODE VAL-2: ------------------------------------------
####screen -dmS "val2" slash-refundd start --home $FLDVAL2
####echo DONE --------------------------------------------------
##### Create validator tx
####PUBKEY=$(slash-refundd tendermint show-validator --home $FLDVAL2)
####
####echo CREATE VAL-2 WITH 10M: ----------------------------------------
####slash-refundd tx staking create-validator -y --from pippo --amount 10000000stake --commission-max-change-rate 1 --commission-max-rate 1 --commission-rate 1 --moniker "validator-2" --home $FLDVAL2 --pubkey ''"$PUBKEY"'' --min-self-delegation 1 \
####    | grep raw_log
####echo DONE --------------------------------------------------
####
####echo DELEGATE 90Mstake FOR VAL-2 FROM BOB: ----------------------------------------
####slash-refundd tx staking delegate cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 90000000stake --from bob --home $FLDVAL2 -y \
####    | grep raw_log
####echo DONE --------------------------------------------------
####
####echo DEPOSIT 10Mstake FOR VAL-2 FROM ALICE: ----------------------------------------
####slash-refundd tx slashrefund deposit cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 10000000stake --from alice --home $FLDVAL2 -y \
####    | grep raw_log
####echo DONE --------------------------------------------------
####
####echo DEPOSIT 10Mstake FOR VAL-2 FROM BOB: ----------------------------------------
####slash-refundd tx slashrefund deposit cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 10000000stake --from alice --home $FLDVAL2 -y \
####    | grep raw_log
####echo DONE --------------------------------------------------