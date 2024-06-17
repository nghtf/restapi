package main

import (
	"log/slog"
	"os"
	"os/signal"
	"restapi"
	"time"
)

func main() {

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	rest := restapi.New(log.With("module", "RestAPI"), true)

	// ROUTE: GET /version (generic handler). Example: curl http://localhost:8080/version
	rest.Router.Get("/version", rest.Generic_GET_handler("API v.1"))

	// ROUTE: POST /upload (generic handler). Example: curl -F "payload=@./example/file.txt" http://localhost:8080/upload
	uploadsChan := make(chan restapi.TFileUpload)
	rest.Router.Post("/upload", rest.Generic_POST_File("payload", "./", "", uploadsChan))
	go manager(log, uploadsChan)

	rest.StartAt(":8080")

	// wait for the signal to quit and shutdown gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	rest.Shutdown(3 * time.Second)

}

func manager(log *slog.Logger, ch chan restapi.TFileUpload) {
	log = log.With(slog.String("module", "manager"))
	for newFile := range ch {
		log.Info("new upload",
			slog.Group("file",
				slog.String("UUID", newFile.ID),
				slog.String("Path", newFile.Path),
				slog.String("Filename", newFile.Header.Filename),
				slog.Int64("Size", newFile.Header.Size),
			),
		)
	}
}
