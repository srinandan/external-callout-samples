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
	"io/ioutil"

	"github.com/open-policy-agent/opa/rego"
	common "github.com/srinandan/sample-apps/common"
)

var query rego.PreparedEvalQuery

func ReadRegoFile(regoFilePath string) (err error) {
	module, err := ioutil.ReadFile(regoFilePath)
	if err != nil {
		common.Error.Println("Rego file was not found")
		return err
	}
	common.Info.Println(string(module))
	return prepareForEvaluate(string(module))
}

func prepareForEvaluate(module string) (err error) {
	r := rego.New(
		rego.Query("data.httpapi.authz.allow"),
		rego.Module("policy,rego", module),
	)
	query, err = r.PrepareForEval(context.Background())
	return err
}

func EvaluatePreparedQuery(input map[string]interface{}) (resultSet rego.ResultSet, err error) {
	common.Info.Printf("input: %+v\n", input)
	resultSet, err = query.Eval(context.Background(), rego.EvalInput(input))
	common.Info.Printf("result: %+v\n", resultSet)
	return resultSet, err
}
