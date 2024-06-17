# restapi
RestAPI wrapper with generic GET/POST handlers and middleware logger. Based on Chi router (https://github.com/go-chi/chi).

    log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // create new server with middlware logger enabled
	rest := restapi.New(log.With("module", "RestAPI"), true)

    // add new generic GET endpoint (route)
	rest.Router.Get("/version", rest.Generic_GET_handler("API v.1"))

    // start server
	rest.StartAt(":8080")

	// as rest.ServerAt() is non blocking, we need to implement a wait loop:
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	rest.Shutdown(3 * time.Second)


Now you can curl http://localhost:8080/version and the output will be:

    {"status":"Ok","data":"API v.1"}

While stdout will provide you with structured logging about request:

```
{"time":"2024-06-17T05:43:57.060262726Z","level":"INFO","msg":"starting server","module":"RestAPI","settings":{"address":":8080","mwlogger":true}}
{"time":"2024-06-17T05:44:01.950194811Z","level":"INFO","msg":"new request","module":"RestAPI","handler":"generic","request":{"id":"71665260ed60/H7dkBWiVKA-000001","route":"/version","method":"GET","remote_addr":"127.0.0.1:36878","user_agent":"curl/7.88.1"}}
{"time":"2024-06-17T05:44:01.95231069Z","level":"INFO","msg":"request completed","module":"RestAPI","middleware":"logger","request":{"id":"71665260ed60/H7dkBWiVKA-000001","endpoint":"/version","method":"GET","remote_addr":"127.0.0.1:36878","user_agent":"curl/7.88.1"},"stats":{"status":200,"bytes":33,"duration":"2.125021ms"}}
```

Module provides generic POST handler with channel-based notification on file upload events and automatic UUID assignment for uploads:

    // create channel and add POST route with destination folder, 
	// naming mask for uploaded files and channel for tracking uploads:

    uploadsChan := make(chan restapi.TFileUpload)
	rest.Router.Post("/upload", rest.Generic_POST_File("payload", "./", "", uploadsChan))

	// start tracking uploads via channel notifications:

	go func(ch chan restapi.TFileUpload) {
		for newFile := range ch {
			fmt.Println("File uploaded:", newFile)
		}
	}(uploadsChan)

Now you can curl curl -F "payload=@./example/file.txt" http://localhost:8080/upload and responsewill be like:

	{"status":"Ok","data":"cdee67b2-3566-46c2-9c83-702b298cb548"}

Handler automatically assignes UUID for the file uploaded. The output on stdout will be like:

```
{"time":"2024-06-17T06:22:32.469050045Z","level":"INFO","msg":"new upload","module":"RestAPI","handler":"generic","request":{"id":"71665260ed60/uGrW4maTks-000001","route":"/upload","method":"POST","remote_addr":"127.0.0.1:49810","user_agent":"curl/7.88.1"},"file":{"UUID":"cdee67b2-3566-46c2-9c83-702b298cb548","Path":"./upload_2458019802","Filename":"file.txt","Size":21}}

File uploaded: {cdee67b2-3566-46c2-9c83-702b298cb548 ./upload_2458019802 0xc000096360}

{"time":"2024-06-17T06:22:32.472910764Z","level":"INFO","msg":"request completed","module":"RestAPI","middleware":"logger","request":{"id":"71665260ed60/uGrW4maTks-000001","endpoint":"/upload","method":"POST","remote_addr":"127.0.0.1:49810","user_agent":"curl/7.88.1"},"stats":{"status":200,"bytes":62,"duration":"7.207908ms"}}
```