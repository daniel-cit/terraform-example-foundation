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

package state

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessState(t *testing.T) {
	jsonFile := filepath.Join(t.TempDir(), "state.json")
	content := fmt.Sprintf("{\"file\": \"%s\", \"steps\": { \"test\": { \"name\": \"test\", \"status\": \"COMPLETE\", \"error\": \"\" }}}", jsonFile)
	err := os.WriteFile(jsonFile, []byte(content), 0644)
	assert.NoError(t, err)

	s, err := LoadState(jsonFile)
	assert.NoError(t, err)
	assert.True(t, s.IsStepComplete("test"), "check if 'test' is 'COMPLETE' should be true")
	assert.False(t, s.IsStepComplete("unit"), "check if 'unit' is 'COMPLETE' should be false")

	s.CompleteStep("unit")
	assert.True(t, s.IsStepComplete("unit"), "check if 'unit' is 'COMPLETE' should be true")

	msg := "step failed"
	s.FailStep("fail", msg)
	assert.False(t, s.IsStepComplete("fail"), "check if 'fail' is 'COMPLETE' should be false")
	assert.Equal(t,s.GetStepError("fail"), msg, "step should have failed")


	assert.False(t, s.IsStepComplete("good"), "check if 'good' is 'COMPLETE' should be false")
	err = RunStepE(s, "good", func() error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, s.IsStepComplete("good"), "check if 'good' is 'COMPLETE' should be true")

	badStepMsg := "bad step"
	assert.False(t, s.IsStepComplete("bad"), "check if 'bad' is 'COMPLETE' should be false")
	err = RunStepE(s, "bad", func() error {
		return fmt.Errorf(badStepMsg)
	})
	assert.Error(t, err)
	assert.False(t, s.IsStepComplete("bad"), "check if 'bad' is 'COMPLETE' should be false")
	assert.Equal(t,s.GetStepError("bad"), badStepMsg, "step 'bad' should have failed")

	// complete states are not executed again
	assert.True(t, s.IsStepComplete("good"), "check if 'good' is 'COMPLETE' should be true")
	err = RunStepE(s, "good", func() error {
		return fmt.Errorf("will fail if executed")
	})
	assert.NoError(t, err)
	assert.True(t, s.IsStepComplete("good"), "check if 'good' is 'COMPLETE' should be true")
}
