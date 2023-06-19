package datax

type DataReader interface {
	ModelSpecList() ([]*ModelSpec, error)
	Search(req *DataRead) (*DataResult, error)
}

type DataWriter interface {
	// Insert(model string, rowFields map[string]interface{}) error
	Upsert(model string, rowFields ...map[string]interface{}) error
	// Delete(model string, filter ...DataRead_Filter) error
}

type DataDriver interface {
	DataReader
	DataWriter
}
