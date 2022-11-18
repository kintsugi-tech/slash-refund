#!/bin/bash

echo "--------------------------------------------------------------------------------"
echo "                                    WARNING                                     "
echo "          check that slash-refundd home is set to \$HOME/.slash-refund          "
echo "--------------------------------------------------------------------------------"

# INITIALIZE
SR_ROOT=$HOME
MAINFOLDER=$SR_ROOT/.slash-refund
FLDVAL2=$SR_ROOT/.sr-node2
FLDVAL3=$SR_ROOT/.sr-node3

# CLEAN
killall screen
rm -rf $FLDVAL2
rm -rf $FLDVAL3


# CREATE VALIDATOR 2
#===================================================================
# Init and run node 2
./init-val-2.sh
echo "Done init val-2."
sleep 2
screen -d -m -S "val2" slash-refundd start --home $FLDVAL2
sleep 2

# Create validator tx
PUBKEY=$(slash-refundd tendermint show-validator --home $FLDVAL2)
slash-refundd tx staking create-validator -y --from carl --amount 10000000stake \
    --commission-max-change-rate 1 --commission-max-rate 1 --commission-rate 1 \
    --moniker "validator-2" --home $FLDVAL2 --pubkey ''"$PUBKEY"'' --min-self-delegation 1 \
    --broadcast-mode block \
    | grep raw_log

# Check node 2 is active
screen -ls
#===================================================================



VALKEY1="alice"
VALKEY2="carl"
valaddr1=$(slash-refundd keys show $VALKEY1 -a --bech val)
valaddr2=$(slash-refundd keys show $VALKEY2 -a --bech val)



# DELEGATE
#===================================================================
slash-refundd tx staking delegate $valaddr2 90000000stake --from bob --home $FLDVAL2 -y \
    --broadcast-mode block \
    | grep raw_log
#===================================================================

slash-refundd q bank balances $(slash-refundd keys show bob -a)

# DEPOSIT 
#===================================================================
slash-refundd tx slashrefund deposit $valaddr2 10000000stake --from alice -y \
    --broadcast-mode block \
    | grep raw_log
#===================================================================



# INIT AND RUN VAL-3 FOR DOUBLESIGN
#===================================================================
./init-val-3-doublesign.sh
sleep 2
screen -d -m -S "val3" slash-refundd start --home $FLDVAL3
sleep 1
#===================================================================



# REPEAT UNDELEGATE 
#===================================================================
# Undelegate are repeated to make at least one undelegate between
# slashing and evidence
for j in {0..6}
do
    slash-refundd tx staking unbond $valaddr2 100000stake --from bob -y \
        --broadcast-mode block \
        | grep raw_log
done
#===================================================================


VALKEY2="carl"
valaddr2=$(slash-refundd keys show $VALKEY2 -a --bech val)
slash-refundd q staking unbonding-delegations-from $valaddr2
slash-refundd q bank balances $(slash-refundd keys show bob -a)