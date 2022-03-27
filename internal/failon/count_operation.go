package failon

import (
	"fmt"
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

	var op_parts []string
	if op_parts = re.FindStringSubmatch(s); s == "" || op_parts == nil {
		return errors.New(
			fmt.Sprintf("count operation (%s) is invalid", s))
	}
	co.num, _ = strconv.Atoi(op_parts[2])
	co.operator = op_parts[1]
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
	}
	return true, errors.New(fmt.Sprintf("count operation (%s) is invalid", co.operator))
}
