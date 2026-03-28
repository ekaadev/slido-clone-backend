package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRoom_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "roomhost", "roomhost@example.com", "password123", "presenter")

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms", map[string]string{
		"title": "My Test Room",
	}, token)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.NotEmpty(t, extractCookieToken(resp))
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	room := data["room"].(map[string]interface{})
	assert.Equal(t, "My Test Room", room["title"])
	assert.NotEmpty(t, room["room_code"])
	assert.Equal(t, "active", room["status"])
}

func TestCreateRoom_Unauthorized(t *testing.T) {
	cleanDB(t)

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms", map[string]string{
		"title": "Unauthorized Room",
	}, "")

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetRoom_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "getroomer", "getroomer@example.com", "password123", "presenter")
	room, _ := createRoom(t, token, "Get Room Test")
	roomCode := room["room_code"].(string)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+roomCode, nil, "")

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	assert.Equal(t, roomCode, data["room_code"])
	assert.Equal(t, "active", data["status"])
}

func TestCloseRoom_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "closehost", "closehost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Close Room Test")
	roomID := room["id"].(float64)

	resp := makeRequest(t, http.MethodPatch, "/api/v1/rooms/"+formatID(roomID)+"/close",
		map[string]string{"status": "closed"}, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "closed", data["status"])
}

func TestCloseRoom_NotOwner(t *testing.T) {
	cleanDB(t)

	// presenter creates room
	presenterToken := registerUser(t, "presenterclose", "presenterclose@example.com", "password123", "presenter")
	room, _ := createRoom(t, presenterToken, "Not Owner Close Test")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	// another user tries to close
	otherToken := registerUser(t, "otherclose", "otherclose@example.com", "password123", "presenter")
	_, otherRoomToken := joinRoom(t, otherToken, roomCode)

	resp := makeRequest(t, http.MethodPatch, "/api/v1/rooms/"+formatID(roomID)+"/close",
		map[string]string{"status": "closed"}, otherRoomToken)

	// system returns 404 (room not found for that presenter_id) rather than 403
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteRoom_MustCloseFirst(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "deleter", "deleter@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Delete Room Test")
	roomID := room["id"].(float64)

	// try to delete an active room — should fail
	resp := makeRequest(t, http.MethodDelete, "/api/v1/rooms/"+formatID(roomID), nil, roomToken)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

