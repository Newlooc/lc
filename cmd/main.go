package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{}
	logger  = log.WithField("app", "fund-tools")
	debug   bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Error("Failed to run fund-tools")
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "debug")
}
