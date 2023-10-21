package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var deleteBucketCmd = &cobra.Command{
	Use:   "delete-bucket",
	Short: "Deletes the S3 bucket.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(cfgAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		err = storage.DeleteBucket(cmd.Flag("bucket").Value.String())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(SuccessfulMessage)
		}
	},
}

func init() {
	deleteBucketCmd.Flags().String("bucket", "", "")

	_ = deleteBucketCmd.MarkFlagRequired("bucket")

	rootCmd.AddCommand(deleteBucketCmd)
}
