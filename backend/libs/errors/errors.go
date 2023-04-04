package errors

import (
	"database/sql"
	"log"

	statusErr "google.golang.org/grpc/status"
)

func HandleDatabaseError(err error) error {
	log.Println(err)

	switch {
	case err == sql.ErrNoRows:
		return statusErr.Error(404, "Record not found")
	case err == sql.ErrTxDone:
		return statusErr.Error(500, "Internal server error has occured")
	case err == sql.ErrConnDone:
		return statusErr.Error(500, "Internal server error has occured")
	default:
		return statusErr.Error(500, "Internal server error has occured")
	}
}

func HandleFatalError(err error, message string) {
	log.Println(err)

	log.Fatal(message)
}

func LogError(err error) error {
	log.Println(err)

	return err
}
