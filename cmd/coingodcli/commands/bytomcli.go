package commands

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"coingod/util"
)

// coingodcli usage template
var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletDisable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

  available with wallet enable:
    {{range .Commands}}{{if (and .IsAvailableCommand (.Name | WalletEnable))}}
    {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

// commandError is an error used to signal different error situations in command handling.
type commandError struct {
	s         string
	userError bool
}

func (c commandError) Error() string {
	return c.s
}

func (c commandError) isUserError() bool {
	return c.userError
}

func newUserError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: true}
}

func newSystemError(a ...interface{}) commandError {
	return commandError{s: fmt.Sprintln(a...), userError: false}
}

func newSystemErrorF(format string, a ...interface{}) commandError {
	return commandError{s: fmt.Sprintf(format, a...), userError: false}
}

// Catch some of the obvious user errors from Cobra.
// We don't want to show the usage message for every error.
// The below may be to generic. Time will show.
var userErrorRegexp = regexp.MustCompile("argument|flag|shorthand")

func isUserError(err error) bool {
	if cErr, ok := err.(commandError); ok && cErr.isUserError() {
		return true
	}

	return userErrorRegexp.MatchString(err.Error())
}

// CoingodcliCmd is Coingodcli's root command.
// Every other command attached to CoingodcliCmd is a child command to it.
var CoingodcliCmd = &cobra.Command{
	Use:   "coingodcli",
	Short: "Coingodcli is a commond line client for coingod core (a.k.a. coingodd)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.SetUsageTemplate(usageTemplate)
			cmd.Usage()
		}
	},
}

// Execute adds all child commands to the root command CoingodcliCmd and sets flags appropriately.
func Execute() {

	AddCommands()
	AddTemplateFunc()

	if _, err := CoingodcliCmd.ExecuteC(); err != nil {
		os.Exit(util.ErrLocalExe)
	}
}

// AddCommands adds child commands to the root command CoingodcliCmd.
func AddCommands() {
	CoingodcliCmd.AddCommand(createAccessTokenCmd)
	CoingodcliCmd.AddCommand(listAccessTokenCmd)
	CoingodcliCmd.AddCommand(deleteAccessTokenCmd)
	CoingodcliCmd.AddCommand(checkAccessTokenCmd)

	CoingodcliCmd.AddCommand(createAccountCmd)
	CoingodcliCmd.AddCommand(deleteAccountCmd)
	CoingodcliCmd.AddCommand(listAccountsCmd)
	CoingodcliCmd.AddCommand(updateAccountAliasCmd)
	CoingodcliCmd.AddCommand(createAccountReceiverCmd)
	CoingodcliCmd.AddCommand(listAddressesCmd)
	CoingodcliCmd.AddCommand(validateAddressCmd)
	CoingodcliCmd.AddCommand(listPubKeysCmd)

	CoingodcliCmd.AddCommand(createAssetCmd)
	CoingodcliCmd.AddCommand(getAssetCmd)
	CoingodcliCmd.AddCommand(listAssetsCmd)
	CoingodcliCmd.AddCommand(updateAssetAliasCmd)

	CoingodcliCmd.AddCommand(getTransactionCmd)
	CoingodcliCmd.AddCommand(listTransactionsCmd)

	CoingodcliCmd.AddCommand(getUnconfirmedTransactionCmd)
	CoingodcliCmd.AddCommand(listUnconfirmedTransactionsCmd)
	CoingodcliCmd.AddCommand(decodeRawTransactionCmd)

	CoingodcliCmd.AddCommand(listUnspentOutputsCmd)
	CoingodcliCmd.AddCommand(listBalancesCmd)

	CoingodcliCmd.AddCommand(rescanWalletCmd)
	CoingodcliCmd.AddCommand(walletInfoCmd)

	CoingodcliCmd.AddCommand(buildTransactionCmd)
	CoingodcliCmd.AddCommand(signTransactionCmd)
	CoingodcliCmd.AddCommand(submitTransactionCmd)
	CoingodcliCmd.AddCommand(estimateTransactionGasCmd)

	CoingodcliCmd.AddCommand(getBlockCountCmd)
	CoingodcliCmd.AddCommand(getBlockHashCmd)
	CoingodcliCmd.AddCommand(getBlockCmd)
	CoingodcliCmd.AddCommand(getBlockHeaderCmd)

	CoingodcliCmd.AddCommand(createKeyCmd)
	CoingodcliCmd.AddCommand(deleteKeyCmd)
	CoingodcliCmd.AddCommand(listKeysCmd)
	CoingodcliCmd.AddCommand(updateKeyAliasCmd)
	CoingodcliCmd.AddCommand(resetKeyPwdCmd)
	CoingodcliCmd.AddCommand(checkKeyPwdCmd)

	CoingodcliCmd.AddCommand(signMsgCmd)
	CoingodcliCmd.AddCommand(verifyMsgCmd)
	CoingodcliCmd.AddCommand(decodeProgCmd)

	CoingodcliCmd.AddCommand(createTransactionFeedCmd)
	CoingodcliCmd.AddCommand(listTransactionFeedsCmd)
	CoingodcliCmd.AddCommand(deleteTransactionFeedCmd)
	CoingodcliCmd.AddCommand(getTransactionFeedCmd)
	CoingodcliCmd.AddCommand(updateTransactionFeedCmd)

	CoingodcliCmd.AddCommand(netInfoCmd)
	CoingodcliCmd.AddCommand(gasRateCmd)

	CoingodcliCmd.AddCommand(versionCmd)
}

// AddTemplateFunc adds usage template to the root command CoingodcliCmd.
func AddTemplateFunc() {
	walletEnableCmd := []string{
		createAccountCmd.Name(),
		listAccountsCmd.Name(),
		deleteAccountCmd.Name(),
		updateAccountAliasCmd.Name(),
		createAccountReceiverCmd.Name(),
		listAddressesCmd.Name(),
		validateAddressCmd.Name(),
		listPubKeysCmd.Name(),

		createAssetCmd.Name(),
		getAssetCmd.Name(),
		listAssetsCmd.Name(),
		updateAssetAliasCmd.Name(),

		createKeyCmd.Name(),
		deleteKeyCmd.Name(),
		listKeysCmd.Name(),
		resetKeyPwdCmd.Name(),
		checkKeyPwdCmd.Name(),
		signMsgCmd.Name(),

		buildTransactionCmd.Name(),
		signTransactionCmd.Name(),

		getTransactionCmd.Name(),
		listTransactionsCmd.Name(),
		listUnspentOutputsCmd.Name(),
		listBalancesCmd.Name(),

		rescanWalletCmd.Name(),
		walletInfoCmd.Name(),
	}

	cobra.AddTemplateFunc("WalletEnable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return true
			}
		}
		return false
	})

	cobra.AddTemplateFunc("WalletDisable", func(cmdName string) bool {
		for _, name := range walletEnableCmd {
			if name == cmdName {
				return false
			}
		}
		return true
	})
}
