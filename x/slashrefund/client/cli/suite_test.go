package cli_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdknetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcliutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/made-in-block/slash-refund/testutil/network"
	"github.com/made-in-block/slash-refund/x/slashrefund/client/cli"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"
)

// UnbondingTime is the value that will be set in the staking module params to override
// the default value in order to allow full withdraw transaction lifecycle testing.
const UnbondingTime = time.Second * 5

type E2ETestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
	ctx     client.Context
	cdc     codec.Codec
}

func NewE2ETestSuite(cfg network.Config) *E2ETestSuite {
	return &E2ETestSuite{cfg: cfg}
}

// This function adds specific deposits and refunds to network configuration.
// This is done to make the transactions' tests independent one to another.
//
// TestCmdClaim needs a refund to be claimed to test a valid transaction.
// Two refunds will be set:
//  1. refund for (acc0,val0), claimed in TestCmdClaim/Valid_transaction
//  2. refund for (acc1,val0), claimed in TestCmdClaim/Valid_transaction_validate
//
// TestCmdWithdraw needs a deposit to be claimed to test a valid transacion.
// Two deposits will be set:
//  1. deposit for (acc0,val2), withdrew in TestCmdWithdraw/Valid_transaction
//  2. deposit for (acc1,val2), withdrew in TestCmdWithdraw/Valid_transaction_validate
func (s *E2ETestSuite) setObjectsToNetworkConfig(l sdknetwork.Logger, config network.Config) network.Config {

	l.Log("setting objects in network config...")

	s.Require().Equal(config.SigningAlgo, string(hd.Secp256k1Type), "error setting network objects: SigningAlgo set in network config must be %s", hd.Secp256k1Type)
	s.Require().GreaterOrEqual(len(config.Mnemonics), 3, "error setting network objects: config.Mnemonics must have at least three mnemonics set")

	derive := hd.Secp256k1.Derive()
	generate := hd.Secp256k1.Generate()

	shares := sdk.NewDec(1000)
	denom := sdk.DefaultBondDenom

	var deposits []types.Deposit
	var depositPools []types.DepositPool
	var refunds []types.Refund
	var refundPools []types.RefundPool

	// Get address and validator address of validator0.
	mnemonic := config.Mnemonics[0]
	bz, err := derive(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	s.Require().NoError(err, "error setting network objects: failed to derive privk from mnemonic 0")
	pubk := generate(bz).PubKey().Address()
	address0 := sdk.AccAddress(pubk).String()
	validator0 := sdk.ValAddress(pubk).String()

	// Get validator address of validator1.
	mnemonic = config.Mnemonics[1]
	bz, err = derive(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	s.Require().NoError(err, "error setting network objects: failed to derive privk from mnemonic 1")
	pubk = generate(bz).PubKey().Address()
	address1 := sdk.AccAddress(pubk).String()

	// Get validator address of validator2.
	mnemonic = config.Mnemonics[2]
	bz, err = derive(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	s.Require().NoError(err, "error setting network objects: failed to derive privk from given mnemonic 2")
	pubk = generate(bz).PubKey().Address()
	validator2 := sdk.ValAddress(pubk).String()

	// Set refund and refund pool (address0,validator0).
	ref00 := types.Refund{
		DelegatorAddress: address0,
		ValidatorAddress: validator0,
		Shares:           shares,
	}
	ref01 := types.Refund{
		DelegatorAddress: address1,
		ValidatorAddress: validator0,
		Shares:           shares,
	}
	refPool0 := types.RefundPool{
		OperatorAddress: validator0,
		Tokens:          sdk.NewCoin(denom, shares.TruncateInt().MulRaw(2)),
		Shares:          shares.MulInt64(2),
	}
	refunds = append(refunds, ref00, ref01)
	refundPools = append(refundPools, refPool0)

	// Set deposit and deposit pool (address0,validator1).
	dep02 := types.Deposit{
		DepositorAddress: address0,
		ValidatorAddress: validator2,
		Shares:           shares,
	}
	dep12 := types.Deposit{
		DepositorAddress: address1,
		ValidatorAddress: validator2,
		Shares:           shares,
	}
	depPool2 := types.DepositPool{
		OperatorAddress: validator2,
		Tokens:          sdk.NewCoin(denom, shares.TruncateInt().MulRaw(2)),
		Shares:          shares.MulInt64(2),
	}
	deposits = append(deposits, dep02, dep12)
	depositPools = append(depositPools, depPool2)

	// Set slashrefund module genesis.
	state := types.GenesisState{}
	state.RefundList = append(state.RefundList, refunds...)
	state.RefundPoolList = append(state.RefundPoolList, refundPools...)
	state.DepositList = append(state.DepositList, deposits...)
	state.DepositPoolList = append(state.DepositPoolList, depositPools...)
	state.Params = types.DefaultParams()
	buf, err := config.Codec.MarshalJSON(&state)
	s.Require().NoError(err)
	config.GenesisState[types.ModuleName] = buf
	l.Log("set refund and deposits in network config")

	// Set staking module genesis: UnbondingTime used in deposit withdraw is the same
	// UnbondingTime of the staking module. It will be set to 5 seconds in order to
	// allow complete testing of withdraw command.
	stateSt := stakingtypes.GenesisState{}
	stateSt.Params = stakingtypes.DefaultParams()
	stateSt.Params.UnbondingTime = UnbondingTime
	buf, err = config.Codec.MarshalJSON(&stateSt)
	s.Require().NoError(err)
	config.GenesisState[stakingtypes.ModuleName] = buf
	l.Logf("set unbonding time in network config (current value: %f seconds)", UnbondingTime.Seconds())

	// Update slashrefund module balance with the deposit and the refund amount in
	// order to have funds to be sent from module account to the depositor address
	// when deposit is withdrawn and to delegator address when refund is claimed.
	var balances []banktypes.Balance
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(types.ModuleName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(denom, shares.TruncateInt().MulRaw(4))},
	})
	bankstate := banktypes.GenesisState{}
	bankstate.Balances = append(bankstate.Balances, balances...)
	buf, err = config.Codec.MarshalJSON(&bankstate)
	s.Require().NoError(err)
	config.GenesisState[banktypes.ModuleName] = buf
	l.Logf("set slashrefund module balance in network config")

	return config
}

// NewNetworkWithObjects creates a new test network with two specific deposits and two
// specific refunds added to network genesis configuration.
func (s *E2ETestSuite) NewNetworkWithObjects(config network.Config) *sdknetwork.Network {
	net, err := sdknetwork.New(s.T(), s.T().TempDir(), s.setObjectsToNetworkConfig(s.T(), config))
	s.Require().NoError(err)

	return net
}

func (s *E2ETestSuite) SetupSuite() {

	s.T().Log("setting up e2e test suite.")

	s.network = s.NewNetworkWithObjects(s.cfg)
	s.cdc = s.network.Config.Codec
	s.cfg = s.network.Config
	s.ctx = s.network.Validators[0].ClientCtx

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	// Import account1 (account of validator1) in the client keyring.
	// For this account a deposit and a refund are available (set in Genesis).
	// This account is already funded during network setup phase.
	s.T().Log("importing account1 in client keyring.")
	_, err = s.ctx.Keyring.NewAccount("account1", s.cfg.Mnemonics[1], keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path(), hd.Secp256k1)
	s.Require().NoError(err)
	s.T().Log("finished setting up suite.")
}

func (s *E2ETestSuite) TearDownSuite() {
	s.T().Log("tearing down e2e test suite")
	s.network.Cleanup()
}

// This test checks for errors during the execution of the deposit command and checks
// also for the correct execution of a valid deposit transaction.
// First, errors and invalid transactions are tested, then in "Valid transaction"
// subtest, a valid deposit transaction from address0 to validator1 is done.
// These tests are used to test the cli responses.
// Eventually, in the "Valid transaction validation" subtest, a new deposit from
// address1 to validator1 is made and checks are performed to validate the result of
// this valid transaction. This subtest is used to check that the deposit and deposit
// pool are updated as expected, and the balance of the depositor decreases as expected
// expected when the valid transaction is processed.
// Checks are performed through queries, as an actual end user would do.
func (s *E2ETestSuite) TestCmdDeposit() {

	denom := sdk.DefaultBondDenom
	idDepositorAddress := s.network.Validators[0].Address.String()
	idValidatorAddress1 := s.network.Validators[1].ValAddress.String()
	feeAmt := sdk.NewInt(10)
	fees := sdk.NewCoins(sdk.NewCoin(denom, feeAmt)).String()
	successCode := sdkerrors.SuccessABCICode
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	testCases := []struct {
		name       string
		args       []string
		expectErr  bool
		txRespCode uint32 //unused if expectErr=true
	}{
		{
			"Error (Without validator address nor amount)",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Without amount)",
			[]string{
				idValidatorAddress1,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Without from-address)",
			[]string{
				idValidatorAddress1,
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Fail decoding validator address)",
			[]string{
				"not-a-validator-address",
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Non valoper validator address)",
			[]string{
				idDepositorAddress,
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Zero amount)",
			[]string{
				idValidatorAddress1,
				sdk.NewCoin(denom, sdk.ZeroInt()).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Negative amount)",
			[]string{
				idValidatorAddress1,
				fmt.Sprintf("-1%s", denom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Invalid (Not found validator address)",
			[]string{
				"cosmosvaloper1uhdmcuszs29hnyqtsjn9cm7cyrmkcnq4undkv5",
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, stakingtypes.ErrNoValidatorFound.ABCICode(),
		},
		{
			"Invalid (Amount higher than actual balance)",
			[]string{
				idValidatorAddress1,
				sdk.NewCoin(denom, sdk.DefaultPowerReduction.MulRaw(999999)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, sdkerrors.ErrInsufficientFunds.ABCICode(),
		},
		{
			"Valid transaction",
			[]string{
				idValidatorAddress1,
				sdk.NewCoin(denom, sdk.NewInt(100)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, successCode,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdDeposit(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				var resp sdk.TxResponse
				s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp), out.String())
				s.Require().Equal(tc.txRespCode, resp.Code, out.String())
			}
		})
	}

	// This subtest must follow TestCmd/Valid_transaction subtest.
	s.Run("Valid transaction validation", func() {

		// A deposit for (acc0,val0) has been created from TestCmdDeposit/Valid_transaction
		// subtest. Require the deposit pool can be found through the query.
		args := []string{idValidatorAddress1, outflag}
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDepositPool(), args)
		s.Require().NoError(err, out.String())
		var resp9 types.QueryGetDepositPoolResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp9), out.String())
		s.Require().NotEmpty(resp9.DepositPool)
		oldPoolShares := resp9.DepositPool.Shares
		oldPoolTokens := resp9.DepositPool.Tokens

		// Get account1 key from the keyring and get its address.
		key, err := s.ctx.Keyring.Key("account1")
		pub, err := key.GetPubKey()
		s.Require().NoError(err)
		depAddr := sdk.AccAddress(pub.Address())
		idDepositorAddress1 := depAddr.String()

		// Get account initial balance.
		out, err = bankcliutil.QueryBalancesExec(s.ctx, depAddr, outflag)
		s.Require().NoError(err)
		var resp banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		amt0 := resp.Balances.AmountOf(denom)

		// Execute valid deposit transaction to validator1 and require it returns the success code.
		depAmt := sdk.NewInt(100)
		depAmtDec := sdk.NewDecFromInt(depAmt)
		args = []string{
			idValidatorAddress1,
			sdk.NewCoin(denom, depAmt).String(),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress1),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdDeposit(), args)
		s.Require().NoError(err, out.String())
		var resp1 sdk.TxResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp1), out.String())
		s.Require().Equal(successCode, resp1.Code, out.String())

		// Get account actual balance.
		depAddr, err = sdk.AccAddressFromBech32(idDepositorAddress1)
		s.Require().NoError(err)
		out, err = bankcliutil.QueryBalancesExec(s.ctx, depAddr, outflag)
		s.Require().NoError(err)
		var resp2 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp2))
		amt1 := resp2.Balances.AmountOf(denom)

		// Require balance of depositor decreased of an amount equal to deposited amount and fees payed.
		s.Require().Equal(depAmt.Add(feeAmt), amt0.Sub(amt1))

		// Require the deposit can be found through the query.
		args = []string{idDepositorAddress1, idValidatorAddress1, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDeposit(), args)
		s.Require().NoError(err, out.String())
		var resp3 types.QueryGetDepositResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp3), out.String())
		s.Require().NotEmpty(resp3.Deposit)
		s.Require().Equal(depAmtDec, resp3.Deposit.Shares, out.String())

		// Require the deposit pool can be found through the query and it matches the two deposits made.
		args = []string{idValidatorAddress1, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDepositPool(), args)
		s.Require().NoError(err, out.String())
		var resp4 types.QueryGetDepositPoolResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp4), out.String())
		s.Require().NotEmpty(resp4.DepositPool)
		s.Require().Equal(depAmt, resp4.DepositPool.Tokens.Amount.Sub(oldPoolTokens.Amount), out.String())
		s.Require().Equal(depAmtDec, resp4.DepositPool.Shares.Sub(oldPoolShares), out.String())
	})
}

// This test checks for errors during the execution of the withdraw command and checks
// also for the correct execution of a valid withdraw transaction.
// In order to isolate withdraw transaction command from deposit transaction command,
// this test transactions target specific deposits set in network genesis configuration.
// First, errors and invalid transactions are tested, then in "Valid transaction"
// subtest, a valid withdraw transaction from address0 to validator2 is done.
// These tests are used to test the cli responses.
// Eventually, in the "Valid transaction validation" subtest, a new withdraw from
// address1 to validator2 is made and checks are performed to validate the result of
// this valid transaction. This subtest is used to check that the deposit and deposit
// pool are updated as expected, and the balance of the depositor increases as expected
// after the unbonding period ends.
// Checks are performed through queries, as an actual end user would do.
func (s *E2ETestSuite) TestCmdWithdraw() {

	denom := sdk.DefaultBondDenom
	idDepositorAddress := s.network.Validators[0].Address.String()
	idValidatorAddress2 := s.network.Validators[2].ValAddress.String()
	idValidatorAddress3 := s.network.Validators[3].ValAddress.String()
	feeAmt := sdk.NewInt(10)
	fees := sdk.NewCoins(sdk.NewCoin(denom, feeAmt)).String()
	successCode := sdkerrors.SuccessABCICode
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	testCases := []struct {
		name       string
		args       []string
		expectErr  bool
		txRespCode uint32 //unused if expectErr=true
	}{
		{
			"Error (Without validator address nor amount)",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Without amount)",
			[]string{
				idValidatorAddress2,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Without from-address)",
			[]string{
				idValidatorAddress2,
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Fail decoding validator address)",
			[]string{
				"not-a-validator-address",
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Non valoper validator address)",
			[]string{
				idDepositorAddress,
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Zero amount)",
			[]string{
				idValidatorAddress2,
				sdk.NewCoin(denom, sdk.ZeroInt()).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Negative amount)",
			[]string{
				idValidatorAddress2,
				fmt.Sprintf("-1%s", denom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Invalid (Not found validator address)",
			[]string{
				"cosmosvaloper1uhdmcuszs29hnyqtsjn9cm7cyrmkcnq4undkv5",
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, stakingtypes.ErrNoValidatorFound.ABCICode(),
		},
		{
			"Invalid (Amount higher than deposited)",
			[]string{
				idValidatorAddress2,
				sdk.NewCoin(denom, sdk.DefaultPowerReduction.MulRaw(999999)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, sdkerrors.ErrInvalidRequest.ABCICode(),
		},
		{
			"Invalid (No deposit for address)",
			[]string{
				idValidatorAddress3,
				sdk.NewCoin(denom, sdk.NewInt(100)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, types.ErrNoDepositForAddress.ABCICode(),
		},
		{
			"Valid transaction",
			[]string{
				idValidatorAddress2,
				sdk.NewCoin(denom, sdk.NewInt(100)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDepositorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, successCode,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdWithdraw(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				var resp sdk.TxResponse
				s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp), out.String())
				s.Require().Equal(tc.txRespCode, resp.Code, out.String())
			}
		})
	}

	// This subtest must follow TestCmdWithdraw/Valid_transaction subtest.
	s.Run("Valid transaction validation", func() {

		// A withdraw for (acc0,val2) has been created from TestCmdWithdraw/
		// Valid_transaction subtest. Require the deposit pool can be found through the
		// query. Actual values of the deposit pool are needed to verify its updating.
		args := []string{idValidatorAddress2, outflag}
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDepositPool(), args)
		s.Require().NoError(err, out.String())
		var resp9 types.QueryGetDepositPoolResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp9), out.String())
		s.Require().NotEmpty(resp9.DepositPool)
		oldPoolShares := resp9.DepositPool.Shares
		oldPoolTokens := resp9.DepositPool.Tokens

		// Get account1 key from the keyring and get its address.
		key, err := s.ctx.Keyring.Key("account1")
		pub, err := key.GetPubKey()
		s.Require().NoError(err)
		depAddr := sdk.AccAddress(pub.Address())
		idDepositorAddress1 := depAddr.String()

		// Require the deposit for (acc1,val2) is correctly set in genesis when network is set up.
		args = []string{idDepositorAddress1, idValidatorAddress2, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDeposit(), args)
		s.Require().NoError(err, out.String())
		var resp99 types.QueryGetDepositResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp99), out.String())
		s.Require().NotEmpty(resp99.Deposit)
		oldDepShares := resp99.Deposit.Shares

		// Get account initial balance.
		out, err = bankcliutil.QueryBalancesExec(s.ctx, depAddr, outflag)
		s.Require().NoError(err)
		var resp banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		amt0 := resp.Balances.AmountOf(denom)

		// Execute valid withdraw transaction to validator2 and require it returns the success code.
		witAmt := sdk.NewInt(100)
		witAmtDec := sdk.NewDecFromInt(witAmt)
		args = []string{
			idValidatorAddress2,
			sdk.NewCoin(denom, witAmt).String(),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, depAddr.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdWithdraw(), args)
		s.Require().NoError(err, out.String())
		var resp1 sdk.TxResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp1), out.String())
		s.Require().Equal(successCode, resp1.Code, out.String())

		// Require the unbonding deposit can be found through a query.
		args = []string{depAddr.String(), idValidatorAddress2, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowUnbondingDeposit(), args)
		s.Require().NoError(err, out.String())
		var resp2 types.QueryGetUnbondingDepositResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp2), out.String())
		s.Require().NotEmpty(resp2.UnbondingDeposit)
		s.Require().Equal(1, len(resp2.UnbondingDeposit.Entries))
		s.Require().Equal(witAmt, resp2.UnbondingDeposit.Entries[0].Balance, out.String())

		// Get account actual balance.
		out, err = bankcliutil.QueryBalancesExec(s.ctx, depAddr, outflag)
		s.Require().NoError(err)
		var resp3 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp3))
		amt1 := resp3.Balances.AmountOf(denom)

		// Require balance has not changed (except for fees payed) since the withdrawn amount is still in its unbonding period.
		s.Require().Equal(amt0.Sub(feeAmt), amt1)

		// Require the deposit still exists and it is decreased by the withdrawn amount.
		args = []string{depAddr.String(), idValidatorAddress2, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDeposit(), args)
		s.Require().NoError(err, out.String())
		var resp4 types.QueryGetDepositResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp4), out.String())
		s.Require().NotEmpty(resp4.Deposit)
		s.Require().Equal(witAmtDec, oldDepShares.Sub(resp4.Deposit.Shares), out.String())

		// Require the deposit pool still exists and matches the withdraw made.
		args = []string{idValidatorAddress2, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowDepositPool(), args)
		s.Require().NoError(err, out.String())
		var resp5 types.QueryGetDepositPoolResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp5), out.String())
		s.Require().NotEmpty(resp5.DepositPool)
		s.Require().Equal(witAmt, oldPoolTokens.Amount.Sub(resp5.DepositPool.Tokens.Amount), out.String())
		s.Require().Equal(witAmtDec, oldPoolShares.Sub(resp5.DepositPool.Shares), out.String())

		// Wait for the unbonding period and an additional block to be sure the
		// unbonding deposit is mature when balance is checked.
		t0 := time.Now()
		t1 := t0.Add(UnbondingTime)
		for t0.Unix() < t1.Unix() {
			s.Assert().NoError(s.network.WaitForNextBlock(), "Failed waiting for next block. Timeout of 10seconds reached.")
			t0 = time.Now()
		}
		s.Require().NoError(s.network.WaitForNextBlock())

		// Get account actual balance.
		out, err = bankcliutil.QueryBalancesExec(s.ctx, depAddr, outflag)
		s.Require().NoError(err)
		var resp7 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp7))
		amt2 := resp7.Balances.AmountOf(denom)

		// Require balance of depositor increased by an amount equal to withrawn amount (minus fees payed).
		s.Require().True(amt2.GT(amt1), "Withdrawn amount did not return to depositor.")
		s.Require().Equal(amt1.Add(witAmt), amt2, "Withdrawn amount returned is different than expected.")

		// Require the unbonding deposit can not be found through a query.
		args = []string{depAddr.String(), idValidatorAddress2, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowUnbondingDeposit(), args)
		s.Require().Error(err)
		s.Require().ErrorContains(err, "key not found", "Unbonding deposit matured but still in queue.")
	})
}

// This test checks for errors during the execution of the claim command and checks
// also for the correct execution of a valid claim transaction.
// In order to isolate claim transaction command from a slash event, this test
// transactions target specific refunds set in network genesis configuration.
// First, errors and invalid transactions are tested, then in "Valid transaction"
// subtest, a valid claim transaction from address0 to validator0 is done.
// These tests are used to test the cli responses.
// Eventually, in the "Valid transaction validation" subtest, a new claim from
// address1 to validator0 is made and checks are performed to validate the result of
// this valid transaction. This subtest is used to check that the refund and refund
// pool are updated as expected, and the balance of the address increases as expected
// after the claim transaction is processed.
// Checks are performed through queries, as an actual end user would do.
func (s *E2ETestSuite) TestCmdClaim() {

	denom := sdk.DefaultBondDenom
	idDelegatorAddress := s.network.Validators[0].Address.String()
	idValidatorAddress := s.network.Validators[0].ValAddress.String()
	feeAmt := sdk.NewInt(10)
	fees := sdk.NewCoins(sdk.NewCoin(denom, feeAmt)).String()
	successCode := sdkerrors.SuccessABCICode
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	testCases := []struct {
		name       string
		args       []string
		expectErr  bool
		txRespCode uint32 //unused if expectErr=true
	}{
		{
			"Error (Without validator address)",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (With amount)",
			[]string{
				idValidatorAddress,
				sdk.NewCoin(denom, sdk.NewInt(1)).String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Without from-address)",
			[]string{
				idValidatorAddress,
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Fail decoding validator address)",
			[]string{
				"not-a-validator-address",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Error (Non valoper validator address)",
			[]string{
				idDelegatorAddress,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			true, 0,
		},
		{
			"Invalid (Not found validator address)",
			[]string{
				"cosmosvaloper1uhdmcuszs29hnyqtsjn9cm7cyrmkcnq4undkv5",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, stakingtypes.ErrNoValidatorFound.ABCICode(),
		},
		{
			"Invalid (No refund for address)",
			[]string{
				s.network.Validators[1].ValAddress.String(),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, types.ErrNoRefundForAddress.ABCICode(),
		},
		{
			"Valid transaction",
			[]string{
				idValidatorAddress,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, idDelegatorAddress),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
			},
			false, successCode,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdClaim(), tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				var resp sdk.TxResponse
				s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp), out.String())
				s.Require().Equal(tc.txRespCode, resp.Code, out.String())
			}
		})
	}

	// This subtest must follow TestCmdRefund/Valid_transaction subtest.
	s.Run("Valid transaction validation", func() {

		// In TestCmdClaim/Valid_transaction the refund for (addr0,val0) has been
		// claimed. This refund must not be returned by a query, since it must be
		// deleted from the store.
		args := []string{idDelegatorAddress, idValidatorAddress, outflag}
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowRefund(), args)
		s.Require().Error(err)
		s.Require().ErrorContains(err, "key not found")

		// Require the refund pool still exists after the first claim, since the
		// refund for (acc1,val0) has not been claimed.
		args = []string{idValidatorAddress, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowRefundPool(), args)
		s.Require().NoError(err)
		var resp1 types.QueryGetRefundPoolResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp1))
		s.Require().NotEmpty(resp1.RefundPool)

		// Get account1 key from the keyring and get its address.
		key, err := s.ctx.Keyring.Key("account1")
		pub, err := key.GetPubKey()
		s.Require().NoError(err)
		addr := sdk.AccAddress(pub.Address())
		idDelegatorAddress1 := addr.String()

		// Require refund for (acc1,val1) has been correctly set in genesis.
		args = []string{idDelegatorAddress1, idValidatorAddress, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowRefund(), args)
		s.Require().NoError(err)
		var resp0 types.QueryGetRefundResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp0))
		s.Require().NotNil(resp0.Refund)

		// Require refund and refund pool are correctly linked, because now only the
		// refund for (acc1,val0) exists in the network.
		s.Require().Equal(resp0.Refund.Shares, resp1.RefundPool.Shares)
		s.Require().Equal(resp0.Refund.Shares, sdk.NewDecFromInt(resp1.RefundPool.Tokens.Amount))
		refAmt := resp1.RefundPool.Tokens.Amount

		// Get account initial balance
		out, err = bankcliutil.QueryBalancesExec(s.ctx, addr, outflag)
		s.Require().NoError(err)
		var resp2 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp2))
		amt0 := resp2.Balances.AmountOf(denom)

		args = []string{
			idValidatorAddress,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, addr.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		}

		// Execute transaction and require it returns success code.
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdClaim(), args)
		s.Require().NoError(err, out.String())
		var resp3 sdk.TxResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp3), out.String())
		s.Require().Equal(successCode, resp3.Code, out.String())

		// Get account actual balance.
		out, err = bankcliutil.QueryBalancesExec(s.ctx, addr, outflag)
		s.Require().NoError(err)
		var resp4 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp4))
		amt1 := resp4.Balances.AmountOf(denom)

		// Require refund has been added to account balance.
		s.Require().Equal(refAmt.Sub(feeAmt), amt1.Sub(amt0))

		// Require refund is no more returned by a query, since it must be removed
		// after claim is correctly executed.
		args = []string{idDelegatorAddress1, idValidatorAddress, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowRefund(), args)
		s.Require().Error(err)
		s.Require().ErrorContains(err, "key not found")

		// Require the refund pool does not exists after the second claim, since all
		// refunds for validator0 have been claimed.
		args = []string{idValidatorAddress, outflag}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdShowRefundPool(), args)
		s.Require().Error(err)
		s.Require().ErrorContains(err, "key not found")

		// Require that a new claim transaction for the already claimed refund returns
		// the "no refund for address" error code.
		args = []string{
			idValidatorAddress,
			fmt.Sprintf("--%s=%s", flags.FlagFrom, addr.String()),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=%s", flags.FlagFees, fees),
		}
		out, err = clitestutil.ExecTestCLICmd(s.ctx, cli.CmdClaim(), args)
		s.Require().NoError(err, out.String())
		var resp5 sdk.TxResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp5), out.String())
		s.Require().Equal(types.ErrNoRefundForAddress.ABCICode(), resp5.Code, out.String())

		// Get account actual balance
		out, err = bankcliutil.QueryBalancesExec(s.ctx, addr, outflag)
		s.Require().NoError(err)
		var resp6 banktypes.QueryAllBalancesResponse
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp6))
		amt2 := resp6.Balances.AmountOf(denom)

		// Require that balance has not changed (only decreased due to fees payed to
		// execute the second claim transaction.
		s.Require().Equal(amt1.Sub(feeAmt), amt2)
	})
}

/*
args := []string{
	fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	fmt.Sprintf("--%s=0", flags.FlagOffset),
	fmt.Sprintf("--%s=%d", flags.FlagLimit, uint64(len(objs))),
	fmt.Sprintf("--%s", flags.FlagCountTotal),
}
*/
