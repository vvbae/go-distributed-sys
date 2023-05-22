package log_v1

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ErrOffsetOutOfRange struct {
	Offset uint64
}

func (e ErrOffsetOutOfRange) GRPCStatus() *status.Status {
	// create a new not-found status
	st := status.New(404,
		fmt.Sprintf("offset out of range: %d", e.Offset),
	)
	msg := fmt.Sprintf(
		"The requested offset is outside the log's range: %d", e.Offset,
	)
	// create the localized message
	d := &errdetails.LocalizedMessage{Locale: "en-US",
		Message: msg,
	}

	// attach the message to status
	std, err := st.WithDetails(d)
	if err != nil {
		return st
	}
	return std
}

func (e ErrOffsetOutOfRange) Error() string {
	return e.GRPCStatus().Err().Error()
}
