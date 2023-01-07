package storage

import "context"

type DeleteOneStreamReqDTO struct {
	Ctx       context.Context
	StreamId string
}

type CreateOneStreamReqDTO struct {
	Ctx       context.Context
	StreamId string
}
