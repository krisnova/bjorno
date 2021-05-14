// Copyright © 2021 Kris Nóva <kris@nivenly.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"sync"

	bjorno "github.com/kris-nova/bjorn"
)

func main() {

	cfg := &bjorno.ServerConfig{
		InterpolateExtensions: []string{
			".html",
		},
		BindAddress:    ":1313",
		ServeDirectory: "../bjorno.com",
		DefaultIndexFiles: []string{
			"index.html",
		},
	}
	bjorno.Runtime(cfg, &BjornoProgram{
		Name:     "Björn",
		Nickname: "Pupperscotch",
	})

}

type BjornoProgram struct {
	Name     string
	Nickname string
	mutex    sync.Mutex
}

func (v *BjornoProgram) Values() interface{} {
	return v
}

func (v *BjornoProgram) Refresh() {
	v.Nickname = "butterscotch"
	v.Name = "björn"
}
func (v *BjornoProgram) Lock() {
	v.mutex.Lock()
}
func (v *BjornoProgram) Unlock() {
	v.mutex.Unlock()
}
