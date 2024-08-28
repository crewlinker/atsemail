package atsemail

import (
	"embed"
	"fmt"
	htemplate "html/template"
	"io"
	ttemplate "text/template"
)

//go:embed exported/html/*.html
var htmlFiles embed.FS

//go:embed  exported/text/*.txt
var textFiles embed.FS

type Render[E any] struct {
	name string
	html *htemplate.Template
	text *ttemplate.Template
}

const (
	leftDelim  = "$"
	rightDelim = "$"
	opts       = "missingkey=error"
)

func New[E any](name string) (r *Render[E], err error) {
	r = &Render[E]{name: name}

	r.html, err = htemplate.New("").
		Delims(leftDelim, rightDelim).
		Option(opts).
		ParseFS(htmlFiles, "exported/html/"+r.name+".html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	r.text, err = ttemplate.New("").
		Delims(leftDelim, rightDelim).
		Option(opts).
		ParseFS(textFiles, "exported/text/"+r.name+".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to parse text: %w", err)
	}

	return r, nil
}

func (r *Render[E]) Render(txtw, htmw io.Writer, data E) error {
	if err := r.text.ExecuteTemplate(txtw, r.name+".txt", data); err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	if err := r.html.ExecuteTemplate(htmw, r.name+".html", data); err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}

	return nil
}
