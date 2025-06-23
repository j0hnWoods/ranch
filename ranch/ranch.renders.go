package ranch

import (
	"fmt"
)

// RanchRender - слой рендеринга (пока заглушка)
type RanchRender struct {
}

func renderFarmTemplate(content map[string]interface{}) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <meta name="description" content="%s">
</head>
<body>
    <h1>%s</h1>
    <div class="content">%s</div>
    <div class="keywords">Domain: %s</div>
</body>
</html>`,
		content["title"], content["description"],
		content["keyword"], content["content"], content["domain"])
}

func renderWhiteTemplate(content map[string]interface{}) string {
	return `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome</title>
</head>
<body>
    <h1>Welcome</h1>
    <p>This is a legitimate website.</p>
    <div class="content">
        <p>Normal content for human visitors.</p>
    </div>
</body>
</html>`
}

func renderError(message string) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body><h1>Error</h1><p>%s</p></body>
</html>`, message)
}

func generateSEOContent() string {
	return `
<h1>SEO Optimized Content</h1>
<p>This is generated SEO content with keywords and optimized structure.</p>
<div class="seo-content">
    <h2>Key Information</h2>
    <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit...</p>
</div>
`
}
