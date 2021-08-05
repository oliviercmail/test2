package qlng

type (
	ASTNode struct {
		pMeta *parserMeta

		Ref  string     `json:"ref,omitempty"`
		Args ASTNodeSet `json:"args,omitempty"`

		Symbol string          `json:"symbol,omitempty"`
		Value  *typedValueWrap `json:"value,omitempty"`

		Raw string `json:"raw,omitempty"`
	}
	ASTNodeSet     []*ASTNode
	typedValueWrap struct {
		Value interface{} `json:"@value"`
		Type  string      `json:"@type"`
	}

	parserMeta struct {
		opDef *opDef
	}
)

func MakeValueOf(t string, v interface{}) *typedValueWrap {
	return &typedValueWrap{
		Type:  t,
		Value: v,
	}
}

// Traverse traverses the AST down to leaf nodes.
//
// If fnc. returns false, the traversal of the current branch ends.
func (n *ASTNode) Traverse(f func(*ASTNode) (bool, *ASTNode, error)) (err error) {
	var ok bool
	var r *ASTNode
	if n == nil {
		return nil
	}

	ok, r, err = f(n)
	if err != nil {
		return err
	}
	*n = *r
	if !ok {
		return
	}

	for _, a := range n.Args {
		if err = a.Traverse(f); err != nil {
			return
		}
	}

	return
}

func (n ASTNode) Clone() *ASTNode {
	aa := n.Args

	if n.Args != nil {
		n.Args = make(ASTNodeSet, len(aa))
		for i, a := range aa {
			n.Args[i] = a.Clone()
		}
	}

	if n.Value != nil {
		n.Value = &typedValueWrap{
			Type:  n.Value.Type,
			Value: n.Value.Value,
		}
	}

	return &n
}

func (n lNull) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: "null",
	}
}

func (n lBoolean) ToAST() (out *ASTNode) {
	return &ASTNode{
		Value: MakeValueOf("Boolean", n.value),
	}
}

func (n lString) ToAST() (out *ASTNode) {
	return &ASTNode{
		Value: MakeValueOf("String", n.value),
	}
}

// @todo differentiate between floats and others
func (n lNumber) ToAST() (out *ASTNode) {
	if isFloaty(n.value) {
		return &ASTNode{
			Value: MakeValueOf("Float", n.value),
		}
	} else {
		return &ASTNode{
			Value: MakeValueOf("Int", n.value),
		}
	}
}

func (n operator) ToAST() (out *ASTNode) {
	op := getOp(n.kind)

	return &ASTNode{
		Ref:   op.name,
		Args:  make(ASTNodeSet, 0, 2),
		pMeta: &parserMeta{opDef: op},
	}
}

func (n Ident) ToAST() (out *ASTNode) {
	return &ASTNode{
		Symbol: n.Value,
	}
}

func (n keyword) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: n.keyword,
	}
}

func (n interval) ToAST() (out *ASTNode) {
	return &ASTNode{
		Ref: "interval",
		Args: ASTNodeSet{
			{Symbol: n.unit},
			{Value: MakeValueOf("Number", n.value)},
		},
	}
}

func (n function) ToAST() (out *ASTNode) {
	auxA := n.arguments.ToAST()

	return &ASTNode{
		Ref:  n.name,
		Args: auxA.Args,
	}
}

func (nn parserNodeSet) ToAST() (out *ASTNode) {
	auxArgs := make(ASTNodeSet, 0, len(nn))

	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	return &ASTNode{
		Ref:  "group",
		Args: auxArgs,
	}
}

func (nn parserNodes) ToAST() (out *ASTNode) {
	// Prep
	auxArgs := make(ASTNodeSet, 0, len(nn))

	// Convert the entire level to AST nodes
	for _, n := range nn {
		auxArgs = append(auxArgs, n.ToAST())
	}

	// In the current level, have the operators consume operands
	// based on their defined weight
	//
	// - find the highest prio. op and have it consume what it wants
	// - repeat until all ops are satisfied
	// -- post optimizations?
	for {
		var bestOp *opDef
		bestOpIx := -1

		// We're done when it is reduced to 1
		if len(auxArgs) <= 1 {
			break
		}

		for _i, _a := range auxArgs {
			i := _i
			a := _a

			// use this as a delimiter for now
			if a.pMeta == nil || a.pMeta.opDef == nil {
				continue
			}

			if bestOp == nil {
				bestOp = a.pMeta.opDef
				bestOpIx = i
				continue
			}

			if a.pMeta.opDef.weight < bestOp.weight {
				bestOp = a.pMeta.opDef
				bestOpIx = i
				continue
			}
		}

		// Have the op consume what it needs.
		// Currently we only have binary operators (the !unary ones) so this is fine.
		arg := auxArgs[bestOpIx]
		arg.Args = append(arg.Args, auxArgs[bestOpIx-1], auxArgs[bestOpIx+1])
		// this is not needed anymore so we can remove it
		arg.pMeta = nil

		// Remove the consumed bits and replace it with the new bit
		aux := auxArgs[0 : bestOpIx-1]
		aux = append(aux, arg)
		// +1 for right side, +1 because the left index is inclusive
		aux = append(aux, auxArgs[bestOpIx+2:]...)
		auxArgs = aux
	}

	return &ASTNode{
		Ref:  "group",
		Args: auxArgs,
	}
}
