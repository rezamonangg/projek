package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestClient struct {
	baseURL    string
	httpClient *http.Client
	jar        *cookiejar.Jar
}

func NewTestClient(baseURL string) *TestClient {
	jar, _ := cookiejar.New(nil)
	return &TestClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Jar: jar,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		jar: jar,
	}
}

func (c *TestClient) Post(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *TestClient) Get(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *TestClient) Put(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *TestClient) Patch(endpoint string, body interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", c.baseURL+endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

func (c *TestClient) Delete(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *TestClient) GetCookies() []*http.Cookie {
	u, _ := url.Parse(c.baseURL)
	return c.jar.Cookies(u)
}

func (c *TestClient) GetSessionCookie() *http.Cookie {
	for _, cookie := range c.GetCookies() {
		if cookie.Name == "session_id" {
			return cookie
		}
	}
	return nil
}

func ParseResponse(t *testing.T, resp *http.Response, target interface{}) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	err = json.Unmarshal(body, target)
	require.NoError(t, err, "Failed to parse JSON response: %s", string(body))
}

func ParseResponseBody(t *testing.T, resp *http.Response) string {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "Failed to read response body")

	return string(body)
}

func AssertStatusCode(t *testing.T, expected, actual int, respBody string) {
	if expected != actual {
		t.Errorf("Expected status %d, got %d. Body: %s", expected, actual, respBody)
	}
}

func PrintResponse(t *testing.T, resp *http.Response, name string) {
	body := ParseResponseBody(t, resp)
	t.Logf("%s: Status=%d, Body=%s", name, resp.StatusCode, body)
}

type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func ParseError(t *testing.T, resp *http.Response) ErrorResponse {
	var errResp ErrorResponse
	ParseResponse(t, resp, &errResp)
	return errResp
}

type URLBuilder struct {
	base string
}

func NewURLBuilder(baseURL string) *URLBuilder {
	return &URLBuilder{base: baseURL}
}

func (b *URLBuilder) Projects() string {
	return "/projects"
}

func (b *URLBuilder) Project(id string) string {
	return fmt.Sprintf("/projects/%s", id)
}

func (b *URLBuilder) ProjectBoards(projectID string) string {
	return fmt.Sprintf("/projects/%s/boards", projectID)
}

func (b *URLBuilder) ProjectEpics(projectID string) string {
	return fmt.Sprintf("/projects/%s/epics", projectID)
}

func (b *URLBuilder) ProjectLabels(projectID string) string {
	return fmt.Sprintf("/projects/%s/labels", projectID)
}

func (b *URLBuilder) Board(boardID string) string {
	return fmt.Sprintf("/boards/%s", boardID)
}

func (b *URLBuilder) BoardTasks(boardID string) string {
	return fmt.Sprintf("/boards/%s/tasks", boardID)
}

func (b *URLBuilder) Task(taskID string) string {
	return fmt.Sprintf("/tasks/%s", taskID)
}

func (b *URLBuilder) TaskMove(taskID string) string {
	return fmt.Sprintf("/tasks/%s/move", taskID)
}

func (b *URLBuilder) TaskLabels(taskID string) string {
	return fmt.Sprintf("/tasks/%s/labels", taskID)
}

func (b *URLBuilder) Epic(epicID string) string {
	return fmt.Sprintf("/epics/%s", epicID)
}

func (b *URLBuilder) Members() string {
	return "/members"
}

func (b *URLBuilder) Member(id string) string {
	return fmt.Sprintf("/members/%s", id)
}

func (b *URLBuilder) AdminStats() string {
	return "/admin/stats"
}

func (b *URLBuilder) AdminSettings() string {
	return "/admin/settings"
}

func (b *URLBuilder) AuthRegister() string {
	return "/auth/register"
}

func (b *URLBuilder) AuthLogin() string {
	return "/auth/login"
}

func (b *URLBuilder) AuthLogout() string {
	return "/auth/logout"
}

func (b *URLBuilder) AuthMe() string {
	return "/auth/me"
}
