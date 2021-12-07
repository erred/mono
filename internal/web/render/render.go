package render

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Render is a wrapper for RenderBytes to work with a Reader.
func Render(o *Options, w io.Writer, r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	return RenderBytes(o, w, b)
}

// RenderBytes renders the markdown in body
func RenderBytes(o *Options, w io.Writer, body []byte) error {
	o.init()

	if len(body) == 0 && len(o.Data.Main) == 0 {
		return fmt.Errorf("no body")
	}

	if len(body) != 0 {
		b, err := o.Data.fromHeader(body)
		if err != nil {
			return fmt.Errorf("extract header: %s", err)
		}

		var buf strings.Builder
		err = o.Markdown.Convert(b, &buf)
		if err != nil {
			return fmt.Errorf("render markdown: %w", err)
		}
		o.Data.Main = buf.String()
	}

	if o.MinifySkip {
		err := o.Template.Execute(w, o.Data)
		if err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
	} else {
		var buf bytes.Buffer
		err := o.Template.Execute(&buf, o.Data)
		if err != nil {
			return fmt.Errorf("execute template: %w", err)
		}
		err = o.Minify.Minify("text/html", w, &buf)
		if err != nil {
			return fmt.Errorf("minify: %w", err)
		}
	}

	return nil
}
