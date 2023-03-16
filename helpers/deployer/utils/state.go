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
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	var state State
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Printf("creating new state file '%s'\n.", file)
		state = State{
			File: file,
		}
	} else {
		f, err := ioutil.ReadFile(file)
		if err != nil {
			return state, err
		}
		err = json.Unmarshal(f, &state)
		if err != nil {
			return state, err
		}
		state.File = file
	}
	if state.Steps == nil {
		state.Steps = map[string]Step{}
	}
	return state, nil
}

func (e State) SaveState() {
	f, _ := json.Marshal(e)
	os.WriteFile(e.File, f, 0644)
}

func (e State) CompleteStep(name string) {
	e.Steps[name] = Step{
		Name:   name,
		Status: completeState,
		Error:  "",
	}
	e.SaveState()
	fmt.Printf("completing step '%s' execution\n", name)
}

func (e State) IsStepComplete(name string) bool {
	val, ok := e.Steps[name]
	if ok {
		return val.Status == completeState
	}
	return false
}

func (e State) FailStep(name string, err string) {
	e.Steps[name] = Step{
		Name:   name,
		Status: errorState,
		Error:  err,
	}
	e.SaveState()
	fmt.Printf("failing step '%s'. Failed with error: %s\n", name, err)
}

func (e State) GetStepError(name string) string {
	val, ok := e.Steps[name]
	if ok {
		return val.Error
	}
	return ""
}

func RunStep(e State, step string, f func() (error)){
	if e.IsStepComplete(step) {
		fmt.Printf("skipping step '%s' execution\n", step)
	} else {
		fmt.Printf("starting step '%s' execution\n", step)
		err := f()
		if err != nil {
			e.FailStep(step, err.Error())
			os.Exit(4)
		}
		e.CompleteStep(step)
	}
}

func RunStepE(e State, step string, f func() (error)) error {
	if e.IsStepComplete(step) {
		fmt.Printf("skipping step '%s' execution\n", step)
	} else {
		fmt.Printf("starting step '%s' execution\n", step)
		err := f()
		if err != nil {
			e.FailStep(step, err.Error())
			return err
		}
		e.CompleteStep(step)
	}
	return nil
}
