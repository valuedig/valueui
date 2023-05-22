package status

import (
	"strings"
	"sync"
	// "github.com/valuedig/valueui/pkg/uiapi"
)

var (
	Assets = sets{
		items: map[string]interface{}{},
	}
)

type sets struct {
	mu    sync.Mutex
	items map[string]interface{}
}

func (it *sets) Sync(name string, v interface{}) {
	it.mu.Lock()
	defer it.mu.Unlock()
	it.items[strings.TrimLeft(name, "/")] = v
}

func (it *sets) Get(name string) interface{} {
	it.mu.Lock()
	defer it.mu.Unlock()
	if v, ok := it.items[strings.TrimLeft(name, "/")]; ok {
		return v
	}
	return nil
}
