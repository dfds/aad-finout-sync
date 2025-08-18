package handler

import (
	"context"
	"fmt"

	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
)

const capabilityTagKey = "capability"

func Capability2FinoutHandler(ctx context.Context) error {
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

	tags, err := finoutClientApp.ApiApp().ListVirtualTags(ctx)
	if err != nil {
		return err
	}

	if tag, exists := tags[capabilityTagKey]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", capabilityTagKey))
		var rules []*finout.CreateVirtualTagRequestRule

		// key: capability id
		var ccRuleMap = make(map[string]*finout.CreateVirtualTagRequestRule)

		for _, capa := range caps {
			if _, ok := ccRuleMap[capa.ID]; !ok {
				ccRuleMap[capa.ID] = &finout.CreateVirtualTagRequestRule{}
				rule := ccRuleMap[capa.ID]
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

		for _, rule := range ccRuleMap {
			rules = append(rules, rule)
		}

		virtualTagRequest := finout.CreateVirtualTagRequest{
			Default: finout.CreateVirtualTagRequestDefault{
				Type:  "string",
				Value: "Untagged",
			},
			Rules: rules,
			Name:  capabilityTagKey,
		}
		_, err := finoutClientApp.ApiApp().CreateVirtualTag(ctx, virtualTagRequest)
		if err != nil {
			return err
		}

	} else {
		util.Logger.Info(fmt.Sprintf("Tag '%s' exists, updating", capabilityTagKey))
		var rules []*finout.UpdateVirtualTagRequestRule

		// key: capability id
		var ccRuleMap = make(map[string]*finout.UpdateVirtualTagRequestRule)

		for _, capa := range caps {
			if _, ok := ccRuleMap[capa.ID]; !ok {
				ccRuleMap[capa.ID] = &finout.UpdateVirtualTagRequestRule{}
				rule := ccRuleMap[capa.ID]
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

		for _, rule := range ccRuleMap {
			rules = append(rules, rule)
		}

		virtualTagUpdateRequest := finout.UpdateVirtualTagRequest{
			Rules:       rules,
			Endpoints:   []string{},
			Name:        capabilityTagKey,
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
