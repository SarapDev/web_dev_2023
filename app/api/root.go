package api

import (
	"fmt"
	"github.kolesa-team.org/backend/go-module/env"
	"net/http"
)

type RootHandler struct {
	env     env.Build
	appName string
	branch  string
}

func NewRootHandler(
	env env.Build,
	appName, branch string,
) *RootHandler {
	return &RootHandler{
		env:     env,
		appName: appName,
		branch:  branch,
	}
}

func (h RootHandler) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(
		[]byte(
			fmt.Sprintf(
				"name: %s\nbranch: %s\nbuild: %s",
				h.appName,
				h.branch,
				h.env.String(),
			),
		),
	)
}
