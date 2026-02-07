package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func newRouter() *httprouter.Router {
	mux := httprouter.New()
	/*  */
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	ytChannelID := os.Getenv("YOUTUBE_CHANNEL_ID")
	if ytApiKey == "" {
		log.Fatalf("youtube API key no provided")
	}
	if ytChannelID == "" {
		log.Fatalf("youtube channel ID no provided")
	}

	mux.GET("/youtube/channel/stats", getChannelStats(ytApiKey, ytChannelID))

	return mux
}

func getChannelStatsFirst() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("response!"))
	}

}

func main() {

	srv := &http.Server{
		Addr:    ":10101",
		Handler: newRouter(),
	}

	idleConsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint
		log.Println("service interrupt received")

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("http server shutdown error: %v", err)

		}
		log.Println("shutdown complete")

		close(idleConsClosed)
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("fatal http server failed to start: %v", err)
		}
	}

	<-idleConsClosed

	log.Println("Service Stop")
}
