package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Check server connection",
	Long:  `Check server connection`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverURL == "" {
			fmt.Println("Please specify server addr")
			return
		}

		resp, err := http.Get("http://" + serverURL + "/ping")
		if err != nil {
			log.Fatalln(err)
			return
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Connection is OK")
		} else {
			fmt.Printf("Bad connection: %v\n", resp.StatusCode)
		}
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
