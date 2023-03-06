package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	asyncCmd = &cobra.Command{
		Use:   "async",
		Short: "Asynchronously sends and receives time synchronization packets",
		Long: "The 'async' command has the same effect as the 'timesync' command, but" +
			"sends and receives packets asynchronously.",
		Run: func(cmd *cobra.Command, args []string) {
			output, err := executeAsync(cmd, args)
			if err != nil {
				handleError(cmd, args, err)
			}
			_, _ = fmt.Fprint(os.Stdout, output)
		},
	}
)

func init() {
	asyncCmd.Flags().IntVarP(&nPrintedHosts, "print", "p", 3,
		"The number of hosts you want to print out the results, no more than 128.")
}
