package registry

import "github.com/MAD-py/go-taskengine/taskengine"

func ProjectQuotaReset(e *taskengine.Engine, task *taskengine.Task) error {
	trigger, err := taskengine.NewCronTrigger("0 0 * * *", true)
	if err != nil {
		return err
	}

	return e.RegisterTask(
		task,
		taskengine.WorkerPolicySerial,
		trigger,
		true,
		0,
	)
}
