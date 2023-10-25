package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var deleteObjectCmd = &cobra.Command{
	Use:   "delete-object",
	Short: `Removes the null version (if there is one) of an object and inserts a delete marker, which becomes the latest version of the object.`,
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(configAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = storage.DeleteObject(cmd.Flag("bucket").Value.String(), cmd.Flag("key").Value.String())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(SuccessfulMessage)
		}
	},
}

func init() {
	deleteObjectCmd.Flags().String("bucket", "", "")
	deleteObjectCmd.Flags().String("key", "", "")

	_ = deleteObjectCmd.MarkFlagRequired("bucket")
	_ = deleteObjectCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(deleteObjectCmd)
}
