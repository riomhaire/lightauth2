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
	"strconv"
	"strings"
	"time"

	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
	"github.com/spf13/cobra"
)

// sessionCmd represents the session command
var sessionCmd = &cobra.Command{
	Use:   "session",
	Short: "Create a session key suitable for accessing a resource.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		username := cmd.Flag("user").Value.String()
		roles := cmd.Flag("roles").Value.String()
		secret := cmd.Flag("secret").Value.String()
		tokenToDecode := cmd.Flag("token").Value.String()
		timeToLive, _ := strconv.Atoi(cmd.Flag("sessionPeriod").Value.String()) // Must be a better way

		if len(tokenToDecode) > 0 {
			// DECODE
			t, err := usecases.DecodeToken(tokenToDecode, secret)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("user    : %v\nexpires : %v\nroles   : %v\n", t.User, time.Unix(t.Expires, 0), t.Roles)
			}
		} else {
			// ENDCODE
			user := entities.User{}
			user.Username = username
			r := strings.Split(roles, ":")
			user.Roles = r

			token, err := usecases.EncodeToken(user, timeToLive, secret)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(token)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(sessionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sessionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sessionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sessionCmd.Flags().IntP("sessionPeriod", "k", 3600, "How long session  will be active for.")
	sessionCmd.Flags().StringP("user", "u", "anonymous", "Username associated with the token.")
	sessionCmd.Flags().StringP("roles", "r", "guest:public", "List of roles separated by ':'.")
	sessionCmd.Flags().StringP("secret", "s", "secret", "Secret used to sign things.")
	sessionCmd.Flags().StringP("token", "t", "", "If populated means decode token")

}
