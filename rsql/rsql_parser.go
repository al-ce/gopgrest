package rsql

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"

	"gopgrest/repatterns"
)

// Separator characters in query
var (
	CLAUSE_ASSIGN   = "=" // assign value `name` to clause `select`: `select=name`
	CLAUSE_SEP      = "&" // separate `select=...` and `where=...`: `select=name&where=name==bob`
	ITEM_SEP        = ";" // separate multiple equality checks: `where=name==bob;age==42`
	ALIAS_SEP       = ":" // separate column name and alias: `select=surname:last_name`
	QUALIFIER_SEP   = "." // qualify column `surname` w/ table `authors`: `select=authors.surname`
	JOIN_ON_ASSIGN  = ":" // `JOIN ON authors WHERE...`: `join=authors:books.author_id==authors.id`
	VALUES_LIST_SEP = "," // separate list of values e.g. `select=surname,forename,died`
)

// newRSQLQuery builds a Pars from the URL
func NewRSQLQuery(url string) (QueryParams, error) {
	// Separate table (t) from URL query params (p)
	pq, err := newPathQuery(url)
	if err != nil || pq == nil {
		return QueryParams{}, err
	}
	query := QueryParams{}
	query.Tables = append(query.Tables, pq.Resource)

	// Example URL query:
	// /authors?where=forename=in=Ann,Anne;surname=Carson&select=forename,surname

	// Parse clauses from URL query params, split at "&"
	// e.g. "where=..." + "select=..."
	for clause := range strings.SplitSeq(pq.Query, CLAUSE_SEP) {
		keyword, namedArgs, err := parseClause(clause)
		if err != nil {
			return QueryParams{}, err
		}

		var clauseErr error
		switch keyword {
		case WHERE: // e.g ?where=
			// WHERE conditions are split at ";" and checked with "==" or "=in="
			conditions, err := NewWhereConditions(namedArgs)
			clauseErr = err
			query.Conditions = conditions
		case SELECT: // e.g. ?select=
			// Select columns are split at ","
			columns, err := newSelect(namedArgs)
			clauseErr = err
			query.Columns = columns
		case JOIN:
			fallthrough
		case INNERJOIN:
			fallthrough
		case LEFTJOIN:
			fallthrough
		case RIGHTJOIN:
			joins, err := newJoins(keyword, namedArgs)
			clauseErr = err
			query.Joins = joins
			// Add all referenced tables to the query
			for _, j := range joins {
				query.Tables = append(query.Tables, j.Table)
			}
		}
		if clauseErr != nil {
			return QueryParams{}, clauseErr
		}

	}
	return query, nil
}

// newPathQuery parses a URL, checking for a table name and an optional query
func newPathQuery(url string) (*PathQuery, error) {
	// Has table in URL but no query
	if repatterns.ReqNoParams.MatchString(url) {
		return nil, nil
	}
	// Has table in URL and query
	queryMatches := repatterns.ReqHasParams.FindStringSubmatch(url)
	if repatterns.ReqHasParams.MatchString(url) {
		return &PathQuery{Resource: queryMatches[1], Query: queryMatches[2]}, nil
	}
	// Bad URL
	return nil, fmt.Errorf("could not parse url %s", url)
}

// parseClause validates the lhs and rhs of a clause string, e.g.
// `select=forename,surname` should split into two substrings at the `=` char
// and the keyword (lhs) should be an implemented clause keyword
func parseClause(clauseStr string) (string, string, error) {
	// Split clause at assignment '=' char, not equality '==' chars
	clause := strings.SplitN(clauseStr, CLAUSE_ASSIGN, 2)
	// Clause requires values on left and right of assignment
	if len(clause) != 2 {
		return "", "", fmt.Errorf("Malformed clause: %s\n", clauseStr)
	}
	// Keyword must be one in a list of implemented clauses
	keyword := clause[0]
	if !slices.Contains(VALIDKEYWORDS, keyword) {
		return "", "", fmt.Errorf("Invalid clause keyword '%s'\n", keyword)
	}
	values := clause[1]
	return keyword, values, nil
}

// NewWhereConditions makes a rsql.Conditions value from the rhs of a URL 'WHERE' query param
// e.g. the rhs of `where=forename=in=Anne,Ann;surname=Carson`
func NewWhereConditions(whereConditions string) ([]Condition, error) {
	conditions := []Condition{}
	// Multiple conditions allowed with ; separator
	for cond := range strings.SplitSeq(whereConditions, ITEM_SEP) {
		f, err := newCondition(cond)
		if err != nil {
			return []Condition{}, err
		}
		conditions = append(conditions, *f)
	}
	return conditions, nil
}

// newCondition creates a new Condition from an item in a ';' separated string of
// `{col}=[values]` WHERE conditions
func newCondition(cond string) (*Condition, error) {
	// Split condition at "==" or "=in=" or "=out=" etc.
	ReConditionOperator := getOperatorSplitRegex()
	splitCondition := ReConditionOperator.Split(cond, -1)
	if len(splitCondition) != 2 {
		return nil, fmt.Errorf("Malformed WHERE clause in url: %s\n", cond)
	}
	operator := ReConditionOperator.FindString(cond)
	conditionCol := splitCondition[0]
	conditionVals := strings.Split(splitCondition[1], VALUES_LIST_SEP)

	// Build column
	column := Column{}
	qualifiedCol := strings.Split(splitCondition[0], QUALIFIER_SEP)
	// Add qualifier if the column was qualified with a table, e.g.
	// `authors.forename`
	if len(qualifiedCol) == 2 {
		column.Qualifier = qualifiedCol[0]
		column.Name = qualifiedCol[1]
	} else {
		column.Name = splitCondition[0]
	}

	nullCheck := hasNullCheck(operator)
	// operator may be empty string and Split always returns array len 1
	// so handle case of no values, except for isnull/isnotnull
	if !nullCheck && len(conditionVals) == 1 && conditionVals[0] == "" {
		return nil, fmt.Errorf("Condition on col %s with no values", conditionCol)
	}
	// null check conditions should not have any rhs values, e.g.
	// `where=born=isnull=` is valid
	// `where=born=isnull=1800` is not valid
	if nullCheck && len(conditionVals) >= 1 && conditionVals[0] != "" {
		return nil, fmt.Errorf("cannot add values to null check conditions")
	}

	// Operator must be implemented
	validOps := slices.Collect(maps.Keys(OperatorToSQLMap))
	if !slices.Contains(validOps, operator) {
		return nil, fmt.Errorf("invalid operator %s on condition %s", operator, cond)
	}

	return &Condition{
		Column:      column,
		Values:      conditionVals,
		SQLOperator: OperatorToSQLMap[operator],
	}, nil
}

// getOperatorSplitRegex builds a regex from the OperatorToSQLMap's keys. The
// regex matches any of the keys in the map
func getOperatorSplitRegex() *regexp.Regexp {
	// Build valid operators to split at
	opRegex := []string{}
	for k := range OperatorToSQLMap {
		opRegex = append(opRegex, fmt.Sprintf("(%s)", k))
	}
	return regexp.MustCompile(fmt.Sprintf("(%s)", strings.Join(opRegex, "|")))
}

func hasNullCheck(operator string) bool {
	return slices.Contains(
		[]string{
			"=isnull=",
			"=na=",
			"=isnotnull=",
			"=notnull",
			"=!null=",
		},
		operator)
}

// newSelect makes a rsql.Columns value from the RHS of a URL select query param
// e.g. the rhs of `select=forename,surename`
func newSelect(selectedColumns string) ([]Column, error) {
	columns := []Column{}
	for sf := range strings.SplitSeq(selectedColumns, VALUES_LIST_SEP) {
		if sf == "" {
			return nil, fmt.Errorf("Empty column in %s", selectedColumns)
		}

		column := Column{}

		// Check for alias indicated by `:` e.g. `genres.name:genre`
		alias := strings.Split(sf, ALIAS_SEP)
		if len(alias) > 2 {
			return nil, fmt.Errorf("Too many alias separators in %s", sf)
		} else if len(alias) == 2 {
			column.Alias = alias[1]
		}
		if slices.Contains(alias, "") {
			return nil, fmt.Errorf("Empty column operand in %s", sf)
		}

		// Check for column quailifier indicated by `.` e.g. `books.author_id`
		col := strings.Split(alias[0], QUALIFIER_SEP)
		if len(col) > 2 {
			return nil, fmt.Errorf("Too many qualifier separators in %s", sf)
		} else if len(col) == 2 {
			column.Qualifier = col[0]
			column.Name = col[1]
		} else {
			column.Name = col[0]
		}
		if slices.Contains(col, "") {
			return nil, fmt.Errorf("Empty column operand in %s", sf)
		}

		columns = append(columns, column)

	}
	return columns, nil
}

// newJoins makes a rsql.Joins value from the RHS of a URL join query param
// e.g. the rhs of `join=authors:books.author_id=authors.id`
func newJoins(joinType, joinRelations string) ([]JoinRelation, error) {
	jr := []JoinRelation{}

	// e.g. transform "inner_join" to INNER JOIN
	joinType = strings.ToUpper(strings.ReplaceAll(joinType, "_", " "))

	// Example:
	// GET /books?join=authors:books.author_id==authors.id;genres:books.genres_id==genres.id
	// Note that this enforces qualified column names in a JOIN statement
	ReJoin := regexp.MustCompile(
		fmt.Sprintf(`(\w+)%s(\w+)\%s(\w+)==(\w+)\%s(\w+)`,
			JOIN_ON_ASSIGN,
			QUALIFIER_SEP,
			QUALIFIER_SEP),
	)

	// Multiple joins allowed with ; separator
	for join := range strings.SplitSeq(joinRelations, ITEM_SEP) {
		matches := ReJoin.FindStringSubmatch(join)
		jr = append(jr, JoinRelation{
			Type:           strings.ToUpper(joinType),
			Table:          matches[1],
			LeftQualifier:  matches[2],
			LeftCol:        matches[3],
			RightQualifier: matches[4],
			RightCol:       matches[5],
		})
	}
	return jr, nil
}
