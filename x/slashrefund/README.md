# Slash Refund Module

## Abstract

This paper describes the `slash-refund` module, how it works and its components.

The SDK module is a component strictly related to the `staking` module and has the purpose of refunding stakers from slashing events.

## How it works

Any user has the possibility to deposit a certain amount of allowed tokens into the module for a particular validator. This funds will be used to repay possible loss deriving from slashing events of the particular validator.

The funds will be maintained in the module and can be withdrawn at any moment. Following the delegated proof of stake philosophy and design pattern, before the tokens claim any withdrawn amount requires an unbonding period. This unbonding time will be equal to the staking unbonding time and is necessary to be able to account for slashing events evidenced in all the allowed time.

## TODO

Here is a list of functionalities or modification arising from the review of the code:

* [] modify the `DepositPool` structure to account for different tokens.
