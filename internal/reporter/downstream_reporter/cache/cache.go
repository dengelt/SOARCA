package cache

import (
	b64 "encoding/base64"
	"errors"
	"reflect"
	"soarca/logger"
	"soarca/models/cacao"
	cache_report "soarca/models/cache"
	itime "soarca/utils/time"
	"time"

	"github.com/google/uuid"
)

var component = reflect.TypeOf(Cache{}).PkgPath()
var log *logger.Log

func init() {
	log = logger.Logger(component, logger.Info, "", logger.Json)
}

const MaxExecutions int = 10

type Cache struct {
	Size         int
	timeUtil     itime.ITime
	Cache        map[string]cache_report.ExecutionEntry // Cached up to max
	fifoRegister []string                               // Used for O(1) FIFO cache management
}

func New(timeUtil itime.ITime, maxExecutions int) *Cache {
	return &Cache{
		Size:     maxExecutions,
		Cache:    make(map[string]cache_report.ExecutionEntry),
		timeUtil: timeUtil,
	}
}

func (cacheReporter *Cache) getExecution(executionKey uuid.UUID) (cache_report.ExecutionEntry, error) {
	executionKeyStr := executionKey.String()
	executionEntry, ok := cacheReporter.Cache[executionKeyStr]
	if !ok {
		err := errors.New("execution is not in cache")
		log.Warning("execution is not in cache. consider increasing cache size.")
		return cache_report.ExecutionEntry{}, err
		// TODO Retrieve from database
	}
	return executionEntry, nil
}
func (cacheReporter *Cache) getExecutionStep(executionKey uuid.UUID, stepKey string) (cache_report.StepResult, error) {
	executionEntry, err := cacheReporter.getExecution(executionKey)
	if err != nil {
		return cache_report.StepResult{}, err
	}
	executionStep, ok := executionEntry.StepResults[stepKey]
	if !ok {
		err := errors.New("execution step is not in cache")
		return cache_report.StepResult{}, err
	}
	return executionStep, nil
}

// Adding executions in FIFO logic
func (cacheReporter *Cache) addExecution(newExecutionEntry cache_report.ExecutionEntry) error {

	if !(len(cacheReporter.fifoRegister) == len(cacheReporter.Cache)) {
		return errors.New("cache fifo register and content are desynchronized")
	}

	newExecutionEntryKey := newExecutionEntry.ExecutionId.String()

	if len(cacheReporter.fifoRegister) >= cacheReporter.Size {
		firstExecution := cacheReporter.fifoRegister[0]
		cacheReporter.fifoRegister = cacheReporter.fifoRegister[1:]
		delete(cacheReporter.Cache, firstExecution)

		cacheReporter.fifoRegister = append(cacheReporter.fifoRegister, newExecutionEntryKey)
		cacheReporter.Cache[newExecutionEntryKey] = newExecutionEntry
		return nil
	}

	cacheReporter.fifoRegister = append(cacheReporter.fifoRegister, newExecutionEntryKey)
	cacheReporter.Cache[newExecutionEntryKey] = newExecutionEntry
	return nil
}

func (cacheReporter *Cache) ReportWorkflowStart(executionId uuid.UUID, playbook cacao.Playbook) error {
	newExecutionEntry := cache_report.ExecutionEntry{
		ExecutionId: executionId,
		PlaybookId:  playbook.ID,
		Started:     cacheReporter.timeUtil.Now(),
		Ended:       time.Time{},
		StepResults: map[string]cache_report.StepResult{},
		Status:      cache_report.Ongoing,
	}
	err := cacheReporter.addExecution(newExecutionEntry)
	if err != nil {
		return err
	}
	return nil
}

func (cacheReporter *Cache) ReportWorkflowEnd(executionId uuid.UUID, playbook cacao.Playbook, workflowError error) error {

	executionEntry, err := cacheReporter.getExecution(executionId)
	if err != nil {
		return err
	}

	if workflowError != nil {
		executionEntry.Error = workflowError
		executionEntry.Status = cache_report.Failed
	} else {
		executionEntry.Status = cache_report.SuccessfullyExecuted
	}
	executionEntry.Ended = cacheReporter.timeUtil.Now()
	cacheReporter.Cache[executionId.String()] = executionEntry

	return nil
}

func (cacheReporter *Cache) ReportStepStart(executionId uuid.UUID, step cacao.Step, variables cacao.Variables) error {
	executionEntry, err := cacheReporter.getExecution(executionId)
	if err != nil {
		return err
	}

	if executionEntry.Status != cache_report.Ongoing {
		return errors.New("trying to report on the execution of a step for an already reported completed or failed execution")
	}

	_, alreadyThere := executionEntry.StepResults[step.ID]
	if alreadyThere {
		log.Warning("a step execution was already reported for this step. overwriting.")
	}

	// TODO: must test
	commandsB64 := []string{}
	isAutomated := true
	for _, cmd := range step.Commands {
		if cmd.Type == cacao.CommandTypeManual {
			isAutomated = false
		}
		if cmd.CommandB64 != "" {
			commandsB64 = append(commandsB64, cmd.CommandB64)
		} else {
			cmdB64 := b64.StdEncoding.EncodeToString([]byte(cmd.Command))
			commandsB64 = append(commandsB64, cmdB64)
		}
	}

	newStepEntry := cache_report.StepResult{
		ExecutionId: executionId,
		StepId:      step.ID,
		Started:     cacheReporter.timeUtil.Now(),
		Ended:       time.Time{},
		Variables:   variables,
		CommandsB64: commandsB64,
		Status:      cache_report.Ongoing,
		Error:       nil,
		IsAutomated: isAutomated,
	}
	executionEntry.StepResults[step.ID] = newStepEntry
	return nil
}

func (cacheReporter *Cache) ReportStepEnd(executionId uuid.UUID, step cacao.Step, returnVars cacao.Variables, stepError error) error {
	executionEntry, err := cacheReporter.getExecution(executionId)
	if err != nil {
		return err
	}

	if executionEntry.Status != cache_report.Ongoing {
		return errors.New("trying to report on the execution of a step for an already reported completed or failed execution")
	}

	executionStepResult, err := cacheReporter.getExecutionStep(executionId, step.ID)
	if err != nil {
		return err
	}

	if executionStepResult.Status != cache_report.Ongoing {
		return errors.New("trying to report on the execution of a step that was already reported completed or failed")
	}

	if stepError != nil {
		executionStepResult.Error = stepError
		executionStepResult.Status = cache_report.ServerSideError
	} else {
		executionStepResult.Status = cache_report.SuccessfullyExecuted
	}
	executionStepResult.Ended = cacheReporter.timeUtil.Now()
	executionStepResult.Variables = returnVars
	executionEntry.StepResults[step.ID] = executionStepResult

	return nil
}

func (cacheReporter *Cache) GetExecutions() ([]cache_report.ExecutionEntry, error) {
	executions := make([]cache_report.ExecutionEntry, 0)
	// NOTE: fetched via fifo register key reference as is ordered array,
	// needed to test and report back ordered executions stored
	for _, executionEntryKey := range cacheReporter.fifoRegister {
		// NOTE: cached executions are passed by reference, so they must not be modified
		entry, present := cacheReporter.Cache[executionEntryKey]
		if !present {
			return []cache_report.ExecutionEntry{}, errors.New("internal error. cache fifo register and cache executions mismatch.")
		}
		executions = append(executions, entry)
	}
	return executions, nil
}

func (cacheReporter *Cache) GetExecutionReport(executionKey uuid.UUID) (cache_report.ExecutionEntry, error) {
	executionEntry, err := cacheReporter.getExecution(executionKey)
	if err != nil {
		return cache_report.ExecutionEntry{}, err
	}
	report := executionEntry

	return report, nil
}
