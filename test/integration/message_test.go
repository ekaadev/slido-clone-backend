package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendMessage_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "msghost", "msghost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Message Test Room")
	roomID := room["id"].(float64)

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/messages",
		map[string]string{"content": "Hello, room!"}, roomToken)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	assert.Equal(t, "Hello, room!", data["content"])
}

func TestListMessages(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "listmsghost", "listmsghost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "List Messages Room")
	roomID := room["id"].(float64)
	path := "/api/v1/rooms/" + formatID(roomID) + "/messages"

	// send two messages
	makeRequest(t, http.MethodPost, path, map[string]string{"content": "First message"}, roomToken)
	makeRequest(t, http.MethodPost, path, map[string]string{"content": "Second message"}, roomToken)

	resp := makeRequest(t, http.MethodGet, path+"?limit=10", nil, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	messages := data["messages"].([]interface{})
	assert.GreaterOrEqual(t, len(messages), 2)
}
