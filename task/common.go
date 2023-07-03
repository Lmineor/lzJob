package task

import "fmt"

const (
	CronType              = "cron"
	IntervalType          = "interval"
	FixedTimeSingleType   = "fixed_time_single"
	DelayedTimeSingleType = "delayed_time_single"
)

func ValidateTaskType(t string) error {
	switch t {
	case CronType, IntervalType, FixedTimeSingleType, DelayedTimeSingleType:
	default:
		return fmt.Errorf("invalid task type, you can select from %s, %s, %s or %s", CronType, IntervalType,
			FixedTimeSingleType, DelayedTimeSingleType)

	}
	return nil
}
