module github.com/pemcne/edi

go 1.17

require (
	github.com/go-joe/file-memory v1.0.0
	github.com/go-joe/joe v0.11.0
	go.uber.org/zap v1.21.0
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.1 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace (
	github.com/go-joe/joe v0.11.0 => github.com/pemcne/joe v0.11.1-0.20220403212347-b9408549999d
	github.com/go-joe/slack-adapter/v2 v2.2.0 => github.com/pemcne/slack-adapter/v2 v2.2.1-0.20220403212237-f6456ddb0c16
)
