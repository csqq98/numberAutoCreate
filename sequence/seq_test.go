package sequence

import (
	"log"
	"testing"
)

type SeqDb struct {
}

var group = map[string]string{
	"20210818": "000001",
	"20210817": "000002",
}

func (s *SeqDb) Data() []*SeqInfo {
	a := []*SeqInfo{
		{
			Module: "1",
			Expr:   "p:1-p:2-d:yyMMdd-i:2:5",
			Group:  group,
			Remark: "模块1",
		},
		{
			Module: "2",
			Expr:   "p:1-d:yyyyMMdd-i:d:5",
			Group:  group,
			Remark: "模块2",
		},
	}
	return a
}
func (s *SeqDb) Save(model *SeqInfo) {
	log.Printf("modelInfo %+v", model.Group)
}

func TestGen(t *testing.T) {
	dbs := SeqDb{}
	s := New(&dbs)
	genTestCreateNo(t, s, &calcCase{
		Model:    "1",
		Str:      s.Modules["1"].Expr,
		Expected: "GTRPSP21081900001",
	}, "GTR", "PSP")
	genTestCreateNo(t, s, &calcCase{
		Model:    "1",
		Str:      s.Modules["1"].Expr,
		Expected: "GTRPSP21081900002",
	}, "GTR", "PSP")
	genTestCreateNo(t, s, &calcCase{
		Model:    "2",
		Str:      s.Modules["2"].Expr,
		Expected: "GTR2021081900001",
	}, "GTR")
	genTestCreateNo(t, s, &calcCase{
		Model:    "1",
		Str:      s.Modules["1"].Expr,
		Expected: "GTRPSP21081900003",
	}, "GTR", "PSP")
}

type calcCase struct {
	Model    string // 模块
	Str      string // 表达式
	Expected string // 测试用例结果
}

func genTestCreateNo(t *testing.T, seq *Sequence, c *calcCase, args ...string) {
	if ans := seq.Gen(c.Model, args...); ans != c.Expected {
		t.Fatalf("Module %s's expression is %s , expected %s, but %s got",
			c.Model, c.Str, c.Expected, ans)
	}
}

func TestIfAccordNumberRule(t *testing.T) {
	cases := []struct {
		Str      string
		Expected bool
	}{
		{"p:1-p:2-d:yyMMdd-i:2:1", true},
		{"p:1-p:2-d:yyMMdd-i:3:2", false},
		{"p:1-p:2-d:yyMMdd-i:4:3", false},
		{"p:1-p:2-d:yyMMdd-i:d:0", false},
		{"p:1-p:2-d:yyMMdd-i:d:4", true},
		{"p:1-d:yyMMdd-i:2:5", false},
		{"p:1-d:yyMMdd-i:1:5", true},
		{"p:1-d:yyMMdd-i:d:5", true},
		{"p:0-d:yyMMdd-i:d:5", false}}
	for _, c := range cases {
		t.Run(c.Str, func(t *testing.T) {
			if ans := IfAccordNumberRule(c.Str); ans != c.Expected {
				t.Fatalf("%s expected %+v, but %+v got",
					c.Str, c.Expected, ans)
			}
		})
	}
}
