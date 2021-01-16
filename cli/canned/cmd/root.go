package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	canVerbose  bool
	canPassword string
	canFile     string

	canFiles = []string{
		expand("~/.config/canned/default.can"),
		"/etc/canned/default.can",
	}

	rootCmd = &cobra.Command{
		Use:   "can",
		Short: "Can stores encrypted goodies",
		Run: func(cmd *cobra.Command, args []string) {
			usage(nil, []string{})
		},
	}
)

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	return nil
}

func usage(cmd *cobra.Command, text []string) {
	fmt.Println(`USAGE
    canned init                # Initializes a new can file
    canned ls                  # Lists all items
    canned set NAME [VALUE]    # Stores an item
    canned get NAME            # Copies the item's content to the clipboard
    canned rm NAME             # Removes an item
    canned tag ls              # Shows all tags
    canned tag add TAG NAME    # Adds the tag TAG to item NAME
    canned tag rm TAG NAME     # Removes the tag TAG from item NAME
    canned random [LENGTH]     # Generates a new random value
    canned env                 # Shows the environment status
    canned version             # Shows the version
    canned help                # Shows this help

GLOBAL OPTIONS
    -f, --file FILE          # Use a custom file
    -v, --verbose            # Shows verbose output

ENVIRONMENT VARIABLES
    canFile                 # Use this file instead of the default
    canPassword             # Use this password (avoids the password prompt)
    canVerbose              # Turns verbosity on
    CAN_AUTO_INIT            # Initializes the can file if it's not ready
    `)

}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.SetHelpFunc(usage)
	rootCmd.PersistentFlags().BoolVarP(&canVerbose, "verbose", "v", false, "Show verbose output")
	rootCmd.PersistentFlags().StringVarP(&canFile, "file", "f", "", "Can file path")
}

func initConfig() {
	// fmt.Println("inside initConfig")
	canFile = env("CAN_FILE", "")
	canPassword = env("CAN_PASSWORD", "")
	canVerbose, _ = strconv.ParseBool(env("CAN_VERBOSE", "false"))
}

func ensureFile() {
	if canFile != "" {
		info, err := os.Stat(canFile)
		if os.IsNotExist(err) {
			bail("Error: Can file %s does not exist.", paint("green", canFile))
		} else if info.IsDir() {
			bail("Error: Can file %s is a directory.", paint("green", canFile))
		} else {
			return
		}
	}

	for _, file := range canFiles {
		info, err := os.Stat(file)
		if os.IsNotExist(err) {
			// continue
		} else if info.IsDir() {
			// continue
		} else {
			canFile = file
			return
		}
	}

	bail("Error: Couldn't find a default can file.")
}

func ensurePassword() {
	if canPassword != "" {
		return
	}

	fmt.Printf("Enter the password: ")
	reader := bufio.NewReader(os.Stdin)
	pass, _ := reader.ReadString('\n')
	pass = strings.TrimSpace(pass)
	if pass == "" {
		fmt.Println("Error: Password can't be empty.")
		os.Exit(1)
	}

	canPassword = pass
}
