package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var createBucketCmd = &cobra.Command{
	Use:   "create-bucket",
	Short: "Creates a new S3 bucket.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(configAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		isPublic, err := cmd.Flags().GetBool("public")
		if err != nil {
			cobra.CheckErr(err)
		}

		output, err := storage.CreateBucket(cmd.Flag("name").Value.String(), isPublic)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s\n", output, SuccessfulMessage)
		}
	},
}

func init() {
	createBucketCmd.Flags().String("name", "", "")
	createBucketCmd.Flags().Bool("public", false, "make bucket a public (read-only)")

	_ = createBucketCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(createBucketCmd)
}
