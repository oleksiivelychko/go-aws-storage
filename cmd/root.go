package cmd

import (
	"github.com/oleksiivelychko/go-aws-storage/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

const SuccessfulMessage = "✅ Operation has been successful!"

var (
	yamlConfig string
	configAWS  *config.AWS
)

var rootCmd = &cobra.Command{
	Short: "Simple Storage Service (S3).",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&yamlConfig, "config", "config.yaml", "YAML config file")
}

func initConfig() {
	viper.SetConfigFile(yamlConfig)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		cobra.CheckErr(err)
	}

	configAWS = &config.AWS{
		Region:             viper.Get("REGION").(string),
		AwsAccessKeyId:     viper.Get("AWS_ACCESS_KEY_ID").(string),
		AwsSecretAccessKey: viper.Get("AWS_SECRET_ACCESS_KEY").(string),
		Endpoint:           viper.Get("ENDPOINT").(string),
		S3ForcePathStyle:   viper.Get("S3_FORCE_PATH_STYLE").(bool),
	}
}
