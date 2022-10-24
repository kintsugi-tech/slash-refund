package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

func CmdListUnbondingDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-unbonding-deposit",
		Short: "list all unbonding_deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllUnbondingDepositRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.UnbondingDepositAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowUnbondingDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-unbonding-deposit [depositor-address] [validator-address]",
		Short: "shows a unbonding_deposit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDepositorAddress := args[0]
			argValidatorAddress := args[1]

			params := &types.QueryGetUnbondingDepositRequest{
				DepositorAddress: argDepositorAddress,
				ValidatorAddress: argValidatorAddress,
			}

			res, err := queryClient.UnbondingDeposit(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
