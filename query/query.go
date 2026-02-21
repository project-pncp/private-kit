package query

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

type Sort struct {
	Key   string `json:"key"`
	Order string `json:"order"`
}

type Match struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}

type Filter struct {
	Rows    int64    `json:"rows"`
	Page    int64    `json:"page"`
	Cursor  string   `json:"cursor"`
	Group   string   `json:"group"`
	Dense   bool     `json:"dense"`
	Sort    Sort     `json:"sort"`
	Matches []Match  `json:"matches"`
	Term    string   `json:"term"`
	Fields  []string `json:"fields"`

	// Deprecated: must use only Sort instead.
	Sorts []Sort `json:"sorts"`
}

func (c *Filter) DecodeGRPC(md metadata.MD) (*Filter, error) {
	var (
		mdRows      = md.Get("rows")
		mdPage      = md.Get("page")
		mdCursor    = md.Get("cursor")
		mdGroup     = md.Get("group")
		mdDense     = md.Get("dense")
		mdSort      = md.Get("sort")
		mdSortOrder = md.Get("sort-order")
		mdTerm      = md.Get("term")
		mdFields    = md.Get("fields")
	)

	if len(mdRows) == 1 {
		rows, _ := strconv.Atoi(mdRows[0])
		c.Rows = int64(rows)
	}
	if c.Rows <= 0 || c.Rows > 100 {
		c.Rows = 10
	}

	if len(mdPage) == 1 {
		page, _ := strconv.Atoi(mdPage[0])
		c.Page = int64(page)
	}
	if c.Page <= 0 {
		c.Page = 1
	}

	if len(mdCursor) == 1 {
		c.Cursor = mdCursor[0]
	}

	if len(mdGroup) == 1 {
		c.Group = mdGroup[0]
	}

	if len(mdDense) == 1 {
		if mdDense[0] == "true" {
			c.Dense = true
		}
	}
	if len(mdTerm) > 0 {
		c.Term = mdTerm[0]
	}
	if len(mdSort) == 1 && len(mdSortOrder) == 1 {
		c.Sort.Key = strings.ReplaceAll(mdSort[0], "_", "")
		c.Sort.Order = mdSortOrder[0]
	}

	for i := 0; ; i++ {
		key := md.Get(fmt.Sprintf("filters-%d-key", i))
		value := md.Get(fmt.Sprintf("filters-%d-value", i))
		op := md.Get(fmt.Sprintf("filters-%d-op", i))
		if len(key) == 0 || len(op) == 0 {
			// no more values
			break
		}

		c.Matches = append(c.Matches, Match{
			Key:   strings.ReplaceAll(key[0], "_", ""),
			Value: value[0],
			Op:    op[0],
		})
	}

	for _, field := range mdFields {
		f := strings.ReplaceAll(field, "_", "")
		c.Fields = append(c.Fields, f)
	}

	return c, nil
}
