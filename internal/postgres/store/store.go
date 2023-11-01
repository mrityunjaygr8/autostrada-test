package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mrityunjaygr8/autostrada-test/assets"
	"github.com/mrityunjaygr8/autostrada-test/internal/postgres/models"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"math"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

var ErrCreatingPostgresPool = errors.New("error creating postgres pool")
var ErrConnectingToPostgres = errors.New("error connecting to postgres")

type transactionFunction func() error

func (p *PostgresStore) createTx() (ctx context.Context, tx pgx.Tx, commit transactionFunction, rollback transactionFunction, err error) {
	ctx = context.Background()
	tx, err = p.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("%w: %q", store.ErrStoreError, err)
	}

	commit = func() error {
		err := tx.Commit(ctx)
		if err != nil {
			return fmt.Errorf("%w: %q", store.ErrStoreError, err)
		}
		return nil
	}
	rollback = func() error {
		err := tx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("%w: %q", store.ErrStoreError, err)
		}
		return nil
	}

	return ctx, tx, commit, rollback, nil

}

func NewPostgresStore(dbString string, autoMigrate bool) (*PostgresStore, func(), error) {
	db, err := pgxpool.New(context.Background(), "postgres://"+dbString)
	if err != nil {
		return nil, nil, ErrCreatingPostgresPool
	}

	err = db.Ping(context.Background())
	if err != nil {
		return nil, nil, ErrConnectingToPostgres
	}

	if autoMigrate {
		iofsDriver, err := iofs.New(assets.EmbeddedFiles, "migrations")
		if err != nil {
			return nil, nil, err
		}

		migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, "pgx5://"+dbString)
		if err != nil {
			return nil, nil, err
		}

		err = migrator.Up()
		switch {
		case errors.Is(err, migrate.ErrNoChange):
			break
		case err != nil:
			return nil, nil, err
		}
	}

	closer := func() {
		db.Close()
	}

	return &PostgresStore{db: db}, closer, nil
}

func (p *PostgresStore) UserInsert(email, password string, id uuid.UUID, admin bool) (*store.User, error) {
	ctx, tx, commit, rollback, err := p.createTx()
	if err != nil {
		return nil, err
	}

	query := models.New(tx)
	params := models.UserInsertParams{
		Email:          email,
		HashedPassword: password,
		ID:             id,
		Admin:          admin,
	}
	dbUser, err := query.UserInsert(ctx, params)
	if err != nil {
		if txErr := rollback(); txErr != nil {
			return nil, txErr
		}
		var pge *pgconn.PgError
		if errors.As(err, &pge) {
			if pge.SQLState() == "23505" {
				return nil, store.ErrUserExists
			}
		}
		return nil, err
	}

	if txErr := commit(); txErr != nil {
		return nil, txErr
	}
	user := &store.User{
		Email:   dbUser.Email,
		ID:      dbUser.ID,
		Admin:   dbUser.Admin,
		Created: dbUser.Created.Time,
	}
	return user, nil
}

func (p *PostgresStore) UserList(pageNumber, pageSize int) (*store.UsersList, error) {
	query := models.New(p.db)
	params := models.UsersListParams{
		Limit:  int32(pageSize),
		Offset: int32((pageNumber - 1) * pageSize),
	}
	dbUsers, err := query.UsersList(context.Background(), params)
	if err != nil {
		return nil, err
	}

	users := make([]store.User, 0)
	totalObjects := 0
	totalPages := 0

	for _, user := range dbUsers {
		users = append(users, store.User{
			Email:   user.Email,
			ID:      user.ID,
			Admin:   user.Admin,
			Created: user.Created.Time,
		})
	}

	if len(dbUsers) > 0 {
		totalObjects = int(dbUsers[0].RowData)
		totalPages = int(math.Ceil(float64(totalObjects) / float64(pageSize)))
	}

	return &store.UsersList{Data: users, TotalObjects: totalObjects, TotalPages: totalPages}, nil
}

func (p *PostgresStore) UserRetrieve(id uuid.UUID) (*store.User, error) {
	query := models.New(p.db)
	user, err := query.UserRetrieve(context.Background(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrUserNotFound
		}
		return nil, err
	}

	return &store.User{
		Email:   user.Email,
		ID:      user.ID,
		Admin:   user.Admin,
		Created: user.Created.Time,
	}, nil
}

func (p *PostgresStore) UserRetrieveByEmail(email string) (*store.User, error) {
	query := models.New(p.db)
	user, err := query.UserRetrieveByEmail(context.Background(), email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrUserNotFound
		}
		return nil, err
	}

	return &store.User{
		Email:   user.Email,
		ID:      user.ID,
		Admin:   user.Admin,
		Created: user.Created.Time,
	}, nil
}
func (p *PostgresStore) UserUpdatePassword(id uuid.UUID, newPassword string) error {
	ctx, tx, commit, rollback, err := p.createTx()
	if err != nil {
		return err
	}
	query := models.New(tx)
	params := models.UserUpdatePasswordParams{
		ID:             id,
		HashedPassword: newPassword,
	}
	res, err := query.UserUpdatePassword(ctx, params)
	if err != nil {
		if txErr := rollback(); txErr != nil {
			return txErr
		}
		return err
	}
	if res.RowsAffected() == 0 {
		if txErr := rollback(); txErr != nil {
			return txErr
		}
		return store.ErrUserNotFound
	}

	if txErr := commit(); txErr != nil {
		return txErr
	}
	return nil
}

func (p *PostgresStore) UserUpdateAdmin(id uuid.UUID, newAdminValue bool) error {
	ctx, tx, commit, rollback, err := p.createTx()
	if err != nil {
		return err
	}
	query := models.New(tx)
	params := models.UserUpdateAdminParams{
		ID:    id,
		Admin: newAdminValue,
	}
	res, err := query.UserUpdateAdmin(ctx, params)
	if err != nil {
		if txErr := rollback(); txErr != nil {
			return txErr
		}
		return err
	}
	if res.RowsAffected() == 0 {
		if txErr := rollback(); txErr != nil {
			return txErr
		}

		return store.ErrUserNotFound
	}
	if txErr := commit(); txErr != nil {
		return txErr
	}
	return nil
}

func (p *PostgresStore) UserDelete(id uuid.UUID) error {
	ctx, tx, commit, rollback, err := p.createTx()
	if err != nil {
		return err
	}

	query := models.New(tx)
	res, err := query.UserDelete(ctx, id)
	if err != nil {
		if txErr := rollback(); txErr != nil {
			return txErr
		}
		return err
	}
	if res.RowsAffected() == 0 {
		if txErr := rollback(); txErr != nil {
			return txErr
		}
		return store.ErrUserNotFound
	}
	if txErr := commit(); txErr != nil {
		return txErr
	}
	return nil
}
