package feed

import (
	"regexp"

	"github.com/PhilippHeuer/rssdownloader/pkg/config"
)

type Filter struct {
	Rules    []*regexp.Regexp
	Excludes []*regexp.Regexp
}

func NewFilter(rules []config.Rule, excludes []config.Rule) (*Filter, error) {
	compiledRules, err := compileRules(rules)
	if err != nil {
		return nil, err
	}

	compiledExcludes, err := compileRules(excludes)
	if err != nil {
		return nil, err
	}

	return &Filter{
		Rules:    compiledRules,
		Excludes: compiledExcludes,
	}, nil
}

func (f *Filter) Matches(title string) bool {
	matchesRules := len(f.Rules) == 0
	matchesExclude := false

	for _, rule := range f.Rules {
		if rule.MatchString(title) {
			matchesRules = true
			break
		}
	}

	for _, exclude := range f.Excludes {
		if exclude.MatchString(title) {
			matchesExclude = true
			break
		}
	}

	return matchesRules && !matchesExclude
}

func compileRules(rules []config.Rule) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, rule := range rules {
		if rule.Type == "regex" {
			re, err := regexp.Compile(rule.Value)
			if err != nil {
				return nil, err
			}
			compiled = append(compiled, re)
		}
	}
	return compiled, nil
}
