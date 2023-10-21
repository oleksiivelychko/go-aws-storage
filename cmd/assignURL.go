package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var assignURLCmd = &cobra.Command{
	Use:   "assign-url",
	Short: "Generate a pre-signed URL for S3 object.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(cfgAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		output, err := storage.AssignURL(cmd.Flag("bucket").Value.String(), cmd.Flag("key").Value.String())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("✅ Signed URL (will expire in %d): %s\n", service.MinutesToExpireSignedURL, output)
		}
	},
}

func init() {
	assignURLCmd.Flags().String("bucket", "", "")
	assignURLCmd.Flags().String("key", "", "")

	_ = assignURLCmd.MarkFlagRequired("bucket")
	_ = assignURLCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(assignURLCmd)
}
