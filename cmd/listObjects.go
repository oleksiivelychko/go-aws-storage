package cmd

import (
	"fmt"
	"github.com/oleksiivelychko/go-aws-storage/service"
	"github.com/spf13/cobra"
)

var listObjectsCmd = &cobra.Command{
	Use:   "list-objects",
	Short: "Returns some or all (up to 1,000) of the objects in a bucket with each request (v2).",
	Run: func(cmd *cobra.Command, args []string) {
		service, err := service.New(configAWS)
		if err != nil {
			cobra.CheckErr(err)
		}

		output, err := service.ListObjects(cmd.Flag("bucket").Value.String())
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("%s\n%s\n", output, SuccessfulMessage)
		}
	},
}

func init() {
	listObjectsCmd.Flags().String("bucket", "", "")

	_ = listObjectsCmd.MarkFlagRequired("bucket")

	rootCmd.AddCommand(listObjectsCmd)
}
