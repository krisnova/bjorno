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
		ServeDirectory: "bjorno.com",
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
