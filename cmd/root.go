/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/waikco/cats-v1/app"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cats-v1",
	Short: "A sample restful API",
	Long:  `Provides cat data through a restful API, backed by multiple storage systems`,
	Run: func(cmd *cobra.Command, args []string) {
		a := &app.App{}

		router := httprouter.New()

		// add actual api routes
		router.GET("/cats/v1/health", a.Health)
		router.GET("/cats/v1/cats/:id", a.GetCat)
		router.GET("/cats/v1/cats", a.GetCats)
		//router.POST("/cats/v1/", a.CreateCat)
		//router.POST("/cats/v1/bulkcatadd", a.MassCreateCat)
		//router.PUT("/cats/v1/{id:[1-9][0-9]*|0}", a.UpdateCat)
		//router.DELETE("/cats/v1/{id:[1-9][0-9]*|0}", a.DeleteCat)

		log.Fatal(http.ListenAndServe(":8080", router))
	},
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cats-v1.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

		// Search config in home directory with name ".cats-v1" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cats-v1")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
