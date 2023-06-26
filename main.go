package main

import (
	"fmt"

	"github.com/netrixframework/redisraft-testing/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}

	rootCmd.AddCommand(cmd.FuzzerCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
