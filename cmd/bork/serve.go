package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"bin/bork/pkg/server/httpserver"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the bork API",
	Long:  `Serve the bork API`,
	Run: func(cmd *cobra.Command, args []string) {
		config := viper.New()
		config.AutomaticEnv()
		fmt.Println("Serving the bork application")
		httpserver.Serve(config)
	},
}
