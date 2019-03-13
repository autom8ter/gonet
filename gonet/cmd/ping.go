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
	"github.com/autom8ter/util/netutil"
	"github.com/spf13/cobra"
	"net"
	"time"
)

// pingCmd represents the ping command
var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "A brief description of your command",
	PreRun: func(cmd *cobra.Command, args []string) {
		p := netutil.Pinger{
			Endpoint: "localhost:3000",
			Do: func() error {
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
				if err != nil {
					return err
				}
				return conn.Close()
			},
		}
		if err := p.Once().Do(); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("ping grpc: success!")
	},
	Run: func(cmd *cobra.Command, args []string) {
		p := netutil.Pinger{
			Endpoint: "localhost:8080",
			Do: func() error {
				conn, err := net.DialTimeout("tcp", addr, 1*time.Second)
				if err != nil {
					return err
				}
				return conn.Close()
			},
		}
		if err := p.Once().Do(); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("ping gateway: success!")
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
