package datax

func NewDataResult() *DataResult {
	return &DataResult{}
}

func NewDataResultError(s string) *DataResult {
	return &DataResult{
		Error: s,
	}
}

// func (it *DataResult) NewRow() *ModelSpec_Row {
// 	return &ModelSpec_Row{
// 		Id: "",
// 	}
// }

func (it *DataResult) SetRow(rows ...*ModelSpec_Row) *DataResult {
	for _, row := range rows {
		it.Rows = append(it.Rows, row)
	}
	return it
}
