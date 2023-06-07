package uiapi

type DataDriver interface {
	Search(*DataRead) (*DataResult, error)
}
