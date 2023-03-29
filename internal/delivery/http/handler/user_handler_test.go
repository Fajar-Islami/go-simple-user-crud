package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/helper"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/dtos"
	repositories "github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/mysql"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
	"github.com/stretchr/testify/assert"
)

type handlerTest struct {
	name          string
	url           string
	status        bool
	token         string
	message       string
	errMessage    string
	isError       bool
	expstatuscode int
}

func TruncateUsersTable(db *sql.DB, t *testing.T) {
	t.Helper()
	t.Log("Truncate users table")
	db.Exec("TRUNCATE users")
}

func HelperInsertUser(db *sql.DB, t *testing.T, args repositories.CreateUserParams) int64 {
	t.Helper()
	t.Log("Insert users table")
	const createUser = `-- name: CreateUser :execlastid
	INSERT INTO users (
	  id,email,first_name,last_name,avatar
	) VALUES ( ?,?,?,?,? )`

	_, err := db.Exec(createUser, args.ID, args.Email, args.FirstName, args.LastName, args.Avatar)
	if err != nil {
		t.Fatal("error when insert data : ", err)
	}

	return args.ID
}

func TestHandlerUsersGet(t *testing.T) {
	e := HelperRouterUser(t)

	tests := []handlerTest{
		{
			name:          "Get Data without Params",
			url:           "/api/v1/user",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with page only",
			url:           "/api/v1/user?page=1",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with limit only",
			url:           "/api/v1/user?limit=1",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with page and limit",
			url:           "/api/v1/user?page=1&limit=1",
			status:        true,
			message:       helper.SUCCEEDGETDATA,
			isError:       false,
			expstatuscode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			response := rec.Result()

			body, _ := io.ReadAll(response.Body)
			// t.Log("body :", string(body))
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			dataJson, _ := json.Marshal(responseBody["data"])
			var responseBodyStruct dtos.ResDataUsers
			json.Unmarshal(dataJson, &responseBodyStruct)

			assert.Equal(t, tt.expstatuscode, rec.Code)
			assert.Equal(t, tt.status, responseBody["status"])
			assert.Equal(t, tt.message, responseBody["message"])
			if tt.isError {
				assert.NotNil(t, responseBody["errors"])
				assert.Contains(t, fmt.Sprint(responseBody["errors"]), tt.errMessage)
				assert.Nil(t, responseBody["data"])
			} else {
				assert.Nil(t, responseBody["errors"])
				assert.NotNil(t, responseBody["data"])
			}

		})
	}
}

func TestHandlerUsersFetch(t *testing.T) {
	e := HelperRouterUser(t)

	tests := []handlerTest{
		{
			name:          "Get Data without Params",
			url:           "/api/v1/user/fetch",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with page only",
			url:           "/api/v1/user/fetch?page=1",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with limit only",
			url:           "/api/v1/user/fetch?limit=1",
			status:        false,
			message:       helper.FAILEDGETDATA,
			isError:       true,
			expstatuscode: http.StatusBadRequest,
			errMessage:    "validation error",
		},
		{
			name:          "Get Data with page and limit",
			url:           "/api/v1/user/fetch?page=1&limit=1",
			status:        true,
			message:       helper.SUCCEEDGETDATA,
			isError:       false,
			expstatuscode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			response := rec.Result()

			body, _ := io.ReadAll(response.Body)
			// t.Log("body :", string(body))
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			dataJson, _ := json.Marshal(responseBody["data"])
			var responseBodyStruct dtos.ResDataUsers
			json.Unmarshal(dataJson, &responseBodyStruct)

			assert.Equal(t, tt.expstatuscode, rec.Code)
			assert.Equal(t, tt.status, responseBody["status"])
			assert.Equal(t, tt.message, responseBody["message"])
			if tt.isError {
				assert.NotNil(t, responseBody["errors"])
				assert.Contains(t, fmt.Sprint(responseBody["errors"]), tt.errMessage)
				assert.Nil(t, responseBody["data"])
			} else {
				assert.Nil(t, responseBody["errors"])
				assert.NotNil(t, responseBody["data"])
			}

		})
	}
}

func TestHandlerUsersCreate(t *testing.T) {
	e := HelperRouterUser(t)
	TruncateUsersTable(contConf.Mysqldb, t)

	type handlerTestCreate struct {
		handlerTest
		body dtos.ReqCreateDataUser
	}

	tests := []handlerTestCreate{
		{
			handlerTest: handlerTest{
				name:          "Insert users without fullfil params",
				url:           "/api/v1/user",
				status:        false,
				message:       helper.FAILEDPOSTDATA,
				isError:       true,
				expstatuscode: http.StatusBadRequest,
				errMessage:    "validation error",
			},
			body: dtos.ReqCreateDataUser{
				Email:     "",
				Firstname: "",
				Lastname:  "",
				Avatar:    "",
			},
		},
		{
			handlerTest: handlerTest{
				name:          "Insert users fullfil params",
				url:           "/api/v1/user",
				status:        true,
				message:       helper.SUCCEEDPOSTDATA,
				isError:       false,
				expstatuscode: http.StatusCreated,
			},
			body: dtos.ReqCreateDataUser{
				Email:     "test_1@mail.com",
				Firstname: "test_1",
				Lastname:  "test_2",
				Avatar:    "",
			},
		},
		{
			handlerTest: handlerTest{
				name:          "Insert users with email that has been used",
				url:           "/api/v1/user",
				status:        false,
				message:       helper.FAILEDPOSTDATA,
				isError:       true,
				expstatuscode: http.StatusBadRequest,
				errMessage:    "Duplicate entry",
			},
			body: dtos.ReqCreateDataUser{
				Email:     "test_1@mail.com",
				Firstname: "new test 1",
				Lastname:  "test 2",
				Avatar:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, tt.url, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			response := rec.Result()

			body, _ := io.ReadAll(response.Body)
			// t.Log("body :", string(body))
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			dataJson, _ := json.Marshal(responseBody["data"])
			var responseBodyStruct dtos.ResDataUsers
			json.Unmarshal(dataJson, &responseBodyStruct)

			assert.Equal(t, tt.expstatuscode, rec.Code)
			assert.Equal(t, tt.status, responseBody["status"])
			assert.Equal(t, tt.message, responseBody["message"])
			if tt.isError {
				assert.NotNil(t, responseBody["errors"])
				assert.Contains(t, fmt.Sprint(responseBody["errors"]), tt.errMessage)
				assert.Nil(t, responseBody["data"])
			} else {
				assert.Nil(t, responseBody["errors"])
				assert.NotNil(t, responseBody["data"])
			}

		})
	}
}

func TestHandlerUsersDelete(t *testing.T) {
	e := HelperRouterUser(t)

	TruncateUsersTable(contConf.Mysqldb, t)

	id := utils.IDGenerator()
	idInserted := HelperInsertUser(contConf.Mysqldb, t, repositories.CreateUserParams{
		ID:        id,
		Email:     "emailtestdelete@mail.com",
		FirstName: "test",
		LastName:  "delete",
	})

	t.Log("idInserted ", idInserted)

	tests := []handlerTest{
		{
			name:          "Delete not exist users",
			url:           fmt.Sprintf("/api/v1/user/%d", 1),
			status:        false,
			message:       helper.FAILEDDELETEDATA,
			isError:       true,
			expstatuscode: http.StatusNotFound,
			errMessage:    "Data Users not found",
			token:         "3cdcnTiBsl",
		},
		{
			name:          "Delete exist users",
			url:           fmt.Sprintf("/api/v1/user/%d", idInserted),
			status:        true,
			message:       helper.SUCCEEDDELETEDATA,
			isError:       false,
			expstatuscode: http.StatusOK,
			token:         "3cdcnTiBsl",
		},
		{
			name:          "Delete without token",
			url:           fmt.Sprintf("/api/v1/user/%d", idInserted),
			status:        false,
			message:       "UNATHORIZED",
			isError:       true,
			expstatuscode: http.StatusUnauthorized,
			errMessage:    "FAILED TO GET TOKEN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			if tt.token != "" {
				token := fmt.Sprintf("Bearer %s", tt.token)
				req.Header.Set("Authorization", token)
			}

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			response := rec.Result()

			body, _ := io.ReadAll(response.Body)
			// t.Log("body :", string(body))
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			dataJson, _ := json.Marshal(responseBody["data"])
			var responseBodyStruct dtos.ResDataUsers
			json.Unmarshal(dataJson, &responseBodyStruct)

			assert.Equal(t, tt.expstatuscode, rec.Code)
			assert.Equal(t, tt.status, responseBody["status"])
			assert.Equal(t, tt.message, responseBody["message"])
			if tt.isError {
				assert.NotNil(t, responseBody["errors"])
				assert.Contains(t, fmt.Sprint(responseBody["errors"]), tt.errMessage)
				assert.Nil(t, responseBody["data"])
			} else {
				assert.Nil(t, responseBody["errors"])
				assert.NotNil(t, responseBody["data"])
			}

		})
	}

}

func TestHandlerUsersUpdate(t *testing.T) {
	e := HelperRouterUser(t)

	TruncateUsersTable(contConf.Mysqldb, t)

	id := utils.IDGenerator()
	idInserted := HelperInsertUser(contConf.Mysqldb, t, repositories.CreateUserParams{
		ID:        id,
		Email:     "emailtestupdate@mail.com",
		FirstName: "test",
		LastName:  "update",
	})

	t.Log("idInserted ", idInserted)

	type handlerTestUpdate struct {
		handlerTest
		body dtos.ReqUpdateDataUser
	}

	tests := []handlerTestUpdate{
		{
			handlerTest: handlerTest{
				name:          "Update not exist users",
				url:           fmt.Sprintf("/api/v1/user/%d", 1),
				status:        false,
				message:       helper.FAILEDUPDATEDATA,
				isError:       true,
				expstatuscode: http.StatusNotFound,
				errMessage:    "Data Users not found",
			},
			body: dtos.ReqUpdateDataUser{
				Firstname: "update firstname",
			},
		},
		{
			handlerTest: handlerTest{
				name:          "Update not correct email format",
				url:           fmt.Sprintf("/api/v1/user/%d", idInserted),
				status:        false,
				message:       helper.FAILEDUPDATEDATA,
				isError:       true,
				expstatuscode: http.StatusBadRequest,
				errMessage:    "validation error for 'Email', Tag: email",
			},
			body: dtos.ReqUpdateDataUser{
				Email: "aaa",
			},
		},
		{
			handlerTest: handlerTest{
				name:          "Update success",
				url:           fmt.Sprintf("/api/v1/user/%d", idInserted),
				status:        true,
				message:       helper.SUCCEEDUPDATEDATA,
				isError:       false,
				expstatuscode: http.StatusOK,
			},
			body: dtos.ReqUpdateDataUser{
				Firstname: "update firstname",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, tt.url, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Accept", "application/json")

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			response := rec.Result()

			body, _ := io.ReadAll(response.Body)
			// t.Log("body :", string(body))
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			dataJson, _ := json.Marshal(responseBody["data"])
			var responseBodyStruct dtos.ResDataUsers
			json.Unmarshal(dataJson, &responseBodyStruct)

			assert.Equal(t, tt.expstatuscode, rec.Code)
			assert.Equal(t, tt.status, responseBody["status"])
			assert.Equal(t, tt.message, responseBody["message"])
			if tt.isError {
				assert.NotNil(t, responseBody["errors"])
				assert.Contains(t, fmt.Sprint(responseBody["errors"]), tt.errMessage)
				assert.Nil(t, responseBody["data"])
			} else {
				assert.Nil(t, responseBody["errors"])
				assert.NotNil(t, responseBody["data"])
			}

		})
	}

}
