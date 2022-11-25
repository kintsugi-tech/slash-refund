package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

func CmdListRefundPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-refund-pool",
		Short: "list all refund_pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllRefundPoolRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RefundPoolAll(context.Background(), params)
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

func CmdShowRefundPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-refund-pool [operator-address]",
		Short: "shows a refund_pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argOperatorAddress := args[0]

			params := &types.QueryGetRefundPoolRequest{
				OperatorAddress: argOperatorAddress,
			}

			res, err := queryClient.RefundPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
