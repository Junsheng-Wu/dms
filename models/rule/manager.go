package rule

type AlertingRules struct {
	Groups []Group `json:"groups,omitempty"`
}

type Group struct {
	Name  string `json:"name,omitempty"`
	Rules []Rule `json:"rules,omitempty"`
}

type Rule struct {
	Alert  string            `json:"alert,omitempty"`
	Expr   string            `json:"expr,omitempty"`
	For    string            `json:"for,omitempty"`
	Labels map[string]string `json:"labels,omitempty"`
}

func (a AlertingRules) ListRules() map[string][]Rule {
	var groupRules map[string][]Rule
	groupRules = make(map[string][]Rule)
	if len(a.Groups) == 0 {
		return map[string][]Rule{}
	}

	for _, group := range a.Groups {
		if len(group.Rules) == 0 {
			continue
		}
		if _, ok := groupRules[group.Name]; ok != true {
			groupRules[group.Name] = []Rule{}
		}
		for _, r := range group.Rules {
			groupRules[group.Name] = append(groupRules[group.Name], r)
		}
	}
	return groupRules
}

func (a AlertingRules) GetRuleByName(name string) *Rule {
	if a.Groups == nil {
		return nil
	}
	for _, group := range a.Groups {
		for _, r := range group.Rules {
			if r.Alert == name {
				return &r
			}

		}
	}

	return nil
}

func (a AlertingRules) GetRuleById(Id string) *Rule {
	if a.Groups == nil {
		return nil
	}
	for _, group := range a.Groups {
		for _, r := range group.Rules {
			if r.Labels["alert_id"] == Id {
				return &r
			}
		}
	}

	return nil
}

func (a AlertingRules) MapGroupRules() map[string][]Rule {
	var mapGroupRules map[string][]Rule
	mapGroupRules = make(map[string][]Rule)
	for _, group := range a.Groups {
		mapGroupRules[group.Name] = group.Rules
	}

	return mapGroupRules
}

//func (a AlertingRules) DeleteRule() map[string][]Rule {
//	var mapGroupRules map[string][]Rule
//	mapGroupRules = make(map[string][]Rule)
//	for _, group := range a.Groups {
//		mapGroupRules[group.Name] = group.Rules
//	}
//
//	return mapGroupRules
//}
