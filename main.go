package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Exercise struct {
	ID       int64
	Title    string
	Level    string
	Week     int
	Skill    string
	Prompt   string
	Code     string
	Hint     string
	Solution string
}

type Attempt struct {
	Student  string
	Title    string
	Mode     string
	Response string
	Feedback string
	Created  string
}

type app struct {
	db        *sql.DB
	templates *template.Template
	anthropic *anthropicClient
}

func main() {
	dbPath := env("DATABASE_PATH", "byggabat.db")
	db, err := openDatabase(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatal(err)
	}

	a := &app{
		db:        db,
		templates: tmpl,
		anthropic: newAnthropicClient(os.Getenv("ANTHROPIC_API_KEY")),
	}
	port := env("PORT", "8080")
	address := env("LISTEN_ADDR", "127.0.0.1:"+port)
	log.Printf("ByggaBat is listening on http://%s", address)
	log.Fatal(http.ListenAndServe(address, a.routes()))
}

func openDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, err
	}
	if err := createSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	if err := seedExercises(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func (a *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.studentHome)
	mux.HandleFunc("GET /exercises/{id}", a.exercisePage)
	mux.HandleFunc("POST /attempts", a.submitAttempt)
	mux.HandleFunc("GET /instructor", a.instructorHome)
	mux.HandleFunc("GET /instructor/exercises/{id}", a.instructorExercise)
	mux.HandleFunc("POST /instructor/generate", a.generateExercise)
	mux.HandleFunc("POST /instructor/import", a.importProblems)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}

func (a *app) studentHome(w http.ResponseWriter, r *http.Request) {
	exercises, err := listExercises(a.db)
	if err != nil {
		http.Error(w, "Could not load exercises", http.StatusInternalServerError)
		return
	}
	a.render(w, "student.html", exercises)
}

func (a *app) exercisePage(w http.ResponseWriter, r *http.Request) {
	exercise, err := getExercise(a.db, r.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Could not load exercise", http.StatusInternalServerError)
		return
	}
	a.render(w, "exercise.html", exercise)
}

func (a *app) submitAttempt(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	exercise, err := getExercise(a.db, r.FormValue("exercise_id"))
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Could not load exercise", http.StatusInternalServerError)
		return
	}
	mode := r.FormValue("mode")
	if mode != "diagnose" && mode != "tutor" && mode != "review" {
		http.Error(w, "Choose diagnose, tutor, or code review mode", http.StatusBadRequest)
		return
	}
	response := strings.TrimSpace(r.FormValue("response"))
	if response == "" || len(response) > 8000 {
		http.Error(w, "Your response must be between 1 and 8,000 characters", http.StatusBadRequest)
		return
	}
	student := strings.TrimSpace(r.FormValue("student"))
	if student == "" {
		student = "Student"
	}
	if len(student) > 80 {
		http.Error(w, "Name must be 80 characters or fewer", http.StatusBadRequest)
		return
	}
	feedback, err := a.anthropic.respond(r.Context(), mode, exercise, response)
	if err != nil {
		http.Error(w, "The AI service could not respond: "+err.Error(), http.StatusBadGateway)
		return
	}
	_, err = a.db.Exec(`INSERT INTO attempts (student, exercise_id, mode, response, feedback)
		VALUES (?, ?, ?, ?, ?)`, student, exercise.ID, mode, response, feedback)
	if err != nil {
		http.Error(w, "Could not save your attempt", http.StatusInternalServerError)
		return
	}
	a.render(w, "feedback.html", feedback)
}

func (a *app) instructorHome(w http.ResponseWriter, r *http.Request) {
	exercises, err := listExercises(a.db)
	if err != nil {
		http.Error(w, "Could not load exercises", http.StatusInternalServerError)
		return
	}
	rows, err := a.db.Query(`SELECT student, exercises.title, mode, response, feedback,
		strftime('%Y-%m-%d %H:%M', attempts.created_at)
		FROM attempts JOIN exercises ON exercises.id = attempts.exercise_id
		ORDER BY attempts.created_at DESC LIMIT 50`)
	if err != nil {
		http.Error(w, "Could not load student activity", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	attempts := make([]Attempt, 0)
	for rows.Next() {
		var attempt Attempt
		if err := rows.Scan(&attempt.Student, &attempt.Title, &attempt.Mode, &attempt.Response, &attempt.Feedback, &attempt.Created); err != nil {
			http.Error(w, "Could not load student activity", http.StatusInternalServerError)
			return
		}
		attempts = append(attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Could not load student activity", http.StatusInternalServerError)
		return
	}
	imported, _ := strconv.Atoi(r.URL.Query().Get("imported"))
	a.render(w, "instructor.html", struct {
		Exercises        []Exercise
		Attempts         []Attempt
		Imported         int
		ShowImportResult bool
	}{exercises, attempts, max(0, imported), r.URL.Query().Has("imported")})
}

func (a *app) instructorExercise(w http.ResponseWriter, r *http.Request) {
	exercise, err := getExercise(a.db, r.PathValue("id"))
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Could not load exercise", http.StatusInternalServerError)
		return
	}
	a.render(w, "instructor_exercise.html", exercise)
}

func (a *app) generateExercise(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}
	reference := strings.TrimSpace(r.FormValue("reference"))
	week := strings.TrimSpace(r.FormValue("week"))
	skill := strings.TrimSpace(r.FormValue("skill"))
	weekNumber, weekErr := strconv.Atoi(week)
	if reference == "" || len(reference) > 12000 || skill == "" || len(skill) > 120 ||
		weekErr != nil || weekNumber < 1 || weekNumber > 12 {
		http.Error(w, "Provide a reference solution (up to 12,000 characters), week, and skill", http.StatusBadRequest)
		return
	}
	generated, err := a.anthropic.generate(r.Context(), week, skill, reference)
	if err != nil {
		http.Error(w, "Could not generate an exercise: "+err.Error(), http.StatusBadGateway)
		return
	}
	if generated.Week != weekNumber {
		http.Error(w, "Generated exercise did not match the requested curriculum week", http.StatusBadGateway)
		return
	}
	result, err := a.db.Exec(`INSERT INTO exercises
		(title, level, week, skill, prompt, code, hint, solution)
		VALUES (?, 'Instructor generated', ?, ?, ?, ?, ?, ?)`,
		generated.Title, generated.Week, generated.Skill, generated.Prompt,
		generated.Code, generated.Hint, generated.Solution)
	if err != nil {
		http.Error(w, "Could not save generated exercise", http.StatusInternalServerError)
		return
	}
	generated.ID, _ = result.LastInsertId()
	w.Header().Set("HX-Redirect", fmt.Sprintf("/exercises/%d", generated.ID))
	http.Redirect(w, r, fmt.Sprintf("/exercises/%d", generated.ID), http.StatusSeeOther)
}

func (a *app) importProblems(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Choose a valid JSON problem file (up to 1 MB)", http.StatusBadRequest)
		return
	}
	defer file.Close()

	problems, err := decodeProblems(file)
	if err != nil {
		http.Error(w, "Invalid problem JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	inserted, err := importProblems(a.db, problems)
	if err != nil {
		http.Error(w, "Could not import problems", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/instructor?imported=%d", inserted), http.StatusSeeOther)
}

func (a *app) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.templates.ExecuteTemplate(w, name, data); err != nil {
		log.Printf("render %s: %v", name, err)
	}
}

func listExercises(db *sql.DB) ([]Exercise, error) {
	rows, err := db.Query(`SELECT id, title, level, week, skill, prompt, code, hint, solution
		FROM exercises ORDER BY week, id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exercises := make([]Exercise, 0)
	for rows.Next() {
		var exercise Exercise
		if err := rows.Scan(&exercise.ID, &exercise.Title, &exercise.Level, &exercise.Week,
			&exercise.Skill, &exercise.Prompt, &exercise.Code, &exercise.Hint, &exercise.Solution); err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}
	return exercises, rows.Err()
}

func getExercise(db *sql.DB, id string) (Exercise, error) {
	var exercise Exercise
	err := db.QueryRow(`SELECT id, title, level, week, skill, prompt, code, hint, solution
		FROM exercises WHERE id = ?`, id).Scan(&exercise.ID, &exercise.Title, &exercise.Level,
		&exercise.Week, &exercise.Skill, &exercise.Prompt, &exercise.Code, &exercise.Hint, &exercise.Solution)
	return exercise, err
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
