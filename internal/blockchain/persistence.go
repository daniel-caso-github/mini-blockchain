package blockchain

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	_ "modernc.org/sqlite"
)

type Store struct {
	db *sql.DB
}

func OpenStore(dsn string) (*Store, error) {
	if dsn == "" {
		return nil, nil
	}

	if dsn != ":memory:" {
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("creating directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if dsn != ":memory:" {
		if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
			db.Close()
			return nil, fmt.Errorf("setting WAL mode: %w", err)
		}
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS metadata (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS blocks (
			idx        INTEGER PRIMARY KEY,
			timestamp  TEXT    NOT NULL,
			data       TEXT    NOT NULL,
			prev_hash  TEXT    NOT NULL,
			hash       TEXT    NOT NULL,
			nonce      INTEGER NOT NULL
		);
	`); err != nil {
		db.Close()
		return nil, fmt.Errorf("creating tables: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) LoadChain() ([]Block, int, error) {
	if s == nil {
		return nil, 0, fmt.Errorf("no store configured")
	}

	difficulty := 0
	row := s.db.QueryRow("SELECT value FROM metadata WHERE key = 'difficulty'")
	var val string
	if err := row.Scan(&val); err == nil {
		difficulty, _ = strconv.Atoi(val)
	}

	rows, err := s.db.Query("SELECT idx, timestamp, data, prev_hash, hash, nonce FROM blocks ORDER BY idx")
	if err != nil {
		return nil, 0, fmt.Errorf("querying blocks: %w", err)
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var b Block
		if err := rows.Scan(&b.Index, &b.Timestamp, &b.Data, &b.PrevHash, &b.Hash, &b.Nonce); err != nil {
			return nil, 0, fmt.Errorf("scanning block: %w", err)
		}
		blocks = append(blocks, b)
	}

	return blocks, difficulty, rows.Err()
}

func (s *Store) SaveBlock(b Block) error {
	if s == nil {
		return nil
	}
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO blocks (idx, timestamp, data, prev_hash, hash, nonce) VALUES (?, ?, ?, ?, ?, ?)",
		b.Index, b.Timestamp, b.Data, b.PrevHash, b.Hash, b.Nonce,
	)
	if err != nil {
		return fmt.Errorf("saving block: %w", err)
	}
	return nil
}

func (s *Store) SaveDifficulty(d int) error {
	if s == nil {
		return nil
	}
	_, err := s.db.Exec(
		"INSERT OR REPLACE INTO metadata (key, value) VALUES ('difficulty', ?)",
		strconv.Itoa(d),
	)
	if err != nil {
		return fmt.Errorf("saving difficulty: %w", err)
	}
	return nil
}
