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

	apigee "github.com/srinandan/external-callout-samples/hello-world-callout/pkg/apigee"
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
	response.Content = "hello world!"

	req.Response = &response

	return req, nil
}
