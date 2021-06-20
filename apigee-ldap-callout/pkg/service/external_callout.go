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
	"context"
	"encoding/json"
	"fmt"

	apigeeresp "github.com/srinandan/external-callout-samples/common/apigeeresp"
	"github.com/srinandan/external-callout-samples/common/pkg/apigee"
	"github.com/srinandan/sample-apps/common"
)

type ldapAuthRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	BaseDN      string `json:"baseDN,omitempty"`
	ObjectClass string `json:"objectClass,omitempty"`
	Scope       string `json:"scope,omitempty"`
}

type ldapSearchRequest struct {
	BaseDN     string   `json:"baseDN,omitempty"`
	Scope      string   `json:"scope,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
	Query      string   `json:"query,omitempty"`
}

type ldapRequest struct {
	LdapAuthRequest   ldapAuthRequest   `json:"ldapAuthRequest,omitempty"`
	LdapSearchRequest ldapSearchRequest `json:"ldapSearchRequest,omitempty"`
}

// server is used to implement ExternalCalloutServiceServer
type ExternalCalloutServiceServer struct {
	apigee.UnimplementedExternalCalloutServiceServer
}

func NewExternalCalloutServiceServer() (apigee.ExternalCalloutServiceServer, error) {
	return &ExternalCalloutServiceServer{}, nil
}

func (c *ExternalCalloutServiceServer) ProcessMessage(ctx context.Context, req *apigee.MessageContext) (*apigee.MessageContext, error) {

	var err error
	response := apigee.Response{}
	ldapRequest := ldapRequest{}

	if len(req.AdditionalFlowVariables) == 0 {
		req.AdditionalFlowVariables = make(map[string]*apigee.FlowVariable)
	}

	if err = json.Unmarshal([]byte(req.Request.Content), &ldapRequest); err != nil {
		common.Error.Println(err)
		req.AdditionalFlowVariables["ldap-result"] = apigeeresp.SetBoolFlowVariable(false)
		response.Content = apigeeresp.GetErrorMessage(err.Error())
	} else {
		req.Request = nil
		if err = validateRequest(ldapRequest); err != nil {
			req.AdditionalFlowVariables["ldap-result"] = apigeeresp.SetBoolFlowVariable(false)
			response.Content = apigeeresp.GetErrorMessage(err.Error())
		} else {
			if ldapRequest.LdapAuthRequest != (ldapAuthRequest{}) {
				result, err := ldapUserAuth(ldapRequest.LdapAuthRequest.Username,
					ldapRequest.LdapAuthRequest.Password,
					ldapRequest.LdapAuthRequest.BaseDN,
					ldapRequest.LdapAuthRequest.ObjectClass,
					ldapRequest.LdapAuthRequest.Scope)
				if err != nil {
					common.Error.Println(err)
					req.AdditionalFlowVariables["ldap-result"] = apigeeresp.SetBoolFlowVariable(false)
					response.Content = apigeeresp.GetErrorMessage(err.Error())
				} else {
					req.AdditionalFlowVariables["ldap-result"] = apigeeresp.SetBoolFlowVariable(result)
				}
			} else {
				sr, err := ldapSearch(ldapRequest.LdapSearchRequest.BaseDN,
					ldapRequest.LdapSearchRequest.Scope,
					ldapRequest.LdapSearchRequest.Query,
					ldapRequest.LdapSearchRequest.Attributes)
				if err != nil {
					common.Error.Println(err)
					req.AdditionalFlowVariables["ldap-search-count"] = setIntFlowVariable(0)
					response.Content = apigeeresp.GetErrorMessage(err.Error())
				} else {
					resultSetBytes, _ := json.Marshal(sr)
					response.Content = string(resultSetBytes)
					req.AdditionalFlowVariables["ldap-search-count"] = setIntFlowVariable(int32(len(sr.Entries)))
				}
			}
		}
	}

	req.Response = &response
	return req, err
}

func setIntFlowVariable(value int32) *apigee.FlowVariable {
	return &apigee.FlowVariable{
		Value: &apigee.FlowVariable_Int32{value},
	}
}

func validateRequest(ldapReq ldapRequest) (err error) {
	if ldapReq.LdapAuthRequest != (ldapAuthRequest{}) && !isEmptySearchRequest(ldapReq.LdapSearchRequest) {
		return fmt.Errorf("Invalid parameters. Send oneOf LdapAuthRequest or LdapSearchRequest")
	}
	if ldapReq.LdapAuthRequest != (ldapAuthRequest{}) {
		if ldapReq.LdapAuthRequest.Username == "" {
			return fmt.Errorf("Username is mandatory for authn")
		}
		if ldapReq.LdapAuthRequest.Password == "" {
			return fmt.Errorf("Password is mandatory for authn")
		}
	}
	if !isEmptySearchRequest(ldapReq.LdapSearchRequest) {
		if len(ldapReq.LdapSearchRequest.Attributes) == 0 {
			return fmt.Errorf("At least one attribute must be sent")
		}
		if ldapReq.LdapSearchRequest.Query == "" {
			return fmt.Errorf("Query cannot be empty")
		}
	}
	return nil
}

func isEmptySearchRequest(ldapSearchRequest ldapSearchRequest) (result bool) {
	if len(ldapSearchRequest.Attributes) > 0 {
		return false
	}
	if ldapSearchRequest.BaseDN != "" {
		return false
	}
	if ldapSearchRequest.Query != "" {
		return false
	}
	if ldapSearchRequest.Scope != "" {
		return false
	}
	return true
}
