package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	cleanDB(t)

	resp := makeRequest(t, http.MethodPost, "/api/v1/users/register", map[string]string{
		"username": "presenter1",
		"email":    "presenter1@example.com",
		"password": "password123",
		"role":     "presenter",
	}, "")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEmpty(t, extractCookieToken(resp))
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	user := data["user"].(map[string]interface{})
	assert.Equal(t, "presenter1", user["username"])
	assert.Equal(t, "presenter", user["role"])
}

func TestRegister_DuplicateUsername(t *testing.T) {
	cleanDB(t)

	payload := map[string]string{
		"username": "dupuser",
		"email":    "dupuser@example.com",
		"password": "password123",
		"role":     "presenter",
	}
	makeRequest(t, http.MethodPost, "/api/v1/users/register", payload, "")

	resp2 := makeRequest(t, http.MethodPost, "/api/v1/users/register", map[string]string{
		"username": "dupuser",
		"email":    "other@example.com",
		"password": "password123",
		"role":     "presenter",
	}, "")

	assert.Equal(t, http.StatusConflict, resp2.StatusCode)
}

func TestRegister_InvalidRequest(t *testing.T) {
	cleanDB(t)

	// missing required fields
	resp := makeRequest(t, http.MethodPost, "/api/v1/users/register", map[string]string{
		"username": "ab", // too short (min 3)
		"email":    "not-an-email",
		"password": "short",
		"role":     "invalid-role",
	}, "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLogin_Success(t *testing.T) {
	cleanDB(t)

	registerUser(t, "loginuser", "loginuser@example.com", "password123", "presenter")

	resp := makeRequest(t, http.MethodPost, "/api/v1/users/login", map[string]string{
		"username": "loginuser",
		"password": "password123",
	}, "")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, extractCookieToken(resp))
}

func TestLogin_WrongPassword(t *testing.T) {
	cleanDB(t)

	registerUser(t, "loginuser2", "loginuser2@example.com", "password123", "presenter")

	resp := makeRequest(t, http.MethodPost, "/api/v1/users/login", map[string]string{
		"username": "loginuser2",
		"password": "wrongpassword",
	}, "")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestAnonymous_Success(t *testing.T) {
	cleanDB(t)

	// presenter creates a room first
	presenterToken := registerUser(t, "presenteranon", "presenteranon@example.com", "password123", "presenter")
	room, _ := createRoom(t, presenterToken, "Anon Test Room")
	roomCode := room["room_code"].(string)

	// anonymous user joins via anon endpoint
	resp := makeRequest(t, http.MethodPost, "/api/v1/users/anonymous", map[string]string{
		"room_code":    roomCode,
		"display_name": "Anonymous Guest",
	}, "")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, extractCookieToken(resp))
}

func TestLogout_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "logoutuser", "logoutuser@example.com", "password123", "presenter")

	resp := makeRequest(t, http.MethodPost, "/api/v1/users/logout", nil, token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// using the same token after logout should fail
	resp2 := makeRequest(t, http.MethodPost, "/api/v1/rooms", map[string]string{"title": "Test Room"}, token)
	assert.Equal(t, http.StatusUnauthorized, resp2.StatusCode)
}
