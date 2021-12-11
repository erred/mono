package render

import (
	"bytes"
	"io"
)

// CompactBytes is a convenience function to render a compact page to a writer
func Compact(dst io.Writer, title, desc, canonicalURL string, body []byte) error {
	o := &Options{
		Data: PageData{
			URLCanonical: canonicalURL,
			Compact:      true,
			Title:        title,
			Description:  desc,
		},
	}

	return Render(o, dst, bytes.NewReader(body))
}

// CompactBytes is a convenience function to render a compact page to a byte slice
func CompactBytes(title, desc, canonicalURL string, body []byte) ([]byte, error) {
	var buf bytes.Buffer
	err := Compact(&buf, title, desc, canonicalURL, body)
	return buf.Bytes(), err
}
