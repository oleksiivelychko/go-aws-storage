package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var listBucketsCmd = &cobra.Command{
	Use:   "list-buckets",
	Short: "Returns a list of all buckets owned by this IAM user.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(cfgAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		output, err := storage.ListBuckets()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s\n", output, SuccessfulMessage)
		}
	},
}

func init() {
	rootCmd.AddCommand(listBucketsCmd)
}
