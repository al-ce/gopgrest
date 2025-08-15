package rsql

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
)

// newRSQLQuery builds a ParsedURL from the URL
func NewRSQLQuery(url string) (*Query, error) {
	// Separate table (t) from URL query params (p)
	rp, err := newPathQuery(url)
	if err != nil || rp == nil {
		return nil, err
	}
	query := Query{}

	// Example URL query:
	// /authors?=filter=forename=in=Ann,Anne;surname=Carson&fields=forename,surname

	// Parse clauses from URL query params, split at "&"
	// e.g. "filter=..." + "fields=..."
	for clause := range strings.SplitSeq(rp.Query, "&") {
		keyword, namedArgs, err := parseClause(clause)
		if err != nil {
			return nil, err
		}

		var clauseErr error
		switch keyword {
		case FILTER: // e.g ?filter=
			// Query filters are split at ";" and checked with "==" or "=in="
			filters, err := newFilters(namedArgs)
			clauseErr = err
			query.Filters = filters
		case FIELDS: // e.g. ?fields=
			// Query fields are split at ","
			fields, err := newFields(namedArgs)
			clauseErr = err
			query.Fields = fields
		case JOIN: // e.g. ?join=
			joins, err := newJoins(JOIN, namedArgs)
			clauseErr = err
			query.Joins = joins
		}
		if clauseErr != nil {
			return nil, clauseErr
		}

	}
	return &query, nil
}

// newPathQuery splits the URL at the query string separator and returns a
// PathQuery value
func newPathQuery(url string) (*PathQuery, error) {
	// Separate table (t) from URL query params (p)
	tp := strings.Split(url, "?")
	if len(tp) < 1 {
		return nil, fmt.Errorf("Could not parse table from URL: %s", url)
	}
	// If there were no queries, return with nil reference
	if len(tp) != 2 {
		return nil, nil
	}
	return &PathQuery{Resource: tp[0], Query: tp[1]}, nil
}

// parseClause validates the lhs and rhs of a clause string, e.g.
// `fields=forename,surname` should split into two substrings at the `=` char
// and the keyword (lhs) should be an implemented clause keyword
func parseClause(clauseStr string) (string, string, error) {
	// Split clause at assignment '=' char, not equality '==' chars
	clause := strings.SplitN(clauseStr, "=", 2)
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

// newFilters makes a rsql.Filters value from the rhs of a URL filter query param
// e.g. the rhs of `filter=forename=in=Anne,Ann;surname=Carson`
func newFilters(filterConditionals string) ([]Filter, error) {
	filters := []Filter{}
	// Multiple filters allowed with ; separator
	for cond := range strings.SplitSeq(filterConditionals, ";") {
		ReFilterOperator := getOperatorSplitRegex()
		f, err := newFilter(cond, ReFilterOperator)
		if err != nil {
			return []Filter{}, err
		}
		filters = append(filters, *f)
	}
	return filters, nil
}

// newFilter creates a new Filter from an item in a ';' separated string of
// `{col}=[values]` filter conditionals
func newFilter(cond string, ReFilterOperator *regexp.Regexp) (*Filter, error) {
	// Split filter at "==" or "=in=" or "=out=" etc.
	splitFilter := ReFilterOperator.Split(cond, -1)
	if len(splitFilter) != 2 {
		return nil, fmt.Errorf("Malformed FILTERS clause in url: %s\n", cond)
	}
	operator := ReFilterOperator.FindString(cond)
	filterCol := splitFilter[0]
	filterVals := strings.Split(splitFilter[1], ",")

	nullCheck := hasNullCheck(operator)
	// operator may be empty string and Split always returns array len 1
	// so handle case of no values, except for isnull/isnotnull
	if !nullCheck && len(filterVals) == 1 && filterVals[0] == "" {
		return nil, fmt.Errorf("attempt to filter on col %s with no values", filterCol)
	}
	// null check conditions should not have any rhs values, e.g.
	// `filter=born=isnull=` is valid
	// `filter=born=isnull=1800` is not valid
	if nullCheck && len(filterVals) >= 1 && filterVals[0] != "" {
		return nil, fmt.Errorf("cannot add values to null check conditions")
	}

	validOps := slices.Collect(maps.Keys(OperatorToSQLMap))
	if !slices.Contains(validOps, operator) {
		return nil, fmt.Errorf("invalid operator %s on filter %s", operator, cond)
	}
	return &Filter{
		Column:      filterCol,
		Values:      filterVals,
		SQLOperator: OperatorToSQLMap[operator],
	}, nil
}

// getOperatorSplitRegex builds a regex from the OperatorToSQLMap's keys. The
// regex matches any of the keys in the map
func getOperatorSplitRegex() *regexp.Regexp {
	// Build valid operators to split at
	opRegex := []string{}
	for _, r := range regexFilterOperators {
		opRegex = append(opRegex, fmt.Sprintf("(%s)", r))
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

// newFields makes a rsql.Fields value from the RHS of a URL fields query param
// e.g. the rhs of `fields=forename,surename`
func newFields(selectedFields string) (Fields, error) {
	fields := strings.Split(selectedFields, ",")
	if slices.Contains(fields, "") {
		return nil, fmt.Errorf("Empty field in %s", selectedFields)
	}
	return fields, nil
}

// newJoins makes a rsql.Joins value from the RHS of a URL join query param
// e.g. the rhs of `join=authors:books.author_id=authors.id`
func newJoins(joinType, joinRelations string) ([]JoinRelation, error) {
	jr := []JoinRelation{}

	// Example:
	// GET /books?join=authors:books.author_id==authors.id;genres:books.genres_id==genres.id
	// Note that this enforces qualified column names in a JOIN statement
	ReJoin := regexp.MustCompile(`(\w+):(\w+)\.(\w+)==(\w+)\.(\w+)`)

	// Multiple joins allowed with ; separator
	for join := range strings.SplitSeq(joinRelations, ";") {
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
