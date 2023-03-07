package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group slashrefund queries under a subcommand.
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListDeposit())
	cmd.AddCommand(CmdShowDeposit())
	cmd.AddCommand(CmdListDepositPool())
	cmd.AddCommand(CmdShowDepositPool())
	cmd.AddCommand(CmdListUnbondingDeposit())
	cmd.AddCommand(CmdShowUnbondingDeposit())
	cmd.AddCommand(CmdListRefund())
	cmd.AddCommand(CmdShowRefund())
	cmd.AddCommand(CmdListRefundPool())
	cmd.AddCommand(CmdShowRefundPool())

	return cmd
}

func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "shows the parameters of the module",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

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

func CmdListDepositPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deposit-pool",
		Short: "list all deposit_pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDepositPoolRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DepositPoolAll(context.Background(), params)
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

func CmdShowDepositPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deposit-pool [operator-address]",
		Short: "shows a deposit_pool",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argOperatorAddress := args[0]

			params := &types.QueryGetDepositPoolRequest{
				OperatorAddress: argOperatorAddress,
			}

			res, err := queryClient.DepositPool(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

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
