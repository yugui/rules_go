/* Copyright 2016 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io/ioutil"
	"os"

	bzl "github.com/bazelbuild/buildifier/core"
	"github.com/bazelbuild/buildifier/differ"
)

func diffFile(file *bzl.File) (err error) {
	f, err := ioutil.TempFile("", "BUILD")
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			if merr := os.Remove(f.Name()); merr != nil {
				err = merr
			}
		}
	}()
	err = func() (err error) {
		defer func() {
			if cerr := f.Close(); cerr != nil {
				if err == nil {
					err = cerr
				}
			}
		}()
		_, err = f.Write(bzl.Format(file))
		return err
	}()
	if err != nil {
		return err
	}

	diff := differ.Find()
	if _, err := os.Stat(file.Path); os.IsNotExist(err) {
		diff.Show(os.DevNull, f.Name())
		return nil
	}
	diff.Show(file.Path, f.Name())
	return nil
}
