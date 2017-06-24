package util

import (
	"fmt"
	"io"
	"text/tabwriter"
)

const (
	tabwriterMinWidth = 10
	tabwriterWidth    = 4
	tabwriterPadding  = 3
	tabwriterPadChar  = ' '
	tabwriterFlags    = 0
)

type TabWriter struct {
	writer *tabwriter.Writer
}

func NewTabWriter(output io.Writer) *TabWriter {
	writer := tabwriter.NewWriter(output, tabwriterMinWidth, tabwriterWidth, tabwriterPadding, tabwriterPadChar, tabwriterFlags)
	return &TabWriter{writer}
}

func (w *TabWriter) Print(str string) {
	fmt.Fprint(w.writer, str)
}

func (w *TabWriter) Nl() {
	w.Print("\n")
}

func (w *TabWriter) Tab() {
	w.Print("\t")
}

func (w *TabWriter) Flush() {
	w.writer.Flush()
}
