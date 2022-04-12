package failon

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const operationRE = `^(>|>=|<|<=|={1,2}|!=)\s*(\d+)$`

type CountOperation struct {
	operator string
	num      int
}

func (co *CountOperation) Parse(s string) error {
	re := regexp.MustCompile(operationRE)

	s = strings.TrimSpace(s)

	var opParts []string
	if opParts = re.FindStringSubmatch(s); s == "" || opParts == nil {
		return errors.Errorf("count operation (%s) is invalid", s)
	}
	co.num, _ = strconv.Atoi(opParts[2])
	co.operator = opParts[1]
	return nil
}

func (co CountOperation) IsFail(count int) (bool, error) {
	switch co.operator {
	case ">":
		return count > co.num, nil
	case ">=":
		return count >= co.num, nil
	case "<":
		return count < co.num, nil
	case "<=":
		return count <= co.num, nil
	case "=", "==":
		return count == co.num, nil
	case "!=":
		return count != co.num, nil
	default:
		return true, errors.Errorf("count operation (%s) is invalid", co.operator)
	}
}
