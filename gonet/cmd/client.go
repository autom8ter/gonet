// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
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
	"github.com/autom8ter/gonet"
	"github.com/autom8ter/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "gonet client cli",
	Run: func(cmd *cobra.Command, args []string) {
		client := gonet.NewClient(uRL, method, user, password, body, headers, formVals)
		resp, err := client.Do()
		util.NewErrCfg("client do", err).FailIfErr()
		body, err := client.ReadBody(resp)
		util.NewErrCfg("client read body", err).FailIfErr()
		fmt.Println(string(body))
	},
}

func init() {
	clientCmd.LocalFlags().StringVar(&uRL, "url", "http://localhost:8080/v1/echo", "request target url")
	clientCmd.LocalFlags().StringVar(&method, "method", "post", "request method")
	clientCmd.LocalFlags().StringToStringVar(&headers, "headers", nil, "request headers")
	clientCmd.LocalFlags().StringToStringVar(&formVals, "form", nil, "request form values")
	clientCmd.LocalFlags().StringVar(&user, "user", "", "request username")
	clientCmd.LocalFlags().StringVar(&password, "pass", "", "request password")
	clientCmd.LocalFlags().StringVar(&body, "body", `{"say": "hello there"}`, "request body")

	_ = viper.BindPFlags(clientCmd.LocalFlags())

}
