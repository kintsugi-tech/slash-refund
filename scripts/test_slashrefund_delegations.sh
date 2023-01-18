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



VALKEY1="alice"  # genesis validator from config.yaml, won't be slashed
VALKEY2="carl"   # new generated validator, will be slashed
valaddr1=$(slash-refundd keys show $VALKEY1 -a --bech val)
valaddr2=$(slash-refundd keys show $VALKEY2 -a --bech val)



# DELEGATE
#===================================================================
slash-refundd tx staking delegate $valaddr2 90000000stake --from bob --home $FLDVAL2 -y \
    --broadcast-mode block \
    | grep raw_log
echo "-------- All delegations for val2:    "
slash-refundd q staking delegations-to $valaddr2
#===================================================================



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



# CONTINUE QUERYING
#===================================================================
for j in {0..19}
do
    echo ========
    echo BLOCK:
    echo $(expr 1 + $(slash-refundd q block | jq '.block.header.height | tonumber'))
    echo ========
    # wrong tx just to sync to wait for block to be committed
    slash-refundd tx staking unbond $valaddr2 1token --from carl -y \
        --broadcast-mode block \
        | grep raw_log
    echo -----------------------------------------------------
done
#===================================================================



# CLAIM 
#===================================================================

# BOB: with slashFactor=0.05 must receive 4'500'000
echo "-------- Balance of bob:" 
slash-refundd q bank balances $(slash-refundd keys show bob -a)
echo "-------- bob claim:"
slash-refundd tx slashrefund claim $valaddr2 --from bob \
    --broadcast-mode block -y | grep raw_log 
echo "-------- Balance of bob:"
slash-refundd q bank balances $(slash-refundd keys show bob -a)

# CARL: with slashFactor=0.05 must receive 500'000
echo "-------- Balance of carl:" 
slash-refundd q bank balances $(slash-refundd keys show carl -a)
echo "-------- carl claim:"
slash-refundd tx slashrefund claim $valaddr2 --from carl \
    --broadcast-mode block -y | grep raw_log 
echo "-------- Balance of bob:"
slash-refundd q bank balances $(slash-refundd keys show carl -a)

# ALICE: must not be present a refund for alice.
echo "-------- Claim (wrong: alice has no deposits):"
slash-refundd tx slashrefund claim $valaddr2 --from alice \
    --broadcast-mode block -y | grep raw_log 

# BOB: must not be present a refund for bob (already claimed)
echo "-------- carl claim:"
slash-refundd tx slashrefund claim $valaddr2 --from carl \
    --broadcast-mode block -y | grep raw_log 

# CARL: must not be present a refund for carl (already claimed)
echo "-------- carl claim:"
slash-refundd tx slashrefund claim $valaddr2 --from carl \
    --broadcast-mode block -y | grep raw_log 
#===================================================================