package build

import (
	"archive/zip"
	"encoding/json"
	"io"
	"text/template"
)

type Package struct {
	Name   string `json:"name"`
	Bin    string `json:"-"`
	Main   string `json:"main"`
	Window Window `json:"window"`
}

type Window struct {
	Title    string `json:"title,omitempty"`
	Toolbar  bool   `json:"toolbar,omitempty"`
	Show     bool   `json:"show,omitempty"`
	Position string `json:"position,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type Templates struct {
	IndexHtml string
	ClientJs  string
	ScriptJs  string
}

var DefaultTemplates = Templates{IndexHtml: index, ClientJs: client, ScriptJs: script}

// CreateNW creates a node-webkit .nw file
func (p Package) CreateNW(zw *zip.Writer, templates Templates, myapp io.Reader) error {
	// Add in a couple of package defaults
	p.Main = "index.html"

	if w, err := zw.Create("package.json"); err != nil {
		return err
	} else {
		if _, err := p.writeJsonTo(w); err != nil {
			return err
		}
	}

	filenameTemplates := map[string]string{
		"index.html": templates.IndexHtml,
		"client.js":  templates.ClientJs,
		"script.js":  templates.ScriptJs}
	for filename, str := range filenameTemplates {
		if w, err := zw.Create(filename); err != nil {
			return err
		} else {
			if t, err := template.New(filename).Parse(str); err != nil {
				return err
			} else {
				if err := t.Execute(w, p); err != nil {
					return err
				}
			}
		}
	}

	if w, err := zw.Create(p.Bin); err != nil {
		return err
	} else {
		if _, err := io.Copy(w, myapp); err != nil {
			return err
		}
	}

	return nil
}

func (p Package) writeJsonTo(w io.Writer) (int64, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(b)
	return int64(n), err
}