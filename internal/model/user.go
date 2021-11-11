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

const (
	userSelect      = "SELECT id, display_name, email, passwd_hash, passwd_salt, created, modified FROM user"
	userSelectID    = userSelect + " WHERE id=?"
	userSelectEmail = userSelect + " WHERE email=?"

	userInsert = "INSERT INTO user (id, display_name, email, passwd_hash, passwd_salt, created, modified) VALUES (?,?,?,?,?,?,?)"
	userUpdate = "UPDATE user SET display_name=?, email=?, passwd_hash=?, passwd_salt=?, modified=? WHERE id=?"
)

type User struct {
	ID          id.ID     `json:"id"`
	DisplayName string    `json:"displayName"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`

	PasswdHash []byte `json:"-"`
	Salt       []byte `json:"-"`
}

func (u *User) scan(row *sql.Row) (*User, error) {
	return u, row.Scan(
		&u.ID,
		&u.DisplayName,
		&u.Email,
		&u.PasswdHash,
		&u.Salt,
		&u.Created,
		&u.Modified,
	)
}

func (u *User) validate() error {
	return u.validatePasswd(false)
}

func (u *User) validatePasswd(allowEmptyPasswd bool) error {
	const (
		maxDisplayName = 72
		maxEmail       = 112

		minPasswd = 8
		maxPasswd = 72

		passwdRequireLetter     = true
		passwdRequireUpperLower = false
		passwdRequireDigit      = true
		passwdRequirePunct      = false
	)

	if err := u.ID.Validate(); err != nil {
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

	// Hash & Salt
	if len(u.PasswdHash) == 0 {
		return fmt.Errorf("user: password hash is empty")
	}
	if len(u.Salt) == 0 {
		return fmt.Errorf("user: salt is empty")
	}

	// Password
	if allowEmptyPasswd && u.Password == "" {
		return nil
	}
	if len(u.Password) < minPasswd {
		return fmt.Errorf("user: password is too short, min %d", minPasswd)
	}
	if len(u.Password) > maxPasswd {
		return fmt.Errorf("user: password is too long, max %d", maxPasswd)
	}
	if !stringMadeOf(u.Password, unicode.Letter, unicode.Digit, unicode.Punct, unicode.Space) {
		return errors.New("user: password contains invalid symbol")
	}

	if passwdRequireLetter && strings.IndexFunc(u.Password, unicode.IsLetter) < 0 {
		return errors.New("user: password must contain at least one letter")
	}
	if passwdRequireUpperLower {
		if strings.IndexFunc(u.Password, unicode.IsUpper) < 0 {
			return errors.New("user: password must contain at least one uppercase letter")
		}
		if strings.IndexFunc(u.Password, unicode.IsLower) < 0 {
			return errors.New("user: password must contain at least one lowercase letter")
		}
	}
	if passwdRequireDigit && strings.IndexFunc(u.Password, unicode.IsDigit) < 0 {
		return errors.New("user: password must contain at least one digit")
	}
	if passwdRequirePunct && strings.IndexFunc(u.Password, unicode.IsPunct) < 0 {
		return errors.New("user: password must contain at least one symbol")
	}
	return nil
}

func (u *User) Copy() *User {
	newUser := new(User)
	*newUser = *u

	newUser.PasswdHash = make([]byte, len(u.PasswdHash))
	copy(newUser.PasswdHash, u.PasswdHash)

	newUser.Salt = make([]byte, len(u.Salt))
	copy(newUser.Salt, u.Salt)
	return newUser
}

func InsertUser(ctx context.Context, u *User) error {
	if err := u.validate(); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, userInsert,
		u.ID,
		u.DisplayName,
		u.Email,
		u.PasswdHash,
		u.Salt,
		u.Created,
		u.Modified,
	)
	if err != nil {
		if isDuplicate(err) {
			return errors.New("email is taken")
		}
		return err
	}
	return nil
}

func GetUser(ctx context.Context, id id.ID) (*User, error) {
	row := db.QueryRowContext(ctx, userSelectID, id)
	return new(User).scan(row)
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row := db.QueryRowContext(ctx, userSelectEmail, email)
	return new(User).scan(row)
}

func UpdateUser(ctx context.Context, uid id.ID, modifyFn func(*User) error) (*User, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	row := tx.QueryRowContext(ctx, userSelectID, uid)
	u, err := new(User).scan(row)
	if err != nil {
		return nil, err
	}

	if err = modifyFn(u); err != nil {
		return nil, err
	}
	if err = u.validatePasswd(true); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, userUpdate,
		u.DisplayName,
		u.Email,
		u.PasswdHash,
		u.Salt,
		u.Modified,
		uid,
	)
	if err != nil {
		if isDuplicate(err) {
			return nil, errors.New("email is taken")
		}
		return nil, err
	}
	return u, tx.Commit()
}

func DeleteUser(ctx context.Context, uid id.ID) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM contributor WHERE user_id=?", uid)
	if err != nil {
		return err
	}
	// Projects.. let's not delete those?...
	res, err := tx.ExecContext(ctx, "DELETE FROM user WHERE id=?", uid)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}

	return tx.Commit()
}
