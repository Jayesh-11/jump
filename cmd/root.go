/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jump",
	Short: "Jump is a CLI tool to help you jump into coding problems : OFFLINE and FREE",
	Long: `Jump is a CLI tool to help you jump into coding problems : OFFLINE and FREE
			You can use it to run your code locally and test it against the test cases.
			It's a great way to practice your coding skills and get better at solving problems.

			Usage:
			jump [command]

			Available Commands:
			start       Start a new coding problem
			help        Help about any command

			Flags:
			-h, --help   help for jump
			-ll, --listlanguages   list all supported languages
			-lp, --listproblems   list all supported problems
			-v, --version   show version
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		listLanguages, _ := cmd.Flags().GetBool("listlanguages")
		listProblems, _ := cmd.Flags().GetBool("listproblems")
		version, _ := cmd.Flags().GetBool("version")
		if listLanguages {
			fmt.Println("Languages:")
			fmt.Println(" - C++")
			fmt.Println(" - Go")
			fmt.Println(" - Java")
			fmt.Println(" - JavaScript")
			fmt.Println(" - Kotlin")
			fmt.Println(" - Python")
			fmt.Println(" - Rust")
			return
		}
		if listProblems {
			fmt.Println("Problems:")
			fmt.Println(" - Two Sum")
			return
		}
		if version {
			fmt.Println("Version: 1.0.0")
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("listlanguages", "l", false, "List all supported languages")
	rootCmd.Flags().BoolP("listproblems", "p", false, "List all supported problems")
	rootCmd.Flags().BoolP("version", "v", false, "Show version")
}


