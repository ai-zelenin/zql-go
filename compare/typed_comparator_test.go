package compare

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestCasterSuite(t *testing.T) {
	suite.Run(t, new(CasterSuite))
}

type CasterSuite struct {
	suite.Suite
	cmp *TypedComparator
}

func (s *CasterSuite) SetupSuite() {
	cmp := NewTypedComparator()
	s.cmp = cmp
}

func (s *CasterSuite) TestCalculateSimple() {
	cases := make([]ctCase, 0)
	cases = s.merge(
		s.makeNilTestcases(),
		s.makeNumberTestCases(),
		s.makeEqTestCases(),
		s.makeNegativeTestCases(),
		s.generateOpTestCases(0.1, 0.1, Eq, true),
		s.generateOpTestCases("a", "a", Eq, true),
		s.generateOpTestCases(true, true, Eq, true),
		s.generateOpTestCases(true, false, Gt, true),
		s.generateOpTestCases(
			time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
			time.Date(1, 1, 1, 1, 1, 1, 1, time.Local),
			Eq,
			true,
		),
	)
	for _, tc := range cases {
		fmt.Println(tc)
		out, err := s.cmp.Op(tc.Op, tc.V1, tc.V2)
		if err != nil {
			s.Require().ErrorIs(err, tc.err)
		} else {
			s.Require().Equal(tc.err, err)
		}
		s.Require().Equal(tc.expected, out)
	}
	fmt.Println(len(cases))
}

func (s *CasterSuite) makeNilTestcases() []ctCase {
	nils := s.allKindsOfNils()
	var testCases = make([]ctCase, 0)
	for _, v1 := range nils {
		testCases = append(testCases, s.generateOpTestCases(v1, v1, Eq, false)...)
		testCases = append(testCases, s.generateOpTestCases(v1, nil, Eq, false)...)
		testCases = append(testCases,
			newCtCase(v1, v1, Gt, false, ErrCompareOperationIsNotAcceptable),
			newCtCase(v1, v1, Gte, true, ErrCompareOperationIsNotAcceptable),
			newCtCase(v1, v1, Lt, false, ErrCompareOperationIsNotAcceptable),
			newCtCase(v1, v1, Lte, true, ErrCompareOperationIsNotAcceptable),
		)
	}
	return testCases
}

func (s *CasterSuite) makeNumberTestCases() []ctCase {
	ones := s.allKindsOfNumbers(1)
	zeros := s.allKindsOfNumbers(0)
	var testCases = make([]ctCase, 0)
	for _, v1 := range ones {
		for _, v2 := range ones {
			testCases = append(testCases, s.generateOpTestCases(v1, v2, Eq, true)...)
		}
	}
	for _, one := range ones {
		for _, zero := range zeros {
			testCases = append(testCases, s.generateOpTestCases(one, zero, Gt, true)...)
		}
	}
	return testCases
}

func (s *CasterSuite) makeNegativeTestCases() []ctCase {
	gen := map[reflect.Kind]func() []any{
		reflect.Complex64: func() []any {
			return []any{
				complex(float32(1), float32(1)),
				complex(float32(2), float32(2)),
			}
		},
		reflect.Complex128: func() []any {
			return []any{
				complex(float64(1), float64(1)),
				complex(float64(2), float64(2)),
			}
		},
		reflect.Array: func() []any {
			return []any{
				[1]int{1},
				[1]string{"1"},
			}
		},
		reflect.UnsafePointer: func() []any {
			var a = bytes.NewBuffer(nil)
			return []any{
				reflect.ValueOf(a).UnsafePointer(),
			}
		},
	}
	s.check(gen)
	var testCases = make([]ctCase, 0)

	ops := []string{
		Eq,
		Neq,
		Gt,
		Gte,
		Lt,
		Lte,
	}
	for _, constructor := range gen {
		vals1 := constructor()
		vals2 := constructor()
		for _, v1 := range vals1 {
			for _, v2 := range vals2 {
				for _, op := range ops {
					testCases = append(testCases, []ctCase{
						{
							V1:       v1,
							V2:       v2,
							Op:       op,
							expected: false,
							err:      ErrTypesIsNotComparable,
						},
					}...)
				}
			}
		}
	}
	return testCases
}

func (s *CasterSuite) makeEqTestCases() []ctCase {
	kinds := map[reflect.Kind]func() []any{
		reflect.Invalid: func() []any {
			return []any{
				nil,
				fmt.Stringer(nil),
			}
		},
		reflect.Bool: func() []any {
			return []any{
				false,
				true,
			}
		},
		reflect.Pointer: func() []any {
			var nilPtr *int = nil
			return []any{
				nilPtr,
			}
		},
		reflect.String: func() []any {
			return []any{
				"a",
				"b",
			}
		},
		reflect.Struct: func() []any {
			return []any{
				time.Date(20, 1, 1, 1, 1, 1, 1, time.Local),
			}
		},
	}

	var testCases = make([]ctCase, 0)
	for k1, constructor := range kinds {
		vals1 := constructor()
		for _, val := range vals1 {
			s.Require().Equal(k1.String(), reflect.ValueOf(val).Kind().String())
		}
	}

	for _, constructor := range kinds {
		vals1 := constructor()
		vals2 := constructor()
		for i := range vals1 {
			testCases = append(testCases, s.generateOpTestCases(vals1[i], vals2[i], Eq, false)...)
		}
	}
	return testCases
}

func (s *CasterSuite) generateOpTestCases(v1, v2 any, op string, full bool) []ctCase {
	var testCases = make([]ctCase, 0)

	switch op {
	case Eq:
		testCases = append(testCases,
			newCtCase(v1, v2, Eq, true, nil),
			newCtCase(v1, v2, Neq, false, nil),
		)
		if full {
			testCases = append(testCases,
				newCtCase(v1, v2, Lt, false, nil),
				newCtCase(v1, v2, Lte, true, nil),
				newCtCase(v1, v2, Gt, false, nil),
				newCtCase(v1, v2, Gte, true, nil),
			)
		}
	case Neq:
		testCases = append(testCases,
			newCtCase(v1, v2, Eq, false, nil),
			newCtCase(v1, v2, Neq, true, nil),
		)
		if full {
			testCases = append(testCases,
				newCtCase(v1, v2, Lte, false, nil),
				newCtCase(v1, v2, Gte, false, nil),
			)
		}
	case Lt:
		testCases = append(testCases,
			newCtCase(v1, v2, Lt, true, nil),
			newCtCase(v1, v2, Lte, true, nil),
			newCtCase(v1, v2, Neq, true, nil),
			newCtCase(v1, v2, Eq, false, nil),
			newCtCase(v1, v2, Gt, false, nil),
			newCtCase(v1, v2, Gte, false, nil),
		)
	case Gt:
		testCases = append(testCases,
			newCtCase(v1, v2, Gt, true, nil),
			newCtCase(v1, v2, Gte, true, nil),
			newCtCase(v1, v2, Neq, true, nil),
			newCtCase(v1, v2, Eq, false, nil),
			newCtCase(v1, v2, Lt, false, nil),
			newCtCase(v1, v2, Lte, false, nil),
		)
	}
	return testCases
}

func (s *CasterSuite) merge(caseList ...[]ctCase) []ctCase {
	var result = make([]ctCase, 0)
	for _, cases := range caseList {
		result = append(result, cases...)
	}
	return result
}

func (s *CasterSuite) allKindsOfNils() []any {
	var pointerNil *int = nil
	var sliceNil []int = nil
	var mapNil map[string]string = nil
	var chanNil chan int = nil
	type testInterface interface {
		A()
	}
	var interfaceNil testInterface = nil
	var funcNil func() = nil
	return []any{
		pointerNil,
		sliceNil,
		mapNil,
		chanNil,
		interfaceNil,
		funcNil,
	}
}

func (s *CasterSuite) allKindsOfNumbers(val uint8) []any {
	var intValue int = int(val)
	var int8Value int8 = int8(val)
	var int16Value int16 = int16(val)
	var int32Value int32 = int32(val)
	var int64Value int64 = int64(val)
	var uintValue uint = uint(val)
	var uint16Value uint16 = uint16(val)
	var uint32Value uint32 = uint32(val)
	var uint64Value uint64 = uint64(val)
	var uintptrValue uintptr = uintptr(val)
	var float32Value float32 = float32(val)
	var float64Value float64 = float64(val)
	var byteValue byte = byte(val)
	return []any{
		intValue,
		int8Value,
		int16Value,
		int32Value,
		int64Value,
		uintValue,
		uint16Value,
		uint32Value,
		uint64Value,
		uintptrValue,
		float32Value,
		float64Value,
		byteValue,
	}
}

func (s *CasterSuite) check(gen map[reflect.Kind]func() []any) {
	for k1, constructor := range gen {
		vals1 := constructor()
		for _, val := range vals1 {
			s.Require().Equal(k1.String(), reflect.ValueOf(val).Kind().String())
		}
	}
}

type ctCase struct {
	V1       any
	V2       any
	Op       string
	expected bool
	err      error
}

func newCtCase(v1 any, v2 any, op string, expected bool, err error) ctCase {
	return ctCase{V1: v1, V2: v2, Op: op, expected: expected, err: err}
}

func (tc ctCase) String() string {
	rv1 := reflect.ValueOf(tc.V1)
	rv2 := reflect.ValueOf(tc.V2)
	return fmt.Sprintf(`%s( %v(%T:%s) , %v(%T:%s) ) = expected:%t err:%v`,
		tc.Op,
		tc.V1, tc.V1, rv1.Kind(),
		tc.V2, tc.V2, rv2.Kind(),
		tc.expected, tc.err)
}
