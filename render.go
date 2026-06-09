package atsemail

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/bufbuild/protovalidate-go"
	emailsv1 "github.com/crewlinker/atsemail/emails/v1"
	"google.golang.org/protobuf/proto"
)

//go:embed exported/html/*.html
var htmlFiles embed.FS

//go:embed  exported/text/*.txt
var textFiles embed.FS

type Render[E EmailData] struct {
	name    string
	htmlSrc string
	txtSrc  string
}

const (
	leftDelim  = "{"
	rightDelim = "}"
)

func protoStringFields(msg proto.Message) map[string]string {
	v := reflect.ValueOf(msg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	m := make(map[string]string, t.NumField())
	for i := range t.NumField() {
		f := t.Field(i)
		fv := v.Field(i)
		if !f.IsExported() || fv.Kind() != reflect.String {
			continue
		}
		m[f.Name] = fv.String()
	}
	return m
}

type EmailData interface {
	proto.Message
	GetThemeOverwrites() *emailsv1.ThemeOverwrites
}

// ResolveVars replaces {variable_name} placeholders in s using vars.
func ResolveVars(s string, vars map[string]string) (string, error) {
	for k, v := range vars {
		s = strings.ReplaceAll(s, leftDelim+k+rightDelim, v)
	}
	return s, nil
}

func New[E EmailData](name string) (r *Render[E], err error) {
	r = &Render[E]{name: name}

	htmlContent, err := fs.ReadFile(htmlFiles, "exported/html/"+r.name+".html")
	if err != nil {
		return nil, fmt.Errorf("failed to read html: %w", err)
	}
	r.htmlSrc = string(htmlContent)

	txtContent, err := fs.ReadFile(textFiles, "exported/text/"+r.name+".txt")
	if err != nil {
		return nil, fmt.Errorf("failed to read text: %w", err)
	}
	r.txtSrc = string(txtContent)

	return r, nil
}

func (r *Render[E]) Render(val *protovalidate.Validator, txtw, htmw io.Writer, data E) error {
	var htmBuf bytes.Buffer

	if err := val.Validate(data); err != nil {
		return fmt.Errorf("invalid email data: %w", err)
	}

	vars := protoStringFields(data)

	txt, err := ResolveVars(r.txtSrc, vars)
	if err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}
	if _, err := io.WriteString(txtw, txt); err != nil {
		return fmt.Errorf("failed to write text: %w", err)
	}

	html, err := ResolveVars(r.htmlSrc, vars)
	if err != nil {
		return fmt.Errorf("failed to render html: %w", err)
	}
	htmBuf.WriteString(html)

	if theme := data.GetThemeOverwrites(); theme != nil {
		if err := ApplyTheme(&htmBuf, theme); err != nil {
			return fmt.Errorf("failed to apply theme: %w", err)
		}
	}

	if _, err := io.Copy(htmw, &htmBuf); err != nil {
		return fmt.Errorf("failed to write to output buffer: %w", err)
	}

	return nil
}

// ThemeOverwritesToCSS defines how we turn theme overwrite data into css styles.
func ThemeOverwritesToCSS(theme *emailsv1.ThemeOverwrites) (
	containerBorderRadius string,
	buttonBorderRadius string,
	buttonBackgroundColor string,
	buttonTextColor string,
	linkTextColor string,
) {
	switch theme.GetBorderRadius() {
	case emailsv1.BorderRadius_BORDER_RADIUS_UNSPECIFIED:
	case emailsv1.BorderRadius_BORDER_RADIUS_NONE:
		containerBorderRadius = `border-radius:0`
		buttonBorderRadius = `border-radius:0`
	case emailsv1.BorderRadius_BORDER_RADIUS_SMALL:
		containerBorderRadius = `border-radius:1%`
		buttonBorderRadius = `border-radius:3px`
	case emailsv1.BorderRadius_BORDER_RADIUS_MEDIUM:
		containerBorderRadius = `border-radius:3%`
		buttonBorderRadius = `border-radius:5px`
	case emailsv1.BorderRadius_BORDER_RADIUS_LARGE:
		containerBorderRadius = `border-radius:5%`
		buttonBorderRadius = `border-radius:10px`
	}

	if theme.GetButtonBackgroundColor() != nil {
		buttonBackgroundColor = fmt.Sprintf(`background-color: rgb(%d,%d,%d)`,
			theme.GetButtonBackgroundColor().GetRed(),
			theme.GetButtonBackgroundColor().GetGreen(),
			theme.GetButtonBackgroundColor().GetBlue())
	}

	if theme.GetButtonTextColor() != nil {
		buttonTextColor = fmt.Sprintf(`color: rgb(%d,%d,%d)`,
			theme.GetButtonTextColor().GetRed(),
			theme.GetButtonTextColor().GetGreen(),
			theme.GetButtonTextColor().GetBlue())
	}

	if theme.GetLinkTextColor() != nil {
		linkTextColor = fmt.Sprintf(`color: rgb(%d,%d,%d)`,
			theme.GetLinkTextColor().GetRed(),
			theme.GetLinkTextColor().GetGreen(),
			theme.GetLinkTextColor().GetBlue())
	}

	return //nolint:nakedret
}

func ApplyTheme(htmBuf *bytes.Buffer, theme *emailsv1.ThemeOverwrites) error {
	doc, err := goquery.NewDocumentFromReader(htmBuf)
	if err != nil {
		return fmt.Errorf("failed to read into document: %w", err)
	}

	borderRadius, buttonBorderRadius, buttonBackgroundColor, buttonTextColor, linkTextColor := ThemeOverwritesToCSS(theme)

	doc.Find("a:not(.sd-theme-button)").Each(func(_ int, s *goquery.Selection) {
		s.Each(func(_ int, s *goquery.Selection) {
			style, _ := s.Attr("style")
			style += `;` + linkTextColor

			s.SetAttr("style", style)
		})
	})

	doc.Find(".sd-theme-button").Each(func(_ int, s *goquery.Selection) {
		s.Each(func(_ int, s *goquery.Selection) {
			style, _ := s.Attr("style")

			style += `;` + buttonTextColor
			style += `;` + buttonBackgroundColor
			style += `;` + buttonBorderRadius

			s.SetAttr("style", style)
		})
	})

	doc.Find(".sd-theme-container").Each(func(_ int, s *goquery.Selection) {
		s.Each(func(_ int, s *goquery.Selection) {
			style, _ := s.Attr("style")
			style += `;` + borderRadius

			s.SetAttr("style", style)
		})
	})

	if theme.GetHeadingImage() != nil {
		doc.Find(".sd-theme-heading-image").Each(func(_ int, s *goquery.Selection) {
			s.Each(func(_ int, s *goquery.Selection) {
				s.SetAttr("src", fmt.Sprintf(`data:%s;base64,%s`,
					theme.GetHeadingImage().GetContentType(),
					base64.StdEncoding.EncodeToString(theme.GetHeadingImage().GetData())))
			})
		})
	}

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
