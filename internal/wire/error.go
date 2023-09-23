package wire

import "fmt"

type DecodError struct {
	Format string
	Reason string
	Err    error
}

func (de *DecodError) Error() string {
	return fmt.Sprintf("failed to decode message. wire format: %s. reason: %s, error: %v", de.Format, de.Reason, de.Err)
}

func ProtoDecodeError(baseErr error, reason string) *DecodError {
	return &DecodError{Format: "proto", Err: baseErr, Reason: reason}
}
