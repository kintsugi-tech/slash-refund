package slashrefund_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/x/slashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/testslashrefund"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func generateRandomAddress() (addres string) {
	hexaddr := secp256k1.GenPrivKey().PubKey().Address()
	address := sdk.AccAddress(hexaddr).String()
	return address
}

func generateRandomOperator() (operator, address string) {
	hexaddr := secp256k1.GenPrivKey().PubKey().Address()
	address = sdk.AccAddress(hexaddr).String()
	operator = sdk.ValAddress(hexaddr).String()
	return operator, address
}

func TestGenesis(t *testing.T) {

	operator1, depositor1 := generateRandomOperator()
	operator2, depositor2 := generateRandomOperator()
	address1 := generateRandomAddress()
	address2 := generateRandomAddress()

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		DepositList: []types.Deposit{
			{
				DepositorAddress: depositor1,
				ValidatorAddress: operator1,
				Shares:           sdk.NewDec(100),
			},
			{
				DepositorAddress: depositor2,
				ValidatorAddress: operator1,
				Shares:           sdk.NewDec(200),
			},
			{
				DepositorAddress: depositor1,
				ValidatorAddress: operator2,
				Shares:           sdk.NewDec(300),
			},
			{
				DepositorAddress: depositor2,
				ValidatorAddress: operator2,
				Shares:           sdk.NewDec(400),
			},
		},
		DepositPoolList: []types.DepositPool{
			{
				OperatorAddress: operator1,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
			{
				OperatorAddress: operator2,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
		},
		UnbondingDepositList: []types.UnbondingDeposit{
			{
				DepositorAddress: depositor1,
				ValidatorAddress: operator1,
				Entries: []types.UnbondingDepositEntry{
					{
						CreationHeight: 0,
						CompletionTime: time.Now().UTC(),
						InitialBalance: sdk.NewInt(100),
						Balance:        sdk.NewInt(100),
					},
					{
						CreationHeight: 0,
						CompletionTime: time.Now().UTC().Add(100),
						InitialBalance: sdk.NewInt(200),
						Balance:        sdk.NewInt(200),
					},
				},
			},
			{
				DepositorAddress: depositor2,
				ValidatorAddress: operator2,
				Entries: []types.UnbondingDepositEntry{
					{
						CreationHeight: 0,
						CompletionTime: time.Now().UTC(),
						InitialBalance: sdk.NewInt(100),
						Balance:        sdk.NewInt(100),
					},
					{
						CreationHeight: 0,
						CompletionTime: time.Now().UTC().Add(100),
						InitialBalance: sdk.NewInt(200),
						Balance:        sdk.NewInt(200),
					},
				},
			},
		},
		RefundPoolList: []types.RefundPool{
			{
				OperatorAddress: operator1,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
			{
				OperatorAddress: operator2,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
		},
		RefundList: []types.Refund{
			{
				DelegatorAddress: address1,
				ValidatorAddress: operator1,
				Shares:           sdk.NewDec(100),
			},
			{
				DelegatorAddress: address2,
				ValidatorAddress: operator2,
				Shares:           sdk.NewDec(100),
			},
		},
	}

	k, ctx := testslashrefund.NewTestKeeper(t)
	slashrefund.InitGenesis(ctx, *k, genesisState)
	got := slashrefund.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	require.Equal(t, genesisState.Params, got.Params)
	require.ElementsMatch(t, genesisState.DepositList, got.DepositList)
	require.ElementsMatch(t, genesisState.DepositPoolList, got.DepositPoolList)
	require.ElementsMatch(t, genesisState.UnbondingDepositList, got.UnbondingDepositList)
	require.ElementsMatch(t, genesisState.RefundPoolList, got.RefundPoolList)
	require.ElementsMatch(t, genesisState.RefundList, got.RefundList)
	// this line is used by starport scaffolding # genesis/test/assert
}
