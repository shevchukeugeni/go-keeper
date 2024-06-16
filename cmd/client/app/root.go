package app

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var login, password, serverURL string

var rootCmd = &cobra.Command{
	Use:   "keeper",
	Short: "keep your secrets safe",
	Long:  `keep your secrets safe`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&login, "l", "", "login for using go-keeper system")
	rootCmd.PersistentFlags().StringVar(&password, "p", "", "password for using go-keeper system")
	rootCmd.PersistentFlags().StringVar(&serverURL, "s", "", "go-keeper server address")
}
