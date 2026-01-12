/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// protocCmd represents the protoc command
var protocCmd = &cobra.Command{
	Use:   "protoc [path]",
	Short: "generate",
	Long:  `simplectl protoc `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("protoc called")
	},
}

func init() {
	rootCmd.AddCommand(protocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// protocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// protocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
