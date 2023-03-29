package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/dtos"
	"github.com/go-playground/assert/v2"
)

func TestHandlerUsersGet(t *testing.T) {
	e := HelperRouterUser(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/user?page=1&limit=1", nil)
	fmt.Println("req.URL", req.URL)

	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	response := rec.Result()
	fmt.Println(rec.Code)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	dataJson, _ := json.Marshal(responseBody["data"])
	var responseBodyStruct dtos.ResDataUsers
	json.Unmarshal(dataJson, &responseBodyStruct)
	t.Logf("responseBodyStruct %#v", responseBodyStruct)

	assert.Equal(t, http.StatusOK, rec.Code)
}
