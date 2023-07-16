package main

import (
	"github.com/spf13/cobra"
)

var alphaCmd = &cobra.Command{
	Use:    "alpha",
	Short:  "Pre-release functionality in alpha stage",
	Hidden: true,
}

func init() {
	rootCmd.AddCommand(alphaCmd)
}
