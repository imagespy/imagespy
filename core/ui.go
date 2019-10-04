package core

import (
	"fmt"
	"net/http"

	"github.com/cbroglie/mustache"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	f *Finder
	t *mustache.Template
}

func (h *Handler) handle(w http.ResponseWriter, r *http.Request) {
	results, err := h.f.Find()
	if err != nil {
		log.Errorf("read result for ui: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "finder returned an error")
		return
	}

	var data []map[string]string
	for _, r := range results {
		rowCSSClass := "table-success"
		if r.Current.Digest != r.Latest.Digest {
			rowCSSClass = "table-danger"
		}

		data = append(data, map[string]string{
			"current_digest": r.Current.Digest[:15],
			"current_tag":    r.Current.Tag,
			"input":          r.Input,
			"instance":       r.Instance,
			"latest_digest":  r.Latest.Digest[:15],
			"latest_tag":     r.Latest.Tag,
			"row_css_class":  rowCSSClass,
			"repository":     r.Current.Repository,
			"source":         r.Current.Source,
		})
	}

	templateData := map[string]interface{}{
		"images": data,
	}
	tpl, err := h.t.Render(templateData)
	if err != nil {
		log.Errorf("render template: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "template could not render")
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tpl))
}

func NewUIHandler(f *Finder, templatePath string) (*Handler, error) {
	t, err := mustache.ParseFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("parse template file: %w", err)
	}

	return &Handler{f, t}, nil
}
