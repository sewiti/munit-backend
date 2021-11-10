package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"

	"github.com/sewiti/munit-backend/pkg/id"
)

type User struct {
	ID          id.ID     `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`

	Hash []byte `json:"-"`
	Salt []byte `json:"-"`
}

func (u *User) scan(row *sql.Row) (*User, error) {
	return u, row.Scan(
		&u.ID,
		&u.DisplayName,
		&u.Email,
		&u.Hash,
		&u.Salt,
		&u.Created,
		&u.Modified,
	)
}

func (u *User) verify() error {
	const (
		maxDisplayName = 72
		maxEmail       = 112

		minPasswd = 8
		maxPasswd = 72

		passwdRequireDigit = true
		passwdRequirePunct = true
	)

	if err := u.ID.Verify(); err != nil {
		return fmt.Errorf("user: %w", err)
	}

	// DisplayName
	if len(u.DisplayName) > maxDisplayName {
		return fmt.Errorf("user: displayname is too long, max %d", maxDisplayName)
	}

	// Email
	if len(u.Email) > maxEmail {
		return fmt.Errorf("user: email is too long, max %d", maxEmail)
	}

	// mail.ParseAddress parses both:
	// - Linus Torvalds <linus@torvalds.com>
	// - linus@torvalds.com
	// requiring us to do an additional check.
	addr, err := mail.ParseAddress(u.Email)
	if err != nil {
		return fmt.Errorf("user: %w", err)
	}
	if addr.Address != u.Email {
		return errors.New("user: email is invalid")
	}

	// Password
	if len(u.Password) < minPasswd {
		return fmt.Errorf("user: password is too short, min %d", minPasswd)
	}
	if len(u.Password) > maxPasswd {
		return fmt.Errorf("user: password is too long, max %d", maxPasswd)
	}
	if !stringMadeOf(u.Password, unicode.Letter, unicode.Digit, unicode.Punct, unicode.Space) {
		return errors.New("user: password contains invalid symbol")
	}

	if passwdRequireDigit && strings.IndexFunc(u.Password, unicode.IsDigit) < 0 {
		return errors.New("user: password must contain at least one digit")
	}
	if passwdRequirePunct && strings.IndexFunc(u.Password, unicode.IsPunct) < 0 {
		return errors.New("user: password must contain at least one symbol")
	}

	// Hash & Salt
	if len(u.Hash) == 0 {
		return fmt.Errorf("user: password hash is empty")
	}
	if len(u.Salt) == 0 {
		return fmt.Errorf("user: salt is empty")
	}
	return nil
}

func InsertUser(ctx context.Context, u *User) error {
	if err := u.verify(); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx,
		"INSERT INTO user (id, display_name, email, passwd_hash, passwd_salt, created, modified) VALUES (?, ?, ?, ?, ?, ?, ?)",
		u.ID,
		u.DisplayName,
		u.Email,
		u.Hash,
		u.Salt,
		u.Created,
		u.Modified,
	)
	if err != nil {
		return errors.New("email already taken")
	}
	return nil
}

func GetUser(ctx context.Context, id id.ID) (*User, error) {
	row := db.QueryRowContext(ctx,
		"SELECT id, display_name, email, passwd_hash, passwd_salt, created, modified FROM user WHERE id=?",
		id,
	)
	return new(User).scan(row)
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := db.QueryRowContext(ctx,
		"SELECT id, display_name, email, passwd_hash, passwd_salt, created, modified FROM user WHERE email=?",
		email,
	)
	return new(User).scan(row)
}
