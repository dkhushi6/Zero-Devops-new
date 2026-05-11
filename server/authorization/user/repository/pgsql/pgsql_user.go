package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"Zero_Devops/server/domain"
)

type pqSqlUserRepository struct {
	Conn *sql.DB
}

func NewPgSqlUserRepository(conn *sql.DB) domain.UserRepository {
	return &pqSqlUserRepository{conn}
}

func (m *pqSqlUserRepository) GetByID (ctx context.Context , id int64) (domain.User, error) {
	query := `SELECT * FROM users WHERE ID = $1`
	row := m.Conn.QueryRowContext(ctx,query,id)

	
	var u domain.User
	err := row.Scan(
		&u.ID,
		&u.ProviderId,
		&u.Provider,
		&u.Username,
		&u.Email,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.RefreshToken,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		logrus.Error(err)
		return domain.User{}, err
	}

	return u, nil
}

func (m *pqSqlUserRepository) GetByUsername (ctx context.Context , username string) (domain.User , error){
	query := `SELECT * FROM users WHERE Username = $1`
	row := m.Conn.QueryRowContext(ctx,query,username)
	
	u := domain.User{}
	err := row.Scan(
		&u.ID,
		&u.ProviderId,
		&u.Provider,
		&u.Username,
		&u.Email,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.RefreshToken,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		logrus.Error(err)
		return domain.User{}, err
	}

	
	return u,nil
}

func (m* pqSqlUserRepository) GetProviderById(ctx context.Context , providerId int64) (domain.User , error){
	query := `SELECT * FROM users WHERE ProviderId = $1`
	row := m.Conn.QueryRowContext(ctx,query,providerId)

	u := domain.User{}
	err := row.Scan(
		&u.ID,
		&u.ProviderId,
		&u.Provider,
		&u.Username,
		&u.Email,
		&u.AvatarURL,
		&u.CreatedAt,
		&u.RefreshToken,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrNotFound
		}
		logrus.Error(err)
		return domain.User{}, err
	}

	return u , nil
}

func (m *pqSqlUserRepository) Store(ctx context.Context , user *domain.User) (error){
	query := `INSERT INTO users (ProviderId , Provider , Username , Email , AvatarURL , CreatedAt) VALUES ($1, $2, $3, $4, $5 , $6 ) RETURNING ID`

	err := m.Conn.QueryRowContext(ctx,query,user.ProviderId,user.Provider, user.Username, user.Email, user.AvatarURL, user.CreatedAt).Scan(&user.ID)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (m* pqSqlUserRepository) Update(ctx context.Context , id int64 , refreshToken string) error {
	query := `UPDATE users SET RefreshToken = $1 WHERE ID = $2`
	
	stmt,err := m.Conn.PrepareContext(ctx,query)

	if err != nil {
		logrus.Error(err)
		return err
	}
	
	defer stmt.Close()

	res,err := stmt.ExecContext(ctx,refreshToken,id)
	if err != nil {
		logrus.Error(err)
		return err
	}


	rowsAffected , err := res.RowsAffected()
	if err != nil {
		logrus.Error(err)
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("weird Behavior. Total Affected: %d", rowsAffected)
	}

	return nil
}