package handlers_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vladkonst/metrics-alerting/app"
	"github.com/vladkonst/metrics-alerting/internal/configs"
)

var a *app.App

func init() {
	cfg := configs.ServerCfg{IntervalsCfg: &configs.ServerIntervalsCfg{}, NetAddressCfg: &configs.NetAddressCfg{}}
	a = app.NewApp(nil, &cfg)

	go func() {
		for range *a.MetricsChan {
			continue
		}
	}()
}

type want struct {
	contentType string
	statusCode  int
	body        string
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, rBody io.Reader) *http.Response {
	req, err := http.NewRequest(method, ts.URL+path, rBody)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	return resp
}

func TestGzipCompression(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		want    want
		body    string
	}{
		{
			name: "compress test",
			want: want{
				statusCode: 200,
				body:       `{"id": "test", "type": "gauge", "value": 1.1}`,
			},
			request: "/update",
			body:    `{"id": "test", "type": "gauge", "value": 1.1}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buff := bytes.NewBuffer(nil)
			zb := gzip.NewWriter(buff)
			_, err := zb.Write([]byte(test.body))
			require.NoError(t, err)
			err = zb.Close()
			require.NoError(t, err)
			req, err := http.NewRequest("POST", ts.URL+test.request, buff)
			require.NoError(t, err)
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Accept-Encoding", "gzip")
			req.Header.Set("Content-Type", "application/json")
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, test.want.statusCode, resp.StatusCode)
			zr, err := gzip.NewReader(resp.Body)
			require.NoError(t, err)
			b, err := io.ReadAll(zr)
			require.NoError(t, err)
			assert.JSONEq(t, test.want.body, string(b))
		})
	}
}

func TestUpdateMetric(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		want    want
		body    string
	}{
		{
			name: "without body test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  500,
				body:        "",
			},
			request: "/update",
			body:    "",
		},
		{
			name: "unsupported type test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  422,
				body:        "",
			},
			request: "/update",
			body:    `{"ID": "test", "MType": "unsupported", "Delta": 0, "Value": 0.0}`,
		},
		{
			name: "success test",
			want: want{
				contentType: "application/json",
				statusCode:  200,
				body:        `{"id": "test", "type": "gauge", "value": 1.1}`,
			},
			request: "/update",
			body:    `{"id": "test", "type": "gauge", "value": 1.1}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buff := bytes.NewBufferString(test.body)
			req, err := http.NewRequest("POST", ts.URL+test.request, buff)
			require.NoError(t, err)
			req.Header.Add("Content-Type", "application/json")
			req.Header.Set("Accept-Encoding", "")
			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
			assert.Equal(t, test.want.statusCode, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
			if test.want.body != "" {
				body := make([]byte, 1024)
				n, _ := resp.Body.Read(body)
				defer resp.Body.Close()
				assert.JSONEq(t, test.want.body, string(body[:n]))
			}
		})
	}
}

func TestGetGaugeMetricValue(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "metric without name test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/value/gauge/",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := testRequest(t, ts, "GET", test.request, nil)
			defer res.Body.Close()
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGetCounterMetricValue(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		want    want
	}{
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
			res := testRequest(t, ts, "GET", test.request, nil)
			defer res.Body.Close()
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestUpdateGaugeMetric(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := testRequest(t, ts, "POST", test.request, nil)
			defer res.Body.Close()
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestUpdateCounterMetric(t *testing.T) {
	ts := httptest.NewServer(a.GetRouter())
	defer ts.Close()
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name: "success test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
			request: "/update/counter/Alloc/1",
		},
		{
			name: "metric without value test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/counter/Alloc/",
		},
		{
			name: "metric without name test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
			request: "/update/counter/",
		},
		{
			name: "incorrect metric value test",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
			request: "/update/counter/Alloc/Alloc",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := testRequest(t, ts, "POST", test.request, nil)
			defer res.Body.Close()
			assert.Equal(t, test.want.statusCode, res.StatusCode)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
