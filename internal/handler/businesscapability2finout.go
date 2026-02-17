package handler

import (
	"context"
	"fmt"
	"strings"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"go.uber.org/zap"
)

const businessCapabilityTagKey = "Business Capability"

// BusinessCapability2FinoutHandler exports dfds.businessCapability tags to Finout.
// Rules are sourced from two places:
//   - SSU metadata (dfds.businessCapability), filtered via the "capability & legacy accounts" virtual tag
//   - mapping.json (TechnicalCapability2BusinessCapability), filtered via the "capability & legacy accounts" virtual tag
//
// SSU metadata takes priority: if a TechnicalCapability value from the JSON matches a capability
// ID already covered by SSU metadata, the JSON entry is skipped.
func BusinessCapability2FinoutHandler(ctx context.Context) error {
	conf, err := config.LoadConfig()
	if err != nil {
		return err
	}

	finoutClientApp := finout.NewFinoutClient()
	finoutClientApp.SetAuthMethod(finout.AuthClientSecretMethod(finout.Config{ClientId: conf.Finout.ClientId, ClientSecret: conf.Finout.ClientSecret}))
	ssuClient := ssu.NewSsuClient(ssu.Config{
		Host:         conf.CapSvc.Host,
		TenantId:     conf.Azure.TenantId,
		ClientId:     conf.CapSvc.ClientId,
		ClientSecret: conf.CapSvc.ClientSecret,
		Scope:        conf.CapSvc.TokenScope,
	})

	caps, err := ssuClient.GetCapabilities()
	if err != nil {
		return err
	}
	util.Logger.Debug("Capabilities retrieved")

	capsTag := make(map[string]string)
	for _, capability := range caps {
		metadata, err := ssuClient.GetCapabilityMetadata(capability.ID)
		if err != nil {
			return err
		}

		var businessCapability string
		if val, exists := metadata["dfds.businessCapability"]; exists {
			businessCapability = val.(string)
		}
		capsTag[capability.ID] = businessCapability
	}
	util.Logger.Debug("Capability metadata retrieved")

	coveredCapIDs := make(map[string]struct{})
	for capID, bc := range capsTag {
		if bc != "" {
			coveredCapIDs[capID] = struct{}{}
		}
	}

	mappings, err := getMappings()
	if err != nil {
		util.Logger.Warn("No manual mappings found, using default values", zap.Error(err))
		mappings = &dataMappings{
			AwsAccountAlias2CostCentre:             []dataMappingsAwsAccountAlias2CostCentre{},
			TechnicalCapability2BusinessCapability: []dataMappingsTechnicalCapability2BusinessCapability{},
		}
	}

	tags, err := finoutClientApp.ApiApp().ListVirtualTags(ctx)
	if err != nil {
		return err
	}

	capAndLegacyTag, exists := tags[strings.ToLower(capabilityAndLegacyAccountsTagKey)]
	if !exists {
		return VirtualTagDoesNotExist.New(VirtualTagDoesNotExistMsg)
	}

	if tag, exists := tags[strings.ToLower(businessCapabilityTagKey)]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", businessCapabilityTagKey))
		var rules []*finout.CreateVirtualTagRequestRule

		capabilityRuleMap := make(map[string]*finout.CreateVirtualTagRequestRule)
		for capID, bc := range capsTag {
			if bc == "" {
				continue
			}
			if _, ok := capabilityRuleMap[bc]; !ok {
				capabilityRuleMap[bc] = &finout.CreateVirtualTagRequestRule{}
				rule := capabilityRuleMap[bc]
				rule.To = bc
				rule.Type = "string"
				rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
					CostCenter: "virtualTag",
					Key:        capAndLegacyTag.ID,
					Type:       "virtual_tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}
			capabilityRuleMap[bc].Filters.Value = append(capabilityRuleMap[bc].Filters.Value.([]string), capID)
		}

		legacyRuleMap := make(map[string]*finout.CreateVirtualTagRequestRule)
		for _, mapping := range mappings.TechnicalCapability2BusinessCapability {
			if mapping.BusinessCapability == "" {
				continue
			}
			if _, covered := coveredCapIDs[mapping.TechnicalCapability]; covered {
				continue
			}
			if _, ok := legacyRuleMap[mapping.BusinessCapability]; !ok {
				legacyRuleMap[mapping.BusinessCapability] = &finout.CreateVirtualTagRequestRule{}
				rule := legacyRuleMap[mapping.BusinessCapability]
				rule.To = mapping.BusinessCapability
				rule.Type = "string"
				rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
					CostCenter: "virtualTag",
					Key:        capAndLegacyTag.ID,
					Type:       "virtual_tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}
			legacyRuleMap[mapping.BusinessCapability].Filters.Value = append(legacyRuleMap[mapping.BusinessCapability].Filters.Value.([]string), mapping.TechnicalCapability)
		}

		for _, rule := range capabilityRuleMap {
			rules = append(rules, rule)
		}
		for _, rule := range legacyRuleMap {
			rules = append(rules, rule)
		}

		virtualTagRequest := finout.CreateVirtualTagRequest{
			Default: finout.CreateVirtualTagRequestDefault{
				Type:  "string",
				Value: "Untagged",
			},
			Rules: rules,
			Name:  businessCapabilityTagKey,
		}
		_, err := finoutClientApp.ApiApp().CreateVirtualTag(ctx, virtualTagRequest)
		if err != nil {
			return err
		}

	} else {
		util.Logger.Info(fmt.Sprintf("Tag '%s' exists, updating", businessCapabilityTagKey))
		var rules []*finout.UpdateVirtualTagRequestRule

		capabilityRuleMap := make(map[string]*finout.UpdateVirtualTagRequestRule)
		for capID, bc := range capsTag {
			if bc == "" {
				continue
			}
			if _, ok := capabilityRuleMap[bc]; !ok {
				capabilityRuleMap[bc] = &finout.UpdateVirtualTagRequestRule{}
				rule := capabilityRuleMap[bc]
				rule.To = bc
				rule.Type = "string"
				rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
					CostCenter: "virtualTag",
					Key:        capAndLegacyTag.ID,
					Type:       "virtual_tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}
			capabilityRuleMap[bc].Filters.Value = append(capabilityRuleMap[bc].Filters.Value.([]string), capID)
		}

		legacyRuleMap := make(map[string]*finout.UpdateVirtualTagRequestRule)
		for _, mapping := range mappings.TechnicalCapability2BusinessCapability {
			if mapping.BusinessCapability == "" {
				continue
			}
			if _, covered := coveredCapIDs[mapping.TechnicalCapability]; covered {
				continue
			}
			if _, ok := legacyRuleMap[mapping.BusinessCapability]; !ok {
				legacyRuleMap[mapping.BusinessCapability] = &finout.UpdateVirtualTagRequestRule{}
				rule := legacyRuleMap[mapping.BusinessCapability]
				rule.To = mapping.BusinessCapability
				rule.Type = "string"
				rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
					CostCenter: "virtualTag",
					Key:        capAndLegacyTag.ID,
					Type:       "virtual_tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}
			legacyRuleMap[mapping.BusinessCapability].Filters.Value = append(legacyRuleMap[mapping.BusinessCapability].Filters.Value.([]string), mapping.TechnicalCapability)
		}

		for _, rule := range capabilityRuleMap {
			rules = append(rules, rule)
		}
		for _, rule := range legacyRuleMap {
			rules = append(rules, rule)
		}

		virtualTagUpdateRequest := finout.UpdateVirtualTagRequest{
			Rules:       rules,
			Endpoints:   []string{},
			Name:        businessCapabilityTagKey,
			Allocations: []string{},
			Default: finout.CreateVirtualTagRequestDefault{
				Type:  "string",
				Value: "Untagged",
			},
		}
		_, err := finoutClientApp.ApiApp().UpdateVirtualTag(ctx, virtualTagUpdateRequest, tag.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
