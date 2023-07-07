package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	l "github.com/kdrkrgz/socalize/pkg/logger"
	"github.com/kdrkrgz/socalize/users"
	"go.uber.org/zap"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type PgxIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
	SendBatch(context.Context, *pgx.Batch) pgx.BatchResults
	Close()
}

type Repository struct {
	pool *pgxpool.Pool
}

type myQueryTracer struct {
	log *zap.SugaredLogger
}

func (tracer *myQueryTracer) TraceQueryStart(
	ctx context.Context,
	_ *pgx.Conn,
	data pgx.TraceQueryStartData) context.Context {
	tracer.log.Debugw("Executing command", "sql", data.SQL, "args", data.Args)

	return ctx
}

func (tracer *myQueryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
}

func New() *Repository {
	dbUrl := "host=0.0.0.0 user=postgres dbname=socialize port=5432 sslmode=disable"
	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.MinConns = 5
	config.MaxConns = 10
	config.ConnConfig.Tracer = &myQueryTracer{
		log: l.Logger().Sugar(),
	}

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return &Repository{
		pool: pool,
	}
}

func (repo *Repository) Close() {
	repo.pool.Close()
}

func (repo *Repository) GetUsers() ([]users.UsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	selectBuilder := psql.Select("users.id",
		"users.created_at",
		"users.updated_at",
		"users.username",
		"users.email",
		"users.phone",
		"users.role",
		"users.password",
		"users.first_name",
		"users.last_name",
		"user_profiles.bio",
		"user_profiles.id",
		"user_profiles.created_at",
		"user_profiles.updated_at",
	).From("users").
		LeftJoin("user_profiles ON users.id = user_profiles.user_id")
	newSql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("could not format query : %w", err)
	}

	query, err := repo.pool.Query(ctx, newSql, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query users: %w", err)
	}
	defer query.Close()

	userCollection := make([]users.UsersResponse, 0)
	for query.Next() {
		var u users.UsersResponse
		err := query.Scan(
			&u.Id,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.Username,
			&u.Email,
			&u.Phone,
			&u.Role,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.Profile.Bio,
			&u.Profile.Id,
			&u.Profile.CreatedAt,
			&u.Profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		userCollection = append(userCollection, u)
	}
	if err := query.Err(); err != nil {
		return nil, err
	}
	return userCollection, nil
}

func (repo *Repository) CreateUser(u *users.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	insertBuilder := psql.Insert("users").
		Columns("username", "email", "phone", "password", "first_name", "last_name").
		Values(u.Username, u.Email, u.Phone, u.Password, u.FirstName, u.LastName).
		Suffix("RETURNING \"id\"")
	newSql, args, err := insertBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("could not format query : %w", err)
	}
	row := repo.pool.QueryRow(ctx, newSql, args...)
	if err := row.Scan(&u.Id); err != nil {
		return fmt.Errorf("could not scan id: %w", err)
	}
	return nil
}

func (repo *Repository) GetUserByEmail(email string) *users.User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	selectBuilder := psql.Select("*").From("users").Where(sq.Eq{"email": email})
	newSql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil
	}
	row := repo.pool.QueryRow(ctx, newSql, args...)
	var u users.User
	if err := row.Scan(
		&u.Id,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Username,
		&u.Email,
		&u.Phone,
		&u.Password,
		&u.Role,
		&u.FirstName,
		&u.LastName); err != nil {
		return nil
	}
	return &u
}

func (repo *Repository) GetUserById(id uint) *users.User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// TODO: check id - check create jwt and lib
	selectBuilder := psql.Select("*").From("users").Where(sq.Eq{"id": id})
	newSql, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil
	}
	row := repo.pool.QueryRow(ctx, newSql, args...)
	var u users.User
	if err := row.Scan(
		&u.Id,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Username,
		&u.Email,
		&u.Phone,
		&u.Role,
		&u.Password,
		&u.FirstName,
		&u.LastName); err != nil {
		return nil
	}
	return &u
}
