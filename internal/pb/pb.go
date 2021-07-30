package progressbar

import (
	"github.com/cheggaaa/pb/v3"
	"os"
)

func NewProcessBarReader(size int64) *pb.ProgressBar {
	// start new bar
	bar := pb.Full.Start64(size)

	// force set io.Writer, by default it's os.Stderr
	bar.SetWriter(os.Stdout)

	// bar will format numbers as bytes (B, KiB, MiB, etc)
	bar.Set(pb.Bytes, true)

	return bar
}
