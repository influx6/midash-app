// Package app is an auto-generated package which exposes the Gu.NApp which 
// can be created to use the constructed views if any. Edit as you see fit.

//go:generate go generate ./components/...
//go:generate go generate ./manifests/...

package app

import (
    "net/http"
	"github.com/gu-io/gu"
	"github.com/gu-io/gu/router"
	"github.com/gu-io/gu/router/cache"
	"github.com/influx6/midash-app/app/manifests"
)

// New returns a new gu.NApp using the provided arguments as needed.
func New(cache router.Cache, server http.Handler) *gu.NApp {
	router := router.NewRouter(server, cache)

	app := gu.App(gu.AppAttr{
		Name:              "app",
		Title:             "app Gu App",
		Manifests:         manifests.Manifests,
		Router: 		   router,
	})

    return app
}
