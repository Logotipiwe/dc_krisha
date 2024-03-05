package pkg

import "time"

func SleepWithInterruption(getSleepTimeFunc func() time.Duration, checkInterval time.Duration) {
	sleepFor := getSleepTimeFunc()
	if sleepFor < checkInterval {
		time.Sleep(sleepFor)
	} else {
		slept := time.Duration(0)
		for {
			left := sleepFor - slept
			if left > checkInterval {
				time.Sleep(checkInterval)
				slept += checkInterval
			} else {
				if left > 0 {
					time.Sleep(left)
				}
				return
			}
			sleepFor = getSleepTimeFunc()
		}
	}
}
