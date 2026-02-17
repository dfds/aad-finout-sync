package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"go.uber.org/zap"
)

const CostCentreToFinoutName = "costCenterToFinout"

const tagKey = "dfds.cost.centre"
const author = "aad-finout-sync"

func CostCentre2FinoutHandler(ctx context.Context) error {
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

		var tribe string = ""
		if val, exists := metadata["dfds.cost.centre"]; exists {
			tribe = val.(string)
		}
		capsTag[capability.ID] = tribe
	}

	util.Logger.Debug("Capability metadata retrieved")

	tags, err := finoutClientApp.ApiApp().ListVirtualTags(ctx)
	if err != nil {
		return err
	}

	capabilityTag, exists := tags["capability"]
	if !exists {
		return VirtualTagDoesNotExist.New(VirtualTagDoesNotExistMsg)
	}

	mappings, err := getMappings()
	if err != nil {
		util.Logger.Warn("No manual mappings found, using default values", zap.Error(err))
		mappings = &dataMappings{
			AwsAccountAlias2CostCentre: []dataMappingsAwsAccountAlias2CostCentre{},
		}
	}

	if tag, exists := tags[tagKey]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", tagKey))
		var rules []*finout.CreateVirtualTagRequestRule

		var ccRuleMapForCapability = make(map[string]*finout.CreateVirtualTagRequestRule)
		var ccRuleMapForAwsAccount = make(map[string]*finout.CreateVirtualTagRequestRule)

		for k, v := range capsTag {
			if v != "" {
				if _, ok := ccRuleMapForCapability[v]; !ok {
					ccRuleMapForCapability[v] = &finout.CreateVirtualTagRequestRule{}
					rule := ccRuleMapForCapability[v]
					rule.To = v
					rule.Type = "string"
					rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
						CostCenter: "virtualTag",
						Key:        capabilityTag.ID,
						Type:       "virtual_tag",
						Operator:   "oneOf",
						Value:      []string{},
					}
				}

				rule := ccRuleMapForCapability[v]
				rule.Filters.Value = append(rule.Filters.Value.([]string), k)
			}
		}

		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			if _, ok := ccRuleMapForAwsAccount[mapping.CostCentre]; !ok {
				ccRuleMapForAwsAccount[mapping.CostCentre] = &finout.CreateVirtualTagRequestRule{}
				rule := ccRuleMapForAwsAccount[mapping.CostCentre]
				rule.To = mapping.CostCentre
				rule.Type = "string"
				rule.Filters = finout.CreateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}

			rule := ccRuleMapForAwsAccount[mapping.CostCentre]
			rule.Filters.Value = append(rule.Filters.Value.([]interface{}), mapping.Alias)
		}

		for _, rule := range ccRuleMapForCapability {
			rules = append(rules, rule)
		}
		for _, rule := range ccRuleMapForAwsAccount {
			rules = append(rules, rule)
		}

		virtualTagRequest := finout.CreateVirtualTagRequest{
			Default: finout.CreateVirtualTagRequestDefault{
				Type:  "string",
				Value: "Untagged",
			},
			Rules: rules,
			Name:  tagKey,
		}
		_, err := finoutClientApp.ApiApp().CreateVirtualTag(ctx, virtualTagRequest)
		if err != nil {
			return err
		}
	} else {
		util.Logger.Info(fmt.Sprintf("Tag '%s' exists, updating", tagKey))

		var rules []*finout.UpdateVirtualTagRequestRule

		var ccRuleMapForCapability = make(map[string]*finout.UpdateVirtualTagRequestRule)
		var ccRuleMapForAwsAccount = make(map[string]*finout.UpdateVirtualTagRequestRule)

		for k, v := range capsTag {
			if v != "" {
				if _, ok := ccRuleMapForCapability[v]; !ok {
					ccRuleMapForCapability[v] = &finout.UpdateVirtualTagRequestRule{}
					rule := ccRuleMapForCapability[v]
					rule.To = v
					rule.Type = "string"
					rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
						CostCenter: "virtualTag",
						Key:        capabilityTag.ID,
						Type:       "virtual_tag",
						Operator:   "oneOf",
						Value:      []string{},
					}
				}

				rule := ccRuleMapForCapability[v]
				rule.Filters.Value = append(rule.Filters.Value.([]string), k)
			}
		}

		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			if _, ok := ccRuleMapForAwsAccount[mapping.CostCentre]; !ok {
				ccRuleMapForAwsAccount[mapping.CostCentre] = &finout.UpdateVirtualTagRequestRule{}
				rule := ccRuleMapForAwsAccount[mapping.CostCentre]
				rule.To = mapping.CostCentre
				rule.Type = "string"
				rule.Filters = finout.UpdateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{},
				}
			}

			rule := ccRuleMapForAwsAccount[mapping.CostCentre]
			rule.Filters.Value = append(rule.Filters.Value.([]string), mapping.Alias)
		}

		for _, rule := range ccRuleMapForCapability {
			rules = append(rules, rule)
		}
		for _, rule := range ccRuleMapForAwsAccount {
			rules = append(rules, rule)
		}

		virtualTagUpdateRequest := finout.UpdateVirtualTagRequest{
			Rules:       rules,
			Endpoints:   []string{},
			Name:        tagKey,
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

type dataMappings struct {
	AwsAccountAlias2CostCentre            []dataMappingsAwsAccountAlias2CostCentre            `json:"awsAccountAlias2CostCentre"`
	TechnicalCapability2BusinessCapability []dataMappingsTechnicalCapability2BusinessCapability `json:"technicalCapability2BusinessCapability"`
}

type dataMappingsAwsAccountAlias2CostCentre struct {
	Alias      string  `json:"alias"`
	MappedName *string `json:"mappedName"`
	CostCentre string  `json:"costCentre"`
}

type dataMappingsTechnicalCapability2BusinessCapability struct {
	TechnicalCapability string `json:"technicalCapability"`
	BusinessCapability  string `json:"businessCapability"`
}

func getMappings() (*dataMappings, error) {
	bytes, err := os.ReadFile("mapping.json")
	if err != nil {
		return nil, err
	}

	var payload *dataMappings

	err = json.Unmarshal(bytes, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
