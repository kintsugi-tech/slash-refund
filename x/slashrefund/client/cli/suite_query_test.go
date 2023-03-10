package cli_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdknetwork "github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/made-in-block/slash-refund/testutil/network"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/made-in-block/slash-refund/x/slashrefund/client/cli"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type E2EQueryTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
	ctx     client.Context
	cdc     codec.Codec

	deposits          []types.Deposit
	depositPools      []types.DepositPool
	params            types.Params
	refunds           []types.Refund
	refundPools       []types.RefundPool
	unbondingDeposits []types.UnbondingDeposit
}

func NewE2EQueryTestSuite(cfg network.Config) *E2EQueryTestSuite {
	return &E2EQueryTestSuite{cfg: cfg}
}

func (s *E2EQueryTestSuite) setObjectsToNetworkConfig(l sdknetwork.Logger, config network.Config) network.Config {

	// Set default params for query params test.
	s.params = types.DefaultParams()

	// Set 25 deposits with 5 associated deposit pools for query deposit and depoist
	// pool tests.
	for i := 0; i < 5; i++ {
		validatorAddress := sample.ValAddress()
		totShares := sdk.ZeroDec()
		for j := 0; j < 5; j++ {
			shares := sdk.NewDec(int64(100 * (j + 1)))
			obj := types.Deposit{
				DepositorAddress: sample.AccAddress(),
				ValidatorAddress: validatorAddress,
				Shares:           shares,
			}
			totShares = totShares.Add(shares)
			s.deposits = append(s.deposits, obj)
		}
		obj := types.DepositPool{
			OperatorAddress: validatorAddress,
			Tokens:          sdk.NewCoin(sdk.DefaultBondDenom, totShares.TruncateInt()),
			Shares:          totShares,
		}
		s.depositPools = append(s.depositPools, obj)
	}

	// Set 5 unbonding deposits for query unbonding deposit tests.
	for i := 0; i < 5; i++ {
		obj := types.UnbondingDeposit{
			DepositorAddress: sample.AccAddress(),
			ValidatorAddress: sample.ValAddress(),
			Entries: []types.UnbondingDepositEntry{
				{
					CreationHeight: int64(i),
					CompletionTime: time.Date(1970, time.January, 1, i, i, i, i, time.UTC),
					InitialBalance: sdk.NewInt(int64(100 * (i + 1))),
					Balance:        sdk.NewInt(int64(100 * i)),
				},
				{
					CreationHeight: int64(i + 1),
					CompletionTime: time.Date(1970, time.January, 1, i+1, i+1, i+1, i+1, time.UTC),
					InitialBalance: sdk.NewInt(int64(100 * (i + 2))),
					Balance:        sdk.NewInt(int64(100 * (i + 1))),
				},
			},
		}
		s.unbondingDeposits = append(s.unbondingDeposits, obj)
	}

	// Set 25 refunds with 5 associated refund pools for query refund and refund pool
	// tests.
	for i := 0; i < 5; i++ {
		validatorAddress := sample.ValAddress()
		totShares := sdk.ZeroDec()
		for j := 0; j < 5; j++ {
			shares := sdk.NewDec(int64(100 * (j + 1)))
			obj := types.Refund{
				DelegatorAddress: sample.AccAddress(),
				ValidatorAddress: validatorAddress,
				Shares:           shares,
			}
			totShares = totShares.Add(shares)
			s.refunds = append(s.refunds, obj)
		}
		obj := types.RefundPool{
			OperatorAddress: validatorAddress,
			Tokens:          sdk.NewCoin(sdk.DefaultBondDenom, totShares.TruncateInt()),
			Shares:          totShares,
		}
		s.refundPools = append(s.refundPools, obj)
	}

	// Set generated objects in slash-refund module genesis state.
	state := types.GenesisState{}
	state.Params = s.params
	state.DepositList = append(state.DepositList, s.deposits...)
	state.DepositPoolList = append(state.DepositPoolList, s.depositPools...)
	state.UnbondingDepositList = append(state.UnbondingDepositList, s.unbondingDeposits...)
	state.RefundList = append(state.RefundList, s.refunds...)
	state.RefundPoolList = append(state.RefundPoolList, s.refundPools...)
	buf, err := config.Codec.MarshalJSON(&state)
	s.Require().NoError(err)

	// Set slash-refund module genesis state in network genesis state configuration.
	config.GenesisState[types.ModuleName] = buf

	return config
}

// NewNetworkWithObjects creates a new test network with slash-refund module objects
// added to network genesis configuration.
func (s *E2EQueryTestSuite) NewNetworkWithObjects(config network.Config) *sdknetwork.Network {
	net, err := sdknetwork.New(s.T(), s.T().TempDir(), s.setObjectsToNetworkConfig(s.T(), config))
	s.Require().NoError(err)

	return net
}

func (s *E2EQueryTestSuite) SetupSuite() {

	s.T().Log("setting up query e2e test suite.")

	s.network = s.NewNetworkWithObjects(s.cfg)

	s.cdc = s.network.Config.Codec
	s.cfg = s.network.Config
	s.ctx = s.network.Validators[0].ClientCtx

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *E2EQueryTestSuite) TearDownSuite() {
	s.T().Log("tearing down query e2e test suite")
	s.network.Cleanup()
}

func (s *E2EQueryTestSuite) TestCmdQueryParams() {

	args := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
	out, err := clitestutil.ExecTestCLICmd(s.ctx, cli.CmdQueryParams(), args)
	s.Require().NoError(err)
	var resp types.QueryParamsResponse
	s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
	s.Require().Equal(s.params, resp.Params)
}

func (s *E2EQueryTestSuite) TestCmdShowDepositPool() {

	objs := s.depositPools
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	s.Run("Errors", func() {
		r := types.QueryGetDepositPoolResponse{}

		for _, tc := range []struct {
			desc               string
			idValidatorAddress string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idValidatorAddress: sample.ValAddress(),
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "FailDecodingValidatorAddress",
				idValidatorAddress: sample.MockAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidValidatorAddress",
				idValidatorAddress: sample.AccAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyValidatorAddress",
				idValidatorAddress: "",
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			s.Run(tc.desc, func() {
				cmd := cli.CmdShowDepositPool()
				args := append([]string{tc.idValidatorAddress}, outflag)
				var resp types.QueryGetDepositPoolResponse

				// Require the command execution returns an error with expErrMsg in its description..
				out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
				s.Require().Contains(err.Error(), tc.expErrMsg)

				// Require the output cannot be unmarshaled.
				s.Require().Error(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().Equal(r.DepositPool, resp.DepositPool)

				// Require the error returned is identified by the expected error code
				// and contains the expected error message in its description.
				stat, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expErrCode, stat.Code())
				s.Require().Contains(stat.Message(), tc.expErrMsg)
			})
		}
	})

	s.Run("Valids", func() {
		for _, obj := range objs {
			cmd := cli.CmdShowDepositPool()
			args := append([]string{obj.OperatorAddress}, outflag)
			var resp types.QueryGetDepositPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().Equal(obj, resp.DepositPool)
		}
	})
}

func (s *E2EQueryTestSuite) TestCmdListDepositPool() {

	objs := s.depositPools

	// These tests analyze the query responses when the following flags are used:
	// --ofset : Offset is used to query objects starting from a given index.
	//           In these tests, offset is incremented at each step by 2.
	// --limit : Limit is used to query objects and limit the number of objects
	//           returned up to the specified value.
	//           In this tests, limit is set equal to 2.
	// A loop is used to analyze all objects, and fails if returned objects does not
	// match the objects set in the genesis configuration.
	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListDepositPool()
			args := s.argsForPaginatedResp(nil, uint64(i), uint64(step), false)
			var resp types.QueryAllDepositPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.DepositPool), step)
			s.Require().Subset(objs, resp.DepositPool)
			s.Require().NotEmpty(resp.DepositPool)
		}
	})

	// These tests analyze the query responses when the following flags are used:
	// --page-key : Page-key is used to query objects starting from a given object key
	//              in the KV-store.
	//              In these tests, the next-key value to set in the next query is
	//              taken from the query response. This tests starts from a query
	//              done without this flag (using --offset=0 instead) in order to
	//              take the first next-key value from the first response.
	// --limit : limit is set equal to 2.
	// A loop is used to analyze all objects, and fails if returned objects does not
	// match the objects set in the genesis configuration.
	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListDepositPool()
			args := s.argsForPaginatedResp(next, 0, uint64(step), false)
			var resp types.QueryAllDepositPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().NotEmpty(resp.DepositPool)
			s.Require().LessOrEqual(len(resp.DepositPool), step)
			s.Require().Subset(objs, resp.DepositPool)
			s.Require().NotEmpty(resp.DepositPool)
			next = resp.Pagination.NextKey
		}
	})

	// These tests analyze the query responses when the following flags are used:
	// --offset : Offset is set to zero, in order to query for all objects.
	// --limit : Limit is set to the total number of requested objects, in order to
	//           query for all objects.
	// --count-total : Count-total is set to true in order to obtain the number of
	//                 objects found by the query.
	// A loop is used to analyze all objects, and fails if returned objects does not
	// match the objects set in the genesis configuration.
	s.Run("Total", func() {
		cmd := cli.CmdListDepositPool()
		args := s.argsForPaginatedResp(nil, 0, uint64(len(objs)), true)
		var resp types.QueryAllDepositPoolResponse
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
		s.Require().NoError(err)
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(objs, resp.DepositPool)
	})
}

func (s *E2EQueryTestSuite) TestCmdShowDeposit() {

	objs := s.deposits
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	s.Run("Errors", func() {
		r := types.QueryGetDepositResponse{}

		for _, tc := range []struct {
			desc               string
			idDepositorAddress string
			idValidatorAddress string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idDepositorAddress: sample.AccAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "FailDecodingAddress",
				idDepositorAddress: sample.MockAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "FailDecodingValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: sample.MockAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidAddress",
				idDepositorAddress: objs[0].ValidatorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[0].DepositorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyAddress",
				idDepositorAddress: "",
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: "",
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			s.Run(tc.desc, func() {
				cmd := cli.CmdShowDeposit()
				args := append([]string{tc.idDepositorAddress, tc.idValidatorAddress}, outflag)
				var resp types.QueryGetDepositResponse

				// Require the command execution returns an error with expErrMsg in its description.
				out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
				s.Require().Contains(err.Error(), tc.expErrMsg)

				// Require the output cannot be unmarshaled.
				s.Require().Error(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().Equal(r.Deposit, resp.Deposit)

				// Require the error returned is identified by the expected error code
				// and contains the expected error message in its description.
				stat, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expErrCode, stat.Code())
				s.Require().Contains(stat.Message(), tc.expErrMsg)
			})
		}
	})

	s.Run("Valids", func() {
		for _, obj := range objs {
			cmd := cli.CmdShowDeposit()
			args := append([]string{obj.DepositorAddress, obj.ValidatorAddress}, outflag)
			var resp types.QueryGetDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().Equal(obj, resp.Deposit)
		}
	})
}

func (s *E2EQueryTestSuite) TestCmdListDeposit() {

	objs := s.deposits

	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListDeposit()
			args := s.argsForPaginatedResp(nil, uint64(i), uint64(step), false)
			var resp types.QueryAllDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.Deposit), step)
			s.Require().Subset(objs, resp.Deposit)
			s.Require().NotEmpty(resp.Deposit)
		}
	})

	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListDeposit()
			args := s.argsForPaginatedResp(next, 0, uint64(step), false)
			var resp types.QueryAllDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.Deposit), step)
			s.Require().Subset(objs, resp.Deposit)
			s.Require().NotEmpty(resp.Deposit)
			next = resp.Pagination.NextKey
		}
	})

	s.Run("Total", func() {
		cmd := cli.CmdListDeposit()
		args := s.argsForPaginatedResp(nil, 0, uint64(len(objs)), true)
		var resp types.QueryAllDepositResponse
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
		s.Require().NoError(err)
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(objs, resp.Deposit)
	})
}

func (s *E2EQueryTestSuite) TestCmdShowRefund() {

	objs := s.refunds
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	s.Run("Errors", func() {
		r := types.QueryGetRefundResponse{}

		for _, tc := range []struct {
			desc               string
			idDelegatorAddress string
			idValidatorAddress string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idDelegatorAddress: sample.AccAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "FailDecodingAddress",
				idDelegatorAddress: sample.MockAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "FailDecodingValidatorAddress",
				idDelegatorAddress: objs[0].DelegatorAddress,
				idValidatorAddress: sample.MockAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidAddress",
				idDelegatorAddress: objs[0].ValidatorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidValidatorAddress",
				idDelegatorAddress: objs[0].DelegatorAddress,
				idValidatorAddress: objs[0].DelegatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyAddress",
				idDelegatorAddress: "",
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyValidatorAddress",
				idDelegatorAddress: objs[0].DelegatorAddress,
				idValidatorAddress: "",
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			s.Run(tc.desc, func() {
				cmd := cli.CmdShowRefund()
				var resp types.QueryGetRefundResponse
				args := append([]string{tc.idDelegatorAddress, tc.idValidatorAddress}, outflag)

				// Require the command execution returns an error with expErrMsg in its description.
				out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
				s.Require().Contains(err.Error(), tc.expErrMsg)

				// Require the output cannot be unmarshaled.
				s.Require().Error(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().Equal(r.Refund, resp.Refund)

				// Require the error returned is identified by the expected error code
				// and contains the expected error message in its description.
				stat, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expErrCode, stat.Code())
				s.Require().Contains(stat.Message(), tc.expErrMsg)
			})
		}
	})

	s.Run("Valids", func() {
		for _, obj := range objs {
			cmd := cli.CmdShowRefund()
			var resp types.QueryGetRefundResponse
			args := append([]string{obj.DelegatorAddress, obj.ValidatorAddress}, outflag)
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().NotNil(resp.Refund)
			s.Require().Equal(obj, resp.Refund)
		}
	})
}

func (s *E2EQueryTestSuite) TestCmdListRefund() {

	objs := s.refunds

	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListRefund()
			args := s.argsForPaginatedResp(nil, uint64(i), uint64(step), false)
			var resp types.QueryAllRefundResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.Refund), step)
			s.Require().Subset(objs, resp.Refund)
			s.Require().NotEmpty(resp.Refund)
		}
	})

	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListRefund()
			args := s.argsForPaginatedResp(next, 0, uint64(step), false)
			var resp types.QueryAllRefundResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.Refund), step)
			s.Require().Subset(objs, resp.Refund)
			s.Require().NotEmpty(resp.Refund)
			next = resp.Pagination.NextKey
		}
	})

	s.Run("Total", func() {
		cmd := cli.CmdListRefund()
		args := s.argsForPaginatedResp(nil, 0, uint64(len(objs)), true)
		var resp types.QueryAllRefundResponse
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
		s.Require().NoError(err)
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(objs, resp.Refund)
	})
}

func (s *E2EQueryTestSuite) TestCmdShowRefundPool() {

	objs := s.refundPools
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	s.Run("Errors", func() {
		r := types.QueryGetRefundPoolResponse{}

		for _, tc := range []struct {
			desc               string
			idValidatorAddress string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idValidatorAddress: sample.ValAddress(),
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "FailDecodingValidatorAddress",
				idValidatorAddress: sample.MockAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidValidatorAddress",
				idValidatorAddress: sample.AccAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyValidatorAddress",
				idValidatorAddress: "",
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			s.Run(tc.desc, func() {
				cmd := cli.CmdShowRefundPool()
				args := append([]string{tc.idValidatorAddress}, outflag)
				var resp types.QueryGetRefundPoolResponse

				// Require the command execution returns an error with expErrMsg in its description.
				out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
				s.Require().Contains(err.Error(), tc.expErrMsg)

				// Require the output cannot be unmarshaled.
				s.Require().Error(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().Equal(r.RefundPool, resp.RefundPool)

				// Require the error returned is identified by the expected error code
				// and contains the expected error message in its description.
				stat, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expErrCode, stat.Code())
				s.Require().Contains(stat.Message(), tc.expErrMsg)
			})
		}
	})

	s.Run("Valids", func() {
		for _, obj := range objs {
			cmd := cli.CmdShowRefundPool()
			args := append([]string{obj.OperatorAddress}, outflag)
			var resp types.QueryGetRefundPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().NotNil(resp.RefundPool)
			s.Require().Equal(obj, resp.RefundPool)
		}
	})
}

func (s *E2EQueryTestSuite) TestCmdListRefundPool() {

	objs := s.refundPools

	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListRefundPool()
			args := s.argsForPaginatedResp(nil, uint64(i), uint64(step), false)
			var resp types.QueryAllRefundPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.RefundPool), step)
			s.Require().Subset(objs, resp.RefundPool)
			s.Require().NotEmpty(resp.RefundPool)
		}
	})

	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListRefundPool()
			args := s.argsForPaginatedResp(next, 0, uint64(step), false)
			var resp types.QueryAllRefundPoolResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.RefundPool), step)
			s.Require().Subset(objs, resp.RefundPool)
			s.Require().NotEmpty(resp.RefundPool)
			next = resp.Pagination.NextKey
		}
	})

	s.Run("Total", func() {
		cmd := cli.CmdListRefundPool()
		args := s.argsForPaginatedResp(nil, 0, uint64(len(objs)), true)
		var resp types.QueryAllRefundPoolResponse
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
		s.Require().NoError(err)
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(objs, resp.RefundPool)
	})
}

func (s *E2EQueryTestSuite) TestCmdShowUnbondingDeposit() {

	objs := s.unbondingDeposits
	outflag := fmt.Sprintf("--%s=json", tmcli.OutputFlag)

	s.Run("Errors", func() {
		r := types.QueryGetUnbondingDepositResponse{}

		for _, tc := range []struct {
			desc               string
			idDepositorAddress string
			idValidatorAddress string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idDepositorAddress: sample.AccAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "FailDecodingAddress",
				idDepositorAddress: sample.MockAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "FailDecodingValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: sample.MockAddress(),
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidAddress",
				idDepositorAddress: objs[0].ValidatorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "InvalidValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[0].DepositorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyAddress",
				idDepositorAddress: "",
				idValidatorAddress: objs[0].ValidatorAddress,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "EmptyValidatorAddress",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: "",
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			s.Run(tc.desc, func() {
				cmd := cli.CmdShowUnbondingDeposit()
				args := append([]string{tc.idDepositorAddress, tc.idValidatorAddress}, outflag)
				var resp types.QueryGetUnbondingDepositResponse

				// Require the command execution returns an error with expErrMsg in its description.
				out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
				s.Require().Contains(err.Error(), tc.expErrMsg)

				// Require the output cannot be unmarshaled.
				s.Require().Error(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
				s.Require().Equal(r.UnbondingDeposit, resp.UnbondingDeposit)

				// Require the error returned is identified by the expected error code
				// and contains the expected error message in its description.
				stat, ok := status.FromError(err)
				s.Require().True(ok)
				s.Require().Equal(tc.expErrCode, stat.Code())
				s.Require().Contains(stat.Message(), tc.expErrMsg)
			})
		}
	})

	s.Run("Valids", func() {
		for _, obj := range objs {
			cmd := cli.CmdShowUnbondingDeposit()
			args := append([]string{obj.DepositorAddress, obj.ValidatorAddress}, outflag)
			var resp types.QueryGetUnbondingDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().NotNil(resp.UnbondingDeposit)
			s.Require().Equal(obj, resp.UnbondingDeposit)
		}
	})
}

func (s *E2EQueryTestSuite) TestCmdListUnbondingDeposit() {

	objs := s.unbondingDeposits

	s.Run("ByOffset", func() {
		step := 2
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListUnbondingDeposit()
			args := s.argsForPaginatedResp(nil, uint64(i), uint64(step), false)
			var resp types.QueryAllUnbondingDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.UnbondingDeposit), step)
			s.Require().Subset(objs, resp.UnbondingDeposit)
			s.Require().NotEmpty(resp.UnbondingDeposit)
		}
	})

	s.Run("ByKey", func() {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			cmd := cli.CmdListUnbondingDeposit()
			args := s.argsForPaginatedResp(next, 0, uint64(step), false)
			var resp types.QueryAllUnbondingDepositResponse
			out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
			s.Require().NoError(err)
			s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
			s.Require().LessOrEqual(len(resp.UnbondingDeposit), step)
			s.Require().Subset(objs, resp.UnbondingDeposit)
			s.Require().NotEmpty(resp.UnbondingDeposit)
			next = resp.Pagination.NextKey
		}
	})

	s.Run("Total", func() {
		cmd := cli.CmdListUnbondingDeposit()
		args := s.argsForPaginatedResp(nil, 0, uint64(len(objs)), true)
		var resp types.QueryAllUnbondingDepositResponse
		out, err := clitestutil.ExecTestCLICmd(s.ctx, cmd, args)
		s.Require().NoError(err)
		s.Require().NoError(s.cdc.UnmarshalJSON(out.Bytes(), &resp))
		s.Require().Equal(len(objs), int(resp.Pagination.Total))
		s.Require().ElementsMatch(objs, resp.UnbondingDeposit)
	})
}

func (s *E2EQueryTestSuite) argsForPaginatedResp(next []byte, offset, limit uint64, total bool) []string {
	args := []string{
		fmt.Sprintf("--%s=json", tmcli.OutputFlag),
	}
	if next == nil {
		args = append(args, fmt.Sprintf("--%s=%d", flags.FlagOffset, offset))
	} else {
		args = append(args, fmt.Sprintf("--%s=%s", flags.FlagPageKey, next))
	}
	args = append(args, fmt.Sprintf("--%s=%d", flags.FlagLimit, limit))
	if total {
		args = append(args, fmt.Sprintf("--%s", flags.FlagCountTotal))
	}
	return args
}
