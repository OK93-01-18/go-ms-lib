package schedule

import (
	"github.com/ok93-01-18/go-ms-lib/log"
)

type Operation struct {
	Name       string
	Interval   string
	Task       TaskInterface
	RunOnInit  bool
	IsNeedFunc func() (bool, error)
}

type CronManager struct {
	appLogger  log.Logger
	cron       Croner
	operations *[]Operation
}

func (m *CronManager) Init() error {

	lenOperations := len(*m.operations)
	if lenOperations == 0 {
		return nil
	}

	for i := 0; i < lenOperations; i++ {
		var err error

		operation := (*m.operations)[i]
		if operation.RunOnInit {
			isNeed := true
			if operation.IsNeedFunc != nil {
				isNeed, err = operation.IsNeedFunc()
				if err != nil {
					return err
				}
			}

			if isNeed {
				err = operation.Task.Do()
				if err != nil {
					return err
				}
			}

		}

		err = m.cron.AddFunc(operation.Interval, func() {
			err = operation.Task.Do()
			if err != nil {
				m.appLogger.Errorf(log.TypeApp, "%s error: %v", operation.Name, err)
			}
		})

		if err != nil {
			return err
		}
	}

	m.cron.Start()

	return nil
}

func NewCronManager(appLogger log.Logger, cron Croner, operations *[]Operation) *CronManager {
	return &CronManager{
		appLogger:  appLogger,
		cron:       cron,
		operations: operations,
	}
}
