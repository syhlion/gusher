package requestworker

import (
	"net/http"
)

type Job struct {
	Resq   *http.Request
	Result chan Result
}

type ErrorResult struct {
	s string
}

func (r *ErrorResult) Error() string {
	return r.s
}

type Result struct {
	Err error
}
