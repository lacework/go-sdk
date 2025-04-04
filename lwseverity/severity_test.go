package lwseverity_test

import (
	"testing"

	"github.com/lacework/go-sdk/v2/lwseverity"
	"github.com/stretchr/testify/assert"
)

type normalizeTest struct {
	Input       string
	ExpectedInt int
	ExpectedStr string
}

var normalizeTests = []normalizeTest{
	{Input: "critical", ExpectedInt: 1, ExpectedStr: "Critical"},
	{Input: "high", ExpectedInt: 2, ExpectedStr: "High"},
	{Input: "medium", ExpectedInt: 3, ExpectedStr: "Medium"},
	{Input: "low", ExpectedInt: 4, ExpectedStr: "Low"},
	{Input: "info", ExpectedInt: 5, ExpectedStr: "Info"},
	{Input: "unknown", ExpectedInt: 0, ExpectedStr: "Unknown"},
	{Input: "foo", ExpectedInt: 0, ExpectedStr: "Unknown"},
}

func TestNormalize(t *testing.T) {
	for _, nt := range normalizeTests {
		t.Run(nt.Input, func(t *testing.T) {
			actualInt, actualStr := lwseverity.Normalize(nt.Input)
			assert.Equal(t, nt.ExpectedInt, actualInt)
			assert.Equal(t, nt.ExpectedStr, actualStr)
		})
	}
}

type notAsCriticalTest struct {
	Name      string
	Severity  string
	Threshold string
	Expected  bool
}

var notAsCriticalTests = []notAsCriticalTest{
	{Name: "less", Severity: "medium", Threshold: "high", Expected: true},
	{Name: "equal", Severity: "high", Threshold: "high", Expected: false},
	{Name: "more", Severity: "critical", Threshold: "high", Expected: false},
	{Name: "unknown-severity", Severity: "fwaasdf", Threshold: "high", Expected: false},
	{Name: "unknown-threshold", Severity: "critical", Threshold: "fwaasdf", Expected: true},
}

func TestNotAsCritical(t *testing.T) {
	for _, nact := range notAsCriticalTests {
		t.Run(nact.Name, func(t *testing.T) {
			assert.Equal(
				t,
				nact.Expected,
				lwseverity.NotAsCritical(nact.Severity, nact.Threshold),
			)
		})
	}
}

type shouldFilterTest struct {
	Name      string
	Severity  string
	Threshold string
	Expected  bool
}

var shouldFilterTests = []shouldFilterTest{
	{Name: "less", Severity: "medium", Threshold: "high", Expected: true},
	{Name: "equal", Severity: "high", Threshold: "high", Expected: false},
	{Name: "greater", Severity: "critical", Threshold: "high", Expected: false},
	{Name: "unknown-severity", Severity: "fwaasdf", Threshold: "high", Expected: false},
	{Name: "unknown-threshold", Severity: "critical", Threshold: "fwaasdf", Expected: false},
}

func TestShouldFilterTest(t *testing.T) {
	for _, sft := range shouldFilterTests {
		t.Run(sft.Name, func(t *testing.T) {
			assert.Equal(
				t,
				sft.Expected,
				lwseverity.ShouldFilter(
					sft.Severity, sft.Threshold),
			)
		})
	}
}

type myStruct struct {
	alertID  int
	severity string
}

func (m myStruct) GetSeverity() string {
	return m.severity
}

type myStructs []myStruct

func TestSort(t *testing.T) {
	m := myStructs{
		myStruct{alertID: 6, severity: "Low"},
		myStruct{alertID: 9, severity: "High"},
		myStruct{alertID: 1, severity: "Low"},
		myStruct{alertID: 3, severity: "Low"},
		myStruct{alertID: 8, severity: "Low"},
	}
	expected := myStructs{
		myStruct{alertID: 9, severity: "High"},
		myStruct{alertID: 6, severity: "Low"},
		myStruct{alertID: 1, severity: "Low"},
		myStruct{alertID: 3, severity: "Low"},
		myStruct{alertID: 8, severity: "Low"},
	}
	lwseverity.SortSlice(m)
	assert.Equal(t, expected, m)

	expected = myStructs{
		myStruct{alertID: 6, severity: "Low"},
		myStruct{alertID: 1, severity: "Low"},
		myStruct{alertID: 3, severity: "Low"},
		myStruct{alertID: 8, severity: "Low"},
		myStruct{alertID: 9, severity: "High"},
	}
	lwseverity.SortSliceA(m)
	assert.Equal(t, expected, m)
}

func TestIsValid(t *testing.T) {
	assert.Equal(t, true, lwseverity.IsValid("Critical"))
	assert.Equal(t, false, lwseverity.IsValid("JackBauer"))
}

func TestValidSeveritiesString(t *testing.T) {
	assert.Equal(t, "critical, high, medium, low, info", lwseverity.ValidSeverities.String())
}
