package filter

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type (
	// Paging is a helper struct that should be embedded in filter types
	// to help with the paging
	Paging struct {
		// How many items per fetch do we want
		Limit uint `json:"limit,omitempty"`

		// PeakNextPage bool `json:"peakNextPage,omitempty"`

		// Opaque pointer to 1st item on page
		// read-only
		PageCursor *PagingCursor `json:"cursor,omitempty"`

		// Opaque cursor that points to the first item on the next page
		// read-only
		NextPage *PagingCursor `json:"nextPage,omitempty"`

		// Opaque cursor that points to the first item on the previous page
		// value of {cursor} will be copied here
		// read-only
		PrevPage *PagingCursor `json:"prevPage,omitempty"`
	}

	PagingCursor struct {
		keys    []string
		values  []interface{}
		desc    []bool
		Reverse bool
	}

	pagingCursorValue struct {
		v interface{}
	}
)

func NewPaging(limit uint, cursor string) (p Paging, err error) {
	p = Paging{Limit: limit}

	if p.PageCursor, err = parseCursor(cursor); err != nil {
		return
	}

	return
}

func (p *PagingCursor) Set(k string, v interface{}, d bool) {
	for i, key := range p.keys {
		if key == k {
			p.values[i] = v
			return
		}
	}

	p.keys = append(p.keys, k)
	p.values = append(p.values, v)
	p.desc = append(p.desc, d)
}

func (p *PagingCursor) Keys() []string {
	return p.keys
}

func (p *PagingCursor) Values() []interface{} {
	return p.values
}

// Stirng to implement Stringer and to get human-readable representation of the cursor
func (p *PagingCursor) String() string {
	var o = "<"

	for i, key := range p.keys {
		o += fmt.Sprintf("%s: %v", key, p.values[i])
		if p.desc[i] {
			o += " DESC"
		}
		o += ", "
	}

	if p.Reverse {
		o += "reverse"
	} else {
		o += "forward"
	}

	return o + ">"
}

// MarshalJSON serializes cursor struct as JSON and encodes it as base64 + adds quotes to be treated as JSON string
func (p *PagingCursor) MarshalJSON() ([]byte, error) {
	buf, err := json.Marshal(struct {
		K []string
		V []interface{}
		D []bool
		R bool
	}{
		p.keys,
		p.values,
		p.desc,
		p.Reverse,
	})

	if err != nil {
		return nil, err
	}

	std := base64.StdEncoding
	dbuf := make([]byte, std.EncodedLen(len(buf)))
	std.Encode(dbuf, buf)

	return append([]byte{'"'}, append(dbuf, '"')...), nil
}

func (p *PagingCursor) Encode() string {
	b, _ := p.MarshalJSON()
	return string(b)
}

func (p *PagingCursor) UnmarshalJSON(in []byte) error {
	var (
		aux struct {
			K []string
			V []pagingCursorValue
			D []bool
			R bool
		}
	)

	if err := json.Unmarshal(in, &aux); err != nil {
		return err
	}

	p.keys = aux.K
	p.desc = aux.D
	p.Reverse = aux.R

	// json.Unmarshal treats uint64 in values ([]interface{}) as float64 and we don't like that.
	p.values = make([]interface{}, len(aux.V))
	for i, v := range aux.V {
		p.values[i] = v.v
	}

	return nil
}

func (p *PagingCursor) Decode(cursor string) error {
	b, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return err
	}

	err = p.UnmarshalJSON(b)
	if err != nil {
		return err
	}

	return err
}

// Sort returns:
//  - sort if cursor is nil
//  - sort from cursor when sort is empty
//  - sort from cursor when sort is compatible with cursor
//  - error if sort & cursor are incompatible
func (p *PagingCursor) Sort(sort SortExprSet) (SortExprSet, error) {
	if p == nil {
		return sort, nil
	}

	if len(sort) == 0 {
		// sort emprt, create it from cursor
		for k := range p.keys {
			sort = append(sort, &SortExpr{
				Column:     p.keys[k],
				Descending: p.desc[k],
			})
		}
		return sort, nil
	}

	// check compatibility
	ss := sort.Columns()
	if len(p.keys) != len(ss) {
		return nil, fmt.Errorf("incompatible sort")
	}

	for k := range p.keys {
		if p.keys[k] != ss[k] {
			return nil, fmt.Errorf("incompatible sort")
		}
	}

	return sort, nil
}

func parseCursor(in string) (p *PagingCursor, err error) {
	if len(in) == 0 {
		return nil, nil
	}

	var buf []byte
	if buf, err = base64.StdEncoding.DecodeString(in); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	p = &PagingCursor{}
	if err = p.UnmarshalJSON(buf); err != nil {
		return nil, fmt.Errorf("could not decode cursor: %w", err)
	}

	return p, nil
}

// Making sure uint64 other int* values are properly unmarshaled
func (v *pagingCursorValue) UnmarshalJSON(in []byte) (err error) {
	var (
		u uint64
		i int64
	)

	if err = json.Unmarshal(in, &u); err == nil {
		v.v = u
		return
	}

	if err = json.Unmarshal(in, &i); err == nil {
		v.v = i
		return
	}

	return json.Unmarshal(in, &v.v)
}
