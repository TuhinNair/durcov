package durcov

import "github.com/jackc/pgx"

// GetPgxPool returns a thread safe postgres connection.
func GetPgxPool(dbURL string) (*pgx.ConnPool, error) {
	pgxcfg, err := pgx.ParseURI(dbURL)
	if err != nil {
		return nil, err
	}

	pgxpool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgxcfg,
	})

	if err != nil {
		return nil, err
	}

	return pgxpool, nil
}
