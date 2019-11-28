package sqlutils

import (
	"database/sql"
	"fmt"
)

type SQLPager struct {
	Query     string
	Arguments []interface{}
	Limit     int
	offset    int
}

func (s *SQLPager) Scroll(db *sql.DB) (*sql.Rows, error) {

	query := fmt.Sprintf(s.Query, s.Limit, s.offset)
	rows, err := db.Query(query, s.Arguments...)

	if err != nil {
		return nil, err
	}

	s.offset += s.Limit

	return rows, nil
}
