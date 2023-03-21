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
	"encoding/json"
	"fmt"
	"os"
)

const (
	completeState = "COMPLETE"
	errorState    = "ERROR"
)

type Step struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error"`
}

type State struct {
	File  string          `json:"file"`
	Steps map[string]Step `json:"steps"`
}

func LoadState(file string) (State, error) {
	var s State
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Printf("creating new state file '%s'\n.", file)
		s = State{
			File: file,
		}
	} else {
		f, err := os.ReadFile(file)
		if err != nil {
			return s, err
		}
		err = json.Unmarshal(f, &s)
		if err != nil {
			return s, err
		}
		s.File = file
	}
	if s.Steps == nil {
		s.Steps = map[string]Step{}
	}
	return s, nil
}

 func (s State) SaveState() {
	f, _ := json.Marshal(s)
	os.WriteFile(s.File, f, 0644)
}

 func (s State) CompleteStep(name string) {
	s.Steps[name] = Step{
		Name:   name,
		Status: completeState,
		Error:  "",
	}
	s.SaveState()
	fmt.Printf("completing step '%s' execution\n", name)
}

 func (s State) IsStepComplete(name string) bool {
	v, ok := s.Steps[name]
	if ok {
		return v.Status == completeState
	}
	return false
}

 func (s State) FailStep(name string, err string) {
	s.Steps[name] = Step{
		Name:   name,
		Status: errorState,
		Error:  err,
	}
	s.SaveState()
	fmt.Printf("failing step '%s'. Failed with error: %s\n", name, err) // TODO make message shorter
}

 func (s State) GetStepError(name string) string {
	v, ok := s.Steps[name]
	if ok {
		return v.Error
	}
	return ""
}

func RunStep(s State, step string, f func() (error)){
	if s.IsStepComplete(step) {
		fmt.Printf("skipping step '%s' execution\n", step)
	} else {
		fmt.Printf("starting step '%s' execution\n", step)
		err := f()
		if err != nil {
			s.FailStep(step, err.Error())
			os.Exit(4)
		}
		s.CompleteStep(step)
	}
}

func RunStepE(s State, step string, f func() (error)) error {
	if s.IsStepComplete(step) {
		fmt.Printf("skipping step '%s' execution\n", step)
	} else {
		fmt.Printf("starting step '%s' execution\n", step)
		err := f()
		if err != nil {
			s.FailStep(step, err.Error())
			return err
		}
		s.CompleteStep(step)
	}
	return nil
}
