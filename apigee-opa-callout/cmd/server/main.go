// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go run cmd/server/main.go
package main

import (
	"flag"
	"fmt"
	"os"

	grpc "github.com/srinandan/external-callout-samples/apigee-opa-callout/pkg/protocol/grpc"
	"github.com/srinandan/external-callout-samples/apigee-opa-callout/pkg/service"
	"github.com/srinandan/sample-apps/common"
)

//Version
var Version, Git string

func usage() {
	fmt.Println("")
	fmt.Println("apigee-opa-callout version ", Version)
	fmt.Println("")
	fmt.Println("Usage: apigee-opa-callout --regopath=<path>")
}

func main() {
	//initialize logging
	common.InitLog()

	var regoPath string
	var help bool

	flag.StringVar(&regoPath, "regopath", "policy.rego", "Rego file path")
	flag.BoolVar(&help, "help", false, "Display usage")

	// Parse commandline parameters
	flag.Parse()

	if help {
		usage()
		os.Exit(1)
	}

	if _, err := os.Stat(regoPath); err != nil {
		usage()
		os.Exit(1)
	}

	if err := service.ReadRegoFile(regoPath); err != nil {
		common.Error.Println(err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = common.GetgRPCPort()
		common.Info.Printf("Defaulting to port %s", port)
	}

	if err := grpc.RunServer(port); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
