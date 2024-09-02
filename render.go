package atsemail

import (
	"bytes"
	"embed"
	"fmt"
	htemplate "html/template"
	"io"
	ttemplate "text/template"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/protobuf/proto"
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

func New[E proto.Message](name string) (r *Render[E], err error) {
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
	var htmBuf bytes.Buffer

	if err := r.text.ExecuteTemplate(txtw, r.name+".txt", data); err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	if err := r.html.ExecuteTemplate(&htmBuf, r.name+".html", data); err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}

	if err := r.ApplyTheme(&htmBuf, data); err != nil {
		return fmt.Errorf("failed to apply theme: %w", err)
	}

	if _, err := io.Copy(htmw, &htmBuf); err != nil {
		return fmt.Errorf("failed to write to output buffer: %w", err)
	}

	return nil
}

func (r *Render[E]) ApplyTheme(htmBuf *bytes.Buffer, _ E) error {
	doc, err := goquery.NewDocumentFromReader(htmBuf)
	if err != nil {
		return fmt.Errorf("failed to read into document: %w", err)
	}

	// @TODO apply theme variables
	// doc.Find(".sd-theme-button").Each(func(_ int, s *goquery.Selection) {
	// 	s.Each(func(_ int, s *goquery.Selection) {
	// 		style, _ := s.Attr("style")

	// 		style += `;background-color: red`
	// 		style += `;border-radius: 0`

	// 		s.SetAttr("style", style)
	// 	})
	// })

	// doc.Find(".sd-theme-container").Each(func(_ int, s *goquery.Selection) {
	// 	s.Each(func(_ int, s *goquery.Selection) {
	// 		style, _ := s.Attr("style")
	// 		style += `;border-radius: 0`

	// 		s.SetAttr("style", style)
	// 	})
	// })

	res, err := doc.Html()
	if err != nil {
		return fmt.Errorf("failed to turn document back into html: %w", err)
	}

	htmBuf.Reset()
	if _, err := htmBuf.WriteString(res); err != nil {
		return fmt.Errorf("failed to write html back into buffer: %w", err)
	}

	return nil
}
