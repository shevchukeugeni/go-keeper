package app

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "sign up command",
	Long:  `sign up command`,
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.AddCommand(registerUserCmd)
}

var registerUserCmd = &cobra.Command{
	Use:   "register [login] [password]",
	Short: "register in go-keeper system",
	Long:  `register in go-keeper system`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		if args[0] == "" || args[1] == "" {
			fmt.Println("Please provide non empty login and password")
			return
		}

		if serverURL == "" {
			fmt.Println("Please specify server addr flag")
			return
		}

		data := fmt.Sprintf("{\"login\": \"%s\", \"password\": \"%s\"}", args[0], args[1])

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(data).
			Post(fmt.Sprintf("http://%s/api/user/register", serverURL))
		if err != nil {
			fmt.Println("Failed to register: ", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to register: %s\n", res.Body())
			return
		}

		fmt.Println("Successfully registered. Now you can use your creds in other commands by setting up --l and --p flags")
	},
}
