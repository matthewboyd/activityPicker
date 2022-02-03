package main

import (
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/matthewboyd/activities"
	"github.com/sony/gobreaker"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	//HostAddress = os.Getenv("HOST_ADDRESS")
	HostAddress   = ":8080"
	hostname      = os.Getenv("POSTGRES_URL")
	username      = os.Getenv("POSTGRES_USER")
	password      = os.Getenv("POSTGRES_PASSWORD")
	database_name = os.Getenv("POSTGRES_DB")
)

const (
	hostPort = 5432
)

var cb *gobreaker.CircuitBreaker

func init() {
	var st gobreaker.Settings
	st.Name = "HTTP GET"
	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return counts.Requests >= 5 && failureRatio >= 0.6
	}
	cb = gobreaker.NewCircuitBreaker(st)
}

func main() {
	logger := log.New(os.Stdout, "matt ", log.LstdFlags|log.Lshortfile)
	pgConString := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostPort, hostname, username, password, database_name)
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})
	database, err := sql.Open("postgres", pgConString)
	database.SetMaxOpenConns(5)
	defer database.Close()
	mux := http.NewServeMux()

	logger.Println("starting now")
	handler := activities.Handler{
		Logger: *logger,
		Db:     *database,
		Redis:  *rdb,
	}
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/{action}", pprof.Index)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/sunny", handler.SunnyEndpoint())
	mux.HandleFunc("/allWeather", handler.NotSunnyEndpoint())
	srv := newServer(mux, HostAddress)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func newServer(mux *http.ServeMux, serverAddress string) *http.Server {
	return &http.Server{

		Addr:         serverAddress,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
