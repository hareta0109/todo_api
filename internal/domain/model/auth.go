package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"todo_api/internal/lib/apperr"
	"unicode/utf8"
)

var (
	errInvalidPasswordLength = errors.New("User Password must be greater than or equal 1 characters")
)

const (
	minUserPasswordLength = 1
)

type Auth struct {
	ID       UserIdentifier
	Name     string
	Hash     string
	Role     UserRole
	UserType UserType
	Company  Company
}

type AuthDescription struct {
	Name     string
	Password *string
	Role     UserRole
	UserType UserType
	Company  *Company
}

func NewAuth(desc AuthDescription) (*Auth, apperr.AppErr) {
	auth := new(Auth)
	if err := auth.Update(desc); err != nil {
		return nil, err
	}

	return auth, nil
}

func (m *Auth) Update(desc AuthDescription) apperr.AppErr {
	if err := desc.validate(); err != nil {
		return apperr.NewBadRequestError().Wrap(err)
	}

	m.Name = desc.Name
	if desc.Password != nil {
		m.Hash = PasswordToHash(*desc.Password)
	}
	m.Role = desc.Role
	m.UserType = desc.UserType
	if desc.Company != nil {
		m.Company = *desc.Company
	}
	return nil
}

func (d *AuthDescription) validate() error {
	nameLength := utf8.RuneCountInString(d.Name)
	if nameLength < minUserNameLength || nameLength > maxUserNameLength {
		return errInvalidUserNameLength
	}

	passwordLength := utf8.RuneCountInString(*d.Password)
	if passwordLength < minUserPasswordLength {
		return errInvalidPasswordLength
	}

	return nil
}

func PasswordToHash(password string) string {
	b := getSHA256Binary(password)
	h := hex.EncodeToString(b)
	return h
}

func getSHA256Binary(s string) []byte {
	r := sha256.Sum256([]byte(s))
	return r[:]
}

// IsSuperAdmin / 管理会社の管理者かどうかを示す。任意の操作が可能。
func (m *Auth) IsSuperAdmin() bool {
	return m.Company.ID == AdminCompanyID &&
		m.UserType == UserTypeAdmin
}

// IsSuperUser / 管理会社に所属しているかどうかを示す。任意の会社の閲覧権限を持つ。
func (m *Auth) IsSuperUser() bool {
	return m.Company.ID == AdminCompanyID
}

// isCompanyUser / 対象の会社に所属しているかどうか
func (m *Auth) isCompanyUser(companyID CompanyIdentifier) bool {
	return m.Company.ID == companyID
}

// isCompanyAdmin / 対象の会社の管理者かどうか
func (m *Auth) isCompanyAdmin(companyID CompanyIdentifier) bool {
	return m.isCompanyUser(companyID) &&
		m.UserType == UserTypeAdmin
}

// CanAdminAction / 会社の管理者としての操作が可能かどうか
func (m *Auth) CanAdminAction(companyID CompanyIdentifier) bool {
	return m.IsSuperAdmin() || m.isCompanyAdmin(companyID)
}

// CanEdit / 対象の会社の編集が可能かどうか
func (m *Auth) CanEdit(companyID CompanyIdentifier) bool {
	return m.isCompanyUser(companyID) &&
		m.Role == UserRoleEditor
}

// CanView / 対象の会社の閲覧が可能かどうか
func (m *Auth) CanView(companyID CompanyIdentifier) bool {
	return m.isCompanyUser(companyID)
}
