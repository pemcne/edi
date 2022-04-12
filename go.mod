module github.com/pemcne/edi

go 1.17

require (
	github.com/go-joe/file-memory v1.0.0
	github.com/go-joe/joe v0.11.0
	github.com/go-joe/slack-adapter/v2 v2.2.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/slack-go/slack v0.10.2 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace (
	github.com/go-joe/joe v0.11.0 => github.com/pemcne/joe v0.11.1-0.20220403212347-b9408549999d
	github.com/go-joe/slack-adapter/v2 v2.2.0 => github.com/pemcne/slack-adapter/v2 v2.2.1-0.20220403212237-f6456ddb0c16
)
