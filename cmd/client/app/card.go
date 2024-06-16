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

var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "easily store your card info",
	Long:  `easily store your card info`,
}

func init() {
	rootCmd.AddCommand(cardCmd)

	cardCmd.AddCommand(cardCreateCmd)
	cardCmd.AddCommand(cardsListCmd)
	cardCmd.AddCommand(cardGetCmd)
	cardCmd.AddCommand(cardDeleteCmd)
	cardCmd.AddCommand(cardUpdateCmd)
}

var cardCreateCmd = &cobra.Command{
	Use:   "create [number] [expiration] [cvv] [metadata]",
	Short: "save card information",
	Long:  `save card information`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		number, err := crypto.Encrypt(password, args[0])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		exp, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		cvv, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[3])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"number\": \"%s\", \"expiration\": \"%s\", \"cvv\": \"%s\", \"metadata\": \"%s\"}",
			number, exp, cvv, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Post(fmt.Sprintf("http://%s/api/secret/card", serverURL))
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

var cardsListCmd = &cobra.Command{
	Use:   "list",
	Short: "get saved cards list",
	Long:  `get saved cards list`,
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
			Get(fmt.Sprintf("http://%s/api/secret/cards", serverURL))
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
			result[i].Key = "*" + decrypted[len(decrypted)-4:]
			fmt.Printf("Number: %s, ID: %s\n", result[i].Key, result[i].Id)
		}
	},
}

var cardGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get card info by id",
	Long:  `get card info by id, you can find ids in list command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result types.CardInfo

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetResult(&result).
			Get(fmt.Sprintf("http://%s/api/secret/card/%s", serverURL, args[0]))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		result.Number, err = crypto.Decrypt(password, result.Number)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.Expiration, err = crypto.Decrypt(password, result.Expiration)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.CVV, err = crypto.Decrypt(password, result.CVV)
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

var cardDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete card info by id",
	Long:  `delete card info by id, you can find ids in list command`,
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
			Delete(fmt.Sprintf("http://%s/api/secret/card/%s", serverURL, args[0]))
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

var cardUpdateCmd = &cobra.Command{
	Use:   "update [id] [number] [expiration] [cvv] [metadata]",
	Short: "update card information",
	Long:  `update card information`,
	Args:  cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		number, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		exp, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		cvv, err := crypto.Encrypt(password, args[3])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[4])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"id\": \"%s\",\"number\": \"%s\", \"expiration\": \"%s\", \"cvv\": \"%s\", \"metadata\": \"%s\"}",
			args[0], number, exp, cvv, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Put(fmt.Sprintf("http://%s/api/secret/card", serverURL))
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
