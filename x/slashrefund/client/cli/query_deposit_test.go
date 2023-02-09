package cli_test

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/made-in-block/slash-refund/testutil/network"
	"github.com/made-in-block/slash-refund/testutil/sample"
	"github.com/made-in-block/slash-refund/x/slashrefund/client/cli"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

func networkWithDepositObjects(t *testing.T, n int) (*network.Network, []types.Deposit) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		deposit := types.Deposit{
			DepositorAddress: sample.AccAddress(),
			ValidatorAddress: sample.ValAddress(),
			Shares:           sdk.NewDec(int64(100 * (i + 1))),
		}
		state.DepositList = append(state.DepositList, deposit)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.DepositList
}

func TestQueryDeposit(t *testing.T) {
	net, objs := networkWithDepositObjects(t, 5)
	ctx := net.Validators[0].ClientCtx

	request := func(next []byte, offset, limit uint64, total bool) []string {
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

	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeposit(), args)
			require.NoError(t, err)
			var resp types.QueryAllDepositResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t, objs, resp.Deposit)
		}
	})

	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeposit(), args)
			require.NoError(t, err)
			var resp types.QueryAllDepositResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.Deposit), step)
			require.Subset(t, objs, resp.Deposit)
			next = resp.Pagination.NextKey
		}
	})

	t.Run("Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListDeposit(), args)
		require.NoError(t, err)
		var resp types.QueryAllDepositResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t, objs, resp.Deposit)
	})

	// TESTS: SHOW SINGLE - VALID
	t.Run("Show-Valid", func(t *testing.T) {
		common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
		for _, tc := range []struct {
			desc               string
			idDepositorAddress string
			idValidatorAddress string
			extraArgs          []string
			obj                types.Deposit
		}{
			{
				desc:               "Found0",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[0],
			},
			{
				desc:               "Found1",
				idDepositorAddress: objs[1].DepositorAddress,
				idValidatorAddress: objs[1].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[1],
			},
			{
				desc:               "Found2",
				idDepositorAddress: objs[2].DepositorAddress,
				idValidatorAddress: objs[2].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[2],
			},
			{
				desc:               "Found3",
				idDepositorAddress: objs[3].DepositorAddress,
				idValidatorAddress: objs[3].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[3],
			},
			{
				desc:               "Found4",
				idDepositorAddress: objs[4].DepositorAddress,
				idValidatorAddress: objs[4].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[4],
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				args := append([]string{tc.idDepositorAddress, tc.idValidatorAddress}, tc.extraArgs...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDeposit(), args)
				require.NoError(t, err)
				var resp types.QueryGetDepositResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.Deposit)
				require.Equal(t, tc.obj, resp.Deposit)
			})
		}
	})

	// TESTS: SHOW SINGLE - ERRORS
	t.Run("Show-Errors", func(t *testing.T) {
		common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
		for _, tc := range []struct {
			desc               string
			idDepositorAddress string
			idValidatorAddress string
			extraArgs          []string
			expErrMsg          string
			expErrCode         codes.Code
		}{
			{
				desc:               "NotFound",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[4].ValidatorAddress,
				extraArgs:          common,
				expErrMsg:          "key not found",
				expErrCode:         codes.NotFound,
			},
			{
				desc:               "Bech32Failed-Decoding",
				idDepositorAddress: sample.MockAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "Bech32Failed-InvalidDelegator",
				idDepositorAddress: objs[0].ValidatorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "Bech32Failed-InvalidValidator",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[0].DepositorAddress,
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "Bech32Failed-EmptyDelegator",
				idDepositorAddress: "",
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "Bech32Failed-EmptyValidator",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: "",
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				args := append([]string{tc.idDepositorAddress, tc.idValidatorAddress}, tc.extraArgs...)

				// ensure the execution returns an error with expErrMsg in its description
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowDeposit(), args)
				require.Contains(t, err.Error(), tc.expErrMsg)

				// ensure the output cannot be unmarshaled
				var resp, nullresp types.QueryGetDepositResponse
				require.Error(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, nullresp.Deposit, resp.Deposit)

				// ensure the error is compatible with package grpc/status, results in status with proper code and with expErrMsg in its description
				stat, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.expErrCode, stat.Code())
				require.Contains(t, stat.Message(), tc.expErrMsg)
			})
		}
	})
}
