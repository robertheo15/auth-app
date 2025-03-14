package repository

import (
	"auth-app/internal/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type UserRepository interface {
	GetAllUsers() ([]model.User, error)
	GetUserByEmail(email string) (model.User, error)
	SetRedis(ctx context.Context, key, value string) error
	CheckRoleRight(roleID int, permission string) (bool, error)
	CreateUser(name, email, password string, roleID int) (int, error)
	UpdateUser(userID int, name, email, password string, roleID int) error
	DeleteUser(userID int) error
}

type UserRepositoryImpl struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) UserRepository {
	return &UserRepositoryImpl{
		DB:    db,
		Redis: redis,
	}
}

func (r *UserRepositoryImpl) SetRedis(ctx context.Context, key, value string) error {
	err := r.Redis.Set(ctx, key, value, 0).Err()
	if err != nil {

	}

	return nil
}

func (r *UserRepositoryImpl) GetUserByEmail(email string) (model.User, error) {
	var u model.User
	var rolesJSON string

	query := `
	SELECT
	    u.id,
	    u.email,
	    u.password,
	    u.last_access,
	    json_agg(json_build_object('role_id', r.id, 'role_name', r.name)) AS roles
	FROM users u
	JOIN user_roles ur ON u.id = ur.user_id
	JOIN roles r ON ur.role_id = r.id
	WHERE u.email = $1
	GROUP BY u.id, u.email, u.last_access;
	`

	row := r.DB.QueryRow(query, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.LastAccess, &rolesJSON); err != nil {
		return u, err
	}

	if err := json.Unmarshal([]byte(rolesJSON), &u.Roles); err != nil {
		return u, fmt.Errorf("failed to parse roles: %w", err)
	}

	return u, nil
}

func (r *UserRepositoryImpl) GetAllUsers() ([]model.User, error) {
	query := `
	SELECT 
	    u.id, 
	    u.email, 
	    u.last_access, 
	    json_agg(json_build_object('role_id', r.id, 'role_name', r.name)) AS roles
	FROM users u
	LEFT JOIN user_roles ur ON u.id = ur.user_id
	LEFT JOIN roles r ON ur.role_id = r.id
	GROUP BY u.id, u.email, u.last_access;
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var u model.User
		var rolesJSON string

		if err := rows.Scan(&u.ID, &u.Email, &u.LastAccess, &rolesJSON); err != nil {
			return nil, err
		}

		if rolesJSON == "" {
			rolesJSON = "[]"
		}

		if err := json.Unmarshal([]byte(rolesJSON), &u.Roles); err != nil {
			return nil, fmt.Errorf("failed to parse roles: %w", err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (r *UserRepositoryImpl) CheckRoleRight(roleID int, permission string) (bool, error) {
	var hasPermission bool

	query := fmt.Sprintf(`SELECT %s FROM role_rights WHERE role_id = $1`, permission)
	err := r.DB.QueryRow(query, roleID).Scan(&hasPermission)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return hasPermission, nil
}

func (r *UserRepositoryImpl) CreateUser(name, email, password string, roleID int) (int, error) {
	var userID int

	query := `
	INSERT INTO users (name, email, password)
	VALUES ($1, $2, $3)
	RETURNING id;
	`

	err := r.DB.QueryRow(query, name, email, password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	_, err = r.DB.Exec(`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, userID, roleID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *UserRepositoryImpl) UpdateUser(userID int, name, email, password string, roleID int) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password = $3
	WHERE id = $4;
	`
	_, err := r.DB.Exec(query, name, email, password, userID)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`UPDATE user_roles SET role_id = $1 WHERE user_id = $2`, roleID, userID)
	return err
}

func (r *UserRepositoryImpl) DeleteUser(userID int) error {
	_, err := r.DB.Exec(`DELETE FROM user_roles WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`DELETE FROM users WHERE id = $1`, userID)
	return err
}
