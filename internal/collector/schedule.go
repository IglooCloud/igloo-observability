package collector

import "github.com/robfig/cron/v3"

type Config struct {
	Schedules []EndpointSchedule
}

type EndpointSchedule struct {
	Endpoint Endpoint
	Schedule string
	Enabled  bool
}

func RequestStream(config Config) chan Endpoint {
	endpointStream := make(chan Endpoint)

	c := cron.New()
	for _, schedule := range config.Schedules {
		if schedule.Enabled {
			endpoint := schedule.Endpoint
			c.AddFunc(schedule.Schedule, func() {
				endpointStream <- endpoint
			})
		}
	}
	c.Start()

	return endpointStream
}
