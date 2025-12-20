package userDomain

import (
	"errors"
	"sort"
	"strings"
	"wimed/internal/domain/userDomain/roles"
)

var (
	ErrIDRequired           = errors.New("user: id is required")
	ErrPasswordHashRequired = errors.New("user: password_hash is required")
	ErrInvalidRole          = errors.New("user: invalid role")
	ErrInvalidEmail         = errors.New("user: invalid email")
	ErrAlreadyActive        = errors.New("user: already active")
	ErrAlreadyInactive      = errors.New("user: already inactive")
	ErrRoleNotAssigned      = errors.New("user: role not assigned")
)

type User struct {
	id           string
	email        Email
	passwordHash string
	active       bool
	roles        map[roles.Role]struct{}
}

func CreateNewUserDomain(id string, email Email, passwordHash string, userRoles ...roles.Role) (*User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrIDRequired
	}
	if email.Value() == "" {
		return nil, ErrInvalidEmail
	}
	passwordHash = strings.TrimSpace(passwordHash)
	if passwordHash == "" {
		return nil, ErrPasswordHashRequired
	}

	u := &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		active:       true,
		roles:        make(map[roles.Role]struct{}),
	}
	for _, r := range userRoles {
		if err := u.AddRole(r); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func RebuildUserDomain(id string, email Email, passwordHash string, active bool, userRoles []roles.Role) (*User, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, ErrIDRequired
	}
	if email.Value() == "" {
		return nil, ErrInvalidEmail
	}
	passwordHash = strings.TrimSpace(passwordHash)
	if passwordHash == "" {
		return nil, ErrPasswordHashRequired
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

func (u *User) ID() string           { return u.id }
func (u *User) Email() Email         { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) IsActive() bool       { return u.active }

func (u *User) Roles() []roles.Role {
	out := make([]roles.Role, 0, len(u.roles))
	for r := range u.roles {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (u *User) AddRole(r roles.Role) error {
	if !r.IsValid() {
		return ErrInvalidRole
	}
	u.roles[r] = struct{}{}
	return nil
}

func (u *User) RemoveRole(r roles.Role) error {
	if !r.IsValid() {

		return ErrInvalidRole
	}
	if !u.HasRole(r) {
		return ErrRoleNotAssigned
	}
	delete(u.roles, r)
	return nil
}

func (u *User) Activate() error {
	if u.active {
		return ErrAlreadyActive
	}
	u.active = true
	return nil
}

func (u *User) Deactivate() error {
	if !u.active {
		return ErrAlreadyInactive
	}
	u.active = false
	return nil
}

func (u *User) SetPasswordHash(passwordHash string) error {
	if strings.TrimSpace(passwordHash) == "" {
		return ErrPasswordHashRequired
	}
	u.passwordHash = passwordHash
	return nil
}

func (u *User) ChangeEmail(email Email) error {
	if email.Value() == "" {
		return ErrInvalidEmail
	}
	u.email = email
	return nil
}

func (u *User) HasRole(r roles.Role) bool {
	_, ok := u.roles[r]
	return ok
}
