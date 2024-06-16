package app

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"

	"keeper-project/internal/crypto"
	"keeper-project/types"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "easily store your files",
	Long:  `easily store your files`,
}

func init() {
	rootCmd.AddCommand(fileCmd)

	fileCmd.AddCommand(fileCreateCmd)
	fileCmd.AddCommand(filesListCmd)
	fileCmd.AddCommand(fileGetCmd)
	fileCmd.AddCommand(fileDeleteCmd)
}

var fileCreateCmd = &cobra.Command{
	Use:   "create [path] [metadata]",
	Short: "save file and metadata",
	Long:  `save file and metadata`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := resty.New()
		token, err := auth(client)
		if err != nil {
			fmt.Println(err)
			return
		}

		md, err := crypto.Encrypt(password, args[1])
		if err != nil {
			fmt.Printf("failed to encrypt: %v\n", err)
			return
		}

		res, err := client.R().
			SetFormData(map[string]string{"Metadata": md}).
			SetHeader("Authorization", token).
			SetFile("file", args[0]).
			Post(fmt.Sprintf("http://%s/api/secret/file", serverURL))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusCreated {
			fmt.Printf("Failed to save: %s\n", res.Body())
			return
		}

		fmt.Println("Successfully saved")
	},
}

var filesListCmd = &cobra.Command{
	Use:   "list",
	Short: "get saved files list",
	Long:  `get saved files list`,
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
			Get(fmt.Sprintf("http://%s/api/secret/files", serverURL))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		for i := range result {
			fmt.Printf("Name: %s, ID: %s\n", result[i].Key, result[i].Id)
		}
	},
}

var fileGetCmd = &cobra.Command{
	Use:   "get [id] [path-to-save]",
	Short: "get file by id",
	Long:  `get file by id, you can find ids in list command`,
	Args:  cobra.ExactArgs(2),
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
			SetOutput(args[1]).
			Get(fmt.Sprintf("http://%s/api/secret/file/%s", serverURL, args[0]))
		if err != nil {
			fmt.Println("Unable to save data", err)
			return
		}

		if res.StatusCode() != http.StatusOK {
			fmt.Printf("Failed to get: %s\n", res.Body())
			return
		}

		metadata, err := crypto.Decrypt(password, res.Header().Get("Meta"))
		if err != nil {
			fmt.Printf("failed to decrypt: %v\n", err)
			return
		}

		fmt.Println("Metadata: ", metadata)
	},
}

var fileDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete file by id",
	Long:  `delete file by id, you can find ids in list command`,
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
			Delete(fmt.Sprintf("http://%s/api/secret/file/%s", serverURL, args[0]))
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
