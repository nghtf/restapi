package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"restapi"
	"time"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rest := restapi.New(log.With("module", "RestAPI"), true)

	// Route (GET, generic handler): /version
	// Example: curl http://localhost:8080/version
	rest.Router.Get("/version", rest.Handler.GET("API v.1"))

	// Route (POST, generic handler): /upload
	uploadsChan := make(chan restapi.TFileUpload)
	go func(ch chan restapi.TFileUpload) {
		for newFile := range ch {
			fmt.Println("File uploaded:", newFile)
		}
	}(uploadsChan)
	// Example: curl -F "payload=@./example/file.txt" http://localhost:8080/upload
	rest.Router.Post("/upload", rest.Handler.POST("payload", "./", restapi.DEFAULT_UPLOAD_PATTERN, uploadsChan))

	// staring the server at localhost:8080
	rest.StartAt(":8080")

	// wait for the signal to quit and shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	rest.Shutdown(3 * time.Second)

}
