package dal

type (
	modelDiffType     string
	ModelModification string
	// ModelDiff defines one identified missmatch between two models
	ModelDiff struct {
		Type         modelDiffType
		Modification ModelModification
		// Original will be nil when a new attribute is being added
		Original *Attribute
		// Inserted will be nil wen an existing attribute is being removed
		Inserted *Attribute
	}

	ModelDiffSet []*ModelDiff
)

const (
	AttributeMissing             modelDiffType     = "attributeMissing"
	AttributeTypeMissmatch       modelDiffType     = "typeMissmatch"
	AttributeSensitivityMismatch modelDiffType     = "sensitivityMismatch"
	AttributeCodecMismatch       modelDiffType     = "codecMismatch"
	AttributeDeleted             ModelModification = "deleted"
	AttributeAdded               ModelModification = "added"
	AttributeChanged             ModelModification = "changed"
)

// Diff calculates the diff between models a and b where a is used as base
func (a *Model) Diff(b *Model) (out ModelDiffSet) {
	if a == nil {
		a = &Model{}
	}
	if b == nil {
		b = &Model{}
	}

	bIndex := make(map[string]struct {
		found bool
		attr  *Attribute
	})
	for _, _attr := range b.Attributes {
		attr := _attr
		bIndex[attr.Ident] = struct {
			found bool
			attr  *Attribute
		}{
			attr: attr,
		}
	}

	aIndex := make(map[string]struct {
		found bool
		attr  *Attribute
	})
	for _, _attr := range a.Attributes {
		attr := _attr
		aIndex[attr.Ident] = struct {
			found bool
			attr  *Attribute
		}{
			attr: attr,
		}
	}

	// Deleted and update ones
	for _, _attrA := range a.Attributes {
		attrA := _attrA
		// store is an interface to something that could be a pointer.
		// we need to copy it to make sure we don't get a nil pointer
		// make sure not to modify this since it would modify the original
		attrA.Store = _attrA.Store

		// Missmatches
		attrBAux, ok := bIndex[attrA.Ident]
		if !ok {
			out = append(out, &ModelDiff{
				Type:         AttributeMissing,
				Modification: AttributeDeleted,
				Original:     attrA,
			})
			continue
		}

		// Typecheck
		if attrA.Type.Type() != attrBAux.attr.Type.Type() {
			out = append(out, &ModelDiff{
				Type:         AttributeTypeMissmatch,
				Modification: AttributeChanged,
				Original:     attrA,
				Inserted:     attrBAux.attr,
			})
		}

		// Other stuff
		// @todo improve; for now it'll do
		if attrA.SensitivityLevelID != attrBAux.attr.SensitivityLevelID {
			out = append(out, &ModelDiff{
				Type:         AttributeSensitivityMismatch,
				Modification: AttributeChanged,
				Original:     attrA,
				Inserted:     attrBAux.attr,
			})
		}
		if attrA.Store.Type() != attrBAux.attr.Store.Type() {
			out = append(out, &ModelDiff{
				Type:         AttributeCodecMismatch,
				Modification: AttributeChanged,
				Original:     attrA,
				Inserted:     attrBAux.attr,
			})
		}
	}

	// New
	for _, _attrB := range b.Attributes {
		attrB := _attrB

		// Missmatches
		_, ok := aIndex[attrB.Ident]
		if !ok {
			out = append(out, &ModelDiff{
				Type:         AttributeMissing,
				Modification: AttributeAdded,
				Original:     nil,
				Inserted:     attrB,
			})
			continue
		}
	}

	return
}

func (dd ModelDiffSet) Alterations() (out []*Alteration) {
	add := func(a *Alteration) {
		out = append(out, a)
	}

	for _, d := range dd {
		switch d.Type {
		case AttributeMissing:
			if d.Inserted == nil {
				// @todo if this was the last attribute we can consider dropping this column
				if d.Original.Store.Type() == AttributeCodecRecordValueSetJSON {
					break
				}

				add(&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: d.Original,
					},
				})
			} else {
				if d.Inserted.Store.Type() == AttributeCodecRecordValueSetJSON {
					add(&Alteration{
						AttributeAdd: &AttributeAdd{
							Attr: &Attribute{
								Ident: d.Inserted.StoreIdent(),
								Type:  &TypeJSON{Nullable: false},
								Store: &CodecPlain{},
							},
						},
					})
				} else {
					add(&Alteration{
						AttributeAdd: &AttributeAdd{
							Attr: d.Inserted,
						},
					})
				}

			}

		case AttributeTypeMissmatch:
			// @todo we might have to do some validation earlier on
			if d.Original.Store.Type() == AttributeCodecRecordValueSetJSON {
				break
			}

			add(&Alteration{
				AttributeReType: &AttributeReType{
					Attr: d.Inserted,
					To:   d.Inserted.Type,
				},
			})

		case AttributeCodecMismatch:
			add(&Alteration{
				AttributeReEncode: &AttributeReEncode{
					Attr: d.Inserted,
					To:   d.Inserted.Store,
				},
			})
		}
	}

	return
}
