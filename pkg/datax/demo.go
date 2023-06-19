package datax

import (
	"time"
)

type helloModelDriver struct {
	specs    []*ModelSpec
	datasets map[string][]*ModelSpec_Row
}

func NewHelloModel() (DataDriver, error) {
	m := &helloModelDriver{
		datasets: map[string][]*ModelSpec_Row{},
	}

	spec := NewModelSpec("events")

	spec.NewField("id", ModelSpec_UINT)
	spec.NewField("updated", ModelSpec_UINT)

	spec.NewField("content", ModelSpec_STRING)

	m.specs = append(m.specs, spec)

	m.datasets[spec.Name] = append(m.datasets[spec.Name], spec.NewRow(map[string]interface{}{
		"id":      1,
		"updated": time.Now().Unix(),
		"content": "hello world",
	}))

	m.datasets[spec.Name] = append(m.datasets[spec.Name], spec.NewRow(map[string]interface{}{
		"id":      2,
		"updated": time.Now().Unix() + 2,
		"content": "hello world 2",
	}))

	return m, nil
}

func (it *helloModelDriver) ModelSpecList() ([]*ModelSpec, error) {
	return it.specs, nil
}

func (it *helloModelDriver) Search(req *DataRead) (*DataResult, error) {
	var spec *ModelSpec
	for _, s := range it.specs {
		if req.Model == s.Name {
			spec = s
			break
		}
	}
	if spec == nil {
		return NewDataResultError("model not found"), nil
	}
	if rows, ok := it.datasets[spec.Name]; ok {
		return NewDataResult().SetRow(rows...), nil
	}
	return NewDataResult(), nil
}

func (it *helloModelDriver) Upsert(model string, rows ...map[string]interface{}) error {
	return nil
}
