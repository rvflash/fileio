package fileio

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var downloadTest = []struct {
	url, key, file string
	err            error
}{
	{url: URL, key: "exists", file: "test.txt"},
	{url: URL, key: "not_exists", file: "test2.txt", err: errors.New("Not Found")},
	{key: "fails", file: "test.txt", err: errors.New("No transport")},
}

var uploadTest = []struct {
	url, key, file, expiry string
	err                    error
	expires                int
}{
	{expires: -1, url: URL, file: "test.txt", key: "2ojE41"},
	{expires: 7, url: URL, file: "test.txt", key: "2ojE41", expiry: "7 days"},
	{expires: 12, url: URL, file: "test.txt", key: "2ojE41", expiry: "12 days"},
	{expires: 31, url: URL, file: "test.txt", key: "2ojE41", expiry: "1 month"},
	{expires: 666, url: URL, file: "test.txt", err: errors.New("Internal error")},
	{expires: 730, url: URL, file: "test.txt", key: "2ojE41", expiry: "2 years"},
	{expires: 999, url: URL, file: "test.txt", err: errors.New("unexpected end of JSON input")},
	{url: URL, file: "test2.txt", err: errors.New("open test2.txt: no such file or directory")},
	{file: "test.txt", err: errors.New("No transport")},
}

// Builds a fake http client by mocking main methods.
type fakeHTTPClient struct{}

// Get mocks the method of same name of the http package.
func (c *fakeHTTPClient) Get(url string) (*http.Response, error) {
	return fakeHTTPHandler(url, "GET")
}

// Post mocks the method of same name of the http package.
func (c *fakeHTTPClient) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return fakeHTTPHandler(url, "POST")
}

func fakeHTTPHandler(url, method string) (*http.Response, error) {
	if !strings.HasPrefix(url, "http") {
		return nil, errors.New("No transport")
	}
	// Mocks responses base on the URL.
	urlHandler := func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if p == "" {
			p = "/"
		}
		if r.URL.RawQuery != "" {
			p += "?" + r.URL.RawQuery
		}
		switch p {
		case "/", "/?expires=14":
			_, _ = io.WriteString(w, `{"success":true,"key":"2ojE41"}`)
		case "/exists":
			_, _ = io.WriteString(w, "This is a test")
		case "/?expires=1w":
			_, _ = io.WriteString(w, `{"success":true,"key":"2ojE41","expiry":"7 days"}`)
		case "/?expires=1m":
			_, _ = io.WriteString(w, `{"success":true,"key":"2ojE41","expiry":"1 month"}`)
		case "/?expires=2y":
			_, _ = io.WriteString(w, `{"success":true,"key":"2ojE41","expiry":"2 years"}`)
		case "/?expires=12":
			_, _ = io.WriteString(w, `{"success":true,"key":"2ojE41","expiry":"12 days"}`)
		case "/?expires=666":
			_, _ = io.WriteString(w, `{"success":false,"error":500,"message":"Internal error"}`)
		case "/?expires=999":
			_, _ = io.WriteString(w, "")
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"success":false,"error":404,"message":"Not Found"}`)
		}
	}

	req := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	urlHandler(w, req)
	return w.Result(), nil
}

// TestDownload tests the method Download.
func TestDownload(t *testing.T) {
	api = &fakeHTTPClient{}

	// Restore http client at the end of the test.
	defer func() { api = http.DefaultClient }()

	// Restore the default url of the API.
	url := URL
	defer func() { URL = url }()

	for _, dt := range downloadTest {
		URL = dt.url
		if err := Download(dt.key, dt.file); err == nil {
			if dt.err != nil {
				t.Errorf("Expected error %v, received no error with /%s", dt.err, dt.key)
			}
		} else if dt.err == nil {
			t.Errorf("Expected no error with /%s, received %v", dt.key, err)
		} else if err.Error() != dt.err.Error() {
			t.Errorf("Expected error %v with /%s, received %v", dt.err, dt.key, err)
		}
	}
}

// TestUpload tests the method Upload.
func TestUpload(t *testing.T) {
	api = &fakeHTTPClient{}

	// Restore http client at the end of the test.
	defer func() { api = http.DefaultClient }()

	// Restore the default url of the API.
	url := URL
	defer func() { URL = url }()

	for _, dt := range uploadTest {
		if dt.expires > 0 {
			continue
		}
		URL = dt.url
		if key, err := Upload(dt.file); err == nil {
			if dt.err != nil {
				t.Errorf("Expected error %v, received no error with file named %s", dt.err, dt.file)
			} else if key != dt.key {
				t.Errorf("Expected key /%s, received /%s with file named %s", dt.key, key, dt.file)
			}
		} else if dt.err == nil {
			t.Errorf("Expected no error with /%s, received %v", dt.file, err)
		} else if err.Error() != dt.err.Error() {
			t.Errorf("Expected error %v with /%s, received %v", dt.err, dt.file, err)
		}
	}
}

// TestUploadWithExpire tests the method UploadWithExpire.
func TestUploadWithExpire(t *testing.T) {
	api = &fakeHTTPClient{}

	// Restore http client at the end of the test.
	defer func() { api = http.DefaultClient }()

	// Restore the default url of the API.s
	url := URL
	defer func() { URL = url }()

	for _, dt := range uploadTest {
		if dt.expires < 0 {
			continue
		}
		URL = dt.url
		if key, expiry, err := UploadWithExpire(dt.file, dt.expires); err == nil {
			if dt.err != nil {
				t.Errorf("Expected error %v, received no error with file named %s", dt.err, dt.file)
			} else if key != dt.key {
				t.Errorf("Expected key /%s, received /%s with file named %s", dt.key, key, dt.file)
			} else if expiry != dt.expiry {
				t.Errorf("Expected expiry in %s, received %s with file named %s", dt.expiry, expiry, dt.file)
			}
		} else if dt.err == nil {
			t.Errorf("Expected no error with /%s, received %v", dt.file, err)
		} else if err.Error() != dt.err.Error() {
			t.Errorf("Expected error %v with /%s, received %v", dt.err, dt.file, err)
		}
	}
}
