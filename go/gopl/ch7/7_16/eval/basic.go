package eval

import (
	"alphas/go/gopl/utility"
	"fmt"
	"math"
	"strconv"
)

type Expr interface {
	Eval(env Env) float64
	fmt.Stringer
}

type Var string

func (v Var) Eval(env Env) float64 {
	return env.Get(v)
}

func (v Var) String() string {
	return fmt.Sprintf("env[%s]", string(v))
}

type literal float64

func (l literal) String() string {
	return fmt.Sprintf("literal[%f]", float64(l))
}

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

type unary struct {
	op rune
	x  Expr
}

func (u unary) String() string {
	return fmt.Sprintf("unary[%q%s]", u.op, u.x)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	default:
		panic(fmt.Errorf("unsupported unary operator: %q", u.op))
	}
}

type binary struct {
	op   rune
	x, y Expr
}

func (b binary) String() string {
	return fmt.Sprintf("binary[%q-%s,%s]", b.op, b.x, b.y)
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	default:
		panic(fmt.Errorf("unsupported binary operator: %q", b.op))
	}
}

type call struct {
	fn   string
	args []Expr
}

func (c call) String() string {
	return fmt.Sprintf("call[%q-%v]", c.fn, c.args)
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	default:
		panic(fmt.Errorf("unsupported call fn: %q", c.fn))
	}
}

type Env map[Var]float64

func (e Env) Get(key Var) float64 {
	val, ok := e[key]
	if ok {
		return val
	}
	val, err := strconv.ParseFloat(utility.RequireInput(string(key)+"="), 64)
	if err != nil {
		panic(err)
	}
	e[key] = val
	return val
}

type min struct {
	args []Expr
}

func (m min) Eval(env Env) float64 {
	if len(m.args) == 0 {
		panic("empty min func input")
	}
	ret := m.args[0].Eval(env)
	for i := 1; i < len(m.args); i++ {
		neo := m.args[i].Eval(env)
		if neo < ret {
			ret = neo
		}
	}
	return ret
}

func (m min) String() string {
	return fmt.Sprintf("min[%v]", m.args)
}
