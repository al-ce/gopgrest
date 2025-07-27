package types

// QueryFilters is map of values to filter by. The expected use is that for a
// row to match the query filter, one of the elements from each key of the
// filter must match the corresponding column value of the row.
type QueryFilters map[string][]string

// RowDataMap is a map used to store data that was retrieved from the database
// or that is to be written to the database
type RowDataMap map[string]any
