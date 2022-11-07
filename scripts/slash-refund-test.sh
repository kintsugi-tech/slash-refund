#/bin/bash

screen -X -S val2 quit
screen -X -S val3 quit 
SR_ROOT=$HOME
MAINFOLDER=$SR_ROOT/.slash-refund
FLDVAL2=$(pwd)/../.sr-node2
FLDVAL3=$(pwd)/../.sr-node3

#===================================================================
#create val2 with 10M
#delegate val2 from bob with 90M
#deposit for val2 from alice 10M
#deposit for val2 from alice 10M
sh run-val-2.sh 
# RUN VALIDATOR 2
echo RUN NODE VAL-2: ------------------------------------------
screen -dmS "val2" slash-refundd start --home $FLDVAL2
echo DONE --------------------------------------------------
# Create validator tx
PUBKEY=$(slash-refundd tendermint show-validator --home $FLDVAL2)

echo CREATE VAL-2 WITH 10M: ----------------------------------------
slash-refundd tx staking create-validator -y --from pippo --amount 10000000stake --commission-max-change-rate 1 --commission-max-rate 1 --commission-rate 1 --moniker "validator-2" --home $FLDVAL2 --pubkey ''"$PUBKEY"'' --min-self-delegation 1 \
    | grep raw_log
echo DONE --------------------------------------------------

echo DELEGATE 90Mstake FOR VAL-2 FROM BOB: ----------------------------------------
slash-refundd tx staking delegate cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 90000000stake --from bob --home $FLDVAL2 -y \
    | grep raw_log
echo DONE --------------------------------------------------

echo DEPOSIT 10Mstake FOR VAL-2 FROM ALICE: ----------------------------------------
slash-refundd tx slashrefund deposit cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 10000000stake --from alice --home $FLDVAL2 -y \
    | grep raw_log
echo DONE --------------------------------------------------

echo DEPOSIT 10Mstake FOR VAL-2 FROM BOB: ----------------------------------------
slash-refundd tx slashrefund deposit cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 10000000stake --from alice --home $FLDVAL2 -y \
    | grep raw_log
echo DONE --------------------------------------------------
# val2 100M
# pool  20M
echo VALIDATORS: ----------------------
slash-refundd q staking validators
echo; 
echo DEPOSIT POOL: ----------------------
slash-refundd q slashrefund list-deposit
echo;
#===================================================================

echo sleeping 5 ...
sleep 5
echo done.

#===================================================================
echo WITHDRAW 1Mstake FOR VAL-2 FROM ALICE: ----------------------------------------
slash-refundd tx slashrefund withdraw cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 1000000stake --from alice --home $FLDVAL2 -y \
    | grep raw_log
echo DONE --------------------------------------------------
#===================================================================

# RUN VALIDATOR 3
#===================================================================
sh run-val-3-doublesign.sh
echo RUN NODE VAL-3: ------------------------------------------
screen -dmS "val3" slash-refundd start --home $FLDVAL3
echo DONE. ----------------------------------------------------
#===================================================================

#===================================================================
echo WITHDRAW 4Mstake FOR VAL-2 FROM ALICE: ----------------------------------------
slash-refundd tx slashrefund withdraw cosmosvaloper1f58m57pcn99r0wdktq07q29uld0uherjatc48z 4000000stake --from alice --home $FLDVAL2 -y \
    | grep raw_log
echo DONE --------------------------------------------------

echo DEPOSIT POOL: ----------------------
slash-refundd q  slashrefund list-deposit
echo; 
#===================================================================

#===================================================================
# wait for slashrefund and then query again:
echo sleeping 10 ...
sleep 10
echo done.
echo DEPOSIT POOL: ----------------------
slash-refundd q  slashrefund list-deposit
echo;
#===================================================================

screen -ls