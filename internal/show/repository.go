package show

import (
	"context"
	"github.com/jackc/pgx/v4"
	"time"
)

type Repository interface {
	Add(ctx context.Context, show Show) (int, error)
	Update(ctx context.Context, show Show) error
	Get(ctx context.Context, id int) (Show, error)
	Delete(ctx context.Context, id int) error
	GetAll(ctx context.Context) ([]Show, error)
	GetShowsWithSameWeekday(ctx context.Context, show Show) ([]Show, error)
}

type PostgresRepo struct {
	conn *pgx.Conn
}

func NewPostgresRepo(conn *pgx.Conn) *PostgresRepo {
	return &PostgresRepo{conn: conn}
}

func (o *PostgresRepo) Add(ctx context.Context, show Show) (int, error) {
	q := "insert into shows (title,weekday,timeslot,description,created_at,updated_at) values($1,$2,$3,$4,$5,$6) returning id"
	row := o.conn.QueryRow(ctx, q, show.Title, show.Weekday, show.Timeslot, show.Description, time.Now(), time.Now())
	var id int
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (o *PostgresRepo) Update(ctx context.Context, show Show) error {
	q := "update shows set title=$1,weekday=$2,timeslot=$3,description=$4,created_at=$5,updated_at=$6"
	_, err := o.conn.Exec(ctx, q, show.Title, show.Weekday, show.Timeslot, show.Description, show.CreatedAt, time.Now())
	return err
}

func (o *PostgresRepo) Get(ctx context.Context, id int) (Show, error) {
	q := "select id,title,weekday,timeslot,description,created_at,updated_at from shows where id=$1"
	row := o.conn.QueryRow(ctx, q, id)
	var s Show
	if err := row.Scan(&s.Id, &s.Title, &s.Weekday, &s.Timeslot, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return s, nil
	}
	return s, nil
}

func (o *PostgresRepo) GetAll(ctx context.Context) ([]Show, error) {
	q := "select id,title,weekday,timeslot,description,created_at,updated_at from shows"
	rows, err := o.conn.Query(ctx, q)
	if err != nil {
		panic(err)
	}
	var result []Show
	for rows.Next() {
		var s Show
		if err := rows.Scan(&s.Id, &s.Title, &s.Weekday, &s.Timeslot, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (o *PostgresRepo) GetShowsWithSameWeekday(ctx context.Context, show Show) ([]Show, error) {
	q := "select id,title,weekday,timeslot,description,created_at,updated_at from shows where weekday=$1"
	rows, err := o.conn.Query(ctx, q, show.Weekday)
	if err != nil {
		return nil, err
	}
	var result []Show
	for rows.Next() {
		var s Show
		if err := rows.Scan(&s.Id, &s.Title, &s.Weekday, &s.Timeslot, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (o *PostgresRepo) Delete(ctx context.Context, id int) error {
	q := "delete from shows where id=$1"
	_, err := o.conn.Exec(ctx, q, id)
	return err
}
