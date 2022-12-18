module github.com/pemcne/edi

go 1.17

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/go-joe/file-memory v1.0.0
	github.com/go-joe/joe v0.11.0
	github.com/go-joe/slack-adapter/v2 v2.2.0
	github.com/lithammer/fuzzysearch v1.1.5
	github.com/notnil/chess v1.8.0
	github.com/robfig/cron/v3 v3.0.1
	go.uber.org/zap v1.21.0
)

require (
	cloud.google.com/go v0.97.0 // indirect
	cloud.google.com/go/firestore v1.6.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/googleapis/gax-go/v2 v2.1.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pemcne/firestore-memory v0.0.0-20220813142323-efff3a4341ca // indirect
	github.com/slack-go/slack v0.11.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220617184016-355a448f1bc9 // indirect
	golang.org/x/oauth2 v0.0.0-20211005180243-6b3c2da341f1 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/api v0.59.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20211028162531-8db9c33dc351 // indirect
	google.golang.org/grpc v1.40.0 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

replace (
	github.com/go-joe/joe v0.11.0 => github.com/pemcne/joe v0.11.1-0.20220403212347-b9408549999d
	github.com/go-joe/slack-adapter/v2 v2.2.0 => github.com/pemcne/slack-adapter/v2 v2.2.1-0.20220618164314-02241dc560f7
)
