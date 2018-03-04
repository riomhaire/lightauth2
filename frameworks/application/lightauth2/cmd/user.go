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

	"github.com/riomhaire/lightauth2/usecases"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Create a user record and signed/encrypted password info for including in your user store.",
	Long:  ` `,
	Run: func(cmd *cobra.Command, args []string) {

		username := cmd.Flag("user").Value.String()
		password := cmd.Flag("password").Value.String()
		roles := cmd.Flag("roles").Value.String()
		claim1 := cmd.Flag("claim1").Value.String()
		claim2 := cmd.Flag("claim2").Value.String()

		if len(password) == 0 {
			fmt.Println("Password Parameter cannot be empty")
			return
		}
		hash := usecases.HashPassword(password, fmt.Sprintf("%v%v", username, password))
		fmt.Printf("%v,%v,true,%v,%v,%v\n", username, hash, roles, claim1, claim2)

	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// userCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	userCmd.Flags().StringP("user", "u", "anonymous", "Username associated with the token.")
	userCmd.Flags().StringP("password", "p", "", "Password (in the raw - will be encoded).")
	userCmd.Flags().StringP("roles", "r", "guest:public", "List of roles separated by ':'.")
	userCmd.Flags().String("claim1", "", "Security Claim - Eg QRCode Hash")
	userCmd.Flags().String("claim2", "", "Security Claim - Eg API Key")
}
