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

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "easily store your notes",
	Long:  `easily store your notes`,
}

func init() {
	rootCmd.AddCommand(noteCmd)

	noteCmd.AddCommand(noteCreateCmd)
	noteCmd.AddCommand(notesListCmd)
	noteCmd.AddCommand(noteGetCmd)
	noteCmd.AddCommand(noteDeleteCmd)
	noteCmd.AddCommand(noteUpdateCmd)
}

var noteCreateCmd = &cobra.Command{
	Use:   "create [title] [text] [metadata]",
	Short: "save notes with title",
	Long:  `save notes with title`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		key, err := crypto.Encrypt(password, args[0])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		text, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"key\": \"%s\", \"data\": \"%s\", \"metadata\": \"%s\"}",
			key, text, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Post(fmt.Sprintf("http://%s/api/secret/text", serverURL))
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

var notesListCmd = &cobra.Command{
	Use:   "list",
	Short: "get saved notes list",
	Long:  `get saved notes list`,
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
			Get(fmt.Sprintf("http://%s/api/secret/texts", serverURL))
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
			fmt.Printf("Title: %s, ID: %s\n", decrypted, result[i].Id)
		}
	},
}

var noteGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get note by id",
	Long:  `get note by id, you can find ids in list command`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		var result types.Note

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetResult(&result).
			Get(fmt.Sprintf("http://%s/api/secret/text/%s", serverURL, args[0]))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		result.Key, err = crypto.Decrypt(password, result.Key)
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}
		result.Text, err = crypto.Decrypt(password, result.Text)
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

var noteDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete note by id",
	Long:  `delete note by id, you can find ids in list command`,
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
			Delete(fmt.Sprintf("http://%s/api/secret/text/%s", serverURL, args[0]))
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

var noteUpdateCmd = &cobra.Command{
	Use:   "update [id] [title] [text] [metadata]",
	Short: "update note",
	Long:  `update note`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		title, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		text, err := crypto.Encrypt(password, args[2])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}
		md, err := crypto.Encrypt(password, args[3])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		data := fmt.Sprintf("{\"id\": \"%s\",\"key\": \"%s\", \"data\": \"%s\", \"metadata\": \"%s\"}",
			args[0], title, text, md)

		res, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Authorization", token).
			SetBody(data).
			Put(fmt.Sprintf("http://%s/api/secret/text", serverURL))
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
