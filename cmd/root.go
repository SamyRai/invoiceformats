package cmd

import (
	"invoiceformats/pkg/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Import subcommands directly
	"invoiceformats/cmd/generate"
	"invoiceformats/cmd/validate"
)

var (
	cfgFile string
	verbose bool
	logger logging.Logger // Use the Logger interface from pkg/logging
)

var rootCmd = &cobra.Command{
	Use:   "invoicegen",
	Short: "A modern PDF invoice generator",
	Long: `InvoiceGen is a production-ready PDF invoice generator written in Go.
	
It supports multiple output formats, template themes, currencies, and can be
used both as a CLI tool and as an HTTP API server.

Features:
• Modern HTML templates with customizable themes
• Multiple currency support with automatic conversions
• Comprehensive validation and error handling
• CLI and API modes
• Docker containerization support`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initLogger()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by the main.main(). The rootCmd will then run the correct
// handler depending on the command line arguments.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

// init function to set up the command line interface
func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Bind flags to viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Register subcommands
	rootCmd.AddCommand(generate.GenerateCmd)
	rootCmd.AddCommand(validate.ValidateCmd)
	// TODO: Add other subcommands here
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	viper.SetEnvPrefix("INVOICEGEN")
	viper.AutomaticEnv()
	// TODO: Set defaults as needed
	_ = viper.ReadInConfig()
}

// initLogger initializes the logger
func initLogger() {
	logger = logging.NewLogger()
}