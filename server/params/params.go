package params

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful"
)

const (
	PagingParam     = "paging"
	OrderByParam    = "orderBy"
	ConditionsParam = "conditions"
	ReverseParam    = "reverse"
)

func ParsePaging(req *restful.Request) (limit, offset int) {
	paging := req.QueryParameter(PagingParam)
	limit = 10
	offset = 0
	if groups := regexp.MustCompile(`^limit=(-?\d+),page=(\d+)$`).FindStringSubmatch(paging); len(groups) == 3 {
		limit, _ = strconv.Atoi(groups[1])
		page, _ := strconv.Atoi(groups[2])
		offset = (page - 1) * limit
	}
	return
}

func ParseConditions(req *restful.Request) (*Conditions, error) {
	conditionsStr := req.QueryParameter(ConditionsParam)

	conditions := &Conditions{Match: make(map[string]string), Fuzzy: make(map[string]string)}

	if conditionsStr == "" {
		return conditions, nil
	}

	// ?conditions=key1=value1,key2~value2,key3=
	for _, item := range strings.Split(conditionsStr, ",") {
		// exact query: key=value, if value is empty means label value must be ""
		// fuzzy query: key~value, if value is empty means label value is "" or label key not exist
		if groups := regexp.MustCompile(`(\S+)([=~])(\S+)?`).FindStringSubmatch(item); len(groups) >= 3 {
			value := ""

			if len(groups) > 3 {
				value = groups[3]
			}

			if groups[2] == "=" {
				conditions.Match[groups[1]] = value
			} else {
				conditions.Fuzzy[groups[1]] = value
			}
		} else {
			return nil, fmt.Errorf("invalid conditions")
		}
	}
	return conditions, nil
}

type Conditions struct {
	Match map[string]string
	Fuzzy map[string]string
}

func GetBoolValueWithDefault(req *restful.Request, name string, dv bool) bool {
	reverse := req.QueryParameter(name)
	if v, err := strconv.ParseBool(reverse); err == nil {
		return v
	}
	return dv
}

func GetStringValueWithDefault(req *restful.Request, name string, dv string) string {
	v := req.QueryParameter(name)
	if v == "" {
		v = dv
	}
	return v
}
