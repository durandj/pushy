package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/durandj/go-pushbullet"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configNamePushbulletAPIKey = "pushbullet_api_key"
	flagNameTitle              = "title"
	flagNameBody               = "body"
)

func setupConfig() error {
	viper.SetConfigName("pushy")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")

	// No need to check the error response when the only error that
	// can be generated for this is a missing argument.
	_ = viper.BindEnv(configNamePushbulletAPIKey, configNamePushbulletAPIKey)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

func main() {
	if err := setupConfig(); err != nil {
		fmt.Printf("Unable to load configuration: %v", err)
		os.Exit(1)
	}

	cmd := cobra.Command{
		Use:   "pushy",
		Short: "Send push notifications via PushBullet",
		Run: func(cmd *cobra.Command, args []string) {
			apiKey := viper.GetString(configNamePushbulletAPIKey)
			if apiKey == "" {
				fmt.Println("No API key was set for Pushbullet")
				os.Exit(1)
			}

			cmdFlags := cmd.Flags()

			title, err := cmdFlags.GetString(flagNameTitle)
			if err != nil {
				fmt.Printf("Unable to determine title for the notification: %v", err)
				os.Exit(1)
			}

			body, err := cmdFlags.GetString(flagNameBody)
			if err != nil {
				fmt.Printf("Unable to determine body for the notification: %v", err)
				os.Exit(1)
			} else if body == "" {
				bodyBytes, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					fmt.Printf("Unable to read notification body from stdin")
					os.Exit(1)
				}

				body = string(bodyBytes)
			}

			client := pushbullet.New(apiKey)

			err = client.PushNote(pushbullet.AllDevices, title, body)
			if err != nil {
				fmt.Printf("Unable to send push notification: %v", err)
				os.Exit(1)
			}
		},
	}

	cmdFlags := cmd.Flags()

	cmdFlags.String(flagNameTitle, "", "The name of the push notification")
	// This would only return an error if the flag name was invalid.
	_ = cmd.MarkFlagRequired(flagNameTitle)

	cmdFlags.String(flagNameBody, "", "The body of the push notification")

	if err := cmd.Execute(); err != nil {
		fmt.Printf("Unable to send push notification: %v", err)
	}
}
