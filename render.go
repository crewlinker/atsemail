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

type EmailTemplate interface {
	Name() string
}

type Render[E EmailTemplate] struct {
	data E
	html *htemplate.Template
	text *ttemplate.Template
}

func New[E EmailTemplate](data E) (r *Render[E], err error) {
	r = &Render[E]{data: data}
	r.html, err = htemplate.ParseFS(htmlFiles, "exported/html/"+data.Name()+".html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	r.text, err = ttemplate.ParseFS(textFiles, "exported/text/"+data.Name()+".txt")
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
