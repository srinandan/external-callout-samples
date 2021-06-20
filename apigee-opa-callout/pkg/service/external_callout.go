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

	apigeeresp "github.com/srinandan/external-callout-samples/common/apigeeresp"
	"github.com/srinandan/external-callout-samples/common/pkg/apigee"
	"github.com/srinandan/sample-apps/common"
)

// server is used to implement ExternalCalloutServiceServer
type ExternalCalloutServiceServer struct {
	apigee.UnimplementedExternalCalloutServiceServer
}

func NewExternalCalloutServiceServer() (apigee.ExternalCalloutServiceServer, error) {
	return &ExternalCalloutServiceServer{}, nil
}

func (c *ExternalCalloutServiceServer) ProcessMessage(ctx context.Context, req *apigee.MessageContext) (*apigee.MessageContext, error) {

	response := apigee.Response{}
	var regoInput map[string]interface{}
	var err error

	if len(req.AdditionalFlowVariables) == 0 {
		req.AdditionalFlowVariables = make(map[string]*apigee.FlowVariable)
	}

	common.Info.Println("Request: ", req.Request.Content)

	if err = json.Unmarshal([]byte(req.Request.Content), &regoInput); err != nil {
		common.Error.Println(err)
		req.Request = nil
		req.AdditionalFlowVariables["opa-result"] = apigeeresp.SetBoolFlowVariable(false)
		response.Content = apigeeresp.GetErrorMessage(err.Error())
	} else {
		req.Request = nil
		rs, err := EvaluatePreparedQuery(regoInput)
		if err != nil {
			common.Error.Println(err)
			req.AdditionalFlowVariables["opa-result"] = apigeeresp.SetBoolFlowVariable(false)
			response.Content = apigeeresp.GetErrorMessage(err.Error())
		} else {
			req.AdditionalFlowVariables["opa-result"] = apigeeresp.SetBoolFlowVariable(rs[0].Expressions[0].Value.(bool))
			resultSetBytes, _ := json.Marshal(rs)
			response.Content = string(resultSetBytes)
		}
	}

	req.Response = &response
	return req, err
}
