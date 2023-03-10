package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/made-in-block/slash-refund/x/slashrefund/types"
)

var (
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdDeposit())
	cmd.AddCommand(CmdWithdraw())
	cmd.AddCommand(CmdClaim())

	return cmd
}

// CmdDeposit implements the command to create and broadcast a MsgDeposit transaction.
func CmdDeposit() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "deposit [validator] [amount]",
		Short: "Deposit tokens for a validator.",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Deposit tokens for a validator that will be used to refund validator's delegators when the validator will be slashed.\n\n"+
				"Example:\n$ %s tx %s deposit %s 1000stake --from mykey",
			appName, types.ModuleName, validatorAddress),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidatorAddress := args[0]
			argAmount := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(argAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeposit(
				clientCtx.GetFromAddress().String(),
				argValidatorAddress,
				amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdWithdraw implements the command to create and broadcast a MsgWithdraw transaction.
func CmdWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw [validator] [amount]",
		Short: "Withdraw tokens from a deposit.",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Withdraw tokens from a deposit.\n\n"+
				"Example:\n$ %s tx %s withdraw %s 1000stake --from mykey",
			appName, types.ModuleName, validatorAddress),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidatorAddress := args[0]
			argAmount := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinNormalized(argAmount)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdraw(
				clientCtx.GetFromAddress().String(),
				argValidatorAddress,
				amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// CmdClaim implements the command to create and broadcast a MsgClaim transaction.
func CmdClaim() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim [validator]",
		Short: "Claim a refund.",
		Long: strings.TrimSpace(fmt.Sprintf(
			"Claim the refund generated for a validator's delegator.\n\n"+
				"Example:\n$ %s tx %s claim %s --from mykey",
			appName, types.ModuleName, validatorAddress),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argValidatorAddress := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgClaim(
				clientCtx.GetFromAddress().String(),
				argValidatorAddress,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
