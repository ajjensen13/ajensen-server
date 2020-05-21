/*
Copyright © 2020 A. Jensen <jensen.aaro@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/ajjensen13/ajensen-server/internal/projects"
	"github.com/ajjensen13/ajensen-server/internal/tags"
)

var (
	cfgFile string
	logger  *log.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ajensen-server",
	Short: `backend server for ajensen-client`,
	Long:  `backend server for ajensen-client`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger = log.New(log.Writer(), log.Prefix(), log.Flags())
		gin.DefaultWriter = logger.Writer()
		if debug {
			gin.SetMode(gin.DebugMode)
			logger.Printf("running in debug mode")
		} else {
			gin.SetMode(gin.ReleaseMode)
			logger.Printf("running in release mode")
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		r := gin.Default()

		err := projects.Init(logger, r, projectsDir)
		if err != nil {
			return err
		}

		r.Use(cors.Default()) // this is public data, allow anyone to access it

		err = tags.Init(logger, r, tagsDir)
		if err != nil {
			return err
		}

		return r.Run(addr)
	},
}

var (
	projectsDir string
	tagsDir     string
	addr        string
	debug       bool
)

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
	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&debug, "debug", "d", false, "true to enable debug mode")
	flags.StringVarP(&projectsDir, "projects", "p", "assets/projects", "directory to load projects from")
	flags.StringVarP(&tagsDir, "tags", "t", "assets/tags", "directory to load tags from")
	flags.StringVar(&addr, "addr", ":http", "address to listen on")
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

		// Search config in home directory with name ".ajensen-projects" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ajensen-projects")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
