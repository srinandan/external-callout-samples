// Copyright 2021 Google LLC
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

package main

import (
	"fmt"
	"log"
)

func main() {
	conn, err := ldapConnectBind("ldap://localhost:389", "cn=admin,dc=example,dc=org", "Adm1n!")
	if err != nil {
		log.Fatalln(err)
	}

	result, err := ldapUserAuth(conn, "billy", "test")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("LDAP Authn Result: ", result)
	_ = ldapDisconnect(conn)
}
