package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

func CmdListDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deposit",
		Short: "list all deposit",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDepositRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DepositAll(context.Background(), params)
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

func CmdShowDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deposit [address] [validator-address]",
		Short: "shows a deposit",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDepositorAddress := args[0]
			argValidatorAddress := args[1]

			params := &types.QueryGetDepositRequest{
				DepositorAddress: argDepositorAddress,
				ValidatorAddress: argValidatorAddress,
			}

			res, err := queryClient.Deposit(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
