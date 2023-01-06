package errors

import "errors"

// This error refers to the moment when no one is connected to the sse server for a particular stream. 
// And the last user, when disconnected, causes the process of deleting this stream
var ErrNoVisitors = errors.New("no visitors")