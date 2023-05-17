package status

import (
	"strings"
	"sync"
	// "github.com/valuedig/valueui/pkg/uiapi"
)

var (
	Layouts = layoutSets{
		items: map[string]interface{}{},
	}
)

type layoutSets struct {
	mu    sync.Mutex
	items map[string]interface{}
}

func (it *layoutSets) Sync(name string, v interface{}) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.items[strings.TrimLeft(name, "/")] = v
}

func (it *layoutSets) Get(name string) interface{} {
	it.mu.Lock()
	defer it.mu.Unlock()
	if v, ok := it.items[strings.TrimLeft(name, "/")]; ok {
		return v
	}
	return nil
}
