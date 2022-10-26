#/bin/bash

# TEST KEYS
VALKEY1="alice"    #VALIDATOR 1
#not set#VALKEY2="bob"      #VALIDATOR 2
DEPKEY1="bob"      #DEPOSITOR 1
DEPKEY2="pippo"    #DEPOSITOR 2
# PROCESS KEYS
valaddr1=$(slash-refundd keys show $VALKEY1 -a --bech val)
#not set#valaddr2=$(slash-refundd keys show $VALKEY2 -a --bech val)
depadrr1=$(slash-refundd keys show $DEPKEY1 -a)
depadrr2=$(slash-refundd keys show $DEPKEY2 -a)
DENOM="stake"

echo ; echo  "||PLAYERS||" ; echo ==========================================
echo validator1 $VALKEY1 $valaddr1
#not set#echo validator2 $VALKEY2 $valaddr2
echo depositor1 $DEPKEY1 $depadrr1
echo depositor2 $DEPKEY2 $depadrr2
export valaddr1=$valaddr1
#not set#export valaddr2=$valaddr2
export depadrr1=$depadrr1
export depadrr2=$depadrr2



#--------------#
#    RECAP     #
#--------------#
## QUERIES
# slash-refundd q slashrefund list-deposit
# slash-refundd q slashrefund list-deposit-pool
# slash-refundd q slashrefund list-unbonding-deposit
# slash-refundd q slashrefund show-deposit $depadrr1 $valaddr1
# slash-refundd q slashrefund show-deposit-pool $valaddr1
## TRANSACTIONS
# slash-refundd tx slashrefund deposit  $valaddr1 $AMTstake --from $DEPKEY1 -y | grep raw_log
# slash-refundd tx slashrefund withdraw $valaddr1 $AMTstake --from $DEPKEY1 -y | grep raw_log



#--------------#
#    TEST1     #
#--------------#

echo ; echo  "[[ TEST WITHDRAW 1 ]]" ; echo ========================================== ; echo


echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


DEPamt1=100
echo  "deposit 1: "$DEPamt1$DENOM" from dep1" ; echo ------------------------------------------
slash-refundd tx slashrefund deposit  $valaddr1 $DEPamt1$DENOM --from $DEPKEY1 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


echo  "TEST: withdraw from depositor 2" ; echo ------------------------------------------ 
slash-refundd tx slashrefund withdraw $valaddr1 1$DENOM --from $DEPKEY2 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


DEPamt2=200
echo  "deposit 2: "$DEPamt2$DENOM" from dep1" ; echo ------------------------------------------
slash-refundd tx slashrefund deposit  $valaddr1 $DEPamt2$DENOM --from $DEPKEY2 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


echo  "withdraw 1: half deposit from dep1" ; echo ------------------------------------------
WITamt1=$(expr $DEPamt1 / 2)
slash-refundd tx slashrefund withdraw $valaddr1 $WITamt1$DENOM --from $DEPKEY1 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


sleep .5s


echo  "withdraw 2: half deposit from dep2" ; echo ------------------------------------------ 
WITamt2=$(expr $DEPamt2 / 2)
slash-refundd tx slashrefund withdraw $valaddr1 $WITamt2$DENOM --from $DEPKEY2 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


sleep .5s


echo  "withdraw 3: 1/4 deposit from dep1" ; echo ------------------------------------------ 
WITamt3=$(expr $DEPamt1 / 4)
slash-refundd tx slashrefund withdraw $valaddr1 $WITamt3$DENOM --from $DEPKEY1 -y | grep raw_log
echo

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo


echo  "sleep 5 seconds . . .    " ; sleep 5s ; echo "done." ; echo # pause 5 (unbonding_time set in config.yml)

echo  "list balances and deposits and deposit pool and unbonding-deposits" ; echo ------------------------------------------
echo "depositor1:" ; slash-refundd q bank balances $depadrr1 ; echo "depositor2:" ; slash-refundd q bank balances $depadrr2 
echo "list all deposits:" ; slash-refundd q  slashrefund list-deposit
echo "deposit-pool:" ; slash-refundd q slashrefund show-deposit-pool $valaddr1
echo "list unbonding deposits:" ; slash-refundd q slashrefund list-unbonding-deposit
echo

#TODO: try withdraw = all
#TODO: try withdraw > all : ErrInvalidRequest, "invalid token amount"