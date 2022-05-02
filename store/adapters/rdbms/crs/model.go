package crs

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza-server/compose/crs"
	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ql"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

type (
	sqlizer interface {
		ToSQL() (string, []interface{}, error)
	}

	queryRunner interface {
		sqlx.QueryerContext
		sqlx.ExecerContext
	}

	attrExpression interface {
		exp.Comparable
		exp.Inable
		exp.Isable
	}

	attributeType interface {
		Type() data.AttributeType
	}

	column struct {
		ident      string
		columnType attributeType
		attributes []*data.Attribute

		encode func(drivers.Dialect, []*data.Attribute, crs.ValueGetter) (any, error)
		decode func(drivers.Dialect, []*data.Attribute, any, crs.ValueSetter) error
	}

	queryParser interface {
		Parse(string) (exp.Expression, error)
	}

	model struct {
		model *data.Model
		conn  queryRunner

		queryParser queryParser
		dialect     drivers.Dialect

		table drivers.TableCodec

		// ID column identifier expression
		//sysColumnID       exp.IdentifierExpression
		//sysPrimaryKeyAttr string

		// optional record fields/columns/expressions
		//sysExprNamespaceID attrExpression
		//sysSoftDeleteAttr  string
		//sysExprModuleID    attrExpression
		//sysModuleAttr      string
		//sysExprDeletedAt   attrExpression
		//sysNamespaceAttr   string

		// all columns we're selecting from when
		// we're selecting from all columns
		//columns []*column
	}
)

// Model returns fully initialized model store
//
// It abstracts database table and its columns and provides unified interface
// for fetching and storing records.
func Model(m *data.Model, c queryRunner, d drivers.Dialect) *model {
	var (
		ms = &model{
			model:       m,
			conn:        c,
			dialect:     d,
			queryParser: ql.Converter(),
			table:       drivers.NewTableCodec(m, d),
		}
	)

	ms.queryParser = ql.Converter(
		ql.SymHandler(func(node *ql.ASTNode) (exp.Expression, error) {
			return ms.table.AttributeExpression(node.Symbol)
		}),
	)

	return ms
}

func (d *model) Truncate(ctx context.Context) error {
	sql, args, err := d.truncateSql().ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Create(ctx context.Context, rr ...crs.ValueGetter) error {
	sql, args, err := d.insertSql(rr...).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Update(ctx context.Context, r crs.ValueGetter) error {
	sql, args, err := d.updateSql(r).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Delete(ctx context.Context, r crs.ValueGetter) error {
	sql, args, err := d.deleteSql(r).ToSQL()
	if err != nil {
		return err
	}

	_, err = d.conn.ExecContext(ctx, sql, args...)
	return err
}

func (d *model) Search(f types.RecordFilter) (i *iterator, err error) {
	if f.PageCursor != nil {
		// Page cursor exists; we need to validate it against used sort
		// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
		// from the cursor.
		// This (extracted sorting info) is then returned as part of response
		if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
			return nil, err
		}
	}

	for _, s := range f.Sort {
		if _, err = d.table.AttributeExpression(s.Column); err != nil {
			return nil, err
		}
	}

	// sanitize filter a bit
	for _, c := range d.table.Columns() {
		if !c.IsPrimaryKey() {
			continue
		}

		attrIdent := c.Attribute().Ident
		if f.Sort.Get(attrIdent) != nil {
			continue
		}

		// Make sure results are always sorted at least by primary key
		f.AppendOrderBy(attrIdent, f.Sort.LastDescending())
	}

	return &iterator{
		ms:      d,
		query:   d.searchSql(f),
		sorting: f.Sorting,
		paging:  f.Paging,
	}, nil
}

func (d *model) Lookup(ctx context.Context, pkv crs.ValueGetter, r crs.ValueSetter) (err error) {
	query, args, err := d.lookupSql(pkv).ToSQL()
	if err != nil {
		return
	}

	// using sql.Rows instead of a row
	// this gives us more control over closing (the rows resource)
	// and ability to use sql.RawBytes
	var rows *sql.Rows
	rows, err = d.conn.QueryContext(ctx, query, args...)
	if err = rows.Err(); err != nil {
		return
	}

	defer rows.Close()
	if !rows.Next() {
		return sql.ErrNoRows
	}

	scanBuf := d.table.MakeScanBuffer()
	if err = rows.Scan(scanBuf...); err != nil {
		return
	}

	if err = d.table.Decode(scanBuf, r); err != nil {
		return
	}

	return rows.Close()
}

// constructs SQL for selecting records from a table,
// converting parts of record filter into conditions
//
// Does not add any limits, sorting or any cursor conditions!
func (d *model) searchSql(f types.RecordFilter) *goqu.SelectDataset {
	var (
		err  error
		base = d.selectSql()
		tmp  exp.Expression
		cnd  []exp.Expression
	)

	{
		// Add model & namespace constraints when model expects (has configured attributes) for them
		//
		// This covers both scenarios:
		// 1) Model is configured to store records in a dedicated table
		//    without model and/or namespace attributes
		//
		// 2) Model has model and/or namespace attribute and saves records
		//    from different modules in the same table

		//if d.sysExprNamespaceID != nil {
		//	cnd = append(cnd, d.sysExprNamespaceID.Eq(f.NamespaceID))
		//} else {
		//	// @todo check if f.NamespaceID is compatible
		//}
		//
		//if d.sysExprModuleID != nil {
		//	cnd = append(cnd, d.sysExprModuleID.Eq(f.ModuleID))
		//} else {
		//	// @todo check if f.ModuleID is compatible
		//}
	}

	{
		// this can no longer be
		//if len(f.LabeledIDs) > 0 {
		//	// Limit by LabeledIDs (list of record IDs)
		//	cnd = append(cnd, d.sysColumnID.In(f.LabeledIDs))
		//}
	}

	{
		// If model supports soft-deletion (= delete-at attribute is present)
		// we need to make sure we respect it
		//if d.sysExprDeletedAt != nil {
		//	switch f.Deleted {
		//	case filter.StateExclusive:
		//		// only not-null values
		//		cnd = append(cnd, d.sysExprDeletedAt.IsNotNull())
		//
		//	case filter.StateExcluded:
		//		// exclude all non-null values
		//		cnd = append(cnd, d.sysExprDeletedAt.IsNull())
		//	}
		//}
	}

	if len(strings.TrimSpace(f.Query)) > 0 {
		if tmp, err = d.queryParser.Parse(f.Query); err != nil {
			return base.SetError(err)
		}

		cnd = append(cnd, tmp)
	}

	return base.Where(cnd...)
}

func (d *model) lookupSql(pkv crs.ValueGetter) *goqu.SelectDataset {
	var (
		sel       = d.selectSql().Limit(1)
		cond, err = d.pkLookupCondition(pkv)
	)

	if err != nil {
		sel = sel.SetError(err)
	}

	return sel.Where(cond)
}

func (d *model) selectSql() *goqu.SelectDataset {
	var (
		cols = d.table.Columns()

		// working around a bug inside goqu lib that adds
		// * to the list of columns to be selected
		// even if we clear the columns first
		q = d.dialect.GOQU().
			From(d.table.Ident()).
			Select(d.table.Ident().Col(cols[0].Name()))
	)

	for _, col := range cols[1:] {
		q = q.SelectAppend(d.table.Ident().Col(col.Name()))
	}

	return q
}

func (d *model) truncateSql() (_ *goqu.TruncateDataset) {
	return d.dialect.GOQU().Truncate(d.table)
}

func (d *model) insertSql(rr ...crs.ValueGetter) (_ *goqu.InsertDataset) {
	var (
		ins = d.dialect.GOQU().Insert(d.table.Ident())
		cc  = d.table.Columns()

		rows = make([][]any, len(rr))
		cols = make([]any, len(cc))

		err error
	)

	for c := range cc {
		cols[c] = cc[c].Name()
	}

	for i, r := range rr {
		rows[i], err = d.table.Encode(r)
		if err != nil {
			return ins.SetError(err)
		}
	}

	return ins.Cols(cols...).Vals(rows...)
}

// updateSql generates SQL command for updating record
func (d *model) updateSql(r crs.ValueGetter) *goqu.UpdateDataset {
	var (
		upd = d.dialect.GOQU().Update(d.table.Ident())

		values    = exp.Record{}
		condition = exp.Ex{}

		encoded, err = d.table.Encode(r)
	)

	if err != nil {
		return upd.SetError(err)
	}

	for i, c := range d.table.Columns() {
		if c.IsPrimaryKey() {
			// values[]
			condition[c.Name()] = encoded[i]
		} else {
			values[c.Name()] = encoded[i]
		}
	}

	return upd.Where(condition).Set(values)
}

func (d *model) deleteSql(pkv crs.ValueGetter) *goqu.DeleteDataset {
	var (
		del       = d.dialect.GOQU().Delete(d.table.Ident())
		cond, err = d.pkLookupCondition(pkv)
	)

	if err != nil {
		del.SetError(err)
	}

	return del.Where(cond)
}

// Constructs primary-key-lookup expression from pk values
func (d *model) pkLookupCondition(pkv crs.ValueGetter) (_ exp.Expression, err error) {
	var (
		cnd = exp.NewExpressionList(exp.AndType)
		val any
	)
	for _, c := range d.table.Columns() {
		if !c.IsPrimaryKey() {
			continue
		}

		val, err = pkv.GetValue(c.Name(), 0)
		if err != nil {
			return nil, fmt.Errorf("could not get value for primary key %q: %w", c.Name(), err)
		}

		cnd.Append(exp.NewBooleanExpression(exp.EqOp, d.table.Ident().Col(c.Name()), val))
	}

	return cnd, nil
}
