package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/made-in-block/slash-refund/x/slashrefund/types"
	"github.com/spf13/cobra"
)

var (
	accountAddress   = "cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p"
	validatorAddress = "cosmosvaloper1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj"
	appName          = "<appd>"
)

// GetQueryCmd returns the cli query commands for this module.
func GetQueryCmd(queryRoute string) *cobra.Command {

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

// CmdQueryParams implements the command to query the module parameters.
func CmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Shows parameters",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show the parameters of the %s module.\n\n"+
				"Example:\n$ %s query %s params\n",
			types.ModuleName, appName, types.ModuleName),
		),
		Args: cobra.NoArgs,
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

// CmdListDeposit implements the command to query all deposits.
func CmdListDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deposit",
		Short: "List all deposits",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show all deposits.\n\n"+
				"Example:\n$ %s query %s list-deposit",
			appName, types.ModuleName),
		),
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

// CmdShowDeposit implements the command to query a single deposit made from an address
// to a validator.
func CmdShowDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deposit [address] [validator]",
		Short: "Show a single deposit",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show a single deposit based on address and validator address.\n\n"+
				"Example:\n$ %s query %s show-deposit %s %s",
			appName, types.ModuleName, accountAddress, validatorAddress),
		),
		Args: cobra.ExactArgs(2),
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

// CmdListDepositPool implements the command to query all deposit pools.
func CmdListDepositPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-deposit-pool",
		Short: "List all deposit pools",
		Long: strings.TrimSpace(fmt.Sprintf(
			"List all deposit pools.\n\n"+
				"Example:\n$ %s query %s list-deposit-pool",
			appName, types.ModuleName),
		),
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

// CmdShowDepositPool implements the command to query a single deposit pool of a
// specific validator.
func CmdShowDepositPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-deposit-pool [validator]",
		Short: "Show a single deposit pool",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show a single deposit pool based validator address.\n\n"+
				"Example:\n$ %s query %s show-deposit-pool %s",
			appName, types.ModuleName, validatorAddress),
		),
		Args: cobra.ExactArgs(1),
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

// CmdListUnbondingDeposit implements the command to query all unbonding deposits.
func CmdListUnbondingDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-unbonding-deposit",
		Short: "List all unbonding deposits",
		Long: strings.TrimSpace(fmt.Sprintf(
			"List all unbonding deposits.\n\n"+
				"Example:\n$ %s query %s list-unbonding-deposit",
			appName, types.ModuleName),
		),
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

// CmdShowUnbondingDeposit implements the command to query a single unbonding deposit
// made from an address to a validator.
func CmdShowUnbondingDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-unbonding-deposit [address] [validator]",
		Short: "Show an unbonding deposit",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show a single unbonding deposit based on address and validator address.\n\n"+
				"Example:\n$ %s query %s show-unbonding-deposit %s %s",
			appName, types.ModuleName, accountAddress, validatorAddress),
		),
		Args: cobra.ExactArgs(2),
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

// CmdListRefund implements the command to query all refunds.
func CmdListRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-refund",
		Short: "List all refunds",
		Long: strings.TrimSpace(fmt.Sprintf(
			"List all refunds.\n\n"+
				"Example:\n$ %s query %s list-refund",
			appName, types.ModuleName),
		),
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

// CmdShowRefund implements the command to query a refund generated for a delegator
// of a validator.
func CmdShowRefund() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-refund [address] [validator]",
		Short: "Show a single refund",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show a single refund based on address and validator address.\n\n"+
				"Example:\n$ %s query %s show-refund %s %s",
			appName, types.ModuleName, accountAddress, validatorAddress),
		),
		Args: cobra.ExactArgs(2),
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

// CmdListRefundPool implements the command to query all refund pools.
func CmdListRefundPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-refund-pool",
		Short: "List all refund pools",
		Long: strings.TrimSpace(fmt.Sprintf(
			"List all refund pools.\n\n"+
				"Example:\n$ %s query %s list-refund-pool",
			appName, types.ModuleName),
		),
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

// CmdShowRefundPool implements the command to query a single refund pool of a specific
// validator.
func CmdShowRefundPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-refund-pool [validator]",
		Short: "Show a single refund pool",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Show a single refund pool based validator address.\n\n"+
				"Example:\n$ %s query %s show-refund-pool %s",
			appName, types.ModuleName, validatorAddress),
		),
		Args: cobra.ExactArgs(1),
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
