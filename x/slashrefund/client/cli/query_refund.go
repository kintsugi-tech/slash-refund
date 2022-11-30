package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

func CmdListRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-refund",
		Short: "list all refund",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllRefundRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.RefundAll(context.Background(), params)
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

func CmdShowRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-refund [delegator] [validator]",
		Short: "shows a refund",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argDelegatorAddress := args[0]
			argValidatorAddress := args[1]

			params := &types.QueryGetRefundRequest{
				Delegator: argDelegatorAddress,
				Validator: argValidatorAddress,
			}

			res, err := queryClient.Refund(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
