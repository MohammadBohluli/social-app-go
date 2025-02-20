package pkg

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginationFeedQuery struct {
	Limit  int      `json:"limit"`
	Offset int      `json:"offset"`
	Sort   string   `json:"sort"`
	Tags   []string `json:"tags"`
	Search string   `json:"search"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	tags := query.Get("tags")
	if tags != "" {
		p.Tags = strings.Split(tags, ",")
	}

	search := query.Get("search")
	if search != "" {
		p.Search = search
	}

	since := query.Get("since")
	if since != "" {
		p.Since = parseTime(since)
	}

	until := query.Get("until")
	if until != "" {
		p.Until = parseTime(since)
	}

	return p, nil
}

func parseTime(s string) string {
	t, err := time.Parse(time.DateTime, s)
	if err != nil {
		return ""
	}

	return t.Format(time.DateTime)
}
