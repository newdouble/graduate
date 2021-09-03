package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/rs/cors"
	"graduate/src/config"
	"graduate/src/controller/fetchactdetailctr"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Version string
	Date string
)

var (
	showVersion bool
)

func main() {
	flag.Parse()
	if showVersion {
		fmt.Println("build by commited:", Version, "date:", Date)
		return
	}
	initServer()

	s := startServer()

	idleConnsClosed := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			fmt.Fprintln(os.Stderr, "server shutdown failed: ", err)
		}
		close(idleConnsClosed)
	}()
	fmt.Println("server starting")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "server start failed: ", err)
		os.Exit(1)
	}
	<-idleConnsClosed
}

func startServer() *http.Server {
	serveMux := http.NewServeMux()
	serveMux.Handle(config.Get().LocationConf.FetchActivityDetailPath, &fetchactdetailctr.FetchActDetailController{})
	serveMux.HandleFunc("/debug/pprof/", pprof.Index)
	serveMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	serveMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	serveMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	serveMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	readTimeout, err := time.ParseDuration(config.Get().ServerConf.ReadTimeout)
	if err != nil {
		readTimeout *= 3000 * time.Millisecond
	}
	writeTimeout, err := time.ParseDuration(config.Get().ServerConf.WriteTimeout)
	if err != nil {
		writeTimeout *= 3000 * time.Millisecond
	}
	cross := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost},
		AllowedOrigins:     []string{"*"},
		OptionsPassthrough: false,
		AllowCredentials:   false,
		Debug:              false,
	})
	handler := cross.Handler(serveMux)
	s := &http.Server{
		Addr: ":" + config.Get().ServerConf.Port,
		Handler: handler,
		ReadTimeout: readTimeout,
		WriteTimeout: writeTimeout,
		MaxHeaderBytes: 1 << 16,
	}
	return s
}

func initServer() {
	config.Parse("")
	os.Setenv("LUA_PATH,", "./lua/?.lua;;")
	//myredis.InitRedis(config.Get().RedisConf.RedisFilePath)
}

func init() {
	flag.BoolVar(&showVersion, "version", false, "get git commit id")
}