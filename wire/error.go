package wire

import "fmt"

type DecodError struct {
	Format string
	Reason error
}

func (de *DecodError) Error() string {
	return fmt.Sprintf("failed to decode message. wire format: %s. reason: %s", de.Format, de.Reason)
}
