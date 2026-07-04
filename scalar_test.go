package gin_swagger_scalar

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

// mockedSwag is a minimal swag.Swagger implementation used to register a doc
// so the "doc.json" route has something to return.
type mockedSwag struct{}

func (s *mockedSwag) ReadDoc() string {
	return `{"swagger":"2.0","info":{"title":"test"}}`
}

func init() {
	gin.SetMode(gin.TestMode)
	swag.Register("swagger", &mockedSwag{})
}

// newRouter mounts the given handler the same way consumers do.
func newRouter(handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET("/scalar/*any", handler)
	return r
}

// perform issues a request against the router and returns the recorder.
func perform(r *gin.Engine, method, target string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestWrapHandlerServesIndex(t *testing.T) {
	r := newRouter(WrapHandler(nil))

	w := perform(r, http.MethodGet, "/scalar/")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("content-type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected html content-type, got %q", ct)
	}
	// The rendered index should reference the standalone script and the doc URL.
	if body := w.Body.String(); !containsAll(body, "/scalar/standalone.js", "doc.json") {
		t.Errorf("index body missing expected references, got:\n%s", body)
	}
}

func TestServesStandaloneJS(t *testing.T) {
	r := newRouter(WrapHandler(nil))

	w := perform(r, http.MethodGet, "/scalar/standalone.js")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("content-type"); ct != "application/javascript" {
		t.Errorf("expected javascript content-type, got %q", ct)
	}
	if w.Body.Len() == 0 {
		t.Error("expected standalone.js body to be non-empty")
	}
}

func TestServesDocJSON(t *testing.T) {
	r := newRouter(WrapHandler(nil))

	w := perform(r, http.MethodGet, "/scalar/doc.json")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("content-type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected json content-type, got %q", ct)
	}
	if body := w.Body.String(); !containsAll(body, "swagger") {
		t.Errorf("expected doc.json body to contain the registered doc, got:\n%s", body)
	}
}

func TestServesIndexHTML(t *testing.T) {
	r := newRouter(WrapHandler(nil))

	w := perform(r, http.MethodGet, "/scalar/index.html")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	if ct := w.Header().Get("content-type"); ct != "text/html; charset=utf-8" {
		t.Errorf("expected html content-type, got %q", ct)
	}
}

func TestUnknownPathReturns404(t *testing.T) {
	r := newRouter(WrapHandler(nil))

	w := perform(r, http.MethodGet, "/scalar/does-not-exist")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestNonGetReturns405(t *testing.T) {
	r := gin.New()
	// Register the handler for POST so it is reached and can reject the method.
	r.POST("/scalar/*any", WrapHandler(nil))

	w := perform(r, http.MethodPost, "/scalar/standalone.js")

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCustomWrapHandlerUsesConfig(t *testing.T) {
	config := &Config{
		URL:          "custom.json",
		InstanceName: "swagger",
		BasePath:     "/docs",
		ProjectName:  "My API",
	}

	r := gin.New()
	r.GET("/docs/*any", CustomWrapHandler(config))

	w := perform(r, http.MethodGet, "/docs/")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w.Code)
	}
	// The custom base path and doc URL should be reflected in the rendered index.
	if body := w.Body.String(); !containsAll(body, "/docs/standalone.js", "custom.json") {
		t.Errorf("index body missing custom config references, got:\n%s", body)
	}
}

func TestConfigOptions(t *testing.T) {
	config := Config{}
	ProjectName("My API")(&config)
	BasePath("/docs")(&config)
	InstanceName("custom")(&config)
	URL("custom.json")(&config)

	if config.ProjectName != "My API" {
		t.Errorf("ProjectName option not applied, got %q", config.ProjectName)
	}
	if config.BasePath != "/docs" {
		t.Errorf("BasePath option not applied, got %q", config.BasePath)
	}
	if config.InstanceName != "custom" {
		t.Errorf("InstanceName option not applied, got %q", config.InstanceName)
	}
	if config.URL != "custom.json" {
		t.Errorf("URL option not applied, got %q", config.URL)
	}
}

func TestDisablingWrapHandler(t *testing.T) {
	const env = "DISABLE_SCALAR"
	t.Setenv(env, "true")

	r := newRouter(DisablingWrapHandler(nil, env))

	w := perform(r, http.MethodGet, "/scalar/")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d when disabled, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDisablingWrapHandlerEnabledWhenEnvUnset(t *testing.T) {
	r := newRouter(DisablingWrapHandler(nil, "DISABLE_SCALAR_UNSET"))

	w := perform(r, http.MethodGet, "/scalar/standalone.js")

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d when enabled, got %d", http.StatusOK, w.Code)
	}
}

func TestDisablingCustomWrapHandler(t *testing.T) {
	const env = "DISABLE_CUSTOM_SCALAR"
	t.Setenv(env, "true")

	config := &Config{
		URL:          "doc.json",
		InstanceName: "swagger",
		BasePath:     "/scalar",
		ProjectName:  "Scalar",
	}

	r := newRouter(DisablingCustomWrapHandler(config, nil, env))

	w := perform(r, http.MethodGet, "/scalar/")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status %d when disabled, got %d", http.StatusNotFound, w.Code)
	}
}

// containsAll reports whether s contains every one of the given substrings.
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
