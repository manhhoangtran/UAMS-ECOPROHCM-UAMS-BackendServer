//go:build integration
// +build integration

package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/ecoprohcm/DMS_BackendServer/handlers"
	"github.com/ecoprohcm/DMS_BackendServer/initializers"
	"github.com/gin-gonic/gin"
)

type TestRouter struct {
	GinRouter *gin.Engine
}

var GlobalTestRouter = &TestRouter{}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	gin.SetMode(gin.TestMode)

	// Change to Work Dir when testing
	_, filename, _, _ := runtime.Caller(0)
	os.Chdir(path.Join(path.Dir(filename), ".."))
	wd, _ := os.Getwd()
	cc, _, err := initializers.InitApplication(fmt.Sprintf("%s/%s", wd, ".env.test"))
	if err != nil {
		fmt.Printf("failed to create event: %s\n", err)
		os.Exit(2)
	}

	// setup router
	router := handlers.SetupRouter(cc.HandlerOptions)
	GlobalTestRouter.GinRouter = router
}

func shutdown() {

}

func DoRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	return w
}

func DoRequestWithBody(r http.Handler, method, path string, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	req.Header.Add("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	return w
}
