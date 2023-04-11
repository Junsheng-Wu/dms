package rule

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	prommodel "github.com/prometheus/common/model"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

const (
	LabelKeyInternalRuleGroup    = "__rule_group__"
	LabelKeyInternalRuleName     = "__rule_name__"
	LabelKeyInternalRuleExpr     = "__rule_query__"
	LabelKeyInternalRuleDuration = "__rule_duration__"
	LabelKeyRuleId               = "alert_id"
)

var (
	ruleNameMatcher              = regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`)
	ErrAlertingRuleNotFound      = errors.New("The alerting rule was not found")
	ErrAlertingRuleAlreadyExists = errors.New("The alerting rule already exists")
	ErrAlertingGroupNotFound     = errors.New("The alerting group was not found")
)

func Validate(r Rule) error {
	var errs []error
	// only one of name and record must be set.
	if r.Alert == "" && r.Expr == "" {
		errs = append(errs, errors.New("one of name and expr must be set"))
	} else if r.Alert == "" {
		if !ruleNameMatcher.MatchString(r.Alert) {
			errs = append(errs, errors.New("rule name must match regular expression ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"))
		}
	}

	if r.For != "" {
		if _, err := prommodel.ParseDuration(r.For); err != nil {
			errs = append(errs, errors.Wrapf(err, "duration is invalid: %s", r.For))
		}
	}

	if len(r.Labels) > 0 {
		for name, v := range r.Labels {
			if !prommodel.LabelName(name).IsValid() || strings.HasPrefix(string(name), "__") {
				errs = append(errs, errors.Errorf(
					"label name (%s) is not valid. The name must match [a-zA-Z_][a-zA-Z0-9_]* and has not the __ prefix (label names with this prefix are for internal use)", name))
			}
			if !prommodel.LabelValue(v).IsValid() {
				errs = append(errs, errors.Errorf("invalid label value: %s", v))
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}
