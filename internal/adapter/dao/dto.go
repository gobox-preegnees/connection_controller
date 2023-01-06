package storage

import "context"

type CreateOneSnapshotReqDTO struct {
	Snapshot []byte
}

type DeleteOneOwnerReqDTO struct {
	Ctx       context.Context
	Usernamme string
	Folder    string
}

type CreateOneOwnerReqDTO struct {
	Ctx       context.Context
	Usernamme string
	Folder    string
}
