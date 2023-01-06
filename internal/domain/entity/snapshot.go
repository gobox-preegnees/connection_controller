package entity

type File struct {
	FileName string `json:"file_name" validate:"required"`
	HashSum  string `json:"hash_sum"`
	SizeFile int    `json:"size_file"`
	ModTime  int    `json:"mod_time" validate:"required"`
}

// Snapshot. This structure must be received by the server from the client.
// This structure is received by the server when the client has a change.
type Snapshot struct {
	// Installed on the server
	RequestId string `json:"request_id" validate:"required"`
	// user = username + folder + client
	StreamId string `json:"stream_id" validate:"required"`
	// Installed on the server
	Timestamp int64 `json:"timestamp" validate:"required"`

	// Installed on the client. Has the values 100 - at start, 300 - at update, 400 - at error
	Action string `json:"on_action" validate:"required,eq=100|eq=300|eq=400"`

	// Username of client
	Username string `json:"username" validate:"required"`
	// Virtual folder, which is synchronizing
	Folder string `json:"folder" validate:"required"`
	// Name of client, which he had set in website
	Client string `json:"client" validate:"required"`
	// File if client directory
	Files []struct {
		File `json:"file" validate:"required"`
	}
}
