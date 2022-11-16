#/bin/bash

clear

# TEST KEYS
VALKEY1="alice"    #VALIDATOR 1
DEPKEY1="bob"      #DEPOSITOR 1
DEPKEY2="carl"     #DEPOSITOR 2
# PROCESS KEYS
valaddr1=$(slash-refundd keys show $VALKEY1 -a --bech val)
depadrr1=$(slash-refundd keys show $DEPKEY1 -a)
depadrr2=$(slash-refundd keys show $DEPKEY2 -a)
DENOM="stake"

echo ; echo  "||PLAYERS||" ; echo ==========================================
echo validator1 $VALKEY1 $valaddr1
echo depositor1 $DEPKEY1 $depadrr1
echo depositor2 $DEPKEY2 $depadrr2
export valaddr1=$valaddr1
export depadrr1=$depadrr1
export depadrr2=$depadrr2


echo "balance of depositor 1:"
slash-refundd q bank balances $depadrr1    ; echo; echo

echo "tx deposit 100stake from depositor1:" 
slash-refundd tx slashrefund deposit  $valaddr1 100stake --from $DEPKEY1 -y | grep raw_log ; echo; echo

echo "balance of depositor 1:"
slash-refundd q bank balances $depadrr1    ; echo; echo

echo "list all deposits:"  
slash-refundd q  slashrefund list-deposit  ; echo; echo

echo "tx withdraw 100stake from depositor1:" 
slash-refundd tx slashrefund withdraw $valaddr1 100stake --from $DEPKEY1 -y | grep raw_log  ; echo; echo

echo "list unbonding deposits:"
slash-refundd q slashrefund list-unbonding-deposit ; echo; echo

sleep 7s

echo "balance of depositor 1:" 
slash-refundd q bank balances $depadrr1    ; echo; echo
