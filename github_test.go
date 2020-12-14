package github

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

// Custom type that allows setting the func that our Mock Do func will run instead
type MockDoType func(req *http.Request) (*http.Response, error)

// MockClient is the mock client
type MockClient struct {
	MockDo MockDoType
}

// Overriding what the Do function should "do" in our MockClient
func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestGitHubCallSuccess(t *testing.T) {

	// build our response JSON
	jsonResponse := `[{
			"full_name": "mock-repo"
		}]`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	result, err := GetRepos("atkinsonbg")
	if err != nil {
		t.Error("TestGitHubCallSuccess failed.")
		return
	}

	if len(result) == 0 {
		t.Error("TestGitHubCallSuccess failed, array was empty.")
		return
	}

	if result[0]["full_name"] != "mock-repo" {
		t.Error("TestGitHubCallSuccess failed, array was not sorted correctly.")
		return
	}
}

func TestGitHubCallFail(t *testing.T) {

	// create a client that throws and returns an error
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 404,
				Body:       nil,
			}, errors.New("Mock Error")
		},
	}

	_, err := GetRepos("atkinsonbgthisusershouldnotexist")
	if err == nil {
		t.Error("TestGitHubCallFail failed.")
		return
	}
}

func TestGitHubCallBadJsonFail(t *testing.T) {

	// build our response JSON
	jsonResponse := `{
		"full_name": "mock-repo"
	}`

	// create a new reader with that JSON
	r := ioutil.NopCloser(bytes.NewReader([]byte(jsonResponse)))

	// create a client that throws and returns an error
	Client = &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	_, err := GetRepos("atkinsonbg")
	if err == nil {
		t.Error("TestGitHubCallBadJsonFail failed.")
		return
	}
}
