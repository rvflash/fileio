package fileio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// DefaultExpires is the default number of days until the file will be deleted by File.io.
const DefaultExpires = 14

// Enable testing by mocking the http client.
type httpClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, bodyType string, body io.Reader) (*http.Response, error)
}

var api httpClient = http.DefaultClient

// Response contains a FileIO response.
type Response struct {
	Success bool   `json:"success"`
	Code    int    `json:"error,omitempty"`
	Err     string `json:"message,omitempty"`
	Expiry  string `json:"expiry,omitempty"`
	Key     string `json:"key,omitempty"`
}

// URL is by default the url of the File.io API.
var URL = "https://file.io"

// Download downloads the file behind the key to the given file.
// An error occurs if the key does not exist or if the file fails to be create.
func Download(key, file string) (err error) {
	// Downloads the file with this key on File.io.
	var resp *http.Response
	if resp, err = api.Get(URL + "/" + key); err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		err = errors.New(resp.Status)
		return
	}

	// Creates a file with the given name.
	var out *os.File
	if out, err = os.Create(file); err != nil {
		return
	}
	defer func() { _ = out.Close() }()

	// Copies the response to the destination.
	_, err = io.Copy(out, resp.Body)

	return
}

// Upload uploads a file to file.io and returns its key or an error if it can not.
// A default expires of 14 days is internally used.
func Upload(file string) (string, error) {
	rs, err := postBody(file, URL)
	if err != nil {
		return "", err
	}
	return rs.Key, nil
}

// UploadWithExpire uploads a file to file.io and sets the expires in days.
// It returns its key, the expiry duration and an error if it can not to get it.
func UploadWithExpire(file string, days int) (string, string, error) {
	rs, err := postBody(file, URL+"/?expires="+expires(days))
	if err != nil {
		return "", "", err
	}
	return rs.Key, rs.Expiry, nil
}

// Converts number of day as expected by the API.
func expires(days int) string {
	if days < 1 {
		return strconv.Itoa(DefaultExpires)
	}
	if days%7 == 0 {
		return strconv.Itoa(days/7) + "w"
	}
	if days%31 == 0 {
		return strconv.Itoa(days/31) + "m"
	}
	if days%365 == 0 {
		return strconv.Itoa(days/365) + "y"
	}
	return strconv.Itoa(days)
}

// {"success":true,"key":"2ojE41"}
// {"success":true,"key":"aQbnDJ","expiry":"7 days"}
// {"success":false,"error":404,"message":"Not Found"}
func parseJSON(data []byte) (*Response, error) {
	res := &Response{}
	if err := json.Unmarshal(data, res); err != nil {
		return res, err
	}
	// The action fails for File.io, deals with it.
	if !res.Success {
		return res, errors.New(res.Err)
	}
	return res, nil
}

func postBody(file, url string) (*Response, error) {
	// Create the form to post.
	body, bodyType, err := createBody(file)
	if err != nil {
		return nil, err
	}
	// Uploads the file.
	resp, err := api.Post(url, bodyType, body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// Gets the response and parse it as JSON.
	rs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseJSON(rs)
}

func createBody(file string) (*bytes.Buffer, string, error) {
	// Create a buffer for the form.
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	defer func() { _ = w.Close() }()

	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return nil, "", err
	}

	// Opens file handle.
	f, err := os.Open(file)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = f.Close() }()

	// Adds the file to the form.
	if _, err = io.Copy(fw, f); err != nil {
		return nil, "", err
	}
	return body, w.FormDataContentType(), nil
}
