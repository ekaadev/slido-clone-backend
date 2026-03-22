package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"slido-clone-backend/internal/config"
)

var (
	testApp   *fiber.App
	testDB    *gorm.DB
	testRedis *redis.Client
	testCfg   *viper.Viper
)

// migrationsPath is relative to the test/integration/ directory at test runtime
const migrationsPath = "../../db/migrations"

func TestMain(m *testing.M) {
	testCfg = loadTestConfig()
	testDB = connectTestDB(testCfg)
	runMigrations(testCfg)

	testRedis = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", testCfg.GetString("REDIS_HOST"), testCfg.GetInt("REDIS_PORT")),
		DB:   testCfg.GetInt("REDIS_DB"),
	})

	log := logrus.New()
	log.SetLevel(logrus.PanicLevel)
	validate := validator.New()

	testApp = fiber.New(fiber.Config{
		ErrorHandler: config.NewErrorHandler(),
	})

	config.Bootstrap(&config.BootstrapConfig{
		DB:        testDB,
		App:       testApp,
		Redis:     testRedis,
		Log:       log,
		Validator: validate,
		Config:    testCfg,
	})

	code := m.Run()

	teardown()
	os.Exit(code)
}

func loadTestConfig() *viper.Viper {
	v := viper.New()

	v.SetConfigFile("../../config.json")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("failed to read config.json: %v", err))
	}

	v.SetConfigFile("../../.env.test")
	v.SetConfigType("env")
	v.AutomaticEnv()
	if err := v.MergeInConfig(); err != nil {
		panic(fmt.Sprintf("failed to read .env.test: %v", err))
	}

	return v
}

func connectTestDB(v *viper.Viper) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		v.GetString("DATABASE_HOST"),
		v.GetString("DATABASE_USERNAME"),
		v.GetString("DATABASE_PASSWORD"),
		v.GetString("DATABASE_NAME"),
		v.GetInt("DATABASE_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect test database: %v", err))
	}
	return db
}

func buildMigrateDSN(v *viper.Viper) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		v.GetString("DATABASE_USERNAME"),
		v.GetString("DATABASE_PASSWORD"),
		v.GetString("DATABASE_HOST"),
		v.GetInt("DATABASE_PORT"),
		v.GetString("DATABASE_NAME"),
	)
}

func runMigrations(v *viper.Viper) {
	m, err := migrate.New("file://"+migrationsPath, buildMigrateDSN(v))
	if err != nil {
		panic(fmt.Sprintf("failed to init migrations: %v", err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}
}

func teardown() {
	m, err := migrate.New("file://"+migrationsPath, buildMigrateDSN(testCfg))
	if err != nil {
		return
	}
	_ = m.Down()
}

// cleanDB truncates all tables in reverse dependency order.
// Should be called at the start of each test function for isolation.
func cleanDB(t *testing.T) {
	t.Helper()
	tables := []string{
		"xp_transactions",
		"poll_responses",
		"poll_options",
		"polls",
		"votes",
		"questions",
		"messages",
		"participants",
		"rooms",
		"users",
	}
	for _, table := range tables {
		if err := testDB.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			t.Logf("warning: failed to truncate %s: %v", table, err)
		}
	}
}

// makeRequest builds an HTTP request and calls testApp.Test().
// body is marshalled to JSON if non-nil.
func makeRequest(t *testing.T, method, path string, body interface{}, token string) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := testApp.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to execute test request: %v", err)
	}
	return resp
}

// readBody reads and returns the response body as a map.
// The response body can only be read once — do not call resp.Body after this.
func readBody(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	defer resp.Body.Close()
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return result
}

// registerUser registers a new user and returns the base JWT token.
func registerUser(t *testing.T, username, email, password, role string) string {
	t.Helper()
	body := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
		"role":     role,
	}
	resp := makeRequest(t, http.MethodPost, "/api/v1/users/register", body, "")
	if resp.StatusCode != http.StatusCreated {
		result := readBody(t, resp)
		t.Fatalf("register failed: status=%d body=%v", resp.StatusCode, result)
	}
	result := readBody(t, resp)
	data := result["data"].(map[string]interface{})
	return data["token"].(string)
}

// loginUser logs in and returns the JWT token.
func loginUser(t *testing.T, username, password string) string {
	t.Helper()
	body := map[string]string{
		"username": username,
		"password": password,
	}
	resp := makeRequest(t, http.MethodPost, "/api/v1/users/login", body, "")
	if resp.StatusCode != http.StatusOK {
		result := readBody(t, resp)
		t.Fatalf("login failed: status=%d body=%v", resp.StatusCode, result)
	}
	result := readBody(t, resp)
	data := result["data"].(map[string]interface{})
	return data["token"].(string)
}

// createRoom creates a room with the given title and returns the room map and room-scoped token.
func createRoom(t *testing.T, token, title string) (map[string]interface{}, string) {
	t.Helper()
	body := map[string]string{"title": title}
	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms", body, token)
	if resp.StatusCode != http.StatusCreated {
		result := readBody(t, resp)
		t.Fatalf("create room failed: status=%d body=%v", resp.StatusCode, result)
	}
	result := readBody(t, resp)
	data := result["data"].(map[string]interface{})
	room := data["room"].(map[string]interface{})
	roomToken := data["token"].(string)
	return room, roomToken
}

// formatID converts a float64 ID from JSON unmarshalling to a URL path string.
func formatID(id float64) string {
	return strconv.FormatInt(int64(id), 10)
}

// joinRoom joins a room by room code and returns the participant map and room-scoped token.
func joinRoom(t *testing.T, token, roomCode string) (map[string]interface{}, string) {
	t.Helper()
	// send an empty JSON body so Fiber's BodyParser can detect the content type
	resp := makeRequest(t, http.MethodPost, "/api/v1/rooms/"+roomCode+"/join", map[string]interface{}{}, token)
	if resp.StatusCode != http.StatusCreated {
		result := readBody(t, resp)
		t.Fatalf("join room failed: status=%d body=%v", resp.StatusCode, result)
	}
	result := readBody(t, resp)
	data := result["data"].(map[string]interface{})
	participant := data["participant"].(map[string]interface{})
	roomToken := data["token"].(string)
	return participant, roomToken
}
