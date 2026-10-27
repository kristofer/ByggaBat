package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type problemDefinition struct {
	Title    string `json:"title"`
	Level    string `json:"level"`
	Week     int    `json:"week"`
	Skill    string `json:"skill"`
	Prompt   string `json:"prompt"`
	Code     string `json:"code"`
	Hint     string `json:"hint"`
	Solution string `json:"solution"`
}

func createSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS exercises (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL UNIQUE,
			level TEXT NOT NULL,
			week INTEGER NOT NULL,
			skill TEXT NOT NULL,
			prompt TEXT NOT NULL,
			code TEXT NOT NULL,
			hint TEXT NOT NULL,
			solution TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS attempts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			student TEXT NOT NULL,
			exercise_id INTEGER NOT NULL REFERENCES exercises(id),
			mode TEXT NOT NULL CHECK (mode IN ('diagnose', 'tutor', 'review')),
			response TEXT NOT NULL,
			feedback TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`)
	return err
}

func seedExercises(db *sql.DB) error {
	seeds := []Exercise{
		{
			Title: "An empty ratings list", Level: "Beginner", Week: 1, Skill: "Lists and edge cases",
			Prompt: "This function should return the average rating, or 0 when there are no ratings. Trace the empty-list case and explain the bug.",
			Code: `def average(ratings):
    total = sum(ratings)
    return total / len(ratings)

print(average([4, 5, 3]))
print(average([]))`,
			Hint:     "What is the length of an empty list, and what happens when you divide by it?",
			Solution: "The empty-list case divides by zero. Check whether ratings is empty before dividing and return 0 for that case.",
		},
		{
			Title: "Count the even numbers", Level: "Beginner", Week: 2, Skill: "Conditionals and modulo",
			Prompt: "The function should count even values in a list. Work through the condition for 2 and 3 and identify why its result is wrong.",
			Code: `def count_even(numbers):
    count = 0
    for number in numbers:
        if number % 2 == 1:
            count += 1
    return count

print(count_even([1, 2, 3, 4, 6]))`,
			Hint:     "An even number has no remainder after division by two.",
			Solution: "The condition counts values with remainder 1, which are odd numbers. It should test for a remainder of 0.",
		},
		{
			Title: "Apply the member discount", Level: "Beginner", Week: 3, Skill: "Functions and comparisons",
			Prompt: "Members spending $100 or more should receive 10% off. Check the exact boundary value and explain what happens.",
			Code: `def total_after_discount(total, is_member):
    if is_member and total > 100:
        return total * 0.9
    return total

print(total_after_discount(100, True))
print(total_after_discount(120, True))`,
			Hint:     "Compare the written requirement with the comparison operator at exactly 100.",
			Solution: "The strict greater-than check excludes a member spending exactly $100. The requirement says $100 or more, so the comparison should include the boundary.",
		},
		{
			Title: "MiniYelp review summaries", Level: "Medium", Week: 5, Skill: "Functions, dictionaries, and empty data",
			Prompt: "Review this small review-summary module. A business with no reviews should have an average of 0, and adding a review should only affect that business. Explain the failure and how you would verify the fix.",
			Code: `# reviews.py
def add_review(businesses, name, rating):
    businesses.setdefault(name, [])
    businesses[name].append(rating)

def average_rating(businesses, name):
    ratings = businesses.get(name, [])
    return sum(ratings) / len(ratings)

def top_businesses(businesses):
    ranked = []
    for name in businesses:
        ranked.append((name, average_rating(businesses, name)))
    ranked.sort(key=lambda row: row[1], reverse=True)
    return ranked

places = {}
add_review(places, "Harbor Cafe", 5)
add_review(places, "Harbor Cafe", 3)
print(average_rating(places, "New Bakery"))
print(top_businesses(places))`,
			Hint:     "Trace the empty business through average_rating. Then think about which list append modifies.",
			Solution: "average_rating divides by zero for a missing business. The prompt also invites verifying that businesses[name] is only changed for the requested business; a regression test should cover both a missing name and two independent businesses.",
		},
		{
			Title: "Warehouse inventory reports", Level: "Medium", Week: 6, Skill: "Classes, iteration, and mutable data",
			Prompt: "The inventory manager should remove every discontinued product, preserve the quantities of active products, and report low stock without changing inventory. Explain why calling discontinue leaves discontinued products in stock.",
			Code: `# inventory.py
class Inventory:
    def __init__(self, quantities):
        self.quantities = quantities

    def receive(self, product, amount):
        self.quantities[product] = self.quantities.get(product, 0) + amount

    def discontinue(self, products):
        for product in products:
            if product not in self.quantities:
                self.quantities.pop(product, None)

    def low_stock(self, threshold):
        return [name for name, count in self.quantities.items()
                if count < threshold]

stock = Inventory({"tea": 4, "mugs": 2, "beans": 8})
stock.discontinue(["tea", "mugs"])
print(stock.quantities)`,
			Hint:     "Check whether the condition in discontinue matches the stated meaning of discontinued.",
			Solution: "The membership condition is reversed: it removes a product only when it is already absent, so the discontinued products stay in stock. Remove the product when it is present and test that active stock remains.",
		},
		{
			Title: "Library overdue notices", Level: "Medium", Week: 7, Skill: "Dates, filtering, and aggregation",
			Prompt: "A loan is overdue when its due date is before today and it has not been returned. Inspect how the module handles today's date and returned loans; describe a test that exposes its error.",
			Code: `# library.py
from datetime import date, timedelta

def overdue_loans(loans, today):
    overdue = []
    for loan in loans:
        due = date.fromisoformat(loan["due"])
        if due <= today:
            overdue.append(loan)
    return overdue

def notice_for(loan, today):
    due = date.fromisoformat(loan["due"])
    days = (today - due).days
    return f'{loan["title"]}: {days} days overdue'

today = date(2025, 4, 10)
loans = [
    {"title": "Python Basics", "due": "2025-04-10", "returned": False},
    {"title": "Clean Code", "due": "2025-04-01", "returned": True},
]
for loan in overdue_loans(loans, today):
    print(notice_for(loan, today))`,
			Hint:     "Check the words 'before today' and 'not returned' against the filtering condition.",
			Solution: "The condition includes loans due today and never checks returned, so it incorrectly reports both sample loans as overdue. Use a strict date comparison and filter returned loans out; test today, yesterday, and a returned overdue loan.",
		},
		{
			Title: "Monthly expense summaries", Level: "Medium", Week: 8, Skill: "Modules, dictionaries, and aggregation",
			Prompt: "The report should total each category for the requested month, excluding refunds and expenses from other months. Explain why the food total is incorrect and describe a test for it.",
			Code: `# expenses.py
def monthly_totals(expenses, month):
    totals = {}
    for expense in expenses:
        if expense["date"][:7] == month:
            category = expense["category"]
            totals[category] = totals.get(category, 0) + expense["amount"]
    return totals

def report(expenses, month, categories):
    totals = monthly_totals(expenses, month)
    for category in categories:
        totals.setdefault(category, 0)
    return totals

records = [
    {"date": "2025-04-02", "category": "travel", "amount": 40},
    {"date": "2025-04-08", "category": "food", "amount": -12},
    {"date": "2025-03-28", "category": "food", "amount": 20},
]
categories = ["travel", "food", "supplies"]
print(report(records, "2025-04", categories))
print(categories)`,
			Hint:     "Trace how each record's amount contributes to its category total.",
			Solution: "The aggregation includes negative amounts, so the April food refund reduces the reported spending. The prompt requires excluding refunds; filter them out and test that the food total excludes the -12 record.",
		},
		{
			Title: "Course enrollment tracker", Level: "Medium", Week: 8, Skill: "Classes, sets, and state changes",
			Prompt: "Students may enroll in multiple courses, but a course's enrollment should count each student once. A dropped course must not be included in the student's active courses. Find the state-management bug and explain its impact.",
			Code: `# enrollment.py
class Student:
    def __init__(self, name):
        self.name = name
        self.courses = set()

    def enroll(self, course):
        self.courses.add(course)
        course.students.append(self)

    def drop(self, course):
        self.courses.discard(course)

class Course:
    def __init__(self, title):
        self.title = title
        self.students = []

    def enrollment_count(self):
        return len(self.students)

python = Course("Python")
maya = Student("Maya")
maya.enroll(python)
maya.enroll(python)
maya.drop(python)
print(python.enrollment_count(), maya.courses)`,
			Hint:     "Trace both sides of the student-course relationship for each operation.",
			Solution: "enroll appends the same student every time, so the course's list counts duplicates. drop updates only the student's set and leaves the course list stale. Keep both sides in sync and prevent duplicate enrollment, then test repeated enroll and drop.",
		},
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, exercise := range seeds {
		_, err := tx.Exec(`INSERT INTO exercises
			(title, level, week, skill, prompt, code, hint, solution)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(title) DO NOTHING`,
			exercise.Title, exercise.Level, exercise.Week, exercise.Skill,
			exercise.Prompt, exercise.Code, exercise.Hint, exercise.Solution)
		if err != nil {
			return fmt.Errorf("seed %q: %w", exercise.Title, err)
		}
	}
	return tx.Commit()
}

func decodeProblems(reader io.Reader) ([]problemDefinition, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	var problems []problemDefinition
	if err := decoder.Decode(&problems); err != nil {
		return nil, fmt.Errorf("decode problem JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return nil, fmt.Errorf("problem JSON must contain a single array")
	}
	if len(problems) == 0 {
		return nil, fmt.Errorf("problem JSON must contain at least one problem")
	}
	for i := range problems {
		problem := &problems[i]
		problem.Title = strings.TrimSpace(problem.Title)
		problem.Level = strings.TrimSpace(problem.Level)
		problem.Skill = strings.TrimSpace(problem.Skill)
		if problem.Title == "" || problem.Level == "" || problem.Week < 1 ||
			problem.Skill == "" || strings.TrimSpace(problem.Prompt) == "" ||
			strings.TrimSpace(problem.Code) == "" || strings.TrimSpace(problem.Hint) == "" ||
			strings.TrimSpace(problem.Solution) == "" {
			return nil, fmt.Errorf("problem %d must include title, level, a positive week, skill, prompt, code, hint, and solution", i+1)
		}
	}
	return problems, nil
}

func importProblems(db *sql.DB, problems []problemDefinition) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	inserted := 0
	for _, problem := range problems {
		result, err := tx.Exec(`INSERT INTO exercises
			(title, level, week, skill, prompt, code, hint, solution)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(title) DO NOTHING`,
			problem.Title, problem.Level, problem.Week, problem.Skill,
			problem.Prompt, problem.Code, problem.Hint, problem.Solution)
		if err != nil {
			return 0, fmt.Errorf("import %q: %w", problem.Title, err)
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		inserted += int(rows)
	}
	return inserted, tx.Commit()
}
