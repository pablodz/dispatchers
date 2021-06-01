package exporter

import (
	"fmt"
	"html/template"
	"os"

	"github.com/CyCoreSystems/dispatchers/v2/sets"
)

// FileExporter is a dispatchers.Exporter which exports formatted dispatcher set data to a file.
type FileExporter struct {
	filename string

	tmpl *template.Template
}

// DefaultFileTemplate is the default file exporter template, suitable for use by the kamailio dispatchers module as a flat file.
var DefaultFileTemplate = `
# Dispatcher sets.
# WARNING: THIS FILE IS AUTOMATICALLY GENERATED.

{{ range . }}
# Dispatcher set {{ .ID }}
{{ _, $ep := range .Endpoints }}
{{ .ID }} sip:{{ $ep }}
{{ end }}
{{ end }}
`

// Export implements dispatchers.Exporter
func (e *FileExporter) Export(sets []*sets.State) error {
	f, err := os.Open(e.filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", e.filename, err)
	}
	defer f.Close()

	if err := e.tmpl.Execute(f, sets); err != nil {
		return fmt.Errorf("failed to write dispatchers to file: %w", err)
	}

	return f.Close()
}

// NewFileExporter creates a new dispatchers.Exporter which writes the dispatcher sets to a file.
// tmpl is optional and if it is set to the empty string, the DefaultFileTemplate will be used, which is compatible with kamailio's dispatchers module as a flat file source.
func NewFileExporter(filename string, tmpl string) (*FileExporter, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename is empty")
	}

	if tmpl == "" {
		tmpl = DefaultFileTemplate
	}

	t, err := template.New("exporter").Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse export template: %w", err)
	}

	return &FileExporter{
		filename: filename,
		tmpl: t,
	}, nil
}
