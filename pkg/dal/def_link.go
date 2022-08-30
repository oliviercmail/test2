package dal

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/filter"
)

type (
	// Link produces a series of left and corresponding right rows based on the
	// provided sources and the LinkPredicate.
	//
	// The Link step produces an SQL left join-like output where left and right
	// rows are served separately and the left rows are not duplicated.
	Link struct {
		Ident    string
		RelLeft  string
		RelRight string
		// @todo allow multiple link predicates; for now (for easier indexing)
		// only allow one (this is the same as we had before)
		On LinkPredicate
		// @todo consider splitting filter into left and right filter
		filter internalFilter
		Filter filter.Filter

		OutLeftAttributes  []AttributeMapping
		OutRightAttributes []AttributeMapping
		LeftAttributes     []AttributeMapping
		RightAttributes    []AttributeMapping

		relLeft  PipelineStep
		relRight PipelineStep

		plan     linkPlan
		analysis stepAnalysis
	}

	// LinkPredicate determines the attributes the two datasets should get joined on
	LinkPredicate struct {
		Left  string
		Right string
	}

	// linkPlan outlines how the optimizer determined the two datasets should be
	// joined on.
	linkPlan struct {
		// @todo add strategy when we have different strategies implemented
		// strategy string

		// partialScan indicates we can partially pull data from the two sources
		// as the data is provided in the correct order.
		partialScan bool
	}

	rowLink struct {
		a, b ValueGetter
	}
)

const (
	LinkRefIdent = "$sys.ref"
)

func (def *Link) Identifier() string {
	return def.Ident
}

func (def *Link) Sources() []string {
	return []string{def.RelLeft, def.RelRight}
}

func (def *Link) Attributes() [][]AttributeMapping {
	return [][]AttributeMapping{def.OutLeftAttributes, def.OutRightAttributes}
}

func (def *Link) Analyze(ctx context.Context) (err error) {
	def.analysis = stepAnalysis{
		scanCost:   costUnknown,
		searchCost: costUnknown,
		filterCost: costUnknown,
		sortCost:   costUnknown,
		outputSize: sizeUnknown,
	}
	return
}

func (def *Link) Analysis() stepAnalysis {
	return def.analysis
}

func (def *Link) Optimize(req internalFilter) (res internalFilter, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (def *Link) init(ctx context.Context) (err error) {
	if len(def.LeftAttributes) == 0 {
		def.LeftAttributes = collectAttributes(def.relLeft)
	}
	if len(def.RightAttributes) == 0 {
		def.RightAttributes = collectAttributes(def.relRight)
	}

	if len(def.OutLeftAttributes) == 0 {
		def.OutLeftAttributes = def.LeftAttributes
	}
	if len(def.OutRightAttributes) == 0 {
		def.OutRightAttributes = def.RightAttributes
	}

	if def.Filter != nil {
		def.filter, err = toInternalFilter(def.Filter)
		if err != nil {
			return
		}
	}

	err = def.validate()
	if err != nil {
		return
	}

	return nil
}

func (def *Link) exec(ctx context.Context, left, right Iterator) (_ Iterator, err error) {
	// @todo adjust the used exec based on other strategies when added
	exec := &linkLeft{
		def:         *def,
		filter:      def.filter,
		leftSource:  left,
		rightSource: right,
	}

	return exec, exec.init(ctx)
}

func (def *Link) validate() (err error) {
	err = func() (err error) {
		if len(def.OutLeftAttributes) == 0 {
			return fmt.Errorf("no left output attributes specified")
		}
		if len(def.OutRightAttributes) == 0 {
			return fmt.Errorf("no right output attributes specified")
		}
		if len(def.LeftAttributes) == 0 {
			return fmt.Errorf("no left attributes specified")
		}
		if len(def.RightAttributes) == 0 {
			return fmt.Errorf("no right attributes specified")
		}

		if def.On.Left == "" {
			return fmt.Errorf("no left attribute in the link predicate specified")
		}
		if def.On.Right == "" {
			return fmt.Errorf("no right attribute in the link predicate specified")
		}

		return
	}()
	if err != nil {
		return fmt.Errorf("invalid definition: %v", err)
	}
	return
}

func (r *rowLink) SelectGVal(ctx context.Context, k string) (interface{}, error) {
	return r.GetValue(k, 0)
}

func (r *rowLink) GetValue(name string, pos uint) (v any, err error) {
	a := r.a.CountValues()
	if cc, ok := a[name]; ok {
		if pos >= cc {
			return nil, nil
		}
		return r.a.GetValue(name, pos)
	}

	b := r.b.CountValues()
	if cc, ok := b[name]; ok {
		if pos >= cc {
			return nil, nil
		}
		return r.b.GetValue(name, pos)
	}

	return
}

func (r *rowLink) CountValues() (out map[string]uint) {
	out = make(map[string]uint)

	for k, c := range r.a.CountValues() {
		out[k] = c
	}
	for k, c := range r.b.CountValues() {
		out[k] = c
	}

	return
}
