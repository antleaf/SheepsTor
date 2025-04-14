package cmd

import (
	"embed"
	"github.com/unrolled/render"
)

// ### CSS
//
//go:embed assets
var embeddedAssets embed.FS

func NewRenderer() *render.Render {
	renderer := render.New(render.Options{
		Layout:          "layout",
		RequirePartials: true,
	})
	return renderer
}
