//go:build integration
// +build integration

package tests

import (
	"net/http"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestFindAllArea(t *testing.T) {
	// reqStr := fmt.Sprintf(`{"name":"test","manager":"test"}`)
	// resp := DoRequestWithBody(GlobalTestRouter.GinRouter, "POST", "/v1/area", reqStr)
	// byteBody, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("ala %+v\n", string(byteBody))
	// assert.Equal(t, http.StatusOK, resp.Code)

	w := DoRequest(GlobalTestRouter.GinRouter, "GET", "/v1/areas")
	// assert if we got 200 on every request
	assert.Equal(t, http.StatusOK, w.Code)

}
