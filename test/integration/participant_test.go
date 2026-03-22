package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinRoom_Success(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "joinhost", "joinhost@example.com", "password123", "presenter")
	room, _ := createRoom(t, presenterToken, "Join Test Room")
	roomCode := room["room_code"].(string)

	otherToken := registerUser(t, "joinuser", "joinuser@example.com", "password123", "presenter")

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+roomCode+"/join", map[string]interface{}{}, otherToken)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	participant := data["participant"].(map[string]interface{})
	assert.Equal(t, "joinuser", participant["display_name"])
}

func TestJoinRoom_Idempotent(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "idempotenthost", "idempotenthost@example.com", "password123", "presenter")
	room, _ := createRoom(t, presenterToken, "Idempotent Room")
	roomCode := room["room_code"].(string)

	otherToken := registerUser(t, "idempotentuser", "idempotentuser@example.com", "password123", "presenter")

	// first join
	resp1 := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+roomCode+"/join", map[string]interface{}{}, otherToken)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	body1 := readBody(t, resp1)
	data1 := body1["data"].(map[string]interface{})
	participant1 := data1["participant"].(map[string]interface{})
	id1 := participant1["id"].(float64)

	// second join — same participant should be returned
	resp2 := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+roomCode+"/join", map[string]interface{}{}, otherToken)
	assert.Equal(t, http.StatusCreated, resp2.StatusCode)
	body2 := readBody(t, resp2)
	data2 := body2["data"].(map[string]interface{})
	participant2 := data2["participant"].(map[string]interface{})
	id2 := participant2["id"].(float64)

	assert.Equal(t, id1, id2, "joining the same room twice should return the same participant ID")
}

func TestListParticipants(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "listhost", "listhost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "List Participants Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	// a second user joins
	user2Token := registerUser(t, "listuser2", "listuser2@example.com", "password123", "presenter")
	joinRoom(t, user2Token, roomCode)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+formatID(roomID)+"/participants?page=1&size=10",
		nil, presenterRoomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	participants := data["participants"].([]interface{})
	// presenter is auto-enrolled + 1 more user
	assert.GreaterOrEqual(t, len(participants), 2)
}
