package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testApp(t *testing.T) (*app, *sql.DB) {
	t.Helper()
	db, err := openDatabase(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		t.Fatal(err)
	}
	a := &app{db: db, templates: tmpl, anthropic: newAnthropicClient("")}
	return a, db
}

func TestSeededExercises(t *testing.T) {
	a, _ := testApp(t)
	exercises, err := listExercises(a.db)
	if err != nil {
		t.Fatal(err)
	}
	beginner, medium := 0, 0
	for _, exercise := range exercises {
		switch exercise.Level {
		case "Beginner":
			beginner++
		case "Medium":
			medium++
		}
	}
	if beginner < 2 || medium != 5 {
		t.Fatalf("seeded %d beginner and %d medium exercises; want at least 2 and exactly 5", beginner, medium)
	}
}

func TestStudentAndInstructorViews(t *testing.T) {
	a, db := testApp(t)
	server := httptest.NewServer(a.routes())
	defer server.Close()

	studentResponse, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer studentResponse.Body.Close()
	studentPage, _ := io.ReadAll(studentResponse.Body)
	if studentResponse.StatusCode != http.StatusOK ||
		!strings.Contains(string(studentPage), "An empty ratings list") ||
		strings.Contains(string(studentPage), "Check whether ratings is empty") {
		t.Fatalf("student page missing exercise or leaked instructor hint")
	}
	exerciseResponse, err := http.Get(server.URL + "/exercises/1")
	if err != nil {
		t.Fatal(err)
	}
	defer exerciseResponse.Body.Close()
	exercisePage, _ := io.ReadAll(exerciseResponse.Body)
	exercise, err := getExercise(db, "1")
	if err != nil {
		t.Fatal(err)
	}
	if exerciseResponse.StatusCode != http.StatusOK ||
		!strings.Contains(string(exercisePage), "def average(ratings):") ||
		strings.Contains(string(exercisePage), "Check whether ratings is empty") ||
		strings.Contains(string(exercisePage), exercise.Solution) {
		t.Fatal("student exercise page is missing its code or leaked instructor hint")
	}
	instructorExercise, err := http.Get(server.URL + "/instructor/exercises/1")
	if err != nil {
		t.Fatal(err)
	}
	defer instructorExercise.Body.Close()
	instructorExercisePage, _ := io.ReadAll(instructorExercise.Body)
	if instructorExercise.StatusCode != http.StatusOK ||
		!strings.Contains(string(instructorExercisePage), exercise.Solution) {
		t.Fatal("instructor exercise page is missing the expected diagnosis")
	}

	instructorResponse, err := http.Get(server.URL + "/instructor")
	if err != nil {
		t.Fatal(err)
	}
	defer instructorResponse.Body.Close()
	instructorPage, _ := io.ReadAll(instructorResponse.Body)
	if instructorResponse.StatusCode != http.StatusOK ||
		!strings.Contains(string(instructorPage), "Generate a flawed variant") ||
		!strings.Contains(string(instructorPage), "Recent student activity") {
		t.Fatal("instructor page is missing its exercise generator or activity view")
	}
}

func TestAttemptUsesAnthropicAndPersists(t *testing.T) {
	anthropic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("unexpected Anthropic request: %s, key header present=%t", r.Method, r.Header.Get("x-api-key") != "")
		}
		var request struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if len(request.Messages) != 1 || !strings.Contains(request.Messages[0].Content, "Assess the student's explanation") {
			t.Error("diagnose mode did not send its mode-specific instruction")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"content":[{"type":"text","text":"Good reasoning. What happens with an empty list?"}]}`)
	}))
	defer anthropic.Close()

	a, db := testApp(t)
	a.anthropic.apiKey = "test-key"
	a.anthropic.endpoint = anthropic.URL
	request := httptest.NewRequest(http.MethodPost, "/attempts", strings.NewReader(
		"exercise_id=1&student=Maya&mode=diagnose&response=It+divides+by+zero"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	a.routes().ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Good reasoning") {
		t.Fatalf("attempt returned %d: %s", response.Code, response.Body.String())
	}
	var student, mode, submission string
	if err := db.QueryRow(`SELECT student, mode, response FROM attempts`).Scan(&student, &mode, &submission); err != nil {
		t.Fatal(err)
	}
	if student != "Maya" || mode != "diagnose" || submission != "It divides by zero" {
		t.Fatalf("stored attempt = (%q, %q, %q)", student, mode, submission)
	}
}

func TestAttemptRejectsInvalidMode(t *testing.T) {
	a, _ := testApp(t)
	request := httptest.NewRequest(http.MethodPost, "/attempts", strings.NewReader(
		"exercise_id=1&mode=answer&response=Here+is+my+answer"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	a.routes().ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("invalid mode returned %d; want 400", response.Code)
	}
}

func TestInstructorCanGenerateStudentExercise(t *testing.T) {
	anthropic := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{"content":[{"type":"text","text":"{\"title\":\"Generated trace\",\"week\":4,\"skill\":\"loops\",\"prompt\":\"Find the flaw.\",\"code\":\"print(1)\",\"hint\":\"Trace it.\",\"solution\":\"Check the loop.\"}"}]}`)
	}))
	defer anthropic.Close()

	a, db := testApp(t)
	a.anthropic.apiKey = "test-key"
	a.anthropic.endpoint = anthropic.URL
	request := httptest.NewRequest(http.MethodPost, "/instructor/generate", strings.NewReader(
		"week=4&skill=loops&reference=for+i+in+range%283%29%3A+print%28i%29"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	a.routes().ServeHTTP(response, request)
	if response.Code != http.StatusSeeOther {
		t.Fatalf("generation returned %d: %s", response.Code, response.Body.String())
	}
	var level string
	if err := db.QueryRow(`SELECT level FROM exercises WHERE title = 'Generated trace'`).Scan(&level); err != nil {
		t.Fatal(err)
	}
	if level != "Instructor generated" {
		t.Fatalf("generated exercise level = %q", level)
	}
	student := httptest.NewRecorder()
	a.routes().ServeHTTP(student, httptest.NewRequest(http.MethodGet, "/", nil))
	if !strings.Contains(student.Body.String(), "Generated trace") {
		t.Fatal("generated exercise did not appear in student practice")
	}
}
