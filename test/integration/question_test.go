package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// submitQuestion is a helper to post a question and return the question map.
func submitQuestion(t *testing.T, roomID float64, content, token string) map[string]interface{} {
	t.Helper()
	path := "/api/v1/rooms/" + formatID(roomID) + "/questions"
	resp := makeRequest(t, http.MethodPost, path, map[string]string{"content": content}, token)
	if resp.StatusCode != http.StatusCreated {
		body := readBody(t, resp)
		t.Fatalf("submit question failed: status=%d body=%v", resp.StatusCode, body)
	}
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	return data["question"].(map[string]interface{})
}

func TestSubmitQuestion_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "qhost", "qhost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Question Room")
	roomID := room["id"].(float64)

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/questions",
		map[string]string{"content": "What is the meaning of life?"}, roomToken)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	q := data["question"].(map[string]interface{})
	assert.Equal(t, "What is the meaning of life?", q["content"])
	assert.Equal(t, "pending", q["status"])

	xp := data["xp_earned"]
	assert.NotNil(t, xp)
}

func TestUpvoteQuestion_Success(t *testing.T) {
	cleanDB(t)

	// presenter creates room and asks a question
	presenterToken := registerUser(t, "upvotehost", "upvotehost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Upvote Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	question := submitQuestion(t, roomID, "Upvote me please", presenterRoomToken)
	questionID := question["id"].(float64)

	// another user joins and upvotes
	user2Token := registerUser(t, "upvoteuser2", "upvoteuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	resp := makeRequest(t, http.MethodPost, "/api/v1/questions/"+formatID(questionID)+"/upvote",
		nil, user2RoomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	q := data["question"].(map[string]interface{})
	assert.Equal(t, float64(1), q["upvote_count"])
}

func TestUpvoteQuestion_OwnQuestion(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "ownupvote", "ownupvote@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Own Upvote Room")
	roomID := room["id"].(float64)

	question := submitQuestion(t, roomID, "My own question", roomToken)
	questionID := question["id"].(float64)

	// upvote own question — should fail
	resp := makeRequest(t, http.MethodPost, "/api/v1/questions/"+formatID(questionID)+"/upvote",
		nil, roomToken)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRemoveUpvote_Success(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "rmupvotehost", "rmupvotehost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Remove Upvote Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	question := submitQuestion(t, roomID, "Remove my upvote", presenterRoomToken)
	questionID := question["id"].(float64)

	user2Token := registerUser(t, "rmupvoteuser2", "rmupvoteuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	// upvote first
	makeRequest(t, http.MethodPost, "/api/v1/questions/"+formatID(questionID)+"/upvote",
		nil, user2RoomToken)

	// then remove upvote
	resp := makeRequest(t, http.MethodDelete, "/api/v1/questions/"+formatID(questionID)+"/upvote",
		nil, user2RoomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	q := data["question"].(map[string]interface{})
	assert.Equal(t, float64(0), q["upvote_count"])
}

func TestValidateQuestion_PresenterOnly(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "validatehost", "validatehost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Validate Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	question := submitQuestion(t, roomID, "Validate this question", presenterRoomToken)
	questionID := question["id"].(float64)

	// non-presenter tries to validate
	user2Token := registerUser(t, "validateuser2", "validateuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	resp := makeRequest(t, http.MethodPatch, "/api/v1/questions/"+formatID(questionID)+"/validate",
		map[string]string{"status": "answered"}, user2RoomToken)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// presenter validates — should succeed
	resp2 := makeRequest(t, http.MethodPatch, "/api/v1/questions/"+formatID(questionID)+"/validate",
		map[string]string{"status": "answered"}, presenterRoomToken)

	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}
