package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"Zero_Devops/server/domain"
)

type pgSqlGithubRepository struct {
	Conn *sql.DB
}

func NewPgSqlGithubRepository(conn *sql.DB) domain.GithubRepository {
	return &pgSqlGithubRepository{conn}
}

func (m *pgSqlGithubRepository) StoreInstallation(ctx context.Context, inst *domain.GithubInstallation) error {
	query := `INSERT INTO Github (UserID , InstallationID , AccountName) VALUES ($1 , $2 , $3)`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		logrus.Error(err)
		return err
	}

	result, err := stmt.ExecContext(ctx, inst.UserID, inst.InstallationID, inst.AccountName)

	if err != nil {
		logrus.Error(err)
		return err
	}

	// This feature is not supported by all the databases sonce here I am using Postgres then I have to check for it
	new_id,err := result.LastInsertId() // Here is a issue I need to check it

	if err != nil{
		logrus.Error(err)
		return err
	}

	inst.ID = new_id

	return nil
}

func (m *pgSqlGithubRepository) GetInstallationByUserID(ctx context.Context, userId int64) (*domain.GithubInstallation, error) {
	query := `SELECT * FROM Github WHERE UserID = $1`
	res := m.Conn.QueryRowContext(ctx, query, userId)

	inst := domain.GithubInstallation{}
	err := res.Scan(
		&inst.ID,
		&inst.UserID,
		&inst.InstallationID,
		&inst.AccountName,
	)

	if err != nil {
		if err == sql.ErrNoRows{
			return nil,domain.ErrNotFound
		}
		logrus.Error(err)
		return nil, err
	}

	return &inst, nil
}

func (m *pgSqlGithubRepository) DeleteInstallationByUserID(ctx context.Context, userId int64) error {
	query := `DELETE FROM Github WHERE UserID = $1`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		logrus.Error(err)
		return err
	}

	res, err := stmt.ExecContext(ctx, userId)
	if err != nil {
		logrus.Error(err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", rowsAffected)
		return err
	}

	return nil
}