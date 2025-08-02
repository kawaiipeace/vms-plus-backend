package funcs

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func InitCronJob() {
	RunCronJobDriverCheckActive()
}

func RunCronJobDriverCheckActive() {
	c := cron.New()

	// Schedule to run every day at midnight (00:00)
	c.AddFunc("0 0 * * *", func() {
		fmt.Println("Running task at:", time.Now().Format("2006-01-02 15:04:05"))
		//JobDriversCheckActive()
	})

	c.Start()
}
