package types_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {

	address0 := sample.AccAddress()
	valAddress0 := sdk.ValAddress(sample.AccAddress()).String()
	address1 := sample.AccAddress()
	valAddress1 := sdk.ValAddress(sample.AccAddress()).String()

	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default genesis state",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc:     "custom valid genesis state",
			genState: validGenesisState(address0, address1, valAddress0, valAddress1),
			valid:    true,
		},
		{
			desc: "no allowed token set",
			genState: &types.GenesisState{
				Params: types.Params{
					AllowedTokens: nil,
				},
			},
			valid: false,
		},
		{
			desc: "duplicated deposit",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositList: []types.Deposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(100),
					},
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(100),
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated unbondingDeposit",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(10),
								Balance:        sdk.NewInt(00),
							},
						},
					},
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 20,
								CompletionTime: time.Unix(20, 0),
								InitialBalance: sdk.NewInt(20),
								Balance:        sdk.NewInt(20),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated depositPool",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(100),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(200),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(200)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated refundPool",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(100),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(200),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(200)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated refund",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundList: []types.Refund{
					{
						DelegatorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(100),
					},
					{
						DelegatorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(200),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset deposit shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositList: []types.Deposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative deposit shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositList: []types.Deposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(-1),
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero deposit shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositList: []types.Deposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(0),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset deposit pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative deposit pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(-1),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero deposit pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(0),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset deposit pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative deposit pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.Coin{Denom: types.DefaultAllowedTokens[0], Amount: sdk.NewInt(-1)},
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero deposit pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid deposit pool denom",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				DepositPoolList: []types.DepositPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.Coin{Denom: "@!z", Amount: sdk.NewInt(1)},
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset refund pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative refund pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(-1),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero refund pool shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(0),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset refund pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative refund pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.Coin{Denom: types.DefaultAllowedTokens[0], Amount: sdk.NewInt(-1)},
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero refund pool tokens",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(0)),
					},
				},
			},
			valid: false,
		},
		{
			desc: "invalid refund pool denom",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundPoolList: []types.RefundPool{
					{
						OperatorAddress: valAddress0,
						Shares:          sdk.NewDec(1),
						Tokens:          sdk.Coin{Denom: "@!z", Amount: sdk.NewInt(1)},
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset refund shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundList: []types.Refund{
					{
						DelegatorAddress: address0,
						ValidatorAddress: valAddress0,
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative refund shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundList: []types.Refund{
					{
						DelegatorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(-1),
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero refund shares",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				RefundList: []types.Refund{
					{
						DelegatorAddress: address0,
						ValidatorAddress: valAddress0,
						Shares:           sdk.NewDec(0),
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset unbonding deposit initial balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								Balance:        sdk.NewInt(00),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative unbonding deposit initial balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(-1),
								Balance:        sdk.NewInt(00),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero unbonding deposit initial balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(0),
								Balance:        sdk.NewInt(00),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "unset unbonding deposit balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(10),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "negative unbonding deposit balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(10),
								Balance:        sdk.NewInt(-1),
							},
						},
					},
				},
			},
			valid: false,
		},
		{
			desc: "zero unbonding deposit balance",
			genState: &types.GenesisState{
				Params: types.DefaultParams(),
				UnbondingDepositList: []types.UnbondingDeposit{
					{
						DepositorAddress: address0,
						ValidatorAddress: valAddress0,
						Entries: []types.UnbondingDepositEntry{
							{
								CreationHeight: 10,
								CompletionTime: time.Unix(10, 0),
								InitialBalance: sdk.NewInt(10),
								Balance:        sdk.NewInt(00),
							},
						},
					},
				},
			},
			valid: true,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func validGenesisState(address0, address1, valAddress0, valAddress1 string) *types.GenesisState {
	return &types.GenesisState{
		Params: types.DefaultParams(),
		DepositList: []types.Deposit{
			{
				DepositorAddress: address0,
				ValidatorAddress: valAddress0,
				Shares:           sdk.NewDec(100),
			},
			{
				DepositorAddress: address1,
				ValidatorAddress: valAddress1,
				Shares:           sdk.NewDec(200),
			},
		},
		UnbondingDepositList: []types.UnbondingDeposit{
			{
				DepositorAddress: address0,
				ValidatorAddress: valAddress0,
				Entries: []types.UnbondingDepositEntry{
					{
						CreationHeight: 0,
						CompletionTime: time.Unix(0, 0),
						InitialBalance: sdk.NewInt(1),
						Balance:        sdk.NewInt(0),
					},
					{
						CreationHeight: 1,
						CompletionTime: time.Unix(1, 0),
						InitialBalance: sdk.NewInt(1),
						Balance:        sdk.NewInt(1),
					},
				},
			},
			{
				DepositorAddress: address1,
				ValidatorAddress: valAddress1,
				Entries: []types.UnbondingDepositEntry{
					{
						CreationHeight: 11,
						CompletionTime: time.Unix(11, 0),
						InitialBalance: sdk.NewInt(11),
						Balance:        sdk.NewInt(11),
					},
					{
						CreationHeight: 12,
						CompletionTime: time.Unix(12, 0),
						InitialBalance: sdk.NewInt(12),
						Balance:        sdk.NewInt(0),
					},
				},
			},
		},
		DepositPoolList: []types.DepositPool{
			{
				OperatorAddress: valAddress0,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
			{
				OperatorAddress: valAddress1,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(200)),
				Shares:          sdk.NewDec(100),
			},
		},
		RefundPoolList: []types.RefundPool{
			{
				OperatorAddress: valAddress0,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(100)),
				Shares:          sdk.NewDec(100),
			},
			{
				OperatorAddress: valAddress1,
				Tokens:          sdk.NewCoin(types.DefaultAllowedTokens[0], sdk.NewInt(200)),
				Shares:          sdk.NewDec(100),
			},
		},
		RefundList: []types.Refund{
			{
				DelegatorAddress: address0,
				ValidatorAddress: valAddress0,
				Shares:           sdk.NewDec(100),
			},
			{
				DelegatorAddress: address1,
				ValidatorAddress: valAddress1,
				Shares:           sdk.NewDec(200),
			},
		},
	}
}
