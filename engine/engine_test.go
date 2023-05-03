package engine

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ai-zelenin/zql-go"
)

func TestEngineSuite(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}

type EngineSuite struct {
	suite.Suite
	engine *ExtendableEngine
}
type testCase struct {
	predicate Predicate
	scope     Scope
	out       bool
	err       error
}

func (s *EngineSuite) SetupSuite() {
	s.engine = NewExtendableEngine(&Config{
		PreHooks: []HookFunc{
			NewVarPrefixHook("$"),
		},
	})
}

func (s *EngineSuite) calc(predicate Predicate, scope Scope) (bool, error) {
	ep, err := s.engine.Compile(predicate)
	if err != nil {
		return false, err
	}
	out, err := ep.Exec(scope)
	fmt.Println(ep.Dump(""))
	return out, err
}

func (s *EngineSuite) TestCalculateComplex() {
	cases := []testCase{
		{
			predicate: zql.And(
				zql.Eq("test1", 1001),
				zql.Eq("test2", 1002),
			),
			scope: NewMapScope(map[string]any{
				"test1": 1001,
				"test2": 1002,
			}),
			out: true,
		},
		{
			predicate: zql.Or(
				zql.Eq("test1", "$test2"), // 1001 == 1002
				zql.Eq("test2", "$test3"), // 1002 == 1003
				zql.Lt("test1", "$test3"), // 1001 < 1003
			),
			scope: NewMapScope(map[string]any{
				"test1": 1001,
				"test2": 1002,
				"test3": 1003,
			}),
			out: true,
		},
		{
			predicate: zql.And(
				zql.Or( // true
					zql.Or( // true
						zql.Eq("test1", 1001), // true
						zql.Eq("test1", 1),    // skip
						zql.Eq("test1", 1),    // skip
					),
					zql.Or( // skip
						zql.Gt("test2", 1001), // skip
						zql.Lt("test2", 1),    // skip
						zql.Lte("test2", 1),   // skip
					),
				),
				zql.And( // false
					zql.And( // false
						zql.Gt("test3", 1001), // true
						zql.Lt("test3", 1),    // false
						zql.Lte("test3", 1),   // skip
					),
					zql.And( //  skip
						zql.Or( // skip
							zql.Eq("test3", 1001), // skip
							zql.Eq("test3", 1),    // skip
							zql.Eq("test3", 1),    // skip
						),
					),
				),
			),
			scope: NewMapScope(map[string]any{
				"test1": 1001,
				"test2": 1002,
				"test3": 1003,
			}),
			out: false,
		},
	}
	for i, tc := range cases {
		name := fmt.Sprintf("%d", i)
		s.Run(name, func() {
			out, err := s.calc(tc.predicate, tc.scope)
			s.Equal(tc.out, out)
			if err != nil && tc.err != nil {
				s.ErrorIs(err, tc.err)
			} else {
				s.Equal(tc.err, err)
			}
		})
	}
}
