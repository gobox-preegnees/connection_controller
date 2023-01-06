package messagebroker

import "time"

type PublishSnapshotReqDTO struct {
	RequestId string
	StreamId  string
	Timestamp time.Time
	Data      []byte
}
