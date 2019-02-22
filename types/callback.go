package types

import (
	"github.com/gowebapi/webidlparser/ast"
)

type Callback struct {
	standardType
	basic      BasicInfo
	Return     TypeRef
	Parameters []*Parameter
	source     *ast.Callback
}

// Callback need to implement Type
var _ Type = &Callback{}

func (t *extractTypes) convertCallback(in *ast.Callback) *Callback {
	params := t.convertParams(in.Parameters)
	ret := &Callback{
		standardType: standardType{
			ref:         createRef(in, t),
			needRelease: false,
		},
		basic:      fromIdlToTypeName(t.main.setup.Package, in.Name, "callback"),
		source:     in,
		Return:     convertType(in.Return, t),
		Parameters: params,
	}
	return ret
}

func (t *Callback) Basic() BasicInfo {
	return TransformBasic(t, t.basic)
}

func (t *Callback) DefaultParam() (info *TypeInfo, inner TypeRef) {
	return t.Param(false, false, false)
}

func (t *Callback) key() string {
	return t.basic.Idl
}

func (t *Callback) lessThan(b *Callback) bool {
	return t.basic.lessThan(&b.basic)
}

func (t *Callback) link(conv *Convert, inuse inuseLogic) TypeRef {
	if t.inuse {
		return t
	}
	t.inuse = true

	t.Return = t.Return.link(conv, make(inuseLogic))
	for i := range t.Parameters {
		t.Parameters[i].Type = t.Parameters[i].Type.link(conv, make(inuseLogic))
	}
	return t
}

func (t *Callback) Param(nullable, option, variadic bool) (info *TypeInfo, inner TypeRef) {
	basic := t.Basic()
	output := newTypeInfo(basic, nullable, option, variadic, false, true, false)
	info, typ := newTypeInfo(basic, nullable, option, variadic, true, false, false), t
	info.Output = output.Output + "Func"
	info.Pointer = true
	return info, typ
}

func (t *Callback) SetBasic(basic BasicInfo) {
	t.basic = basic
}

func (t *Callback) TypeID() TypeID {
	return TypeCallback
}
