package storage

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/GrosbergKirr/WeatherApp/internal/models"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func (s *Storage) RegisterUser(log *slog.Logger, user models.User) (int, error) {
	const path = "/storage/auth"
	_, err := s.UserExistence(log, user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("Login is free")
		} else {
			log.Error("User checking error")
			return http.StatusInternalServerError, err
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("internal error (password)", slog.String("path", path))
		return http.StatusInternalServerError, err
	}
	user.Password = string(hashedPassword)

	tx, err := s.Db.Beginx()
	if err != nil {
		log.Error("Could not begin transaction", slog.Any("err", err))
		return http.StatusInternalServerError, err
	}
	query := "INSERT INTO users (login, password) VALUES(:login, :password)"
	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		log.Error("Failed to prepare cities transactions to DB", slog.String("path", path))
		return http.StatusInternalServerError, err
	}
	defer func(stmt *sqlx.NamedStmt) {
		err := stmt.Close()
		if err != nil {
			log.Error("failed to close statement", slog.String("path", path))
		}
	}(stmt)
	_, err = stmt.Exec(user)
	if err != nil {
		log.Error("Failed to insert into DB", slog.String("path", path))
		return http.StatusInternalServerError, err
	}
	if err = tx.Commit(); err != nil {
		log.Error("Failed to commit transaction", slog.String("path", path))
		return http.StatusInternalServerError, err
	}
	log.Debug("Save new user data to db successfully", slog.String("user", user.Login))
	return http.StatusCreated, nil
}
func (s *Storage) Login(log *slog.Logger, user models.User) (int, error) {
	const path = "/storage/auth"
	userData, err := s.UserExistence(log, user.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Invalid login or password")
			return http.StatusUnauthorized, err
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
			if err != nil {
				log.Error("Invalid login or password")
				return http.StatusUnauthorized, err
			}
		}
	}
	log.Info("Login & password is valid")
	return http.StatusAccepted, nil
}
func (s *Storage) UserExistence(log *slog.Logger, loginReq string) (models.User, error) {
	const path = "/storage/auth"
	var user models.User
	queryLogin := "SELECT login, password FROM users where login = $1"
	if err := s.Db.Get(&user, queryLogin, loginReq); err != nil {
		log.Error("Checking user", slog.String("path", path))
		return models.User{}, err
	}
	return user, nil
}
