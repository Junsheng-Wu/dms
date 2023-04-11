package rule

import (
	"encoding/json"
	"github.com/pkg/errors"
	prommodel "github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	HttpPATH = "http://"
	ReqPATH  = "/apis/alerting/v1/rules/update"
)

type Operator interface {
	ListRules() (map[string][]Rule, error)
	AddRule(group string, rule Rule) error
	UpdateRule(group, id string, rule Rule) error
	DeleteRule(group string, ruleId string) error
	UpdateAllRules(groups AlertingRules) error
}

type operator struct {
	RulePath           string
	Urls               []string // peer
	PrometheusEndpoint string
}

func NewOperator(path, prometheusEndpoint, peers string) Operator {
	var urls []string
	ps := strings.Split(peers, ",")
	for _, p := range ps {
		url := HttpPATH + p + ReqPATH
		urls = append(urls, url)
	}
	return &operator{
		RulePath:           path,
		Urls:               urls,
		PrometheusEndpoint: prometheusEndpoint,
	}
}

func (o *operator) ListRules() (map[string][]Rule, error) {
	var groups AlertingRules

	rulesData, err := ioutil.ReadFile(o.RulePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(rulesData, &groups)
	if err != nil {
		return nil, err
	}

	return groups.ListRules(), nil
}

func (o *operator) AddRule(group string, rule Rule) error {
	var (
		groups AlertingRules
	)

	err := Validate(rule)
	if err != nil {
		return err
	}
	setRuleId(rule)

	rulesData, err := ioutil.ReadFile(o.RulePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(rulesData, &groups)
	if err != nil {
		return err
	}

	ruleExit := groups.GetRuleByName(rule.Alert)
	if ruleExit != nil {
		return ErrAlertingRuleAlreadyExists
	}

	mapGroupRules := groups.MapGroupRules()

	if _, ok := mapGroupRules[group]; ok != true {
		mapGroupRules[group] = []Rule{}
	}

	for g, _ := range mapGroupRules {
		if g == group {
			mapGroupRules[group] = append(mapGroupRules[group], rule)
			break
		}
	}

	gs := mapToRules(mapGroupRules)
	data, err := yaml.Marshal(gs)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(o.RulePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(o.RulePath, data, 0644)
	if err != nil {
		return err
	}

	err = o.ReloadPrometheusConfig()
	if err != nil {
		return err
	}

	if len(o.Urls) > 0 {
		err = o.SyncAlertingRules(gs, o.Urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *operator) UpdateRule(group, ruleId string, rule Rule) error {
	var (
		groups AlertingRules
	)

	err := Validate(rule)
	if err != nil {
		return err
	}

	rulesData, err := ioutil.ReadFile(o.RulePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(rulesData, &groups)
	if err != nil {
		return err
	}

	ruleExit := groups.GetRuleById(ruleId)
	if ruleExit == nil {
		return ErrAlertingRuleNotFound
	}

	rule.Labels["alert_id"] = ruleId

	mapGroupRules := groups.MapGroupRules()
Loop:
	for g, rules := range mapGroupRules {
		if g == group {
			for i, r := range rules {
				if r.Labels["alert_id"] == ruleId {
					rules := append(rules[:i], rules[i+1:]...)
					mapGroupRules[group] = rules
					mapGroupRules[group] = append(mapGroupRules[group], rule)
					break Loop
				}
			}
		}
	}

	gs := mapToRules(mapGroupRules)
	data, err := yaml.Marshal(gs)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(o.RulePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(o.RulePath, data, 0644)
	if err != nil {
		return err
	}

	err = o.ReloadPrometheusConfig()
	if err != nil {
		return err
	}

	if len(o.Urls) > 0 {
		err = o.SyncAlertingRules(gs, o.Urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *operator) DeleteRule(group string, ruleId string) error {
	var (
		groups AlertingRules
	)

	rulesData, err := ioutil.ReadFile(o.RulePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(rulesData, &groups)
	if err != nil {
		return err
	}

	ruleExit := groups.GetRuleById(ruleId)
	if ruleExit == nil {
		return ErrAlertingRuleNotFound
	}

	mapGroupRules := groups.MapGroupRules()
Loop:
	for g, rules := range mapGroupRules {
		if g == group {
			for i, rule := range rules {
				if rule.Labels["alert_id"] == ruleId {
					rules := append(rules[:i], rules[i+1:]...)
					mapGroupRules[group] = rules
					break Loop
				}
			}
		}
	}

	gs := mapToRules(mapGroupRules)
	data, err := yaml.Marshal(gs)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(o.RulePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(o.RulePath, data, 0644)
	if err != nil {
		return err
	}

	err = o.ReloadPrometheusConfig()
	if err != nil {
		return err
	}

	if len(o.Urls) > 0 {
		err = o.SyncAlertingRules(gs, o.Urls)
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *operator) SyncAlertingRules(groups AlertingRules, urls []string) error {
	data, err := json.Marshal(groups)
	if err != nil {
		return err
	}

	for _, url := range urls {
		_, err := http.Post(url, "text/html", strings.NewReader(string(data)))
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *operator) ReloadPrometheusConfig() error {
	url := "http://" + o.PrometheusEndpoint + "/-/reload"

	_, err := http.Post(url, "text/html", nil)
	if err != nil {
		return err
	}

	return nil
}

func (o *operator) UpdateAllRules(groups AlertingRules) error {
	if reflect.DeepEqual(groups, AlertingRules{}) {
		return ErrAlertingGroupNotFound
	}

	data, err := yaml.Marshal(groups)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(o.RulePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(o.RulePath, data, 0644)
	if err != nil {
		return err
	}

	err = o.ReloadPrometheusConfig()
	if err != nil {
		return err
	}

	return nil
}

func setRuleId(rule Rule) {
	if rule.Labels == nil {
		rule.Labels = make(map[string]string)
	}
	id := GenResourceRuleIdIgnoreFormat(rule)
	rule.Labels[LabelKeyRuleId] = id
}

func GenResourceRuleIdIgnoreFormat(rule Rule) string {
	lbls := make(map[string]string)
	for k, v := range rule.Labels {
		if k == LabelKeyRuleId {
			continue
		}
		lbls[k] = v
	}

	lbls[LabelKeyInternalRuleName] = rule.Alert
	lbls[LabelKeyInternalRuleDuration] = rule.For
	lbls[LabelKeyInternalRuleExpr] = rule.Expr

	return prommodel.Fingerprint(prommodel.LabelsToSignature(lbls)).String()
}

func GetRuleId(rule Rule) (string, error) {
	id, ok := rule.Labels["alert_id"]
	if !ok {
		return "", errors.New("rule id not found")
	}
	return id, nil
}

func mapToRules(mg map[string][]Rule) AlertingRules {
	var groups AlertingRules
	for g, rs := range mg {
		groups.Groups = append(groups.Groups, Group{
			Name:  g,
			Rules: rs,
		})
	}
	return groups
}
