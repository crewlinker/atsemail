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

type EmailData interface {
	Name() string
}

type Render[E EmailData] struct {
	data E
	html *htemplate.Template
	text *ttemplate.Template
}

const (
	leftDelim  = "$"
	rightDelim = "$"
	opts       = "missingkey=error"
)

func New[E EmailData](data E) (r *Render[E], err error) {
	r = &Render[E]{data: data}

	r.html, err = htemplate.New("").
		Delims(leftDelim, rightDelim).
		Option(opts).
		ParseFS(htmlFiles, "exported/html/"+data.Name()+".html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	r.text, err = ttemplate.New("").
		Delims(leftDelim, rightDelim).
		Option(opts).
		ParseFS(textFiles, "exported/text/"+data.Name()+".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to parse text: %w", err)
	}

	return r, nil
}

func (r *Render[E]) Render(txtw, htmw io.Writer) error {
	if err := r.text.ExecuteTemplate(txtw, r.data.Name()+".txt", r.data); err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	if err := r.html.ExecuteTemplate(htmw, r.data.Name()+".html", r.data); err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}

	return nil
}
