package users

import (
	"bookstore_users-api/datasources/mysql/users_db"
	"bookstore_users-api/logger"
	"bookstore_users-api/utils/date_utils"
	"bookstore_users-api/utils/mysql_utils"
	"github.com/IrinaVolosevich/bookstore_utils-go/rest_errors"
	"fmt"
	"strings"
)

const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, password, status) VALUES(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE id=?"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=?;"
	queryDeleteUser             = "DELETE FROM users WHERE id=?"
	queryFindByStatus           = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? and password=? and status=?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Id)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password); getErr != nil {
		logger.Error("error when trying to get user by Id", getErr)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	user.DateCreated = date_utils.GetNowString()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Password, user.Status)
	if saveErr != nil {
		logger.Error("error when trying to save user statement", saveErr)
		return rest_errors.NewInternalServerError("database error", err)
	}

	userId, err := insertResult.LastInsertId()

	if err != nil {
		logger.Error("error when trying to get last insert Id", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	user.Id = userId
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to update user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	_, updateErr := stmt.Exec(user.FirstName, user.LastName, user.Email)
	if updateErr != nil {
		logger.Error("error when trying to update user statement", updateErr)
		return rest_errors.NewInternalServerError("database error", err)
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to delete user statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	if _, deleteErr := stmt.Exec(user.Id); deleteErr != nil {
		logger.Error("error when trying to delete user statement", deleteErr)
		return rest_errors.NewInternalServerError("database error", deleteErr)
	}

	return nil
}

func (user *User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := users_db.Client.Prepare(queryFindByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("database error", err)
	}
	defer rows.Close()

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan row to find users by status statement", err)
			return nil, rest_errors.NewInternalServerError("database error", err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by email and password statement", err)
		return rest_errors.NewInternalServerError("database error", err)
	}

	defer stmt.Close()

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status, &user.Password); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return rest_errors.NewNotFoundError("no user found")
		}

		logger.Error("error when trying to get user by email and password", getErr)
		return rest_errors.NewInternalServerError("database error", getErr)
	}

	return nil
}