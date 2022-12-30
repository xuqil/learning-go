package orm

// 衍生类型
type op string

// 别名
//type op = string

const (
	opEq  op = "="
	opLT  op = "<"
	opNot op = "NOT"
	opAnd op = "AND"
	opOr  op = "OR"
)

func (o op) String() string {
	return string(o)
}

type Predicate struct {
	left  Expression
	op    op
	right Expression
}

// EQ EQ("id", 12)
// EQ(sub, "id, 12)
// EQ(sub.id, 12)
//func EQ(column string, right any) Predicate {
//	return Predicate{
//		Column: column,
//		Op:     "=",
//		Arg:    right,
//	}
//}

type Column struct {
	name string
}

func C(name string) Column {
	return Column{name: name}
}

// Eq 例如 C("id").EQ(12)
// sub.C("id").EQ(12)
func (c Column) Eq(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opEq,
		right: value{val: arg},
	}
}

func (c Column) LT(arg any) Predicate {
	return Predicate{
		left:  c,
		op:    opLT,
		right: value{val: arg},
	}
}

func (c Column) expr() {}

// Not  Not(C("name").EQ("Tom"))
func Not(p Predicate) Predicate {
	return Predicate{
		op:    opNot,
		right: p,
	}
}

// And  C("id").EQ(12).And(C("name").EQ("Tom"))
func (left Predicate) And(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opAnd,
		right: right,
	}
}

// Or  C("id").EQ(12).Or(C("name").EQ("Tom"))
func (left Predicate) Or(right Predicate) Predicate {
	return Predicate{
		left:  left,
		op:    opOr,
		right: right,
	}
}

func (Predicate) expr() {}

type value struct {
	val any
}

func (value) expr() {}

// Expression 是一个标记接口，代表表达式
type Expression interface {
	expr()
}
