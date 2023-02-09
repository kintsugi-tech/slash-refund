package cli_test

import (
	"fmt"
	"testing"

	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/made-in-block/slash-refund/x/slashrefund/client/cli"
	"github.com/stretchr/testify/require"

	"github.com/made-in-block/slash-refund/testutil/network"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"

	tmcli "github.com/tendermint/tendermint/libs/cli"
)

func networkWithParams(t *testing.T, params types.Params) (*network.Network, types.Params) {
	t.Helper()
	cfg := network.DefaultConfig()
	state := types.GenesisState{}
	require.NoError(t, cfg.Codec.UnmarshalJSON(cfg.GenesisState[types.ModuleName], &state))

	state.Params = params
	buf, err := cfg.Codec.MarshalJSON(&state)
	require.NoError(t, err)
	cfg.GenesisState[types.ModuleName] = buf
	return network.New(t, cfg), state.Params
}

func TestQueryParams(t *testing.T) {
	net, params := networkWithParams(t,
		types.Params{
			AllowedTokens: []string{"foo"},
		})
	ctx := net.Validators[0].ClientCtx

	t.Run("Query-Valid", func(t *testing.T) {

		args := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
		out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdQueryParams(), args)
		require.NoError(t, err)

		var resp types.QueryParamsResponse
		require.NoError(t, net.Config.Codec.UnmarshalJSON(out.Bytes(), &resp))
		require.NotNil(t, resp.Params)
		require.Equal(t, params, resp.Params)
	})
}
