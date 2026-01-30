package handler

import (
	"context"
	"fmt"
	"strings"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
)

const businessCapabilityTagKey = "dfds.businessCapability"

// BusinessCapability2FinoutHandler exports dfds.businessCapability tags to Finout
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

		var businessCapability string = ""
		if val, exists := metadata["dfds.businessCapability"]; exists {
			businessCapability = val.(string)
		}
		capsTag[capability.ID] = businessCapability
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

	if tag, exists := tags[strings.ToLower(businessCapabilityTagKey)]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", businessCapabilityTagKey))
		var rules []*finout.CreateVirtualTagRequestRule

		var bcRuleMapForCapability = make(map[string]*finout.CreateVirtualTagRequestRule)

		for k, v := range capsTag {
			if v != "" {
				if _, ok := bcRuleMapForCapability[v]; !ok {
					bcRuleMapForCapability[v] = &finout.CreateVirtualTagRequestRule{}
					rule := bcRuleMapForCapability[v]
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

				rule := bcRuleMapForCapability[v]
				rule.Filters.Value = append(rule.Filters.Value.([]string), k)
			}
		}

		for _, rule := range bcRuleMapForCapability {
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
		var bcRuleMapForCapability = make(map[string]*finout.UpdateVirtualTagRequestRule)

		for k, v := range capsTag {
			if v != "" {
				if _, ok := bcRuleMapForCapability[v]; !ok {
					bcRuleMapForCapability[v] = &finout.UpdateVirtualTagRequestRule{}
					rule := bcRuleMapForCapability[v]
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

				rule := bcRuleMapForCapability[v]
				rule.Filters.Value = append(rule.Filters.Value.([]string), k)
			}
		}

		for _, rule := range bcRuleMapForCapability {
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
