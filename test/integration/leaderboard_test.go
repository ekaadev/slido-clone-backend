package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeaderboard_AfterActions(t *testing.T) {
	cleanDB(t)

	// presenter creates room and earns XP through a question
	presenterToken := registerUser(t, "lbhost", "lbhost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Leaderboard Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	// submit a question (+10 XP for presenter)
	path := "/api/v1/rooms/" + formatID(roomID) + "/questions"
	makeRequest(t, http.MethodPost, path, map[string]string{"content": "Leaderboard question"}, presenterRoomToken)

	// second user joins and sends a message (+1 XP)
	user2Token := registerUser(t, "lbuser2", "lbuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)
	makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/messages",
		map[string]string{"content": "Hi!"}, user2RoomToken)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+formatID(roomID)+"/leaderboard",
		nil, presenterRoomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	leaderboard := data["leaderboard"].([]interface{})
	assert.GreaterOrEqual(t, len(leaderboard), 1)

	// presenter (10 XP) should rank above user2 (1 XP)
	first := leaderboard[0].(map[string]interface{})
	assert.GreaterOrEqual(t, first["xp_score"].(float64), float64(10))
}

func TestXPTransactions_List(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "xphost", "xphost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "XP Transactions Room")
	roomID := room["id"].(float64)

	// earn XP by submitting a question
	path := "/api/v1/rooms/" + formatID(roomID) + "/questions"
	makeRequest(t, http.MethodPost, path, map[string]string{"content": "XP question"}, roomToken)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+formatID(roomID)+"/xp-transactions",
		nil, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	transactions := data["transactions"].([]interface{})
	assert.GreaterOrEqual(t, len(transactions), 1)
}
