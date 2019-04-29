package plugin

import (
	"github.com/thundra-io/thundra-lambda-agent-go/application"
	"github.com/thundra-io/thundra-lambda-agent-go/config"
	"github.com/thundra-io/thundra-lambda-agent-go/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/utils"
)

type CompositeDataModel struct {
	BaseDataModel
	ID                string      `json:"id"`
	Type              string      `json:"type"`
	AllMonitoringData interface{} `json:"allMonitoringData"`
}

type BaseDataModel struct {
	AgentVersion              *string                 `json:"agentVersion,omitempty"`
	DataModelVersion          *string                 `json:"dataModelVersion,omitempty"`
	ApplicationID             *string                 `json:"applicationId,omitempty"`
	ApplicationDomainName     *string                 `json:"applicationDomainName,omitempty"`
	ApplicationClassName      *string                 `json:"applicationClassName,omitempty"`
	ApplicationName           *string                 `json:"applicationName,omitempty"`
	ApplicationVersion        *string                 `json:"applicationVersion,omitempty"`
	ApplicationStage          *string                 `json:"applicationStage,omitempty"`
	ApplicationRuntime        *string                 `json:"applicationRuntime,omitempty"`
	ApplicationRuntimeVersion *string                 `json:"applicationRuntimeVersion,omitempty"`
	ApplicationTags           *map[string]interface{} `json:"applicationTags,omitempty"`
}

func PrepareCompositeData(baseDataModel BaseDataModel, allData []MonitoringDataWrapper) CompositeDataModel {

	var allDataUnwrapped []Data
	for i := range allData {
		allDataUnwrapped = append(allDataUnwrapped, allData[i].Data)
	}

	return CompositeDataModel{
		BaseDataModel:     baseDataModel,
		ID:                utils.GenerateNewID(),
		Type:              "Composite",
		AllMonitoringData: allDataUnwrapped,
	}
}

func PrepareBaseData() BaseDataModel {
	agentVersion := constants.AgentVersion
	dataModelVersion := constants.DataModelVersion
	applicationRuntime := application.ApplicationRuntime
	applicationRuntimeVersion := application.ApplicationRuntimeVersion
	return BaseDataModel{
		AgentVersion:              &agentVersion,
		DataModelVersion:          &dataModelVersion,
		ApplicationID:             &application.ApplicationID,
		ApplicationDomainName:     &application.ApplicationDomainName,
		ApplicationClassName:      &application.ApplicationClassName,
		ApplicationName:           &application.ApplicationName,
		ApplicationVersion:        &application.ApplicationVersion,
		ApplicationStage:          &application.ApplicationStage,
		ApplicationRuntime:        &applicationRuntime,
		ApplicationRuntimeVersion: &applicationRuntimeVersion,
		ApplicationTags:           &application.ApplicationTags,
	}
}

func GetBaseData() BaseDataModel {
	if (config.ReportRestCompositeDataEnabled && !config.ReportPublishCloudwatchEnabled) ||
		(config.ReportPublishCloudwatchEnabled && config.ReportCloudwatchCompositeDataEnabled) {
		return BaseDataModel{}
	}
	return PrepareBaseData()
}
