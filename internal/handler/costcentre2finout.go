package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go.dfds.cloud/aad-finout-sync/internal/config"
	"go.dfds.cloud/aad-finout-sync/internal/finout"
	"go.dfds.cloud/aad-finout-sync/internal/ssu"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"go.uber.org/zap"
	"os"
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

	mappings, err := getMappings()
	if err != nil {
		util.Logger.Warn("No manual mappings found, using default values", zap.Error(err))
		mappings = &dataMappings{
			AwsAccountAlias2CostCentre: []dataMappingsAwsAccountAlias2CostCentre{},
		}
	}

	if tag, exists := tags[tagKey]; !exists {
		util.Logger.Info(fmt.Sprintf("Tag '%s' doesn't exist, creating", tagKey))
		var rules []finout.CreateVirtualTagRequestRule

		for k, v := range capsTag {
			if v != "" {
				rules = append(rules, finout.CreateVirtualTagRequestRule{
					To: v,
					Filters: finout.CreateVirtualTagRequestRuleFilter{
						CostCenter: "virtualTag",
						Key:        "52c02d7e-093a-42b7-bf06-eb13050a8687", //id of capability virtual tag, retrieve this dynamically later
						Type:       "virtual_tag",
						Operator:   "oneOf",
						Value:      []string{k},
					},
				})
			}
		}

		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			rules = append(rules, finout.CreateVirtualTagRequestRule{
				To: mapping.CostCentre,
				Filters: finout.CreateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{mapping.Alias},
				},
			})
		}

		virtualTagRequest := finout.CreateVirtualTagRequest{
			Default: finout.CreateVirtualTagRequestDefault{
				Type:  "string",
				Value: "Untagged",
			},
			Rules:     rules,
			Category:  "Project",
			Name:      tagKey,
			UpdatedBy: author,
			CreatedBy: author,
		}
		_, err := finoutClientApp.ApiApp().CreateVirtualTag(ctx, virtualTagRequest)
		if err != nil {
			return err
		}
	} else {
		util.Logger.Info(fmt.Sprintf("Tag '%s' exists, updating", tagKey))

		var rules []finout.UpdateVirtualTagRequestRule

		for k, v := range capsTag {
			if v != "" {
				rules = append(rules, finout.UpdateVirtualTagRequestRule{
					To: v,
					Filters: finout.UpdateVirtualTagRequestRuleFilter{
						CostCenter: "virtualTag",
						Key:        "52c02d7e-093a-42b7-bf06-eb13050a8687", //id of capability virtual tag, retrieve this dynamically later
						Type:       "virtual_tag",
						Operator:   "oneOf",
						Value:      []string{k},
					},
				})
			}
		}

		for _, mapping := range mappings.AwsAccountAlias2CostCentre {
			rules = append(rules, finout.UpdateVirtualTagRequestRule{
				To: mapping.CostCentre,
				Filters: finout.UpdateVirtualTagRequestRuleFilter{
					CostCenter: "amazon-cur",
					Key:        "aws_account_name",
					Type:       "tag",
					Operator:   "oneOf",
					Value:      []string{mapping.Alias},
				},
			})
		}

		virtualTagUpdateRequest := finout.UpdateVirtualTagRequest{
			Rules:     rules,
			Category:  "Project",
			Endpoints: []string{},
			Name:      tagKey,
			UpdatedBy: author,
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
	AwsAccountAlias2CostCentre []dataMappingsAwsAccountAlias2CostCentre `json:"awsAccountAlias2CostCentre"`
}

type dataMappingsAwsAccountAlias2CostCentre struct {
	Alias      string `json:"alias"`
	CostCentre string `json:"costCentre"`
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
