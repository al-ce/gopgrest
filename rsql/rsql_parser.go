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
		keyword, conditionals, err := parseClause(url, clause)
		if err != nil {
			return nil, err
		}

		switch keyword {
		case FILTER:
			// Query filters are split at ";" and checked with "==" or "=in="
			filters, err := newFilters(conditionals)
			if err != nil {
				return nil, err
			}
			query.Filters = filters
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
func parseClause(url, clauseStr string) (string, string, error) {
	// Split clause at assignment '=' char, not equality '==' chars
	clause := strings.SplitN(clauseStr, "=", 2)
	// Clause requires values on left and right of assignment
	if len(clause) != 2 {
		return "", "", fmt.Errorf("Malformed clause: %s in url %s\n", clauseStr, url)
	}
	// Keyword must be one in a list of implemented clauses
	keyword := clause[0]
	if !slices.Contains(VALIDKEYWORDS, keyword) {
		return "", "", fmt.Errorf("Invalid clause keyword '%s' in url %s\n", keyword, url)
	}
	values := clause[1]
	return keyword, values, nil
}

// newFilters makes a rsql.Filters value from the rhs of a URL query param
// e.g. the rhs of `filters=forename=in=Anne,Ann;surname=Carson“
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
