package aws

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"regexp"
	"strings"
)

//go:generate curl -o reference.json https://raw.githubusercontent.com/fluggo/aws-service-auth-reference/master/service-auth.json
// bash -c "curl https://raw.githubusercontent.com/fluggo/aws-service-auth-reference/master/service-auth.json | jq -c > reference.json"

//go:embed reference.json
var referenceRaw []byte
var reference []ServiceAuthorizationReference
var actions []string

func init() {
	err := json.Unmarshal(referenceRaw, &reference)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load AWS reference")
	}
	for _, service := range reference {
		for _, action := range service.Actions {
			actions = append(actions, fmt.Sprintf("%s:%s", service.ServicePrefix, action.Name))
		}
	}
}

// TODO don't copy the following structs, the original library does not contain a package
// https://github.com/fluggo/aws-service-auth-reference

type ServiceAuthorizationReference struct {
	Name              string          `json:"name"`
	ServicePrefix     string          `json:"servicePrefix"`
	AuthReferenceHref string          `json:"authReferenceHref"`
	ApiReferenceHref  string          `json:"apiReferenceHref,omitempty"`
	Actions           []*Action       `json:"actions"`
	ResourceTypes     []*ResourceType `json:"resourceTypes"`
	ConditionKeys     []*ConditionKey `json:"conditionKeys"`
}

type ActionResourceType struct {
	ResourceType     string   `json:"resourceType"`
	Required         bool     `json:"required"`
	ConditionKeys    []string `json:"conditionKeys"`
	DependentActions []string `json:"dependentActions"`
}

type Action struct {
	Name           string               `json:"name"`
	PermissionOnly bool                 `json:"permissionOnly"`
	ReferenceHref  string               `json:"referenceHref,omitempty"`
	Description    string               `json:"description"`
	AccessLevel    string               `json:"accessLevel"`
	ResourceTypes  []ActionResourceType `json:"resourceTypes"`
}

type ResourceType struct {
	Name          string   `json:"name"`
	ReferenceHref string   `json:"referenceHref,omitempty"`
	ArnPattern    string   `json:"arnPattern"`
	ConditionKeys []string `json:"conditionKeys"`
}

type ConditionKey struct {
	Name          string `json:"name"`
	ReferenceHref string `json:"referenceHref,omitempty"`
	Description   string `json:"description"`
	Type          string `json:"type"`
}

func GetAllActions() []string {
	return actions
}

func ExpandActions(action string) ([]string, error) {
	log.Debug().Str("action", action)

	// Replace *
	action = strings.Replace(action, "*", ".*", -1)

	// Replace ?
	action = strings.Replace(action, "?", ".{1}", -1)

	patternStr := fmt.Sprintf("(?i)^%s$", action)
	log.Debug().Str("pattern", patternStr)

	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		return nil, err
	}

	return lo.Filter[string](actions, func(action string, _ int) bool {
		return pattern.MatchString(action)
	}), nil
}
