package output

import (
	_ "embed"

	"encoding/base64"
)

var (
	//go:embed templates/core.html
	coreHtml string
)

var (
	//go:embed templates/favicon.svg
	favicon []byte
)

var (
	faviconBase64 = base64.StdEncoding.EncodeToString(favicon)
)

type data struct {
	Title   string
	Favicon string
	Payload any
}
