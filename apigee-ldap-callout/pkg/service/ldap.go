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

package service

import (
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap"
)

var bindUser, bindPassword string
var conn *ldap.Conn

const defaultBaseDN = "dc=example,dc=org"
const defaultObjectClass = "inetOrgPerson"

func LdapConnectBind(host string, tlsConfig bool, binduser string, bindpassword string) (err error) {
	if err = ldapConnect(host, tlsConfig); err != nil {
		return err
	}
	bindUser = binduser
	bindPassword = bindpassword
	return ldapBind(bindUser, bindPassword)
}

func ldapBind(binduser string, bindpassword string) (err error) {
	if conn == nil {
		return fmt.Errorf("ldap connection not set")
	}
	return conn.Bind(binduser, bindpassword)
}

func ldapConnect(host string, tlsConfig bool) (err error) {
	if conn, err = ldap.DialURL(host); err != nil {
		return err
	}
	if tlsConfig {
		if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return err
		}
	}
	return err
}

func LdapDisconnect() (err error) {
	if conn == nil {
		return fmt.Errorf("ldap connection not set")
	}
	conn.Close()
	return nil
}

func ldapUserAuth(username string, password string, baseDN string, objectClass string, scope string) (result bool, err error) {

	if objectClass == "" {
		objectClass = defaultObjectClass
	}

	if baseDN == "" {
		baseDN = defaultBaseDN
	}

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		getScope(scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=%s)(uid=%s))", objectClass, ldap.EscapeFilter(username)),
		[]string{"dn"},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return false, err
	}

	if len(sr.Entries) != 1 {
		return false, nil
	}

	userdn := sr.Entries[0].DN

	// Bind as the user to verify their password
	if err = ldapBind(userdn, password); err != nil {
		return false, err
	}

	// Rebind as the read only user for any further queries
	if err = ldapBind(bindUser, bindPassword); err != nil {
		return false, err
	}

	return true, nil
}

func ldapSearch(baseDN string, scope string, query string, attributes []string) (sr *ldap.SearchResult, err error) {

	if baseDN == "" {
		baseDN = defaultBaseDN
	}

	if len(attributes) == 0 {
		attributes = []string{"dn"}
	}

	// Search based on attributes
	searchRequest := ldap.NewSearchRequest(
		baseDN,
		getScope(scope), ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(%s))", query),
		attributes,
		nil,
	)

	sr, err = conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	return sr, nil
}

func getScope(scope string) int {
	switch ldapScope := scope; ldapScope {
	case "object":
		return ldap.ScopeBaseObject
	case "onelevel":
		return ldap.ScopeSingleLevel
	case "subtree":
		return ldap.ScopeWholeSubtree
	default:
		return ldap.ScopeWholeSubtree
	}
}
