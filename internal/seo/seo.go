package seo

import (
	"html"
	"net/url"

	"portaljuridico/internal/content"
	"portaljuridico/internal/editorial"
)

func IsAbsoluteHTTPSURL(value string) bool {
	parsed, err := url.Parse(value)
	return err == nil && parsed.Scheme == "https" && parsed.Host != ""
}

func RenderHead(page content.Page) string {
	return `<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>` + html.EscapeString(page.Title) + `</title>
<meta name="description" content="` + html.EscapeString(page.MetaDescription) + `">
<link rel="canonical" href="` + html.EscapeString(page.CanonicalURL) + `">
<meta name="robots" content="` + editorial.RobotsDirective(page) + `">`
}
