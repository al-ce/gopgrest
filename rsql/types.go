package rsql

// VALIDKEYWORDS are valid clause keywords for a URL query
var VALIDKEYWORDS = []string{
	WHERE,
	SELECT,
	JOIN,
	INNERJOIN,
	LEFTJOIN,
	RIGHTJOIN,
}

const (
	WHERE     = "where"
	SELECT    = "select"
	JOIN      = "join"
	INNERJOIN = "inner_join"
	LEFTJOIN  = "left_join"
	RIGHTJOIN = "right_join"
)

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

// PathQuery holds the parts of a RESTful GET request's URL that are translated
// to a SQL SELECT query. In the example:
//
// `GET /authors?where=forename=in=Ann,Anne;surname=Carson&select=forename,surname`
//
// the Resource and Query are split at the `?` character
type PathQuery struct {
	Resource string
	Query    string
}

// QueryParams is a parsed PathQuery that is used to build a SQL query.
type QueryParams struct {
	Tables     []string       // Tables to SELECT in either FROM or JOIN
	Columns    []Column       // Columns to return in SELECT query
	Conditions []Condition    // Conditionals for WHERE clause
	Joins      []JoinRelation // Relations for JOIN clauses
}

// Condition is the parsed result of one of any `;` separated 'where' conditions in
// a URL query. The example:
//
// `GET /authors?where=forename=in=Ann,Anne;surname=!=Carson`
//
// would be parsed as two separate Condition values:
// {Column: "forename", Values: []string{"Ann", "Anne"}, SQLOperator: "IN" }
// {Column: "surname", Values: []string{"Carson"}, SQLOperator: "IN" }
type Condition struct {
	Column      Column
	Values      []string
	SQLOperator string
}

type Column struct {
	Qualifier string
	Name      string
	Alias     string
}

type JoinRelation struct {
	Type           string
	Table          string
	LeftQualifier  string
	LeftCol        string
	RightQualifier string
	RightCol       string
}
