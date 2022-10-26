#/bin/bash
# TEST KEYS
VALKEY1="alice"    #VALIDATOR 1
##VALKEY2="bob"      #VALIDATOR 2
DEPKEY1="bob"      #DEPOSITOR 1
DEPKEY2="pippo"    #DEPOSITOR 2
# PROCESS KEYS
valaddr1=$(slash-refundd keys show $VALKEY1 -a --bech val)
##valaddr2=$(slash-refundd keys show $VALKEY2 -a --bech val)
depadrr1=$(slash-refundd keys show $DEPKEY1 -a)
depadrr2=$(slash-refundd keys show $DEPKEY2 -a)
echo ; echo  "||PLAYERS||" ; echo ================================
echo validator1 $VALKEY1 $valaddr1
##echo validator2 $VALKEY2 $valaddr2
echo depositor1 $DEPKEY1 $depadrr1
echo depositor2 $DEPKEY2 $depadrr2
export valaddr1=$valaddr1
##export valaddr2=$valaddr2
export depadrr1=$depadrr1
export depadrr2=$depadrr2
#
#
# SYNTAXES:
#  slash-refundd tx slashrefund deposit [validator-address] [amount] [flags --from [depositor-address]]
#  slash-refundd query slashrefund show-deposit [address] [validator-address] [flags]
#
#
# WRONG TRANSACTION
echo ; echo  "||TEST WRONG DENOM||" ; echo ================================
echo ; echo TEST WRONG DENOM: TX deposit banane ; echo --------
slash-refundd tx slashrefund deposit $valaddr1 2000000banane --from $DEPKEY1 -y | grep raw_log
#
# CORRECT TRANSACTION: 2M from depositor1 to validator1
echo ; echo  "||TEST DEPOSIT||" ; echo ================================
echo ; echo  TEST CORRECT DENOM: TX deposit 2Mstake from depositor1 to validator1 ; echo --------
slash-refundd tx slashrefund deposit $valaddr1 2000000stake --from $DEPKEY1 -y | grep raw_log
#
# CORRECT TRANSACTION: 9M from depositor1 to validator2
##echo ; echo  "||TEST DEPOSIT||" ; echo ================================
##echo ; echo  TEST CORRECT DENOM: TX deposit 9Mstake from depositor1 to validator2 ; echo --------
##slash-refundd tx slashrefund deposit $valaddr2 9000000stake --from $DEPKEY1 -y | grep raw_log
#
#
# TEST QUERY
echo ; echo  "||TEST DEPOSIT QUERIES||" ; echo ================================
echo ; echo  "TEST: QUERY list-deposit" ; echo --------
slash-refundd q  slashrefund list-deposit
echo ; echo  "TEST: QUERY list-deposit-pool" ; echo --------
slash-refundd q  slashrefund list-deposit-pool
echo ; echo  "TEST: QUERY list-unbonding-deposit" ; echo --------
slash-refundd q  slashrefund list-unbonding-deposit
echo ; echo  "TEST: QUERY show-deposit depadrr1 valaddr1" ; echo --------
slash-refundd q  slashrefund show-deposit $depadrr1 $valaddr1
echo ; echo  "TEST: QUERY show-deposit-pool valaddr1" ; echo --------
slash-refundd q  slashrefund show-deposit-pool $valaddr1
##echo ; echo  "TEST: QUERY show-deposit-pool valaddr1" ; echo --------
##slash-refundd q  slashrefund show-deposit-pool $valaddr2
#
#
# TEST: MORE DEPOSITS
# CORRECT TRANSACTION: 5M from depositor2
echo ; echo  "||TEST CONSECUTIVE DEPOSITS||" ; echo ================================
echo ; echo  "TEST: TX deposit 5Mstake from depositor2 to validator1" ; echo --------
slash-refundd tx slashrefund deposit $valaddr1 5000000stake --from $DEPKEY2 -y | grep raw_log
echo ; echo  "TEST: TX deposit 1Mstake from depositor1 to validator1" ; echo --------
slash-refundd tx slashrefund deposit $valaddr1 1000000stake --from $DEPKEY1 -y | grep raw_log
echo ; echo  "TEST: QUERY list-deposit" ; echo -------- 
slash-refundd q  slashrefund list-deposit
echo ; echo  "TEST: QUERY list-deposit-pool" ; echo --------
slash-refundd q  slashrefund list-deposit-pool
#
#
# RECAP
echo ; echo  "||PLAYERS||" ; echo ================================
echo validator1 $VALKEY1 $valaddr1
echo validator2 $VALKEY1 $valaddr2
echo depositor1 $DEPKEY1 $depadrr1
echo depositor2 $DEPKEY2 $depadrr2
