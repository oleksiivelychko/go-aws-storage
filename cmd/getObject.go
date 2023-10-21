package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var getObjectCmd = &cobra.Command{
	Use:   "get-object",
	Short: "Retrieves objects from the S3 bucket.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(cfgAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = storage.GetObject(
			cmd.Flag("bucket").Value.String(),
			cmd.Flag("key").Value.String(),
			cmd.Flag("path").Value.String(),
		)

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(SuccessfulMessage)
		}
	},
}

func init() {
	getObjectCmd.Flags().String("bucket", "", "")
	getObjectCmd.Flags().String("key", "", "")
	getObjectCmd.Flags().String("path", "", "path to save")

	_ = getObjectCmd.MarkFlagRequired("bucket")
	_ = getObjectCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(getObjectCmd)
}
