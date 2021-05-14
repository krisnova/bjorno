package bjorno

import "sync"

type EmptyProgram struct {
	mutex sync.Mutex
}

func (v *EmptyProgram) Values() interface{} {
	return v
}

func (v *EmptyProgram) Refresh() {
}
func (v *EmptyProgram) Lock() {
	v.mutex.Lock()
}
func (v *EmptyProgram) Unlock() {
	v.mutex.Unlock()
}
