package rdbms

// This file is an auto-generated file
//
// Template:    pkg/codegen/assets/store_rdbms.gen.go.tpl
// Definitions: store/resource_activity_log.yaml
//
// Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated.

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/cortezaproject/corteza-server/pkg/discovery/types"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/store/rdbms/builders"
)

var _ = errors.Is

// SearchResourceActivityLogs returns all matching rows
//
// This function calls convertResourceActivityLogFilter with the given
// types.ResourceActivityFilter and expects to receive a working squirrel.SelectBuilder
func (s Store) SearchResourceActivityLogs(ctx context.Context, f types.ResourceActivityFilter) (types.ResourceActivitySet, types.ResourceActivityFilter, error) {
	var (
		err error
		set []*types.ResourceActivity
		q   squirrel.SelectBuilder
	)

	return set, f, func() error {
		q, err = s.convertResourceActivityLogFilter(f)
		if err != nil {
			return err
		}

		// Paging enabled
		// {search: {enablePaging:true}}
		// Cleanup unwanted cursor values (only relevant is f.PageCursor, next&prev are reset and returned)
		f.PrevPage, f.NextPage = nil, nil

		if f.PageCursor != nil {
			// Page cursor exists so we need to validate it against used sort
			// To cover the case when paging cursor is set but sorting is empty, we collect the sorting instructions
			// from the cursor.
			// This (extracted sorting info) is then returned as part of response
			if f.Sort, err = f.PageCursor.Sort(f.Sort); err != nil {
				return err
			}
		}

		// Make sure results are always sorted at least by primary keys
		if f.Sort.Get("id") == nil {
			f.Sort = append(f.Sort, &filter.SortExpr{
				Column:     "id",
				Descending: f.Sort.LastDescending(),
			})
		}

		// Cloned sorting instructions for the actual sorting
		// Original are passed to the fetchFullPageOfUsers fn used for cursor creation so it MUST keep the initial
		// direction information
		sort := f.Sort.Clone()

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		if f.PageCursor != nil && f.PageCursor.ROrder {
			sort.Reverse()
		}

		// Apply sorting expr from filter to query
		if q, err = setOrderBy(q, sort, s.sortableResourceActivityLogColumns(), s.Config().SqlSortHandler); err != nil {
			return err
		}

		set, f.PrevPage, f.NextPage, err = s.fetchFullPageOfResourceActivityLogs(
			ctx,
			q, f.Sort, f.PageCursor,
			f.Limit,
			nil,
			func(cur *filter.PagingCursor) squirrel.Sqlizer {
				return builders.CursorCondition(cur, nil)
			},
		)

		if err != nil {
			return err
		}

		f.PageCursor = nil
		return nil
	}()
}

// fetchFullPageOfResourceActivityLogs collects all requested results.
//
// Function applies:
//  - cursor conditions (where ...)
//  - limit
//
// Main responsibility of this function is to perform additional sequential queries in case when not enough results
// are collected due to failed check on a specific row (by check fn).
//
// Function then moves cursor to the last item fetched
func (s Store) fetchFullPageOfResourceActivityLogs(
	ctx context.Context,
	q squirrel.SelectBuilder,
	sort filter.SortExprSet,
	cursor *filter.PagingCursor,
	reqItems uint,
	check func(*types.ResourceActivity) (bool, error),
	cursorCond func(*filter.PagingCursor) squirrel.Sqlizer,
) (set []*types.ResourceActivity, prev, next *filter.PagingCursor, err error) {
	var (
		aux []*types.ResourceActivity

		// When cursor for a previous page is used it's marked as reversed
		// This tells us to flip the descending flag on all used sort keys
		reversedOrder = cursor != nil && cursor.ROrder

		// copy of the select builder
		tryQuery squirrel.SelectBuilder

		// Copy no. of required items to limit
		// Limit will change when doing subsequent queries to fill
		// the set with all required items
		limit = reqItems

		// cursor to prev. page is only calculated when cursor is used
		hasPrev = cursor != nil

		// next cursor is calculated when there are more pages to come
		hasNext bool
	)

	set = make([]*types.ResourceActivity, 0, DefaultSliceCapacity)

	for try := 0; try < MaxRefetches; try++ {
		if cursor != nil {
			tryQuery = q.Where(cursorCond(cursor))
		} else {
			tryQuery = q
		}

		if limit > 0 {
			// fetching + 1 so we know if there are more items
			// we can fetch (next-page cursor)
			tryQuery = tryQuery.Limit(uint64(limit + 1))
		}

		if aux, err = s.QueryResourceActivityLogs(ctx, tryQuery, check); err != nil {
			return nil, nil, nil, err
		}

		if len(aux) == 0 {
			// nothing fetched
			break
		}

		// append fetched items
		set = append(set, aux...)

		if reqItems == 0 {
			// no max requested items specified, break out
			break
		}

		collected := uint(len(set))

		if reqItems > collected {
			// not enough items fetched, try again with adjusted limit
			limit = reqItems - collected

			if limit < MinEnsureFetchLimit {
				// In case limit is set very low and we've missed records in the first fetch,
				// make sure next fetch limit is a bit higher
				limit = MinEnsureFetchLimit
			}

			// Update cursor so that it points to the last item fetched
			cursor = s.collectResourceActivityLogCursorValues(set[collected-1], sort...)

			// Copy reverse flag from sorting
			cursor.LThen = sort.Reversed()
			continue
		}

		if reqItems < collected {
			set = set[:reqItems]
			hasNext = true
		}

		break
	}

	collected := len(set)

	if collected == 0 {
		return nil, nil, nil, nil
	}

	if reversedOrder {
		// Fetched set needs to be reversed because we've forced a descending order to get the previous page
		for i, j := 0, collected-1; i < j; i, j = i+1, j-1 {
			set[i], set[j] = set[j], set[i]
		}

		// when in reverse-order rules on what cursor to return change
		hasPrev, hasNext = hasNext, hasPrev
	}

	if hasPrev {
		prev = s.collectResourceActivityLogCursorValues(set[0], sort...)
		prev.ROrder = true
		prev.LThen = !sort.Reversed()
	}

	if hasNext {
		next = s.collectResourceActivityLogCursorValues(set[collected-1], sort...)
		next.LThen = sort.Reversed()
	}

	return set, prev, next, nil
}

// QueryResourceActivityLogs queries the database, converts and checks each row and
// returns collected set
//
// Fn also returns total number of fetched items and last fetched item so that the caller can construct cursor
// for next page of results
func (s Store) QueryResourceActivityLogs(
	ctx context.Context,
	q squirrel.Sqlizer,
	check func(*types.ResourceActivity) (bool, error),
) ([]*types.ResourceActivity, error) {
	var (
		tmp = make([]*types.ResourceActivity, 0, DefaultSliceCapacity)
		set = make([]*types.ResourceActivity, 0, DefaultSliceCapacity)
		res *types.ResourceActivity

		// Query rows with
		rows, err = s.Query(ctx, q)
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		if err = rows.Err(); err == nil {
			res, err = s.internalResourceActivityLogRowScanner(rows)
		}

		if err != nil {
			return nil, err
		}

		tmp = append(tmp, res)
	}

	for _, res = range tmp {

		set = append(set, res)
	}

	return set, nil
}

// LookupResourceActivityLogByID searches for corteza resource activity by ID
// It returns corteza resource activity even if deleted
func (s Store) LookupResourceActivityLogByID(ctx context.Context, id uint64) (*types.ResourceActivity, error) {
	return s.execLookupResourceActivityLog(ctx, squirrel.Eq{
		s.preprocessColumn("ral.id", ""): store.PreprocessValue(id, ""),
	})
}

// CreateResourceActivityLog creates one or more rows in resource_activity_log table
func (s Store) CreateResourceActivityLog(ctx context.Context, rr ...*types.ResourceActivity) (err error) {
	for _, res := range rr {
		err = s.checkResourceActivityLogConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execCreateResourceActivityLogs(ctx, s.internalResourceActivityLogEncoder(res))
		if err != nil {
			return err
		}
	}

	return
}

// UpdateResourceActivityLog updates one or more existing rows in resource_activity_log
func (s Store) UpdateResourceActivityLog(ctx context.Context, rr ...*types.ResourceActivity) error {
	return s.partialResourceActivityLogUpdate(ctx, nil, rr...)
}

// partialResourceActivityLogUpdate updates one or more existing rows in resource_activity_log
func (s Store) partialResourceActivityLogUpdate(ctx context.Context, onlyColumns []string, rr ...*types.ResourceActivity) (err error) {
	for _, res := range rr {
		err = s.checkResourceActivityLogConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpdateResourceActivityLogs(
			ctx,
			squirrel.Eq{
				s.preprocessColumn("ral.id", ""): store.PreprocessValue(res.ID, ""),
			},
			s.internalResourceActivityLogEncoder(res).Skip("id").Only(onlyColumns...))
		if err != nil {
			return err
		}
	}

	return
}

// UpsertResourceActivityLog updates one or more existing rows in resource_activity_log
func (s Store) UpsertResourceActivityLog(ctx context.Context, rr ...*types.ResourceActivity) (err error) {
	for _, res := range rr {
		err = s.checkResourceActivityLogConstraints(ctx, res)
		if err != nil {
			return err
		}

		err = s.execUpsertResourceActivityLogs(ctx, s.internalResourceActivityLogEncoder(res))
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteResourceActivityLog Deletes one or more rows from resource_activity_log table
func (s Store) DeleteResourceActivityLog(ctx context.Context, rr ...*types.ResourceActivity) (err error) {
	for _, res := range rr {

		err = s.execDeleteResourceActivityLogs(ctx, squirrel.Eq{
			s.preprocessColumn("ral.id", ""): store.PreprocessValue(res.ID, ""),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// DeleteResourceActivityLogByID Deletes row from the resource_activity_log table
func (s Store) DeleteResourceActivityLogByID(ctx context.Context, ID uint64) error {
	return s.execDeleteResourceActivityLogs(ctx, squirrel.Eq{
		s.preprocessColumn("ral.id", ""): store.PreprocessValue(ID, ""),
	})
}

// TruncateResourceActivityLogs Deletes all rows from the resource_activity_log table
func (s Store) TruncateResourceActivityLogs(ctx context.Context) error {
	return s.Truncate(ctx, s.resourceActivityLogTable())
}

// execLookupResourceActivityLog prepares ResourceActivityLog query and executes it,
// returning types.ResourceActivity (or error)
func (s Store) execLookupResourceActivityLog(ctx context.Context, cnd squirrel.Sqlizer) (res *types.ResourceActivity, err error) {
	var (
		row rowScanner
	)

	row, err = s.QueryRow(ctx, s.resourceActivityLogsSelectBuilder().Where(cnd))
	if err != nil {
		return
	}

	res, err = s.internalResourceActivityLogRowScanner(row)
	if err != nil {
		return
	}

	return res, nil
}

// execCreateResourceActivityLogs updates all matched (by cnd) rows in resource_activity_log with given data
func (s Store) execCreateResourceActivityLogs(ctx context.Context, payload store.Payload) error {
	return s.Exec(ctx, s.InsertBuilder(s.resourceActivityLogTable()).SetMap(payload))
}

// execUpdateResourceActivityLogs updates all matched (by cnd) rows in resource_activity_log with given data
func (s Store) execUpdateResourceActivityLogs(ctx context.Context, cnd squirrel.Sqlizer, set store.Payload) error {
	return s.Exec(ctx, s.UpdateBuilder(s.resourceActivityLogTable("ral")).Where(cnd).SetMap(set))
}

// execUpsertResourceActivityLogs inserts new or updates matching (by-primary-key) rows in resource_activity_log with given data
func (s Store) execUpsertResourceActivityLogs(ctx context.Context, set store.Payload) error {
	upsert, err := s.config.UpsertBuilder(
		s.config,
		s.resourceActivityLogTable(),
		set,
		s.preprocessColumn("id", ""),
	)

	if err != nil {
		return err
	}

	return s.Exec(ctx, upsert)
}

// execDeleteResourceActivityLogs Deletes all matched (by cnd) rows in resource_activity_log with given data
func (s Store) execDeleteResourceActivityLogs(ctx context.Context, cnd squirrel.Sqlizer) error {
	return s.Exec(ctx, s.DeleteBuilder(s.resourceActivityLogTable("ral")).Where(cnd))
}

func (s Store) internalResourceActivityLogRowScanner(row rowScanner) (res *types.ResourceActivity, err error) {
	res = &types.ResourceActivity{}

	err = row.Scan(
		&res.ID,
		&res.ResourceID,
		&res.ResourceType,
		&res.ResourceAction,
		&res.Timestamp,
		&res.Meta,
	)

	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound.Stack(1)
	}

	if err != nil {
		return nil, errors.Store("could not scan resourceActivityLog db row: %s", err).Wrap(err)
	} else {
		return res, nil
	}
}

// QueryResourceActivityLogs returns squirrel.SelectBuilder with set table and all columns
func (s Store) resourceActivityLogsSelectBuilder() squirrel.SelectBuilder {
	return s.SelectBuilder(s.resourceActivityLogTable("ral"), s.resourceActivityLogColumns("ral")...)
}

// resourceActivityLogTable name of the db table
func (Store) resourceActivityLogTable(aa ...string) string {
	var alias string
	if len(aa) > 0 {
		alias = " AS " + aa[0]
	}

	return "resource_activity_log" + alias
}

// ResourceActivityLogColumns returns all defined table columns
//
// With optional string arg, all columns are returned aliased
func (Store) resourceActivityLogColumns(aa ...string) []string {
	var alias string
	if len(aa) > 0 {
		alias = aa[0] + "."
	}

	return []string{
		alias + "id",
		alias + "rel_resource",
		alias + "resource_type",
		alias + "resource_action",
		alias + "ts",
		alias + "meta",
	}
}

// {true true false true true false}

// sortableResourceActivityLogColumns returns all ResourceActivityLog columns flagged as sortable
//
// With optional string arg, all columns are returned aliased
func (Store) sortableResourceActivityLogColumns() map[string]string {
	return map[string]string{
		"id": "id",
	}
}

// internalResourceActivityLogEncoder encodes fields from types.ResourceActivity to store.Payload (map)
//
// Encoding is done by using generic approach or by calling encodeResourceActivityLog
// func when rdbms.customEncoder=true
func (s Store) internalResourceActivityLogEncoder(res *types.ResourceActivity) store.Payload {
	return store.Payload{
		"id":              res.ID,
		"rel_resource":    res.ResourceID,
		"resource_type":   res.ResourceType,
		"resource_action": res.ResourceAction,
		"ts":              res.Timestamp,
		"meta":            res.Meta,
	}
}

// collectResourceActivityLogCursorValues collects values from the given resource that and sets them to the cursor
// to be used for pagination
//
// Values that are collected must come from sortable, unique or primary columns/fields
// At least one of the collected columns must be flagged as unique, otherwise fn appends primary keys at the end
//
// Known issue:
//   when collecting cursor values for query that sorts by unique column with partial index (ie: unique handle on
//   undeleted items)
func (s Store) collectResourceActivityLogCursorValues(res *types.ResourceActivity, cc ...*filter.SortExpr) *filter.PagingCursor {
	var (
		cursor = &filter.PagingCursor{LThen: filter.SortExprSet(cc).Reversed()}

		hasUnique bool

		// All known primary key columns

		pkId bool

		collect = func(cc ...*filter.SortExpr) {
			for _, c := range cc {
				switch c.Column {
				case "id":
					cursor.Set(c.Column, res.ID, c.Descending)

					pkId = true

				}
			}
		}
	)

	collect(cc...)
	if !hasUnique || !(pkId && true) {
		collect(&filter.SortExpr{Column: "id", Descending: false})
	}

	return cursor
}

// checkResourceActivityLogConstraints performs lookups (on valid) resource to check if any of the values on unique fields
// already exists in the store
//
// Using built-in constraint checking would be more performant but unfortunately we cannot rely
// on the full support (MySQL does not support conditional indexes)
func (s *Store) checkResourceActivityLogConstraints(ctx context.Context, res *types.ResourceActivity) error {
	// Consider resource valid when all fields in unique constraint check lookups
	// have valid (non-empty) value
	//
	// Only string and uint64 are supported for now
	// feel free to add additional types if needed
	var valid = true

	if !valid {
		return nil
	}

	var checks = make([]func() error, 0)

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}
