package storage

import "context"

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
