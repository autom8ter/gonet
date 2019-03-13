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
	"context"
	"github.com/autom8ter/gonet"
	"github.com/autom8ter/gonet/config"
	"github.com/autom8ter/gonet/testing/gen/echo"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
)

var grpcGatewayConfig = &gonet.GrpcGatewayConfig{
	EnvPrefix:    "",
	DialOptions:  []grpc.DialOption{grpc.WithInsecure()},
	RegisterFunc: echopb.RegisterEchoServiceHandlerFromEndpoint,
}

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A brief description of your command",
	PreRun: func(cmd *cobra.Command, args []string) {
		go func() {
			if err := echopb.NewEchoServer().Serve(grpc.NewServer()); err != nil {
				log.Fatal(err.Error())
			}
		}()
		config.SetupViper("")
	},
	Run: func(cmd *cobra.Command, args []string) {
		gw := gonet.NewGrpcGateway(context.Background(), grpcGatewayConfig, gonet.NewRouter(addr))
		gw.WithDebug()
		gw.Serve()
	},
}

func init() {
	grpcCmd.LocalFlags().StringVar(&addr, "address", ":8080", "address to listen on")
}
