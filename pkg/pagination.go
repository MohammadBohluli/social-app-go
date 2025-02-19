package pkg

import (
	"net/http"
	"strconv"
)

type PaginationFeedQuery struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

func (p PaginationFeedQuery) Parse(r *http.Request) (PaginationFeedQuery, error) {
	query := r.URL.Query()

	limit := query.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return p, nil
		}

		p.Limit = l
	}

	offset := query.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return p, nil
		}

		p.Offset = o
	}

	sort := query.Get("sort")
	if sort != "" {
		p.Sort = sort
	}

	return p, nil
}
