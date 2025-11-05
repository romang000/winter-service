package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	"github.com/romang000/winter-service/internal/client/http/geocoding"
	"github.com/romang000/winter-service/internal/client/http/open_meteo"
	"github.com/romang000/winter-service/internal/config"
	"github.com/romang000/winter-service/internal/db"
	"github.com/romang000/winter-service/internal/handler"
	"github.com/romang000/winter-service/internal/repository"
	"github.com/romang000/winter-service/internal/service"
	"log"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	city = "moscow"
)

type Reading struct {
	Timestamp   time.Time
	Temperature float64
}

type Storage struct {
	data map[string][]Reading
	mu   sync.RWMutex
}

func main() {
	configPath := flag.String("config", "./config", "Path to config file")
	flag.Parse()
	
	var logger *slog.Logger
	
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}
	
	newDb, dsn, err := db.NewDatabaseConnection(cfg.Database)
	
	if err != nil {
		panic(err)
	}
	
	repo, err := repository.New(newDb)
	if err != nil {
		logger.Error(fmt.Sprintf("Error creating repository: %v", err))
	}
	
	logger.Info(fmt.Sprintf("Using database %s", dsn))
	
	srv := service.NewWeatherService(repo, logger)
	
	router := gin.Default()
	
	handler.RegisterRoutes(router, logger, srv)
	
	logger.Info(fmt.Sprintf("starting app on %s:%s", cfg.Service.Address, cfg.Service.ServerPort))
	if err = router.Run(fmt.Sprintf("%s:%s", cfg.Service.Address, cfg.Service.ServerPort)); err != nil {
		panic(err)
	}
	
	//r := chi.NewRouter()
	//r.Use(middleware.Logger)
	//
	//storage := &Storage{
	//	data: make(map[string][]Reading),
	//}
	//
	//r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
	//	cityName := chi.URLParam(r, "city")
	//
	//	fmt.Println(cityName)
	//
	//	storage.mu.RLock()
	//	defer storage.mu.RUnlock()
	//
	//	reading, ok := storage.data[cityName]
	//	if !ok {
	//		w.WriteHeader(http.StatusNotFound)
	//		w.Write([]byte("not found"))
	//		return
	//	}
	//
	//	raw, err := json.Marshal(reading)
	//	if err != nil {
	//		log.Println("json marshal error:", err)
	//	}
	//
	//	_, err = w.Write(raw)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//
	//})
	//
	//s, err := gocron.NewScheduler()
	//if err != nil {
	//	panic(err)
	//}
	//
	//jobs, err := initJobs(s, storage)
	//if err != nil {
	//	panic(err)
	//}
	//
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//
	//go func() {
	//	defer wg.Done()
	//
	//	fmt.Println("starting server on port " + httpPort)
	//	err = http.ListenAndServe(httpPort, r)
	//	if err != nil {
	//		panic(err)
	//	}
	//}()
	//
	//go func() {
	//	defer wg.Done()
	//
	//	fmt.Printf("starting job: %v\n", jobs[0].ID())
	//	s.Start()
	//}()
	//
	//wg.Wait()
}

func initJobs(scheduler gocron.Scheduler, storage *Storage) ([]gocron.Job, error) {
	
	httpClient := &http.Client{Timeout: 10 * time.Second}
	
	geoClient := geocoding.NewClient(httpClient)
	
	openMeteoClient := open_meteo.NewClient(httpClient)
	
	j, err := scheduler.NewJob(
		gocron.DurationJob(
			5*time.Second,
		),
		gocron.NewTask(
			func() {
				geoResp, err := geoClient.GetCoords(city)
				if err != nil {
					log.Println(err)
					return
				}
				
				openMeteoResponse, err := openMeteoClient.GetTemperature(geoResp.Latitude, geoResp.Longitude)
				if err != nil {
					log.Println(err)
					return
				}
				storage.mu.Lock()
				defer storage.mu.Unlock()
				
				timestamp, err := time.Parse("2006-01-2T15:04", openMeteoResponse.Current.Time)
				
				if err != nil {
					log.Println(err)
					return
				}
				
				storage.data[city] = append(storage.data[city], Reading{
					timestamp,
					openMeteoResponse.Current.Temperature2m,
				})
				
				fmt.Printf("%v update data for city: %s", time.Now(), city)
			},
		),
	)
	
	if err != nil {
		return nil, err
	}
	
	return []gocron.Job{j}, nil
}
