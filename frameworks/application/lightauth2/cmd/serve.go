// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/riomhaire/lightauth2/frameworks/application/lightauth2/bootstrap"
	"github.com/spf13/cobra"
)

var cfgFile string
var projectBase string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts an authentication server",
	Long: `Starts an authentication server on the given port and using given 
	       secret based on users stored in a predefined location - by default a csv file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		application := bootstrap.Application{}

		application.Initialize(cmd, args)
		application.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	//	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lightauth2.yaml)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serveCmd.Flags().IntP("port", "p", 3030, "Default Port to Listen to.")
	serveCmd.Flags().IntP("sessionPeriod", "k", 3600, "How long session returned will be active for.")
	serveCmd.Flags().StringP("sessionSecret", "s", "secret", "Secret used to sign things.")
	serveCmd.Flags().StringP("usersFile", "u", "users.csv", "If User File used this is the one to use.")
	serveCmd.Flags().Bool("profile", false, "Enable Profiling.")
	serveCmd.Flags().String("serverCert", "server.crt", "Server SSL Cert File ")
	serveCmd.Flags().String("serverKey", "server.key", "Server SSL Key File ")
	serveCmd.Flags().Bool("useSSL", false, "Use SSL")
	serveCmd.Flags().Bool("useUserAPI", false, "Use User API (see https://githib.com/riomhaire/lightauthuserapi")
	serveCmd.Flags().String("userAPIKey", "", "The API Access Token Needed for User API")
	serveCmd.Flags().String("userAPIHost", "", "The base endpoint where service resides")
	serveCmd.Flags().StringP("loggingLevel", "l", "Debug", "Logging Level: Trace,Debug,Info,Warn,Error")
}
