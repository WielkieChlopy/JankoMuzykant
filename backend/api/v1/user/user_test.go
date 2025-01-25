package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"

	"backend/testutils"
)

func TestSignup(t *testing.T) {

	table := []struct {
		name        string
		payload     string
		assert_func func(*testutils.Test, *httptest.ResponseRecorder)
	}{
		{
			"good", `{"user":{"username":"test1","password":"secret"}}`,
			func(test *testutils.Test, rec *httptest.ResponseRecorder) {
				var m map[string]map[string]interface{}
				if !assert.Equal(t, http.StatusCreated, rec.Code) {
					log.Error(rec.Body.String())
				}
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m))
				assert.Equal(t, "test1", m["user"]["username"])
				assert.Empty(t, m["user"]["password"])
				assert.NotEmpty(t, m["user"]["id"])
				assert.NotEmpty(t, m["user"]["token"])
			},
		},
		{
			"missing username",
			`{"user":{"password":"secret"}}`,
			func(test *testutils.Test, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			},
		},
		{
			"missing password",
			`{"user":{"username":"secret"}}`,
			func(test *testutils.Test, rec *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
			},
		},
	}
	for _, tc := range table {
		t.Run(tc.name, func(t *testing.T) {
			test, err := testutils.SetupTest()
			assert.NoError(t, err)

			req := httptest.NewRequest(echo.POST, "/api/v1/singup", strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := test.Router.NewContext(req, rec)

			assert.NoError(t, test.Handler.UserHandler.SignUp(c))
			tc.assert_func(test, rec)
		})
	}
}
