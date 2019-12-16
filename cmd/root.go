package cmd

import (
	"fmt"
	"os"

	"github.com/jaxenlau/pagenor-go/log"
	"github.com/jaxenlau/pagenor-go/services"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "pagenor",
	Short: "pagenor",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pagenor.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".pagenor" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".pagenor")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

type ApplicationOptions struct {
	Pagenor services.PagenorOptions `mapstructure:"pagenor" yaml:"pagenor"`
}

func (s *ApplicationOptions) Load() {
	err := viper.Unmarshal(s)
	if err != nil {
		log.DefaultLogger.WithError(err).Fatal("failed to parse config file")
	}
}

func loadApplicationOptions() ApplicationOptions {
	opts := ApplicationOptions{}
	opts.Load()
	return opts
}

func handleInitError(module string, err error) {
	if err == nil {
		return
	}
	log.DefaultLogger.WithError(err).Fatalf("init %s failed", module)
}
