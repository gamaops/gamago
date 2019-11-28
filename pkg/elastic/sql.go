package elastic

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type ElasticBulkResponse struct {
	Errors bool `json:"errors"`
	Items  []struct {
		Index struct {
			ID     string `json:"_id"`
			Result string `json:"result"`
			Status int    `json:"status"`
			Error  struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
				Cause  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
				} `json:"caused_by"`
			} `json:"error"`
		} `json:"index"`
	} `json:"items"`
	Count int
}

var ErrElasticBulk = errors.New("errors on Elasticsearch bulk operation")

func (e *Elastic) IndexFromSQLRows(index string, rows *sql.Rows) (*ElasticBulkResponse, error) {

	var buf bytes.Buffer
	blk := new(ElasticBulkResponse)

	for rows.Next() {
		var id string
		var data []byte
		if err := rows.Scan(&id, &data); err != nil {
			return nil, err
		}
		meta := []byte(fmt.Sprintf(`{"index":{"_id" :"%s"}}\n`, id))
		data = append(data, '\n')
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)
		blk.Count++
	}
	defer buf.Reset()
	res, err := e.Client.Bulk(bytes.NewReader(buf.Bytes()), e.Client.Bulk.WithIndex(index))
	if err != nil {
		return blk, err
	}
	defer res.Body.Close()
	var raw map[string]interface{}

	if res.IsError() {
		if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
			return blk, err
		}
		return blk, &ErrElasticResponse{
			Response: res,
			Raw:      raw,
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
		return blk, err
	}
	for _, d := range blk.Items {
		if d.Index.Status > 201 {
			return blk, ErrElasticBulk
		}
	}

	return blk, nil
}
