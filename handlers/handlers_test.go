package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vladkonst/metrics-alerting/routers"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

// func TestUpdateAndGetCounterMetric(t *testing.T) {
// 	ts := httptest.NewServer(routers.GetRouter())
// 	defer ts.Close()

// 	t.Run("update and get gauge test", func(t *testing.T) {
// 		testRequest(t, ts, "POST", "/update/counter/Alloc/1")
// 		res := testRequest(t, ts, "GET", "/value/counter/Alloc/")
// 		var body []byte
// 		res.Body.Read(body)
// 		defer res.Body.Close()
// 		assert.Equal(t, 200, res.StatusCode)
// 		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
// 		assert.Equal(t, "1", string(body))
// 	})
// }

// func TestUpdateAndGetGaugeMetric(t *testing.T) {
// 	ts := httptest.NewServer(routers.GetRouter())
// 	defer ts.Close()

// 	t.Run("update and get gauge test", func(t *testing.T) {
// 		testRequest(t, ts, "POST", "/update/gauge/Alloc/1.0")
// 		res := testRequest(t, ts, "GET", "/value/gauge/Alloc/")
// 		var body []byte
// 		res.Body.Read(body)
// 		defer res.Body.Close()
// 		assert.Equal(t, 200, res.StatusCode)
// 		assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
// 		assert.Equal(t, "1.0", string(body))
// 	})
// }

func TestGetCurrentMetricValue(t *testing.T) {
	ts := httptest.NewServer(routers.GetRouter())
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "metric without name gauge test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/value/gauge/",
		},
		{
			name: "metric without name counter test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/value/counter/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := testRequest(t, ts, "GET", test.request)
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestUpdateMetric(t *testing.T) {
	ts := httptest.NewServer(routers.GetRouter())
	defer ts.Close()
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{{
		name: "success gauge test",
		want: want{
			contentType: "text/plain; charset=utf-8",
			statusCode:  200,
		},
		request: "/update/gauge/Alloc/1.0",
	},
		{
			name: "metric without value gauge test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/gauge/Alloc/",
		},
		{
			name: "metric without name gauge test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/gauge/",
		},
		{
			name: "incorrect metric value gauge test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			request: "/update/gauge/Alloc/Alloc",
		},
		{
			name: "success counter test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
			request: "/update/counter/Alloc/1",
		},
		{
			name: "metric without value counter test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/counter/Alloc/",
		},
		{
			name: "metric without name counter test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/counter/",
		},
		{
			name: "incorrect metric value counter test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			request: "/update/counter/Alloc/Alloc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := testRequest(t, ts, "POST", test.request)
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
