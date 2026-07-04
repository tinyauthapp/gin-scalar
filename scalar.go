package gin_swagger_scalar

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "embed"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

//go:embed standalone.js
var ScalarStandaloneJS string

//go:embed scalar.gohtml
var ScalarGoHTML string

// Config stores the Scalar configuration options
type Config struct {
	// The URL pointing to API definition (normally swagger.json or swagger.yaml), default is `doc.json`
	URL string
	// The swagger instance name, default is swag
	InstanceName string
	// The base path of Scalar UI, default is /scalar
	BasePath string
	// Project name, default is "Scalar"
	ProjectName string
}

// The configuration for the Scalar template
type scalarConfig struct {
	ScalarScript string
	DocPath      string
	ProjectName  string
}

// ProjectName sets the title of the Scalar UI.
// Defaults to "Scalar".
func ProjectName(name string) func(*Config) {
	return func(c *Config) {
		c.ProjectName = name
	}
}

// BasePath sets the base path in which the Scalar UI will be served.
// Defaults to "/scalar".
func BasePath(basePath string) func(*Config) {
	return func(c *Config) {
		c.BasePath = basePath
	}
}

// InstanceName sets the instance name used to generate the swagger documents.
// Defaults to swag.Name ("swagger").
func InstanceName(name string) func(*Config) {
	return func(c *Config) {
		c.InstanceName = name
	}
}

// URL sets the URL pointing to API definition (normally swagger.json or swagger.yaml).
// Defaults to "doc.json".
func URL(url string) func(*Config) {
	return func(c *Config) {
		c.URL = url
	}
}

// WrapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func WrapHandler(_ *http.Handler, options ...func(*Config)) gin.HandlerFunc {
	var config = Config{
		URL:          "doc.json",
		InstanceName: "swagger",
		BasePath:     "/scalar",
		ProjectName:  "Scalar",
	}

	for _, opt := range options {
		opt(&config)
	}

	return CustomWrapHandler(&config)
}

// CustomWrapHandler wraps `http.Handler` into `gin.HandlerFunc`.
func CustomWrapHandler(config *Config) gin.HandlerFunc {
	// create the scalar config
	scalarCfg := scalarConfig{
		ScalarScript: config.BasePath + "/standalone.js",
		DocPath:      config.URL,
	}

	// create the scalar template
	index := template.Must(template.New("index").Parse(ScalarGoHTML))

	return func(ctx *gin.Context) {
		if ctx.Request.Method != http.MethodGet {
			ctx.AbortWithStatus(http.StatusMethodNotAllowed)
			return
		}

		// return index.html content for /scalar or /scalar/
		if strings.TrimSuffix(ctx.Request.RequestURI, "/") == config.BasePath {
			ctx.Header("content-type", "text/html; charset=utf-8")
			buf := bytes.Buffer{}
			err := index.Execute(&buf, scalarCfg)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			ctx.String(http.StatusOK, buf.String())
			return
		}

		path := strings.TrimPrefix(ctx.Request.RequestURI, config.BasePath+"/")

		switch filepath.Ext(path) {
		case ".html":
			ctx.Header("content-type", "text/html; charset=utf-8")
		case ".js":
			ctx.Header("content-type", "application/javascript")
		case ".json":
			ctx.Header("content-type", "application/json; charset=utf-8")
		}

		switch path {
		case "index.html":
			buf := bytes.Buffer{}
			err := index.Execute(&buf, scalarCfg)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			ctx.String(http.StatusOK, buf.String())
			return
		case "standalone.js":
			ctx.String(http.StatusOK, ScalarStandaloneJS)
			return
		case "doc.json":
			doc, err := swag.ReadDoc(config.InstanceName)
			if err != nil {
				ctx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			ctx.String(http.StatusOK, doc)
		default:
			ctx.AbortWithStatus(http.StatusNotFound)
		}
	}
}

// DisablingWrapHandler turns handler off
// if specified environment variable passed.
func DisablingWrapHandler(handler *http.Handler, env string) gin.HandlerFunc {
	if os.Getenv(env) != "" {
		return func(c *gin.Context) {
			// Simulate behavior when route unspecified and
			// return 404 HTTP code
			c.String(http.StatusNotFound, "")
		}
	}

	return WrapHandler(handler)
}

// DisablingCustomWrapHandler turn handler off
// if specified environment variable passed.
func DisablingCustomWrapHandler(config *Config, handler *http.Handler, env string) gin.HandlerFunc {
	if os.Getenv(env) != "" {
		return func(c *gin.Context) {
			// Simulate behavior when route unspecified and
			// return 404 HTTP code
			c.String(http.StatusNotFound, "")
		}
	}

	return CustomWrapHandler(config)
}
