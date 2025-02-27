package pessoa

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joaovrivero/rinha-backend"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	Conn      *pgxpool.Pool
	Cache     *redis.Client
	ChPessoas chan rinha.Pessoa
}

var Repo *Repository

func NewRepository(Conn *pgxpool.Pool, Cache *redis.Client) *Repository {
	if Repo == nil {
		Repo = &Repository{Conn: Conn, Cache: Cache, ChPessoas: make(chan rinha.Pessoa)}
	}

	return Repo
}

func (r *Repository) Create(ctx context.Context, pessoa rinha.Pessoa) error {
	if v, _ := r.Cache.Get(ctx, pessoa.Apelido).Result(); v != "" {
		return rinha.ErrApelidoJaExiste
	}

	j, err := json.Marshal(pessoa)
	if err != nil {
		return err
	}

	pipe := r.Cache.Pipeline()
	pipe.Set(ctx, pessoa.Apelido, "t", 0)
	pipe.Set(ctx, pessoa.ID.String(), j, 24*time.Hour)
	if _, err := pipe.Exec(ctx); err != nil {
		slog.Error(err.Error())
		return err
	}

	r.ChPessoas <- pessoa

	return nil
}

func (r *Repository) Insert(pessoas []rinha.Pessoa) error {
	if len(pessoas) == 0 {
		return nil
	}

	_, err := r.Conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"pessoas"},
		[]string{"id", "apelido", "nome", "nascimento", "stack", "search_index"},
		pgx.CopyFromSlice(len(pessoas), func(i int) ([]any, error) {
			p := pessoas[i]
			index := fmt.Sprintf("%s %s %s", strings.ToLower(p.Apelido), strings.ToLower(p.Nome), strings.ToLower(strings.Join(p.Stack, " ")))
			return []any{p.ID, p.Apelido, p.Nome, p.Nascimento.Time, p.Stack, index}, nil
		}),
	)

	if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.ConstraintName == "pessoas_apelido_key" {
		// @TODO how to deal with conflicts on database
		slog.Error("algum apelido ja existe")
		return pgErr
	}

	return err
}

func (r *Repository) Count(ctx context.Context) (int, error) {
	var count int

	err := r.Conn.QueryRow(ctx, `
		SELECT COUNT(id)
		FROM pessoas
	`).Scan(&count)

	return count, err
}

func (r *Repository) FindOne(ctx context.Context, id uuid.UUID) (rinha.Pessoa, error) {
	var pessoa rinha.Pessoa

	entry, err := r.Cache.Get(ctx, id.String()).Result()
	if err != nil && err != redis.Nil {
		slog.Error(err.Error())
	}

	if entry != "" {
		err := json.Unmarshal([]byte(entry), &pessoa)
		if err == nil {
			return pessoa, nil
		}

		slog.Error(err.Error())
	}

	var nascimento time.Time
	err = r.Conn.QueryRow(ctx, `
		SELECT id, apelido, nome, nascimento, stack
		FROM pessoas
		WHERE id = $1
		LIMIT 1
	`, id).Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &nascimento, &pessoa.Stack)

	pessoa.Nascimento = rinha.Date{Time: nascimento}

	if err == pgx.ErrNoRows {
		return pessoa, nil
	}

	return pessoa, err
}

func (r *Repository) FindByTermo(ctx context.Context, t string) ([]rinha.Pessoa, error) {
	pessoas := []rinha.Pessoa{}

	rows, err := r.Conn.Query(ctx, `
		SELECT id, apelido, nome, nascimento, stack
		FROM pessoas
		WHERE search_index ILIKE '%' || $1 || '%'
		LIMIT 50
	`, strings.ToLower(t))

	if err != nil {
		return pessoas, err
	}

	defer rows.Close()

	for rows.Next() {
		var pessoa rinha.Pessoa
		var nascimento time.Time
		err := rows.Scan(&pessoa.ID, &pessoa.Apelido, &pessoa.Nome, &nascimento, &pessoa.Stack)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		pessoa.Nascimento = rinha.Date{Time: nascimento}
		pessoas = append(pessoas, pessoa)
	}

	return pessoas, err
}
