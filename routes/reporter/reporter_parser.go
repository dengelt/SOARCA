package reporter

import (
	api_model "soarca/models/api"
	cache_model "soarca/models/cache"
)

const defaultRequestInterval int = 5

func parseCachePlaybookEntry(cacheEntry cache_model.ExecutionEntry) (api_model.PlaybookExecutionReport, error) {
	playbookStatus, err := api_model.CacheStatusEnum2String(cacheEntry.Status)
	if err != nil {
		return api_model.PlaybookExecutionReport{}, err
	}

	playbookStatusText, err := api_model.GetCacheStatusText(playbookStatus, api_model.ReportLevelPlaybook)
	if err != nil {
		return api_model.PlaybookExecutionReport{}, err
	}
	if cacheEntry.Error != nil {
		playbookStatusText = playbookStatusText + " - error: " + cacheEntry.Error.Error()
	}

	stepResults, err := parseCacheStepEntries(cacheEntry.StepResults)
	if err != nil {
		return api_model.PlaybookExecutionReport{}, err
	}

	executionReport := api_model.PlaybookExecutionReport{
		Type:            "execution_status",
		ExecutionId:     cacheEntry.ExecutionId.String(),
		PlaybookId:      cacheEntry.PlaybookId,
		Started:         cacheEntry.Started,
		Ended:           cacheEntry.Ended,
		Status:          playbookStatus,
		StatusText:      playbookStatusText,
		StepResults:     stepResults,
		RequestInterval: defaultRequestInterval,
	}
	return executionReport, nil
}

func parseCacheStepEntries(cacheStepEntries map[string]cache_model.StepResult) (map[string]api_model.StepExecutionReport, error) {
	parsedEntries := map[string]api_model.StepExecutionReport{}
	for stepId, stepEntry := range cacheStepEntries {

		stepStatus, err := api_model.CacheStatusEnum2String(stepEntry.Status)
		if err != nil {
			return map[string]api_model.StepExecutionReport{}, err
		}
		stepStatusText, err := api_model.GetCacheStatusText(stepStatus, api_model.ReportLevelStep)
		if err != nil {
			return map[string]api_model.StepExecutionReport{}, err
		}

		if stepEntry.Error != nil {
			stepStatusText = stepStatusText + " - error: " + stepEntry.Error.Error()
		}

		parsedEntries[stepId] = api_model.StepExecutionReport{
			ExecutionId:        stepEntry.ExecutionId.String(),
			StepId:             stepEntry.StepId,
			Started:            stepEntry.Started,
			Ended:              stepEntry.Ended,
			Status:             stepStatus,
			StatusText:         stepStatusText,
			ExecutedBy:         "soarca",
			CommandsB64:        stepEntry.CommandsB64,
			Variables:          stepEntry.Variables,
			AutomatedExecution: stepEntry.IsAutomated,
		}
	}
	return parsedEntries, nil
}
