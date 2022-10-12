package cli

import (
	"context"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

func CmdListUnbondingDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-unbonding-deposit",
		Short: "list all unbonding-deposit",
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
		Use:   "show-unbonding-deposit [id]",
		Short: "shows a unbonding-deposit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			params := &types.QueryGetUnbondingDepositRequest{
				Id: id,
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
