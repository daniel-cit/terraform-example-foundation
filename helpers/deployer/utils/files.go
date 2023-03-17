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
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CopyFile(src string, dest string) error {
	s, err := os.Stat(src)
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dest, buf, s.Mode())
}

func CopyDirectory(src string, dest string) error {
	err := os.MkdirAll(dest, 0755)
	if err != nil {
		return err
	}
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			if f.Name() != ".terraform" {
				err := CopyDirectory(filepath.Join(src, f.Name()), filepath.Join(dest, f.Name()))
				if err != nil {
					return err
				}
			}
		} else {
			if f.Name() != ".terraform.lock.hcl" {
				err := CopyFile(filepath.Join(src, f.Name()), filepath.Join(dest, f.Name()))
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func ReplaceStringInFile(filename, old, new string) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, bytes.Replace(f, []byte(old), []byte(new), -1), 0644)
}

func FindFiles(dir, filename string) ([]string, error) {
	found := []string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if d.Name() == filename && !strings.Contains(path, ".terraform") {
			found = append(found, path)
		}
		return nil
	})
	return found, err
}
