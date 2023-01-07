package entity

import (
	"time"
)

// Consistency. Эти данные пересылаются извне, из другого микросервиса
type Consistency struct {
	RequestId string 
	StreamId  string
	Timestamp time.Time
	Data      []byte
}
