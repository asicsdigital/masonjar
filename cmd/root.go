// Copyright © 2018 Steve Huff <steve.huff@asics.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile, logFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "masonjar",
	Short: "A tool for provisioning canned workflows",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/masonjar/masonjar.yaml)")

	rootCmd.PersistentFlags().StringVar(&logFile, "logfile", "", "log file (default is $HOME/.config/masonjar/masonjar.log)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	viper.BindPFlag("IsVerbose", rootCmd.PersistentFlags().Lookup("verbose"))
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

		// Search config in home directory with name ".masonjar" (without extension).
		cfgPath := filepath.Join(home, ".config", "masonjar")
		cfgFile = filepath.Join(cfgPath, "masonjar.yaml")
		viper.AddConfigPath(cfgPath)
		viper.SetConfigName("masonjar")
	}

	viper.BindPFlag("CfgFile", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("LogFile", rootCmd.PersistentFlags().Lookup("logfile"))

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// now that the config is read, derive the homedir
	setMasonjarHomedir(viper.GetString("CfgFile"))

	// configure logging
	viper.SetDefault("LogFile", FilenameInHomedir("masonjar.log"))
	initLogging(viper.GetString("LogFile"))

	// derive repository path
	viper.Set("RepoDir", FilenameInHomedir("repo"))
}

func initLogging(logFile string) {
	jww.SetLogOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    8,
		MaxBackups: 3,
		Compress:   true,
	})

	jww.INFO.Printf("configured logging to LogFile: %v", logFile)

	if viper.GetBool("IsVerbose") {
		jww.SetLogThreshold(jww.LevelTrace)
		jww.SetStdoutThreshold(jww.LevelDebug)
	}
}

func FilenameInHomedir(fileName string) string {
	homeDir := viper.GetString("HomeDir")

	if len(homeDir) == 0 {
		jww.ERROR.Panic("FilenameInHomedir() called before setMasonjarHomedir()")
	}

	return filepath.Join(homeDir, fileName)
}

func setMasonjarHomedir(cfgFile string) string {
	dir, _ := filepath.Split(cfgFile)
	viper.Set("HomeDir", dir)
	return viper.GetString("HomeDir")
}
