module github.com/pemcne/edi

go 1.17

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/go-joe/file-memory v1.0.0
	github.com/go-joe/joe v0.11.0
	github.com/go-joe/slack-adapter/v2 v2.2.0
	github.com/lithammer/fuzzysearch v1.1.4
	github.com/robfig/cron/v3 v3.0.1
	go.uber.org/zap v1.21.0
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/slack-go/slack v0.10.3 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace (
	github.com/go-joe/joe v0.11.0 => github.com/pemcne/joe v0.11.1-0.20220403212347-b9408549999d
	github.com/go-joe/slack-adapter/v2 v2.2.0 => github.com/pemcne/slack-adapter/v2 v2.2.1-0.20220403212237-f6456ddb0c16
)
