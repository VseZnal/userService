package repository

import (
	"database/sql"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"userService/libs/errors"
	"userService/services/user-service/config"
	user_service "userService/services/user-service/proto/user-service"
)

type Database interface {
	SignIn(user *user_service.User) (*user_service.User, error)
	SignUp(user *user_service.User) (*user_service.User, error)
}

type DatabaseConn struct {
	conn *sql.DB
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
		return nil, err
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
		return nil, errors.HandleDatabaseError(err)
	}

	return userOut, err
}

func NewDatabase() (*DatabaseConn, error) {
	conf := config.UserConfig()

	connStr := conf.PgConnString

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		errors.HandleFatalError(err, "failed to connect to postgres")
	}

	err = db.Ping()

	if err != nil {
		errors.HandleFatalError(err, "failed to connect to postgres")
	}

	return &DatabaseConn{conn: db}, err
}
