package repository

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"time"
)

// ActiveUserCount struct to hold the result of count queries
type ActiveUserCount struct {
	Count int `db:"count"`
}

// make a struct for user
type User struct {
	App core.App
}

// NewUserRepository creates a new User repository
func NewUserRepository(app core.App) *User {
	return &User{
		App: app,
	}
}

// get statistic user active daily
// ambil count aja sih , jadi bisa user_history itu diambil grub by user_id count gitu aja
// GetDailyActiveUsers returns the count of daily active users
func (r *User) GetDailyActiveUsers(date time.Time) (int, error) {

	// If date is zero value, use current date
	targetDate := date
	if targetDate.IsZero() {
		targetDate = time.Now()
	}

	// Format date as YYYY-MM-DD for SQL
	dateStr := targetDate.Format("2006-01-02")

	var result ActiveUserCount
	query := "SELECT COUNT(DISTINCT user_id) AS count FROM user_history WHERE date(created) = date({:date})"

	err := r.App.DB().
		NewQuery(query).
		Bind(dbx.Params{
			"date": dateStr,
		}).
		One(&result)

	if err != nil {
		return 0, err
	}

	return result.Count, nil
}
