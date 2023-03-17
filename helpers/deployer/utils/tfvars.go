// Copyright 2023 Google LLC
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

package utils

import (
	"errors"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func ReadTfvars(filename string, val interface{}) error {
	data, diagnostic := hclparse.NewParser().ParseHCLFile(filename)
	if diagnostic.HasErrors() {
		return errors.New(diagnostic.Error())
	}
	decoded := gohcl.DecodeBody(data.Body, nil, val)
	if decoded.HasErrors() {
		return errors.New(decoded.Error())
	}
	return nil
}

func WriteTfvars(filename string, val interface{}) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(val, f.Body())
	return os.WriteFile(filename, f.Bytes(), 0644)
}
