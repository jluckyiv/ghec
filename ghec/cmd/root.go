/*
Copyright © 2023 Jackson Lucky <jack@jacksonlucky.net>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jluckyiv/ghec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile              string
	numTargets           int
	level                int
	previousEnhancements int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghec",
	Short: "Gloomhaven Enhancement Calculator",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

// run is a helper function for the subcommands, which are similar.
func run(be ghec.BaseEnhancement, desc string) {
	l := ghec.Level(level)
	pe := ghec.PreviousEnhancements(previousEnhancements)
	e := ghec.
		NewEnhancement(be,
			ghec.OptionWithLevel(l),
			ghec.OptionWithMultipleTarget(numTargets),
			ghec.OptionWithPreviousEnhancements(ghec.PreviousEnhancements(pe)),
		)
	cost, err := e.Cost()
	cobra.CheckErr(err)
	fmt.Printf("%s costs %d", desc, cost)
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
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ghec.yaml)")
	rootCmd.PersistentFlags().IntVarP(&numTargets, "targets", "t", 1, "number of current targets")
	rootCmd.PersistentFlags().IntVarP(&level, "level", "l", 1, "ability card level")
	rootCmd.PersistentFlags().IntVarP(&previousEnhancements, "previous", "p", 0, "number of previous enhancements")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".ghec" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".ghec")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
