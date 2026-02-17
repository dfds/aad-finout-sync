package handler

import (
	"context"
	"fmt"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"go.uber.org/zap"
)

const capabilityAndLegacyAccountsTagKey = "capability & legacy accounts"

func CapabilityAndLegacyAccounts2FinoutHandler(ctx context.Context) error {
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

	mappings, err := getMappings()
	if err != nil {
		util.Logger.Warn("No manual mappings found, using default values", zap.Error(err))
		mappings = &dataMappings{
			AwsAccountAlias2CostCentre: []dataMappingsAwsAccountAlias2CostCentre{},
		}
	}

	tags, err := finoutClientApp.ApiApp().ListVirtualTags(ctx)
	if err != nil {
		return err
	}

	if tag, exists := tags[capabilityAndLegacyAccountsTagKey]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", capabilityAndLegacyAccountsTagKey))
		var rules []*finout.CreateVirtualTagRequestRule

		// Capability rules: key = capability ID
		var capRuleMap = make(map[string]*finout.CreateVirtualTagRequestRule)
		for _, capa := range caps {
			if _, ok := capRuleMap[capa.ID]; !ok {
				capRuleMap[capa.ID] = &finout.CreateVirtualTagRequestRule{}
				rule := capRuleMap[capa.ID]
				rule.To = capa.ID
				rule.Type = "string"
				rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
					Or: []finout.CreateVirtualTagRequestRuleFilter{
						{
							CostCenter: "amazon-cur",
							Key:        "aws_account_name",
							Type:       "tag",
							Operator:   "oneOf",
							Value:      []string{fmt.Sprintf("dfds-%s", capa.RootID)},
						},
						{
							CostCenter: "kubernetes",
							Key:        "k8s_namespace",
							Type:       "kubernetesResources",
							Operator:   "oneOf",
							Value:      []string{capa.RootID},
						},
						{
							CostCenter: "Azure",
							Key:        "resourcegroup",
							Type:       "tag",
							Operator:   "contains",
							Value:      capa.RootID,
						},
					},
				}
			}
		}

		// Legacy account rules: one rule per alias, To = alias
		var legacyRuleMap = make(map[string]*finout.CreateVirtualTagRequestRule)
		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			if _, ok := legacyRuleMap[mapping.Alias]; !ok {
				mapTo := ""
				if mapping.MappedName != nil {
					mapTo = *mapping.MappedName
				} else {
					mapTo = mapping.Alias
				}
				legacyRuleMap[mapping.Alias] = &finout.CreateVirtualTagRequestRule{}
				rule := legacyRuleMap[mapping.Alias]
				rule.To = mapTo
				rule.Type = "string"
				rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{mapping.Alias},
				}
			}
		}

		for _, rule := range capRuleMap {
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
			Name:  capabilityAndLegacyAccountsTagKey,
		}
		_, err := finoutClientApp.ApiApp().CreateVirtualTag(ctx, virtualTagRequest)
		if err != nil {
			return err
		}

	} else {
		util.Logger.Info(fmt.Sprintf("Tag '%s' exists, updating", capabilityAndLegacyAccountsTagKey))
		var rules []*finout.UpdateVirtualTagRequestRule

		// Capability rules: key = capability ID
		var capRuleMap = make(map[string]*finout.UpdateVirtualTagRequestRule)
		for _, capa := range caps {
			if _, ok := capRuleMap[capa.ID]; !ok {
				capRuleMap[capa.ID] = &finout.UpdateVirtualTagRequestRule{}
				rule := capRuleMap[capa.ID]
				rule.To = capa.ID
				rule.Type = "string"
				rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
					Or: []finout.UpdateVirtualTagRequestRuleFilter{
						{
							CostCenter: "amazon-cur",
							Key:        "aws_account_name",
							Type:       "tag",
							Operator:   "oneOf",
							Value:      []string{fmt.Sprintf("dfds-%s", capa.RootID)},
						},
						{
							CostCenter: "kubernetes",
							Key:        "k8s_namespace",
							Type:       "kubernetesResources",
							Operator:   "oneOf",
							Value:      []string{capa.RootID},
						},
						{
							CostCenter: "Azure",
							Key:        "resourcegroup",
							Type:       "tag",
							Operator:   "contains",
							Value:      capa.RootID,
						},
					},
				}
			}
		}

		// Legacy account rules: one rule per alias, To = alias
		var legacyRuleMap = make(map[string]*finout.UpdateVirtualTagRequestRule)
		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			if _, ok := legacyRuleMap[mapping.Alias]; !ok {
				legacyRuleMap[mapping.Alias] = &finout.UpdateVirtualTagRequestRule{}
				mapTo := ""
				if mapping.MappedName != nil {
					mapTo = *mapping.MappedName
				} else {
					mapTo = mapping.Alias
				}
				rule := legacyRuleMap[mapping.Alias]
				rule.To = mapTo
				rule.Type = "string"
				rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{mapping.Alias},
				}
			}
		}

		for _, rule := range capRuleMap {
			rules = append(rules, rule)
		}
		for _, rule := range legacyRuleMap {
			rules = append(rules, rule)
		}

		virtualTagUpdateRequest := finout.UpdateVirtualTagRequest{
			Rules:       rules,
			Endpoints:   []string{},
			Name:        capabilityAndLegacyAccountsTagKey,
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
