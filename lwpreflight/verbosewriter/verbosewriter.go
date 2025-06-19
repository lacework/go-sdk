package verbosewriter

import (
	"fmt"
	"io"

	"github.com/abiosoft/colima/util/terminal"
)

type WriteCloser interface {
	Write(s string)
	Close()
}

type VerboseWriter struct {
	vw io.WriteCloser
}

func New() *VerboseWriter {
	return &VerboseWriter{
		vw: terminal.NewVerboseWriter(10),
	}
}

func (v *VerboseWriter) Write(s string) {
	v.vw.Write(fmt.Appendf(nil, "%s\n", s))
}

func (v *VerboseWriter) Close() {
	v.vw.Close()
}
