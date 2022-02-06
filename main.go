package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/matthewboyd/activities"
	"log"
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	//HostAddress = os.Getenv("HOST_ADDRESS")
	HostAddress  = ":8080"
	hostname     = os.Getenv("POSTGRES_URL")
	username     = os.Getenv("POSTGRES_USER")
	password     = os.Getenv("POSTGRES_PASSWORD")
	databaseName = os.Getenv("POSTGRES_DB")
)

const (
	hostPort = 5432
)

func main() {
	logger := log.New(os.Stdout, "matt ", log.LstdFlags|log.Lshortfile)
	pgConString := fmt.Sprintf("port=%d host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		hostPort, hostname, username, password, databaseName)
	redisAddress := fmt.Sprintf("%s:6379", os.Getenv("REDIS_URL"))
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	database, err := pgxpool.Connect(context.Background(), pgConString)
	if err != nil {
		log.Fatalf("Failed to connect to postgres db", err)
	}
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

		Addr:              serverAddress,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
}
