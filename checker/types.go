package checker

import (
	"text/template/parse"
	"sync"
	"reflect"
	"text/template"
	"io"
	"github.com/vinkdong/gox/log"
	"sort"
	"strings"
	"fmt"
	"errors"
)

var builtins = template.FuncMap{
	"and": and,
	"eq":  eq,
}


type kind int

const (
	invalidKind kind = iota
	boolKind
	complexKind
	intKind
	floatKind
	stringKind
	uintKind
)
var (
	errBadComparisonType = errors.New("invalid type for comparison")
	errBadComparison     = errors.New("incompatible types for comparison")
	errNoComparison      = errors.New("missing argument for comparison")
)

func basicKind(v reflect.Value) (kind, error) {
	switch v.Kind() {
	case reflect.Bool:
		return boolKind, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intKind, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintKind, nil
	case reflect.Float32, reflect.Float64:
		return floatKind, nil
	case reflect.Complex64, reflect.Complex128:
		return complexKind, nil
	case reflect.String:
		return stringKind, nil
	}
	return invalidKind, errBadComparisonType
}

func eq(arg1 reflect.Value, arg2 ...reflect.Value) (bool, error) {
	v1 := indirectInterface(arg1)
	k1, err := basicKind(v1)
	if err != nil {
		return false, err
	}
	if len(arg2) == 0 {
		return false, errNoComparison
	}
	for _, arg := range arg2 {
		v2 := indirectInterface(arg)
		k2, err := basicKind(v2)
		if err != nil {
			return false, err
		}
		truth := false
		if k1 != k2 {
			// Special case: Can compare integer values regardless of type's sign.
			switch {
			case k1 == intKind && k2 == uintKind:
				truth = v1.Int() >= 0 && uint64(v1.Int()) == v2.Uint()
			case k1 == uintKind && k2 == intKind:
				truth = v2.Int() >= 0 && v1.Uint() == uint64(v2.Int())
			default:
				return false, errBadComparison
			}
		} else {
			switch k1 {
			case boolKind:
				truth = v1.Bool() == v2.Bool()
			case complexKind:
				truth = v1.Complex() == v2.Complex()
			case floatKind:
				truth = v1.Float() == v2.Float()
			case intKind:
				truth = v1.Int() == v2.Int()
			case stringKind:
				truth = v1.String() == v2.String()
			case uintKind:
				truth = v1.Uint() == v2.Uint()
			default:
				panic("invalid kind")
			}
		}
		if truth {
			return true, nil
		}
	}
	return false, nil
}

func truth(arg reflect.Value) bool {
	t, _ := isTrue(indirectInterface(arg))
	return t
}

type Checker struct {
	name string
	*parse.Tree
	*common
	leftDelim  string
	rightDelim string
}

type common struct {
	tmpl   map[string]*Checker // Map from name to defined templates.
	// We use two maps, one for parsing and one for execution.
	// This separation makes the API cleaner since it doesn't
	// expose reflection to the client.
	muFuncs    sync.RWMutex // protects parseFuncs and execFuncs
	parseFuncs template.FuncMap
	execFuncs  map[string]reflect.Value
}

func  (c *Checker) init()  {
	if c.common == nil {
		comm := new(common)
		comm.tmpl = make(map[string]*Checker)
		comm.parseFuncs = make(template.FuncMap)
		comm.execFuncs = make(map[string]reflect.Value)
		c.common = comm
	}
}

func New(name string) *Checker {
	c := &Checker{
		name: name,
	}
	c.init()
	return c
}

func (c *Checker) New(name string) *Checker {
	c.init()
	nt := &Checker{
		name:       name,
		common:     c.common,
		leftDelim:  c.leftDelim,
		rightDelim: c.rightDelim,
	}
	return nt
}

var zero reflect.Value

func (c *Checker) Parse(text string) (*Checker, error) {
	c.init()
	c.muFuncs.RLock()
	trees, err := parse.Parse(c.name, text, c.leftDelim, c.rightDelim, c.parseFuncs, builtins)
	c.muFuncs.RUnlock()
	if err != nil {
		return nil, err
	}
	// Add the newly parsed trees, including the one for t, into our common structure.
	for name, tree := range trees {
		if _, err := c.AddParseTree(name, tree); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Checker) AddParseTree(name string, tree *parse.Tree) (*Checker, error) {
	c.init()
	// If the name is the name of this template, overwrite this template.
	nt := c
	if name != c.name {
		nt = c.New(name)
	}
	// Even if nt == t, we need to install it in the common.tmpl map.
	if replace, err := c.associate(nt, tree); err != nil {
		return nil, err
	} else if replace || nt.Tree == nil {
		nt.Tree = tree
	}
	return nt, nil
}

func (c *Checker) associate(new *Checker, tree *parse.Tree) (bool, error) {
	if new.common != c.common {
		panic("internal error: associate not common")
	}
	if old := c.tmpl[new.name]; old != nil && parse.IsEmptyTree(tree.Root) && old.Tree != nil {
		// If a template by that name exists,
		// don't replace it with an empty template.
		return false, nil
	}
	c.tmpl[new.name] = new
	return true, nil
}


func (c *Checker) execute(wr io.Writer, data interface{}) (err error) {
	defer recover()
	value, ok := data.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(data)
	}
	state := &state{
		tmpl: c,
		wr:   wr,
		vars: []variable{{"$", value}},
	}
	if c.Tree == nil || c.Root == nil {
		panic(c.name)
	}
	state.walk(value, c.Root)
	return
}

type state struct {
	tmpl  *Checker
	wr    io.Writer
	node  parse.Node // current node, for errors
	vars  []variable // push-down stack of variable values.
	depth int        // the height of the stack of executing templates.
}

func (s *state) walk(dot reflect.Value, node parse.Node) {
	s.at(node)
	switch node := node.(type) {
	case *parse.IfNode:
		s.walkIfOrWith(parse.NodeIf, dot, node.Pipe, node.List, node.ElseList)
	case *parse.ListNode:
		for _, node := range node.Nodes {
			s.walk(dot, node)
		}
	case *parse.TextNode:
		if _, err := s.wr.Write(node.Text); err != nil {
			log.Errorf("LAL")
		}
	case *parse.WithNode:
		s.walkIfOrWith(parse.NodeWith, dot, node.Pipe, node.List, node.ElseList)
	default:
		log.Errorf("unknown node: %s", node)
	}
}

func (s *state) at(node parse.Node) {
	s.node = node
}

func (s *state) walkIfOrWith(typ parse.NodeType, dot reflect.Value, pipe *parse.PipeNode, list, elseList *parse.ListNode) {
	defer s.pop(s.mark())
	val := s.evalPipeline(dot, pipe)
	truth, ok := isTrue(val)
	if !ok {
		log.Errorf("if/with can't use %v", val)
	}
	if truth {
		if typ == parse.NodeWith {
			s.walk(val, list)
		} else {
			s.walk(dot, list)
		}
	} else if elseList != nil {
		s.walk(dot, elseList)
	}
}

func (s *state) pop(mark int) {
	s.vars = s.vars[0:mark]
}

func (s *state) mark() int {
	return len(s.vars)
}

type variable struct {
	name  string
	value reflect.Value
}

func (s *state) evalPipeline(dot reflect.Value, pipe *parse.PipeNode) (value reflect.Value) {
	if pipe == nil {
		return
	}
	s.at(pipe)
	for _, cmd := range pipe.Cmds {
		value = s.evalCommand(dot, cmd, value) // previous value is this one's final arg.
		// If the object has type interface{}, dig down one level to the thing inside.
		if value.Kind() == reflect.Interface && value.Type().NumMethod() == 0 {
			value = reflect.ValueOf(value.Interface()) // lovely!
		}
	}
	for _, variable := range pipe.Decl {
		s.push(variable.Ident[0], value)
	}
	return value
}

func (s *state) push(name string, value reflect.Value) {
	s.vars = append(s.vars, variable{name, value})
}

func (s *state) evalCommand(dot reflect.Value, cmd *parse.CommandNode, final reflect.Value) reflect.Value {
	firstWord := cmd.Args[0]
	switch n := firstWord.(type) {
	case *parse.FieldNode:
		return s.evalFieldNode(dot, n, cmd.Args, final)
	case *parse.ChainNode:
		return s.evalChainNode(dot, n, cmd.Args, final)
	case *parse.PipeNode:
		// Parenthesized pipeline. The arguments are all inside the pipeline; final is ignored.
		return s.evalPipeline(dot, n)
	case *parse.VariableNode:
		return s.evalVariableNode(dot, n, cmd.Args, final)
	}
	s.at(firstWord)
	s.notAFunction(cmd.Args, final)
	switch word := firstWord.(type) {
	case *parse.BoolNode:
		return reflect.ValueOf(word.True)
	case *parse.DotNode:
		return dot
	case *parse.NilNode:
		s.errorf("nil is not a command")
	case *parse.NumberNode:
		return s.idealConstant(word)
	case *parse.StringNode:
		return reflect.ValueOf(word.Text)
	}
	s.errorf("can't evaluate command %q", firstWord)
	panic("not reached")
}

// IsTrue reports whether the value is 'true', in the sense of not the zero of its type,
// and whether the value has a meaningful truth value. This is the definition of
// truth used by if and other such actions.
func IsTrue(val interface{}) (truth, ok bool) {
	return isTrue(reflect.ValueOf(val))
}

func (s *state) setVar(n int, value reflect.Value) {
	s.vars[len(s.vars)-n].value = value
}

func (s *state) walkRange(dot reflect.Value, r *parse.RangeNode) {
	s.at(r)
	defer s.pop(s.mark())
	val, _ := indirect(s.evalPipeline(dot, r.Pipe))
	// mark top of stack before any variables in the body are pushed.
	mark := s.mark()
	oneIteration := func(index, elem reflect.Value) {
		// Set top var (lexically the second if there are two) to the element.
		if len(r.Pipe.Decl) > 0 {
			s.setVar(1, elem)
		}
		// Set next var (lexically the first if there are two) to the index.
		if len(r.Pipe.Decl) > 1 {
			s.setVar(2, index)
		}
		s.walk(elem, r.List)
		s.pop(mark)
	}
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() == 0 {
			break
		}
		for i := 0; i < val.Len(); i++ {
			oneIteration(reflect.ValueOf(i), val.Index(i))
		}
		return
	case reflect.Map:
		if val.Len() == 0 {
			break
		}
		for _, key := range sortKeys(val.MapKeys()) {
			oneIteration(key, val.MapIndex(key))
		}
		return
	case reflect.Chan:
		if val.IsNil() {
			break
		}
		i := 0
		for ; ; i++ {
			elem, ok := val.Recv()
			if !ok {
				break
			}
			oneIteration(reflect.ValueOf(i), elem)
		}
		if i == 0 {
			break
		}
		return
	case reflect.Invalid:
		break // An invalid value is likely a nil map, etc. and acts like an empty map.
	default:
		log.Errorf("range can't iterate over %v", val)
	}
	if r.ElseList != nil {
		s.walk(dot, r.ElseList)
	}
}

func (s *state) walkTemplate(dot reflect.Value, t *parse.TemplateNode) {
	s.at(t)
	tmpl := s.tmpl.tmpl[t.Name]
	if tmpl == nil {
		log.Errorf("template %q not defined", t.Name)
	}
	// Variables declared by the pipeline persist.
	dot = s.evalPipeline(dot, t.Pipe)
	newState := *s
	newState.depth++
	newState.tmpl = tmpl
	// No dynamic scoping: template invocations inherit no variables.
	newState.vars = []variable{{"$", dot}}
	newState.walk(dot, tmpl.Root)
}

// Eval functions evaluate pipelines, commands, and their elements and extract
// values from the data structure by examining fields, calling methods, and so on.
// The printing of those values happens only through walk functions.

func (s *state) notAFunction(args []parse.Node, final reflect.Value) {
	if len(args) > 1 || final.IsValid() {
		log.Errorf("can't give argument to non-function %s", args[0])
	}
}

// idealConstant is called to return the value of a number in a context where
// we don't know the type. In that case, the syntax of the number tells us
// its type, and we use Go rules to resolve. Note there is no such thing as
// a uint ideal constant in this situation - the value must be of int type.
func (s *state) idealConstant(constant *parse.NumberNode) reflect.Value {
	// These are ideal constants but we don't know the type
	// and we have no context.  (If it was a method argument,
	// we'd know what we need.) The syntax guides us to some extent.
	s.at(constant)
	switch {
	case constant.IsComplex:
		return reflect.ValueOf(constant.Complex128) // incontrovertible.
	case constant.IsFloat && !isHexConstant(constant.Text) && strings.ContainsAny(constant.Text, ".eE"):
		return reflect.ValueOf(constant.Float64)
	case constant.IsInt:
		n := int(constant.Int64)
		if int64(n) != constant.Int64 {
			log.Errorf("%s overflows int", constant.Text)
		}
		return reflect.ValueOf(n)
	case constant.IsUint:
		log.Errorf("%s overflows int", constant.Text)
	}
	return zero
}

func isHexConstant(s string) bool {
	return len(s) > 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

func (s *state) evalFieldNode(dot reflect.Value, field *parse.FieldNode, args []parse.Node, final reflect.Value) reflect.Value {
	s.at(field)
	return s.evalFieldChain(dot, dot, field, field.Ident, args, final)
}

func (s *state) evalChainNode(dot reflect.Value, chain *parse.ChainNode, args []parse.Node, final reflect.Value) reflect.Value {
	s.at(chain)
	if len(chain.Field) == 0 {
		s.errorf("internal error: no fields in evalChainNode")
	}
	if chain.Node.Type() == parse.NodeNil {
		s.errorf("indirection through explicit nil in %s", chain)
	}
	// (pipe).Field1.Field2 has pipe as .Node, fields as .Field. Eval the pipeline, then the fields.
	pipe := s.evalArg(dot, nil, chain.Node)
	return s.evalFieldChain(dot, pipe, chain, chain.Field, args, final)
}

func (s *state) varValue(name string) reflect.Value {
	for i := s.mark() - 1; i >= 0; i-- {
		if s.vars[i].name == name {
			return s.vars[i].value
		}
	}
	s.errorf("undefined variable: %s", name)
	return zero
}


func (s *state) evalVariableNode(dot reflect.Value, variable *parse.VariableNode, args []parse.Node, final reflect.Value) reflect.Value {
	// $x.Field has $x as the first ident, Field as the second. Eval the var, then the fields.
	s.at(variable)
	value := s.varValue(variable.Ident[0])
	if len(variable.Ident) == 1 {
		s.notAFunction(args, final)
		return value
	}
	return s.evalFieldChain(dot, value, variable, variable.Ident[1:], args, final)
}

// evalFieldChain evaluates .X.Y.Z possibly followed by arguments.
// dot is the environment in which to evaluate arguments, while
// receiver is the value being walked along the chain.
func (s *state) evalFieldChain(dot, receiver reflect.Value, node parse.Node, ident []string, args []parse.Node, final reflect.Value) reflect.Value {
	n := len(ident)
	for i := 0; i < n-1; i++ {
		receiver = s.evalField(dot, ident[i], node, nil, zero, receiver)
	}
	// Now if it's a method, it gets the arguments.
	return s.evalField(dot, ident[n-1], node, args, final, receiver)
}

// evalField evaluates an expression like (.Field) or (.Field arg1 arg2).
// The 'final' argument represents the return value from the preceding
// value of the pipeline, if any.
func (s *state) evalField(dot reflect.Value, fieldName string, node parse.Node, args []parse.Node, final, receiver reflect.Value) reflect.Value {
	if !receiver.IsValid() {
		return zero
	}
	typ := receiver.Type()
	receiver, isNil := indirect(receiver)
	// Unless it's an interface, need to get to a value of type *T to guarantee
	// we see all methods of T and *T.
	ptr := receiver
	if ptr.Kind() != reflect.Interface && ptr.Kind() != reflect.Ptr && ptr.CanAddr() {
		ptr = ptr.Addr()
	}
	if method := ptr.MethodByName(fieldName); method.IsValid() {
		return s.evalCall(dot, method, node, fieldName, args, final)
	}
	hasArgs := len(args) > 1 || final.IsValid()
	// It's not a method; must be a field of a struct or an element of a map.
	switch receiver.Kind() {
	case reflect.Struct:
		tField, ok := receiver.Type().FieldByName(fieldName)
		if ok {
			if isNil {
				s.errorf("nil pointer evaluating %s.%s", typ, fieldName)
			}
			field := receiver.FieldByIndex(tField.Index)
			if tField.PkgPath != "" { // field is unexported
				s.errorf("%s is an unexported field of struct type %s", fieldName, typ)
			}
			// If it's a function, we must call it.
			if hasArgs {
				s.errorf("%s has arguments but cannot be invoked as function", fieldName)
			}
			return field
		}
	case reflect.Map:
		if isNil {
			s.errorf("nil pointer evaluating %s.%s", typ, fieldName)
		}
		// If it's a map, attempt to use the field name as a key.
		nameVal := reflect.ValueOf(fieldName)
		if nameVal.Type().AssignableTo(receiver.Type().Key()) {
			if hasArgs {
				s.errorf("%s is not a method but has arguments", fieldName)
			}
			result := receiver.MapIndex(nameVal)
			return result
		}
	}
	s.errorf("can't evaluate field %s in type %s", fieldName, typ)
	panic("not reached")
}

var (
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType  = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
	reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

// evalCall executes a function or method call. If it's a method, fun already has the receiver bound, so
// it looks just like a function call. The arg list, if non-nil, includes (in the manner of the shell), arg[0]
// as the function itself.
func (s *state) evalCall(dot, fun reflect.Value, node parse.Node, name string, args []parse.Node, final reflect.Value) reflect.Value {
	if args != nil {
		args = args[1:] // Zeroth arg is function name/node; not passed to function.
	}
	typ := fun.Type()
	numIn := len(args)
	if final.IsValid() {
		numIn++
	}
	numFixed := len(args)
	if typ.IsVariadic() {
		numFixed = typ.NumIn() - 1 // last arg is the variadic one.
		if numIn < numFixed {
			s.errorf("wrong number of args for %s: want at least %d got %d", name, typ.NumIn()-1, len(args))
		}
	} else if numIn != typ.NumIn() {
		s.errorf("wrong number of args for %s: want %d got %d", name, typ.NumIn(), len(args))
	}
	// Build the arg list.
	argv := make([]reflect.Value, numIn)
	// Args must be evaluated. Fixed args first.
	i := 0
	for ; i < numFixed && i < len(args); i++ {
		argv[i] = s.evalArg(dot, typ.In(i), args[i])
	}
	// Now the ... args.
	if typ.IsVariadic() {
		argType := typ.In(typ.NumIn() - 1).Elem() // Argument is a slice.
		for ; i < len(args); i++ {
			argv[i] = s.evalArg(dot, argType, args[i])
		}
	}
	// Add final value if necessary.
	if final.IsValid() {
		t := typ.In(typ.NumIn() - 1)
		if typ.IsVariadic() {
			if numIn-1 < numFixed {
				// The added final argument corresponds to a fixed parameter of the function.
				// Validate against the type of the actual parameter.
				t = typ.In(numIn - 1)
			} else {
				// The added final argument corresponds to the variadic part.
				// Validate against the type of the elements of the variadic slice.
				t = t.Elem()
			}
		}
		argv[i] = s.validateType(final, t)
	}
	result := fun.Call(argv)
	// If we have an error that is not nil, stop execution and return that error to the caller.
	if len(result) == 2 && !result[1].IsNil() {
		s.at(node)
		s.errorf("error calling %s: %s", name, result[1].Interface().(error))
	}
	v := result[0]
	if v.Type() == reflectValueType {
		v = v.Interface().(reflect.Value)
	}
	return v
}

// canBeNil reports whether an untyped nil can be assigned to the type. See reflect.Zero.
func canBeNil(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return true
	case reflect.Struct:
		return typ == reflectValueType
	}
	return false
}

// validateType guarantees that the value is valid and assignable to the type.
func (s *state) validateType(value reflect.Value, typ reflect.Type) reflect.Value {
	if !value.IsValid() {
		if typ == nil || canBeNil(typ) {
			// An untyped nil interface{}. Accept as a proper nil value.
			return reflect.Zero(typ)
		}
		s.errorf("invalid value; expected %s", typ)
	}
	if typ == reflectValueType && value.Type() != typ {
		return reflect.ValueOf(value)
	}
	if typ != nil && !value.Type().AssignableTo(typ) {
		if value.Kind() == reflect.Interface && !value.IsNil() {
			value = value.Elem()
			if value.Type().AssignableTo(typ) {
				return value
			}
			// fallthrough
		}
		// Does one dereference or indirection work? We could do more, as we
		// do with method receivers, but that gets messy and method receivers
		// are much more constrained, so it makes more sense there than here.
		// Besides, one is almost always all you need.
		switch {
		case value.Kind() == reflect.Ptr && value.Type().Elem().AssignableTo(typ):
			value = value.Elem()
			if !value.IsValid() {
				s.errorf("dereference of nil pointer of type %s", typ)
			}
		case reflect.PtrTo(value.Type()).AssignableTo(typ) && value.CanAddr():
			value = value.Addr()
		default:
			s.errorf("wrong type for value; expected %s; got %s", typ, value.Type())
		}
	}
	return value
}

func (s *state) evalArg(dot reflect.Value, typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	switch arg := n.(type) {
	case *parse.DotNode:
		return s.validateType(dot, typ)
	case *parse.NilNode:
		if canBeNil(typ) {
			return reflect.Zero(typ)
		}
		s.errorf("cannot assign nil to %s", typ)
	case *parse.FieldNode:
		return s.validateType(s.evalFieldNode(dot, arg, []parse.Node{n}, zero), typ)
	case *parse.VariableNode:
		return s.validateType(s.evalVariableNode(dot, arg, nil, zero), typ)
	case *parse.PipeNode:
		return s.validateType(s.evalPipeline(dot, arg), typ)
	case *parse.ChainNode:
		return s.validateType(s.evalChainNode(dot, arg, nil, zero), typ)
	}
	switch typ.Kind() {
	case reflect.Bool:
		return s.evalBool(typ, n)
	case reflect.Complex64, reflect.Complex128:
		return s.evalComplex(typ, n)
	case reflect.Float32, reflect.Float64:
		return s.evalFloat(typ, n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return s.evalInteger(typ, n)
	case reflect.Interface:
		if typ.NumMethod() == 0 {
			return s.evalEmptyInterface(dot, n)
		}
	case reflect.Struct:
		if typ == reflectValueType {
			return reflect.ValueOf(s.evalEmptyInterface(dot, n))
		}
	case reflect.String:
		return s.evalString(typ, n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return s.evalUnsignedInteger(typ, n)
	}
	s.errorf("can't handle %s for arg of type %s", n, typ)
	panic("not reached")
}

func (s *state) evalBool(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.BoolNode); ok {
		value := reflect.New(typ).Elem()
		value.SetBool(n.True)
		return value
	}
	s.errorf("expected bool; found %s", n)
	panic("not reached")
}

func (s *state) evalString(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.StringNode); ok {
		value := reflect.New(typ).Elem()
		value.SetString(n.Text)
		return value
	}
	s.errorf("expected string; found %s", n)
	panic("not reached")
}

func (s *state) evalInteger(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsInt {
		value := reflect.New(typ).Elem()
		value.SetInt(n.Int64)
		return value
	}
	s.errorf("expected integer; found %s", n)
	panic("not reached")
}

func (s *state) evalUnsignedInteger(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsUint {
		value := reflect.New(typ).Elem()
		value.SetUint(n.Uint64)
		return value
	}
	s.errorf("expected unsigned integer; found %s", n)
	panic("not reached")
}

func (s *state) evalFloat(typ reflect.Type, n parse.Node) reflect.Value {
	s.at(n)
	if n, ok := n.(*parse.NumberNode); ok && n.IsFloat {
		value := reflect.New(typ).Elem()
		value.SetFloat(n.Float64)
		return value
	}
	s.errorf("expected float; found %s", n)
	panic("not reached")
}

func (s *state) evalComplex(typ reflect.Type, n parse.Node) reflect.Value {
	if n, ok := n.(*parse.NumberNode); ok && n.IsComplex {
		value := reflect.New(typ).Elem()
		value.SetComplex(n.Complex128)
		return value
	}
	s.errorf("expected complex; found %s", n)
	panic("not reached")
}

func (s *state) evalEmptyInterface(dot reflect.Value, n parse.Node) reflect.Value {
	s.at(n)
	switch n := n.(type) {
	case *parse.BoolNode:
		return reflect.ValueOf(n.True)
	case *parse.DotNode:
		return dot
	case *parse.FieldNode:
		return s.evalFieldNode(dot, n, nil, zero)
	case *parse.NilNode:
		// NilNode is handled in evalArg, the only place that calls here.
		s.errorf("evalEmptyInterface: nil (can't happen)")
	case *parse.NumberNode:
		return s.idealConstant(n)
	case *parse.StringNode:
		return reflect.ValueOf(n.Text)
	case *parse.VariableNode:
		return s.evalVariableNode(dot, n, nil, zero)
	case *parse.PipeNode:
		return s.evalPipeline(dot, n)
	}
	log.Errorf("can't handle assignment of %s to empty interface argument", n)
	panic("not reached")
}

// indirect returns the item at the end of indirection, and a bool to indicate if it's nil.
func indirect(v reflect.Value) (rv reflect.Value, isNil bool) {
	for ; v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface; v = v.Elem() {
		if v.IsNil() {
			return v, true
		}
	}
	return v, false
}

// printValue writes the textual representation of the value to the output of
// the template.
func (s *state) printValue(n parse.Node, v reflect.Value) {
	s.at(n)
	iface, ok := printableValue(v)
	if !ok {
		log.Errorf("can't print %s of type %s", n, v.Type())
	}
	_, err := fmt.Fprint(s.wr, iface)
	if err != nil {

	}
}

// printableValue returns the, possibly indirected, interface value inside v that
// is best for a call to formatted printer.
func printableValue(v reflect.Value) (interface{}, bool) {
	if v.Kind() == reflect.Ptr {
		v, _ = indirect(v) // fmt.Fprint handles nil.
	}
	if !v.IsValid() {
		return "<no value>", true
	}

	if !v.Type().Implements(errorType) && !v.Type().Implements(fmtStringerType) {
		if v.CanAddr() && (reflect.PtrTo(v.Type()).Implements(errorType) || reflect.PtrTo(v.Type()).Implements(fmtStringerType)) {
			v = v.Addr()
		} else {
			switch v.Kind() {
			case reflect.Chan, reflect.Func:
				return nil, false
			}
		}
	}
	return v.Interface(), true
}

// Types to help sort the keys in a map for reproducible output.

type rvs []reflect.Value

func (x rvs) Len() int      { return len(x) }
func (x rvs) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type rvInts struct{ rvs }

func (x rvInts) Less(i, j int) bool { return x.rvs[i].Int() < x.rvs[j].Int() }

type rvUints struct{ rvs }

func (x rvUints) Less(i, j int) bool { return x.rvs[i].Uint() < x.rvs[j].Uint() }

type rvFloats struct{ rvs }

func (x rvFloats) Less(i, j int) bool { return x.rvs[i].Float() < x.rvs[j].Float() }

type rvStrings struct{ rvs }

func (x rvStrings) Less(i, j int) bool { return x.rvs[i].String() < x.rvs[j].String() }

// sortKeys sorts (if it can) the slice of reflect.Values, which is a slice of map keys.
func sortKeys(v []reflect.Value) []reflect.Value {
	if len(v) <= 1 {
		return v
	}
	switch v[0].Kind() {
	case reflect.Float32, reflect.Float64:
		sort.Sort(rvFloats{v})
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sort.Sort(rvInts{v})
	case reflect.String:
		sort.Sort(rvStrings{v})
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		sort.Sort(rvUints{v})
	}
	return v
}

func (s *state) errorf(format string, args ...interface{}) {
	log.Errorf(s.tmpl.name)
}