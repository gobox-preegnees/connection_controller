package entity

// Consistency. This describes the actions that need to be performed
// by the client in order to bring his directory into a consistent form.
// This structure is sent by the server every time the consistent folder on the server changes.
// If the client that received this structure has the same "Client" and "Update=true",
// then it does nothing, and if after "Client" is different from the real client,
// then the client accepts the changes
type Consistency struct {
	// Installed on the server from field of client
	ExternalId int `json:"external_id"`
	// Installed on the server
	InternalId int `json:"internal_id"`
	// Installed on the server
	Timestamp int `json:"timestamp" validate:"required"`

	Update bool `json:"update" validate:"required"`

	Username string `json:"username" validate:"required"`
	Folder   string `json:"folder" validate:"required"`
	Client   string `json:"client" validate:"required"`

	NeedToRemove []struct {
		File `json:"file"`
	}
	NeedToBlock []struct {
		File `json:"file"`
	}
	NeedToUpload []struct {
		File `json:"file"`
		// This field stores a jwt token,
		// with which the client will later contact the server to send or download a file.
		// Fields: folder, username, client, file_name, size_file, mod_time, hash_sum, action:upload|download|download
		ActionToken string `json:"action_token"`
	}
	NeedToDownload []struct {
		File `json:"file"`
		// Fields: Virtual_file_name, folder, username, client, file_name, size_file, mod_time, hash_sum, action:upload|download|download
		ActionToken string `json:"action_token"`
	}
}
