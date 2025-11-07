package worker

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/romang000/winter-service/internal/client/http/geocoding"
	"github.com/romang000/winter-service/internal/client/http/open_meteo"
	"github.com/romang000/winter-service/internal/repository"
	"log/slog"
	"time"
)

const durationJob = 5 * time.Second

type JobScheduler struct {
	ctx             context.Context
	logger          *slog.Logger
	geoClient       *geocoding.Client
	openMeteoClient *open_meteo.Client
	Scheduler       gocron.Scheduler
	cache           *repository.CacheRepository
	city            string
}

func NewJobScheduler(
	ctx context.Context,
	logger *slog.Logger,
	geoClient *geocoding.Client,
	openMeteoClient *open_meteo.Client,
	cache *repository.CacheRepository,
	city string,
) (*JobScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	
	return &JobScheduler{
		ctx:             ctx,
		logger:          logger,
		geoClient:       geoClient,
		openMeteoClient: openMeteoClient,
		Scheduler:       s,
		cache:           cache,
		city:            city,
	}, nil
}

func (j *JobScheduler) InitJobs() ([]gocron.Job, error) {
	const f = "worker.weather"
	job, err := j.Scheduler.NewJob(
		gocron.DurationJob(durationJob),
		gocron.NewTask(
			func() {
				j.logger.Info("Starting job")
				ctx, cancel := context.WithTimeout(j.ctx, 30*time.Second)
				defer cancel()
				
				geoResp, err := j.geoClient.GetCoords(j.city)
				if err != nil {
					j.logger.Error(fmt.Sprintf(" %v: %s", f, err))
					return
				}
				j.logger.Info(fmt.Sprintf(" long: %f, lat: %f", geoResp.Longitude, geoResp.Latitude))
				
				openMeteoResp, err := j.openMeteoClient.GetTemperature(geoResp.Longitude, geoResp.Latitude)
				if err != nil {
					j.logger.Error(fmt.Sprintf("error get temperature %v: %s", f, err))
					return
				}
				var timestamp time.Time
				
				if openMeteoResp.Current.Time != "" {
					timestamp, err = time.Parse("2006-01-2T15:04", openMeteoResp.Current.Time)
					if err != nil {
						j.logger.Error(fmt.Sprintf("error parse timestamp: %v: %s", f, err))
						return
					}
				}
				
				reading := repository.Reading{
					Timestamp:   timestamp,
					Temperature: openMeteoResp.Current.Temperature2m,
				}
				
				if err = j.cache.Set(ctx, j.city, reading); err != nil {
					j.logger.Error(fmt.Sprintf("error set cache: %v:%s", f, err))
					return
				}
				j.logger.Info("done job")
			}),
	)
	if err != nil {
		return nil, err
	}
	return []gocron.Job{job}, nil
}
