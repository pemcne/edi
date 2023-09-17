module github.com/pemcne/edi

go 1.20

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/go-joe/joe v0.11.0
	github.com/go-joe/slack-adapter/v2 v2.2.0
	github.com/lithammer/fuzzysearch v1.1.5
	github.com/notnil/chess v1.9.0
	github.com/pemcne/firestore-memory v0.0.0-20221218153031-323b330c3e9d
	github.com/robfig/cron/v3 v3.0.1
	github.com/slack-go/slack v0.12.0
)

require (
	cloud.google.com/go v0.107.0 // indirect
	cloud.google.com/go/compute v1.13.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.2 // indirect
	cloud.google.com/go/firestore v1.9.0 // indirect
	cloud.google.com/go/longrunning v0.3.0 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/oauth2 v0.3.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.105.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221207170731-23e4bf6bdc37 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace (
	github.com/go-joe/joe v0.11.0 => github.com/pemcne/joe v0.11.1-0.20220403212347-b9408549999d
	github.com/go-joe/slack-adapter/v2 v2.2.0 => github.com/pemcne/slack-adapter/v2 v2.2.1-0.20220618164314-02241dc560f7
)
