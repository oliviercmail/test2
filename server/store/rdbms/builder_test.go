package rdbms

import (
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza/server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuilder(t *testing.T) {
	var (
		req = require.New(t)
		cfg = &Config{PlaceholderFormat: squirrel.Question}
	)

	upsert, err := UpsertBuilder(cfg, "tbl", store.Payload{"c1": "v1", "c2": "v2"}, "c1")
	req.NoError(err)
	sql, args, err := upsert.ToSql()
	req.NoError(err)
	req.Contains(sql, "ON CONFLICT (c1) DO UPDATE SET")
	req.Contains(sql, "INSERT INTO tbl")
	req.Equal([]interface{}{"v1", "v2", "v2"}, args)

}
