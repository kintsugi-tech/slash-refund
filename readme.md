# Slash Refund Module

## Abstract

This paper describes the `slash-refund` module, how it works and its components.

The SDK module is a component strictly related to the `staking` module and has the purpose of refunding stakers from slashing events.

## How it works

Any user has the possibility to deposit a certain amount of allowed tokens into the module for a particular validator. This funds will be used to repay possible loss deriving from slashing events of the particular validator.

The funds will be maintained in the module and can be withdrawn at any moment. Following the delegated proof of stake philosophy and design pattern, before the tokens claim any withdrawn amount requires an unbonding period. This unbonding time will be equal to the staking unbonding time and is necessary to be able to account for slashing events evidenced in all the allowed time.

## `types`

Protobuf defines the following types:

### `Deposit`

```go
type Deposit struct {
    DepositorAddress string                                 `protobuf:"bytes,1,opt,name=depositor_address,json=depositorAddress,proto3" json:"depositor_address,omitempty"`
    ValidatorAddress string                                 `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
    Shares           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares"`
}
```

where `DepositorAddress` is address of a user who deposited tokens, `ValidatorAddress` is the validator for which tokens are provided, and `Shares` are the pool's partial shares associated to the depositor.

### `DepositPool`

```go
type DepositPool struct {
    OperatorAddress string                                 `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
    Tokens          types.Coin                             `protobuf:"bytes,2,opt,name=tokens,proto3" json:"tokens"`
    Shares          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares" yaml:"depositor_shares"`
}
```

where `OperatorAddress` is the operator address of an active, or previously active, validator, `Tokens` are the token available to refund the associated validator slashes, and `Shares` are the actual circulating shares associated to this pool.

### `Refund`

```go
type Refund struct {
    DelegatorAddress string                                 `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
    ValidatorAddress string                                 `protobuf:"bytes,2,opt,name=validator_address,json=validatorAddress,proto3" json:"validator_address,omitempty"`
    Shares           github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares"`
}
```

This type is the analogous of `Deposit` but in this case `Shares` represents the portion of associated pool tokens the user can withdraw as result of a refund.

### `RefundPool`

```go
type RefundPool struct {
    OperatorAddress string                                 `protobuf:"bytes,1,opt,name=operator_address,json=operatorAddress,proto3" json:"operator_address,omitempty"`
    Tokens          types.Coin                             `protobuf:"bytes,2,opt,name=tokens,proto3" json:"tokens"`
    Shares          github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"shares" yaml:"refund_shares"`
}
```

This type is the anologous of `DepositPool` but is used to manage refunded tokens.

### `Params`

```go
type Params struct {
    AllowedTokens []string `protobuf:"bytes,1,rep,name=allowedTokens,proto3" json:"allowedTokens,omitempty" yaml:"allowed_tokens"`
}
```

where `allowedTokens` are the native tokens that a user can deposit to help refunding validators slahes.

## `msgServer`

The module's `msgServer` is composed of the following methods:

* `Deposit`: handles the deposit of funds from a generic address into the module. This method performs basic checks on the received message: depositor address, validator address, valid validtor and non zero deposit. After that, it calls the `Deposit` methos of the keeper.

* `Withdraw`: handles the withdraw of previously deposited funds if still available. This method checks the validaity of the depositor address, the one of the validator and the validaity of the requested withdraw. No check is performed on the bonding status of the validator since a deposit could have been made to a validator that has become unbonded. After that, it calls the `Withdraw` method of the keeper.

* `Claim`: ...

## `Keeper`

The slash refund module's staking keeper expects the following Cosmos SDK keeper:

* `types.BankKeeper`;

* `stakingKeeper`;

* `slashingKeeper`;

Here are described the provided methods:

* `Deposit`: this methods implements the state transition logic of funds added to a pool and shares emitted for the sender. The function creates a new deposit pool if there isn't already one for the target validator. Shares to be emitted are computed from deposited tokens. After computing shares the `types.Deposit` and `types.DepositPool` are updated accordingly. Given the deposit pool total tokens $poolTokens$, deposit pool total emitted shares $poolShares$, and the input $inputTokens$, the emitted shares $emittedShares$ are computed using the following equation:

$$ emittedShares = \frac{poolShares}{inputTokens} \cdot poolTokens $$

* `Withdraw`: this method deefines the state transition logic for a user that wants to withdraw its previously deposited tokens if still available. The user inputs the amount of shares it would like to withdraw. Given an amount of shares to be withdrawn, it insert the associated tokens into an unbonding deposit queue.

* `ComputeAssociatedShares`: this method is used to compute the amount of tokens associated to a specific amount of shares for a particular `types.DepositPool`. This method is almost entirely copied from the [staking module](https://github.com/cosmos/cosmos-sdk/blob/d74d0e4e8cd57d77de9590892ef89584765251c8/x/staking/keeper/delegation.go#L995).

## Important note

The module is based on a slightly modified version of the Cosmos SDK v0.46.8 in which we added the infraction height to slash events. This is imposed in the `go.mod` as a replacement:

```go
replace github.com/cosmos/cosmos-sdk v0.46.1 => github.com/made-in-block/cosmos-sdk v0.46.8-infraction-height
```
