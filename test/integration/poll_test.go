package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// createPoll is a helper to create a poll in the given room.
// Returns the poll map. Uses the provided token (must be room-scoped presenter token).
func createPoll(t *testing.T, roomID float64, question string, options []string, token string) map[string]interface{} {
	t.Helper()
	body := map[string]interface{}{
		"question": question,
		"options":  options,
	}
	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/polls", body, token)
	if resp.StatusCode != http.StatusCreated {
		rb := readBody(t, resp)
		t.Fatalf("create poll failed: status=%d body=%v", resp.StatusCode, rb)
	}
	rb := readBody(t, resp)
	data := rb["data"].(map[string]interface{})
	return data["poll"].(map[string]interface{})
}

func TestCreatePoll_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "pollhost", "pollhost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Poll Room")
	roomID := room["id"].(float64)

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/polls", map[string]interface{}{
		"question": "What is your favorite color?",
		"options":  []string{"Red", "Green", "Blue"},
	}, roomToken)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	poll := data["poll"].(map[string]interface{})
	assert.Equal(t, "What is your favorite color?", poll["question"])
	assert.Equal(t, "active", poll["status"])
	options := poll["options"].([]interface{})
	assert.Len(t, options, 3)
}

func TestCreatePoll_NotPresenter(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "pollhost2", "pollhost2@example.com", "password123", "presenter")
	room, _ := createRoom(t, presenterToken, "Poll Not Presenter")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	user2Token := registerUser(t, "polluser2", "polluser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+formatID(roomID)+"/polls", map[string]interface{}{
		"question": "Unauthorized poll",
		"options":  []string{"Yes", "No"},
	}, user2RoomToken)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestVotePoll_Success(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "votehost", "votehost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Vote Poll Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	poll := createPoll(t, roomID, "Vote question", []string{"Option A", "Option B"}, presenterRoomToken)
	pollID := poll["id"].(float64)
	options := poll["options"].([]interface{})
	optionID := options[0].(map[string]interface{})["id"].(float64)

	user2Token := registerUser(t, "voteuser2", "voteuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	resp := makeRequest(t, http.MethodPost, "/api/v1/polls/"+formatID(pollID)+"/vote",
		map[string]interface{}{"option_id": int(optionID)}, user2RoomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	results := data["updated_results"].(map[string]interface{})
	assert.Equal(t, float64(1), results["total_votes"])
}

func TestVotePoll_DuplicateVote(t *testing.T) {
	cleanDB(t)

	presenterToken := registerUser(t, "dupvotehost", "dupvotehost@example.com", "password123", "presenter")
	room, presenterRoomToken := createRoom(t, presenterToken, "Duplicate Vote Room")
	roomCode := room["room_code"].(string)
	roomID := room["id"].(float64)

	poll := createPoll(t, roomID, "Dup vote question", []string{"Yes", "No"}, presenterRoomToken)
	pollID := poll["id"].(float64)
	options := poll["options"].([]interface{})
	optionID := options[0].(map[string]interface{})["id"].(float64)

	user2Token := registerUser(t, "dupvoteuser2", "dupvoteuser2@example.com", "password123", "presenter")
	_, user2RoomToken := joinRoom(t, user2Token, roomCode)

	// first vote
	makeRequest(t, http.MethodPost, "/api/v1/polls/"+formatID(pollID)+"/vote",
		map[string]interface{}{"option_id": int(optionID)}, user2RoomToken)

	// duplicate vote
	resp := makeRequest(t, http.MethodPost, "/api/v1/polls/"+formatID(pollID)+"/vote",
		map[string]interface{}{"option_id": int(optionID)}, user2RoomToken)

	assert.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestClosePoll_Success(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "closepollhost", "closepollhost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Close Poll Room")
	roomID := room["id"].(float64)

	poll := createPoll(t, roomID, "Close this poll", []string{"A", "B"}, roomToken)
	pollID := poll["id"].(float64)

	resp := makeRequest(t, http.MethodPatch, "/api/v1/polls/"+formatID(pollID)+"/close", nil, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	closedPoll := data["poll"].(map[string]interface{})
	assert.Equal(t, "closed", closedPoll["status"])
}

func TestGetActivePolls(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "activehost", "activehost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Active Polls Room")
	roomID := room["id"].(float64)

	createPoll(t, roomID, "Active poll 1", []string{"A", "B"}, roomToken)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+formatID(roomID)+"/polls/active", nil, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	polls := data["polls"].([]interface{})
	assert.GreaterOrEqual(t, len(polls), 1)
}

func TestGetPollHistory(t *testing.T) {
	cleanDB(t)

	token := registerUser(t, "historyhost", "historyhost@example.com", "password123", "presenter")
	room, roomToken := createRoom(t, token, "Poll History Room")
	roomID := room["id"].(float64)

	createPoll(t, roomID, "History poll 1", []string{"A", "B"}, roomToken)

	resp := makeRequest(t, http.MethodGet, "/api/v1/rooms/"+formatID(roomID)+"/polls?limit=10", nil, roomToken)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body := readBody(t, resp)
	data := body["data"].(map[string]interface{})
	polls := data["polls"].([]interface{})
	assert.GreaterOrEqual(t, len(polls), 1)
}
