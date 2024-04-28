package web

import (
	"net/http"

	"github.com/chheller/go-htmx-todo/modules/config"
)

type ErrorPageHandlers struct {
}

func (h *ErrorPageHandlers) Init(router *http.ServeMux) {
	router.HandleFunc("GET /404", h.handle404)
	router.HandleFunc("GET /403", h.handle403)
	router.HandleFunc("/", h.handle404)
}

func (h *ErrorPageHandlers) handle404(res http.ResponseWriter, req *http.Request) {
	Templates.WriteTemplateResponse(res, "/pages", "error_404_page", struct {
		InjectBrowserReloadScript bool
	}{
		InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
	})
}

func (h *ErrorPageHandlers) handle403(res http.ResponseWriter, req *http.Request) {
	Templates.WriteTemplateResponse(res, "/pages", "error_403_page", struct {
		InjectBrowserReloadScript bool
		ErrorMsg                  string
		HttpPrintDebugError       bool
	}{
		InjectBrowserReloadScript: config.GetEnvironment().InjectBrowserReload,
		ErrorMsg:                  "Forbidden",
		HttpPrintDebugError:       config.GetEnvironment().ApplicationConfiguration.HttpPrintDebugError,
	})
}
