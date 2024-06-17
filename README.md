# restapi
RestAPI wrapper for Chi router with generic GET/POST handlers and middleware logger (slog-based). Helps to simplify http API implementation.

### Generic GET handler

Boilerplate to run a server with GET route:


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


Now you can `curl http://localhost:8080/version` and the output will be:

    {"status":"Ok","data":"API v.1"}

While stdout will provide you with structured logging about request:

```
{"time":"2024-06-17T05:43:57.060262726Z","level":"INFO","msg":"starting server","module":"RestAPI","settings":{"address":":8080","mwlogger":true}}
{"time":"2024-06-17T05:44:01.950194811Z","level":"INFO","msg":"new request","module":"RestAPI","handler":"generic","request":{"id":"71665260ed60/H7dkBWiVKA-000001","route":"/version","method":"GET","remote_addr":"127.0.0.1:36878","user_agent":"curl/7.88.1"}}
{"time":"2024-06-17T05:44:01.95231069Z","level":"INFO","msg":"request completed","module":"RestAPI","middleware":"logger","request":{"id":"71665260ed60/H7dkBWiVKA-000001","endpoint":"/version","method":"GET","remote_addr":"127.0.0.1:36878","user_agent":"curl/7.88.1"},"stats":{"status":200,"bytes":33,"duration":"2.125021ms"}}
```

You can supply rest.Generic_GET_handler(data interface{}) with any sort of the data that can be marshalled to the client.

### Generic POST handler

Module provides generic POST handler with channel-based notification on file upload events and automatic UUID assignment for files uploaded:

    // create channel and add POST route with destination folder, 
	// naming mask for uploaded files and channel for tracking uploads:

    uploadsChan := make(chan restapi.TFileUpload)
	rest.Router.Post("/upload", rest.Generic_POST_File("payload", "./upload/dir", "upload_*", uploadsChan))

	// start tracking uploads via channel notifications:

	go func(ch chan restapi.TFileUpload) {
		for newFile := range ch {
			fmt.Println("File uploaded:", newFile)
		}
	}(uploadsChan)

Now you can `curl -F "payload=@./example/file.txt" http://localhost:8080/upload` and response will contain new UUID for the file uploaded:

	{"status":"Ok","data":"cdee67b2-3566-46c2-9c83-702b298cb548"}

The output on stdout provides detailed logging for the event:

```
{"time":"2024-06-17T06:22:32.469050045Z","level":"INFO","msg":"new upload","module":"RestAPI","handler":"generic","request":{"id":"71665260ed60/uGrW4maTks-000001","route":"/upload","method":"POST","remote_addr":"127.0.0.1:49810","user_agent":"curl/7.88.1"},"file":{"UUID":"cdee67b2-3566-46c2-9c83-702b298cb548","Path":"./upload/dir/upload_2458019802","Filename":"file.txt","Size":21}}

File uploaded: {cdee67b2-3566-46c2-9c83-702b298cb548 ./upload/dir/upload_2458019802 0xc000096360}

{"time":"2024-06-17T06:22:32.472910764Z","level":"INFO","msg":"request completed","module":"RestAPI","middleware":"logger","request":{"id":"71665260ed60/uGrW4maTks-000001","endpoint":"/upload","method":"POST","remote_addr":"127.0.0.1:49810","user_agent":"curl/7.88.1"},"stats":{"status":200,"bytes":62,"duration":"7.207908ms"}}
```

Naming convention for uploads follows https://pkg.go.dev/os#CreateTemp. Original file name is stored in restapi.TFileUpload.