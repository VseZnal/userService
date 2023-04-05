package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	customErr "userService/libs/errors"
	"userService/services/user-service/config"
	user_service "userService/services/user-service/proto/user-service"
)

type Database interface {
	SignIn(user *user_service.User) (*user_service.User, error)
	SignUp(user *user_service.User) (*user_service.User, error)
	CheckUser(user *user_service.User) error
	CheckCode(code string, user *user_service.User) error
	RefreshPassword(user *user_service.User) error
	CreateNotification(user *user_service.User, code string) error
}

type DatabaseConn struct {
	conn *sql.DB
}

func (db DatabaseConn) CreateNotification(user *user_service.User, code string) error {
	query := `
				INSERT INTO notification
				    (username, code)
				VALUES
				    ($1, $2)
				RETURNING username, code
				`

	var usernameCheck, codeCheck string

	err := db.conn.QueryRow(query, &user.Username, code).Scan(&usernameCheck, &codeCheck)
	if err != nil {
		return customErr.HandleDatabaseError(err)
	}

	return err

}

func (db DatabaseConn) RefreshPassword(user *user_service.User) error {
	query := `
		UPDATE users
		SET password = $2
		WHERE username = $1
		RETURNING id, username
			`

	userOut := &user_service.User{}

	err := db.conn.QueryRow(query, &user.Username, &user.Password).Scan(&userOut.Id, &userOut.Username)
	if err != nil {
		return customErr.HandleDatabaseError(err)
	}

	return err
}

func (db DatabaseConn) CheckCode(code string, user *user_service.User) error {
	query := `
			SELECT username, code
			FROM notification 
			WHERE code = $1 and username = $2
			`

	var usernameCheck, codeCheck string

	err := db.conn.QueryRow(query, code, &user.Username).Scan(&usernameCheck, &codeCheck)
	if err != nil {
		return errors.New("code error")
	}

	return err
}

func (db DatabaseConn) CheckUser(user *user_service.User) error {
	query := `
			SELECT id
			FROM users 
			WHERE username = $1
			`

	var usernameCheck int

	err := db.conn.QueryRow(query, &user.Username).Scan(&usernameCheck)
	if err != nil {
		return errors.New("user not found")
	}

	return err
}

func (db DatabaseConn) SignUp(user *user_service.User) (*user_service.User, error) {
	query := `
				INSERT INTO users
				    (username, password)
				VALUES
				    ($1, $2)
				RETURNING id
				`

	rows := db.conn.QueryRow(query, &user.Username, &user.Password)

	userOut := &user_service.User{}

	err := rows.Scan(
		&userOut.Id,
	)
	if err != nil {
		return nil, customErr.HandleDatabaseError(err)
	}

	return userOut, err
}

func (db DatabaseConn) SignIn(user *user_service.User) (*user_service.User, error) {
	query := `
			SELECT id, username
			FROM users 
			WHERE username = $1 and password = $2
			`

	userOut := &user_service.User{}

	rows := db.conn.QueryRow(query, &user.Username, &user.Password)

	err := rows.Scan(
		&userOut.Id,
		&userOut.Username,
	)
	if err != nil {
		return nil, customErr.HandleDatabaseError(err)
	}

	return userOut, err
}

func NewDatabase() (*DatabaseConn, error) {
	conf := config.UserConfig()

	connStr := conf.PgConnString

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		customErr.HandleFatalError(err, "failed to connect to postgres")
	}

	err = db.Ping()

	if err != nil {
		customErr.HandleFatalError(err, "failed to connect to postgres")
	}

	return &DatabaseConn{conn: db}, err
}
