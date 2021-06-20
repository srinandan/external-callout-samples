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

	grpc "github.com/srinandan/external-callout-samples/apigee-ldap-callout/pkg/protocol/grpc"
	"github.com/srinandan/external-callout-samples/apigee-ldap-callout/pkg/service"
	"github.com/srinandan/sample-apps/common"
)

//Version
var Version, Git string

func usage() {
	fmt.Println("")
	fmt.Println("apigee-ldap-callout version ", Version)
	fmt.Println("")
	fmt.Println("Usage: apigee-ldap-callout --host=<ldap url> --binduser=<user> --bindpass=<passwd>")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("--baseDN=\"dc=example,dc=org\"  The base level of LDAP under which all of your data exists (default: dc=example,dc=org)")
	fmt.Println("--tls=true  Enable TLS, default is false")
}

func main() {
	//initialize logging
	common.InitLog()

	var host, binduser, bindpasswd, baseDN string
	var help, tlsConfig bool

	flag.StringVar(&host, "host", "", "LDAP URL (ex: ldap://localhost:389)")
	flag.StringVar(&binduser, "binduser", "", "LDAP bind username")
	flag.StringVar(&bindpasswd, "bindpasswd", "", "LDAP bind password")
	flag.StringVar(&baseDN, "baseDN", "", "The base level of LDAP under which all of your data exists (default: dc=example,dc=org)")
	flag.BoolVar(&tlsConfig, "tls", false, "Enable TLS")
	flag.BoolVar(&help, "help", false, "Display usage")

	// Parse commandline parameters
	flag.Parse()

	if help {
		usage()
		os.Exit(1)
	}

	if host == "" && os.Getenv("LDAP_HOST") == "" {
		usage()
		os.Exit(1)
	}

	if os.Getenv("LDAP_HOST") != "" {
		host = os.Getenv("LDAP_HOST")
	}

	if binduser == "" && os.Getenv("LDAP_BIND_USER") == "" {
		usage()
		os.Exit(1)
	}

	if os.Getenv("LDAP_BIND_USER") != "" {
		binduser = os.Getenv("LDAP_BIND_USER")
	}

	if bindpasswd == "" && os.Getenv("LDAP_BIND_PASSWD") == "" {
		usage()
		os.Exit(1)
	}

	if os.Getenv("LDAP_BIND_PASSWD") != "" {
		bindpasswd = os.Getenv("LDAP_BIND_PASSWD")
	}

	if os.Getenv("TLS") == "true" {
		tlsConfig = true
	}

	if err := service.LdapConnectBind(host, tlsConfig, binduser, bindpasswd); err != nil {
		common.Error.Println(err)
		os.Exit(1)
	}

	if err := grpc.RunServer(common.GetgRPCPort()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
