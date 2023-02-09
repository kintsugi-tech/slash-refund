package cli_test

import (
	"fmt"
	"testing"
	"time"

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

func networkWithUnbondingDepositObjects(t *testing.T, n int) (*network.Network, []types.UnbondingDeposit) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	for i := 0; i < n; i++ {
		unbondingDeposit := types.UnbondingDeposit{
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
		state.UnbondingDepositList = append(state.UnbondingDepositList, unbondingDeposit)
	}
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.UnbondingDepositList
}

func TestQueryUnbondingDeposit(t *testing.T) {
	net, objs := networkWithUnbondingDepositObjects(t, 5)
	ctx := net.Validators[0].ClientCtx

	// TESTS: LIST ALL
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
	t.Run("List-ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(objs); i += step {
			args := request(nil, uint64(i), uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListUnbondingDeposit(), args)
			require.NoError(t, err)
			var resp types.QueryAllUnbondingDepositResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t, objs, resp.UnbondingDeposit)
		}
	})
	t.Run("List-ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(objs); i += step {
			args := request(next, 0, uint64(step), false)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListUnbondingDeposit(), args)
			require.NoError(t, err)
			var resp types.QueryAllUnbondingDepositResponse
			require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
			require.LessOrEqual(t, len(resp.UnbondingDeposit), step)
			require.Subset(t, objs, resp.UnbondingDeposit)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("List-Total", func(t *testing.T) {
		args := request(nil, 0, uint64(len(objs)), true)
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdListUnbondingDeposit(), args)
		require.NoError(t, err)
		var resp types.QueryAllUnbondingDepositResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NoError(t, err)
		require.Equal(t, len(objs), int(resp.Pagination.Total))
		require.ElementsMatch(t, objs, resp.UnbondingDeposit)
	})

	// TESTS: SHOW SINGLE - VALID
	t.Run("Show-Valid", func(t *testing.T) {
		common := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
		for _, tc := range []struct {
			desc               string
			idDepositorAddress string
			idValidatorAddress string
			extraArgs          []string
			obj                types.UnbondingDeposit
		}{
			{
				desc:               "Single-Found0",
				idDepositorAddress: objs[0].DepositorAddress,
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[0],
			},
			{
				desc:               "Single-Found1",
				idDepositorAddress: objs[1].DepositorAddress,
				idValidatorAddress: objs[1].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[1],
			},
			{
				desc:               "Single-Found2",
				idDepositorAddress: objs[2].DepositorAddress,
				idValidatorAddress: objs[2].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[2],
			},
			{
				desc:               "Single-Found3",
				idDepositorAddress: objs[3].DepositorAddress,
				idValidatorAddress: objs[3].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[3],
			},
			{
				desc:               "Single-Found4",
				idDepositorAddress: objs[4].DepositorAddress,
				idValidatorAddress: objs[4].ValidatorAddress,
				extraArgs:          common,
				obj:                objs[4],
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				args := append([]string{tc.idDepositorAddress, tc.idValidatorAddress}, tc.extraArgs...)
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowUnbondingDeposit(), args)
				require.NoError(t, err)
				var resp types.QueryGetUnbondingDepositResponse
				require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.NotNil(t, resp.UnbondingDeposit)
				require.Equal(t, tc.obj, resp.UnbondingDeposit)
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
				desc:               "Bech32Failed",
				idDepositorAddress: sample.MockAddress(),
				idValidatorAddress: objs[0].ValidatorAddress,
				extraArgs:          common,
				expErrMsg:          "invalid request",
				expErrCode:         codes.InvalidArgument,
			},
			{
				desc:               "Bech32Failed-InvalidDepositor",
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
				desc:               "Bech32Failed-EmptyDepositor",
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
				out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdShowUnbondingDeposit(), args)
				require.Contains(t, err.Error(), tc.expErrMsg)

				// ensure the output cannot be unmarshaled
				var resp, nullresp types.QueryGetUnbondingDepositResponse
				require.Error(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, nullresp.UnbondingDeposit, resp.UnbondingDeposit)

				// ensure the error is compatible with package grpc/status, results in status with proper code and with expErrMsg in its description
				stat, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.expErrCode, stat.Code())
				require.Contains(t, stat.Message(), tc.expErrMsg)
			})
		}
	})
}
