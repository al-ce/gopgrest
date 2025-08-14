package rsql

// PathQuery holds the parts of a RESTful GET request's URL that are translated
// to a SQL SELECT query. In the example:
//
// `GET /authors?filter=forename=in=Ann,Anne;surname=Carson&fields=forename,surname`
//
// the Resource and Query are split at the `?` character
type PathQuery struct {
	Resource string
	Query    string
}

// Query is a parsed PathQuery that is used to build a SQL query.
type Query struct {
	Fields  Fields
	Filters []Filter
	Join    Joins
}

// FiltersMap is map of values to filter by. The expected use is that for a
// row to match the query filter, one of the elements from each key of the
// filter must match the corresponding column value of the row.
type FiltersMap map[string][]string

// Filter is the parsed result of one of any `;` separated filter conditions in
// a URL query. The example:
//
// `GET /authors?filter=forename=in=Ann,Anne;surname=!=Carson`
//
// would be parsed as two seaparte Filter values:
// {Column: "forename", Values: []string{"Ann", "Anne"}, SQLOperator: "IN" }
// {Column: "surname", Values: []string{"Carson"}, SQLOperator: "IN" }
type Filter struct {
	Column      string
	Values      []string
	SQLOperator string
}

type Joins map[string]JoinRelation

type Fields []string

type JoinRelation struct {
	LeftTable  string
	LeftCol    string
	RightTable string
	RightCol   string
}

// VALIDKEYWORDS are valid clause keywords for a URL query
var VALIDKEYWORDS = []string{FILTER, FIELDS, JOIN}

const (
	FILTER = "filter"
	FIELDS = "fields"
	JOIN   = "join"
)

// regexFilterOperators is an array of valid operators for a filter condition. We
// use this to guarantee the order of the subexpressions, so that "<=" is
// checked before "<" etc.
var regexFilterOperators = []string{
	"==",
	"!=",
	"=in=",
	"=out=",
	"=like=",
	"=!like=",
	"=notlike=",
	"=nk=",
	"=isnull=",
	"=na=",
	"=isnotnull=",
	"=notnull=",
	"=nn=",
	"=!null=",
	"=le=",
	"=ge=",
	"<=",
	">=",
	"=lt=",
	"=gt=",
	"<", // KEEP THIS AFTER <=
	">", // KEEP THIS AFTER >=
}

// OperatorToSQLMap is a map of RSQL operators to their SQL counterpart
var OperatorToSQLMap = map[string]string{
	"==":          "=",
	"!=":          "!=",
	"=in=":        "IN",
	"=out=":       "NOT IN",
	"=like=":      "LIKE",
	"=!like=":     "NOT LIKE",
	"=notlike=":   "NOT LIKE",
	"=nk=":        "NOT LIKE",
	"=isnull=":    "IS NULL",
	"=na=":        "IS NULL",
	"=isnotnull=": "IS NOT NULL",
	"=notnull=":   "IS NOT NULL",
	"=nn=":        "IS NOT NULL",
	"=!null=":     "IS NOT NULL",
	"=le=":        "<=",
	"=ge=":        ">=",
	"<=":          "<=",
	">=":          ">=",
	"=lt=":        "<",
	"=gt=":        ">",
	"<":           "<",
	">":           ">",
}
