package userDomain

import (
	"errors"
	"strings"
	"wimed/internal/domain/userDomain/roles"
)

type User struct {
	id           string
	email        Email
	passwordHash string
	active       bool
	roles        map[roles.Role]struct{}
}

func CreateNewUserDomain(id string, email Email, passwordHash string, active bool, userRoles ...roles.Role) (*User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("id is required")
	}
	if strings.TrimSpace(passwordHash) == "" {
		return nil, errors.New("password is required")
	}

	u := &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		active:       active,
		roles:        make(map[roles.Role]struct{}),
	}
	for _, r := range userRoles {
		if err := u.AddRole(r); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func RebuildUserDomain(id string, email Email, passwordHash string, active bool, userRoles []roles.Role) *User {
	u := &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		active:       active,
		roles:        make(map[roles.Role]struct{}),
	}

	for _, r := range userRoles {
		u.roles[r] = struct{}{}
	}
	return u
}

func (u *User) ID() string           { return u.id }
func (u *User) EmailValue() string   { return u.email.Value() }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) IsActive() bool       { return u.active }

func (u *User) Roles() []roles.Role {
	out := make([]roles.Role, 0, len(u.roles))
	for r := range u.roles {
		out = append(out, r)
	}
	return out
}

func (u *User) AddRole(r roles.Role) error {
	if !r.IsValid() {
		return errors.New("invalid role")
	}
	u.roles[r] = struct{}{}
	return nil
}

func (u *User) RemoveRole(r roles.Role) error {
	if !r.IsValid() {
		return errors.New("invalid role")
	}
	delete(u.roles, r)
	return nil
}

func (u *User) Activate() error {
	if u.active {
		return errors.New("userDomain is already active")
	}
	u.active = true
	return nil
}

func (u *User) Deactivate() error {
	if !u.active {
		return errors.New("userDomain as already inactive")
	}
	u.active = false
	return nil
}

func (u *User) SetPasswordHash(passwordHash string) error {
	if strings.TrimSpace(passwordHash) == "" {
		return errors.New("password is required")
	}
	u.passwordHash = passwordHash
	return nil
}

func (u *User) ChangeEmail(email Email) {
	u.email = email
}
