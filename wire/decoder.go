package wire

import "github.com/thediveo/enumflag"

type Format enumflag.Flag

const (
	PlainText Format = iota
	Proto
)

var FormatIDs = map[Format][]string{
	PlainText: {"plaintext"},
	Proto:     {"proto"},
}

type Decoder interface {
	Decode([]byte) (string, error)
}
