package main

import (
	"bytes"
	"encoding/json"
	"linkshortener/internal/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	godotenv.Load()
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}
	return db
}

func removeDbData(db *gorm.DB) {
	db.Exec("DELETE FROM links")
	db.Exec("DELETE FROM stats")
	db.Exec("DELETE FROM users")
}

func TestRegisterSuccess(t *testing.T) {
	db := initDb()
	defer removeDbData(db)

	ts := httptest.NewServer(appInit())
	defer ts.Close()

	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@test.com",
		Password: "password123!",
		Name:     "test",
	})

	res, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status Created, got %d", res.StatusCode)
	}

	response := auth.RegisterResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.Email != "test@test.com" {
		t.Fatalf("expected email test@test.com, got %s", response.Email)
	}
}

type testCase struct {
	request  auth.RegisterRequest
	expected int
}

func TestRegisterFailed(t *testing.T) {
	db := initDb()
	defer removeDbData(db)

	ts := httptest.NewServer(appInit())
	defer ts.Close()

	testCases := []testCase{
		{
			request: auth.RegisterRequest{
				Email:    "test@test.com",
				Password: "123",
			},
			expected: http.StatusBadRequest,
		},
		{
			request: auth.RegisterRequest{
				Email:    "a@test.com",
				Password: "password123!",
				Name:     "a",
			},
			expected: http.StatusBadRequest,
		},
		{
			request: auth.RegisterRequest{
				Email:    "valid@test.com",
				Password: "password123!",
				Name:     "validuser",
			},
			expected: http.StatusCreated,
		},
	}

	for _, testCase := range testCases {
		data, _ := json.Marshal(&testCase.request)

		res, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBuffer(data))
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != testCase.expected {
			t.Fatalf("expected status %d, got %d", testCase.expected, res.StatusCode)
		}
	}
}

func TestLoginSuccess(t *testing.T) {
	db := initDb()
	defer removeDbData(db)

	ts := httptest.NewServer(appInit())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "password123!",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status OK, got %d", res.StatusCode)
	}

	response := auth.LoginResponse{}
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if response.AccessToken == "" {
		t.Fatalf("expected access token, got %s", response.AccessToken)
	}

	if response.RefreshToken == "" {
		t.Fatalf("expected refresh token, got %s", response.RefreshToken)
	}
}

func TestLoginFailed(t *testing.T) {
	db := initDb()
	defer removeDbData(db)

	ts := httptest.NewServer(appInit())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@test.com",
		Password: "wrong",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		t.Fatalf("expected status not OK, got %d", res.StatusCode)
	}
}
