package utils

import (
	"net/http"
	"strconv"
)

var (
	defaultLimit  = 10
	defaultOffset = 0
)

type Pagination struct {
	Limit  int
	Offset int
}

func GetPagination(r *http.Request) Pagination {

	sLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(sLimit)
	if err != nil || limit <= 0 {
		limit = defaultLimit
	}

	// prevent excessively large limits
	if limit > 100 {
		limit = 100
	}

	sOffset := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(sOffset)
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	return Pagination{
		Limit:  limit,
		Offset: offset,
	}
}
