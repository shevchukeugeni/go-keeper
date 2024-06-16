package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"keeper-project/internal/crypto"
	"keeper-project/types"
)

var credCmd = &cobra.Command{
	Use:   "credentials",
	Short: "easily store your credentials",
	Long:  `easily store your credentials`,
}

func init() {
	rootCmd.AddCommand(credCmd)

	credCmd.AddCommand(credCreateCmd)
	credCmd.AddCommand(credsListCmd)
	credCmd.AddCommand(credGetCmd)
	credCmd.AddCommand(credDeleteCmd)
	credCmd.AddCommand(credUpdateCmd)
}

var credCreateCmd = &cobra.Command{
	Use:   "create [site] [login] [password] [metadata]",
	Short: "save credentials",
	Long:  `save credentials`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		site, err := crypto.Encrypt(password, args[0])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		lgn, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		pass, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[3])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"site\": \"%s\", \"login\": \"%s\", \"password\": \"%s\", \"metadata\": \"%s\"}",
			site, lgn, pass, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Post(fmt.Sprintf("http://%s/api/secret/cred", serverURL))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusAccepted {
			fmt.Printf("Failed to save: %s\n", res.Body())
			return
		}

		fmt.Println("Successfully saved")
	},
}

var credsListCmd = &cobra.Command{
	Use:   "list",
	Short: "get saved credentials list",
	Long:  `get saved credentials list`,
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()

		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result []*types.Key

		res, err := client.R().
			SetHeader("Authorization", token).
			SetResult(&result).
			Get(fmt.Sprintf("http://%s/api/secret/creds", serverURL))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		for i := range result {
			decrypted, err := crypto.Decrypt(password, result[i].Key)
			if err != nil {
				fmt.Printf("failed to decrypt: %v\n", err)
				return
			}
			fmt.Printf("Site: %s, ID: %s\n", decrypted, result[i].Id)
		}
	},
}

var credGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get credentials by id",
	Long:  `get credentials by id, you can find ids in list command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result types.Credentials

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetResult(&result).
			Get(fmt.Sprintf("http://%s/api/secret/cred/%s", serverURL, args[0]))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		result.Site, err = crypto.Decrypt(password, result.Site)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.Login, err = crypto.Decrypt(password, result.Login)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.Password, err = crypto.Decrypt(password, result.Password)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.Metadata, err = crypto.Decrypt(password, result.Metadata)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}

		s, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			fmt.Println("failed to print: ", err)
		}
		fmt.Println(string(s))
	},
}

var credDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete credentials by id",
	Long:  `delete credentials by id, you can find ids in list command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			Delete(fmt.Sprintf("http://%s/api/secret/cred/%s", serverURL, args[0]))
		if err != nil {
			fmt.Println("Unable to delete data", err)
			return
		}

		if res.StatusCode() != http.StatusNoContent {
			fmt.Printf("Failed to delete: %s\n", res.Body())
			return
		}

		fmt.Println("Successfully deleted")
	},
}

var credUpdateCmd = &cobra.Command{
	Use:   "update [id] [site] [login] [password] [metadata]",
	Short: "update credentials",
	Long:  `update credentials`,
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		site, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		lgn, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		pass, err := crypto.Encrypt(password, args[3])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[4])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"id\": \"%s\",\"site\": \"%s\", \"login\": \"%s\", \"password\": \"%s\", \"metadata\": \"%s\"}",
			args[0], site, lgn, pass, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Put(fmt.Sprintf("http://%s/api/secret/cred", serverURL))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to save: %s\n", res.Body())
			return
		}

		fmt.Println("Successfully updated")
	},
}
