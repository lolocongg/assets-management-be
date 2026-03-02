package email

import (
	"bytes"
	"embed"
	"html/template"
)

//go:embed files/overdue_template.html
var templateFS embed.FS

type Renderer struct {
	overdueTmpl *template.Template
}

func NewRenderer() (*Renderer, error) {
	tmpl, err := template.ParseFS(templateFS, "files/overdue_template.html")
	if err != nil {
		return nil, err
	}

	return &Renderer{
		overdueTmpl: tmpl,
	}, nil
}

func (r *Renderer) RenderOverdue(data OverdueEmailData) (string, error) {
	var buf bytes.Buffer

	if err := r.overdueTmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
