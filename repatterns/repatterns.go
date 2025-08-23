package repatterns

import "regexp"

var (
	ReqWithId         = regexp.MustCompile(`^/(\w+)/([0-9]+)/?\??$`)
	ReqOptionalParams = regexp.MustCompile(`^/(\w+)(\?.*)?/?\??$`)
	ReqNoParams       = regexp.MustCompile(`^/(\w+)/?\??$`)
	ReqHasParams      = regexp.MustCompile(`^/(\w+)\?(.*)$`)
)
