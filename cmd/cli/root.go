package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// The Root Cli Handler
var rootCmd = &cobra.Command{
	Version: "1",
}

// Execute starts the program
func Execute() {
	// Run the program
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}
