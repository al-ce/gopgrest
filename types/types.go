package types

// RowData is a map used to store data that was retrieved from the database
// or that is to be written to the database
type RowData map[string]any

// RowDataIdMap is a map of RowData with the row id as the key
type RowDataIdMap map[int64]RowData
