package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var keys []string

var putObjectsCmd = &cobra.Command{
	Use:   "put-objects",
	Short: "Adds an object(-s) to the S3 bucket.",
	Run: func(cmd *cobra.Command, args []string) {
		storage, err := service.New(cfgAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		outCh := make(chan string, len(keys))
		errCh := make(chan error, len(keys))

		storage.PutObjectsAsync(cmd.Flag("bucket").Value.String(), keys, outCh, errCh)

		if len(errCh) > 0 {
			for err := range errCh {
				fmt.Printf("⛔️ %s", err)
			}
		}

		if len(outCh) > 0 {
			for output := range outCh {
				fmt.Printf("✅ %s", output)
			}
		}
	},
}

func init() {
	putObjectsCmd.Flags().String("bucket", "", "")
	putObjectsCmd.Flags().StringArrayVar(&keys, "key", nil, "allows multiple values")

	_ = putObjectsCmd.MarkFlagRequired("bucket")
	_ = putObjectsCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(putObjectsCmd)
}
