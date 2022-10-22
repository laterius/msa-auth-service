package service

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/laterius/service_architecture_hw3/app/internal/domain"
	"github.com/laterius/service_architecture_hw3/app/internal/repo"
	"github.com/laterius/service_architecture_hw3/app/modules/hash"
	"github.com/laterius/service_architecture_hw3/app/modules/rand"
	"github.com/laterius/service_architecture_hw3/app/pkg/nullable"
	"github.com/laterius/service_architecture_hw3/app/pkg/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	// UserPwPepper Adding the Pepper value
	UserPwPepper = "secret-random-string"

	// HmacSecret for creating the HMAC
	HmacSecret = "secret-hmac-key"
)

type UserData struct {
	Username     string `json:"username" schema:"username"`
	FirstName    string `json:"firstName" schema:"firstName"`
	LastName     string `json:"lastName" schema:"lastName"`
	Email        string `json:"email" schema:"email"`
	Phone        string `json:"phone" schema:"phone"`
	Password     string `json:"password" schema:"password"`
	PasswordHash string `json:"passwordHash" schema:"passwordHash"`
	Remember     string `json:"remember" schema:"remember"`
	RememberHash string `json:"rememberHash" schema:"rememberHash"`
}

type UserLogin struct {
	Username string `json:"username" schema:"username"`
	Password string `json:"password" schema:"password"`
}

type User struct {
	Id int64 `json:"id"`
	UserData
}

func (u *User) FromDomain(d *domain.User) *User {
	u.Id = int64(d.Id)
	u.Username = d.Username
	u.FirstName = d.FirstName
	u.LastName = d.LastName
	u.Email = d.Email
	u.Phone = d.Phone
	u.Password = d.Password
	u.PasswordHash = d.PasswordHash
	u.Remember = d.Remember
	u.RememberHash = d.RememberHash

	return u
}

type UserCreate UserData

func (u UserCreate) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(2, 32)),
		validation.Field(&u.FirstName, validation.Required, validation.Length(1, 32)),
		validation.Field(&u.LastName, validation.Required, validation.Length(1, 32)),
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Phone, validation.Required, is.E164),
		validation.Field(&u.Phone, validation.Required, validation.Length(5, 16)),
	)
}

func (u UserCreate) ToDomain() *domain.User {
	return &domain.User{
		Id:        0,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Password:  u.Password,
	}
}

type UserUpdate UserData

func (u UserUpdate) ToDomain() *domain.User {
	return &domain.User{
		Id:        0,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Phone:     u.Phone,
		Password:  u.Password,
	}
}

type UserPartialUpdate struct {
	Username  nullable.String `json:"username"`
	FirstName nullable.String `json:"firstName"`
	LastName  nullable.String `json:"lastName"`
	Email     nullable.String `json:"email"`
	Phone     nullable.String `json:"phone"`
	Remember  nullable.String `json:"token"`
}

func (pu UserPartialUpdate) ToDomain() *domain.UserPartialData {
	d := types.NewKv()
	if pu.Username.Set {
		d.Set("username", pu.Username.Value)
	}
	if pu.FirstName.Set {
		d.Set("firstName", pu.FirstName.Value)
	}
	if pu.LastName.Set {
		d.Set("lastName", pu.LastName.Value)
	}
	if pu.Email.Set {
		d.Set("email", pu.Email.Value)
	}
	if pu.Phone.Set {
		d.Set("phone", pu.Phone.Value)
	}

	if pu.Remember.Set {
		d.Set("password", pu.Remember.Value)
	}

	return d
}

type UserReader interface {
	Get(domain.UserId) (*domain.User, error)
}

type UserLoginReader interface {
	Login(domain.Username, domain.Password) (*domain.User, error)
}

type UserRememberReader interface {
	ByRemember(token string) (u *domain.User, err error)
}

type UserCreator interface {
	Create(*UserCreate) (*domain.User, error)
}

type UserUpdater interface {
	Update(domain.UserId, *UserUpdate) (*domain.User, error)
}

type UserPartialUpdater interface {
	PartialUpdate(domain.UserId, *domain.UserPartialData) (*domain.User, error)
}

type UserDeleter interface {
	Delete(domain.UserId) error
}

type UserService interface {
	UserReader
	UserCreator
	UserUpdater
	UserPartialUpdater
	UserDeleter
	UserLoginReader
	UserRememberReader
}

type userService struct {
	reader         repo.UserReader
	observer       repo.UserObserver
	creator        repo.UserCreator
	updater        repo.UserUpdater
	partialUpdater repo.UserPartialUpdater
	deleter        repo.UserDeleter
	loginReader    repo.UserLoginReader
	rememberReader repo.UserRememberReader
	hmac           hash.HMAC
}

func NewUserService(repo repo.UserRepo) UserService {
	return &userService{
		reader:         repo,
		observer:       repo,
		creator:        repo,
		updater:        repo,
		partialUpdater: repo,
		deleter:        repo,
		loginReader:    repo,
		rememberReader: repo,
		hmac:           hash.NewHMAC(HmacSecret),
	}
}

func (s *userService) Get(id domain.UserId) (*domain.User, error) {
	err := id.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.reader.Get(id)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, domain.ErrUserNotFound
		}
	}
	return user, err
}

func (s *userService) Login(username domain.Username, password domain.Password) (*domain.User, error) {
	user, err := s.loginReader.Login(username, password)
	// Validate if the user is existed in the database or no
	if err != nil {
		if err.Error() == "record not found" {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	//Compare the login based in the Hash value
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(string(password)+UserPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, domain.ErrInvalidUserId
		case nil:
			return nil, err
		default:
			return nil, err
		}
	}

	err = s.signIn(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *userService) ByRemember(token string) (*domain.User, error) {
	newuser := domain.User{
		Remember: token,
	}
	// Validating and Normalizing then creating the hash
	s.hmacRememberToken(&newuser)

	user, err := s.rememberReader.ByRemember(newuser.RememberHash)
	//TODO надо протестировать, может лишнее
	user.Remember = token
	// Validate if the user is existed in the database or no
	if err != nil {
		if err.Error() == "record not found" {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, err
}

func (s *userService) Create(req *UserCreate) (*domain.User, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	newUser := req.ToDomain()

	err = bcryptPassword(newUser)
	if err != nil {
		return nil, err
	}

	err = passwordHashRequired(newUser)
	if err != nil {
		return nil, err
	}

	err = setRememberIfUnset(newUser)
	if err != nil {
		return nil, err
	}

	err = s.hmacRememberToken(newUser)
	if err != nil {
		return nil, err
	}

	user, err := s.creator.Create(newUser)

	err = s.signIn(user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *userService) Update(id domain.UserId, req *UserUpdate) (*domain.User, error) {
	exists, err := s.observer.Exists(id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return s.updater.Update(id, req.ToDomain())
}

func (s *userService) PartialUpdate(id domain.UserId, data *domain.UserPartialData) (*domain.User, error) {
	return s.partialUpdater.PartialUpdate(id, data)
}

func (s *userService) Delete(id domain.UserId) error {
	return s.deleter.Delete(id)
}

func (s *userService) signIn(user *domain.User) error {
	// Making sure to Remember the token
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token

		s.hmacRememberToken(user)
		_, err = s.updater.Update(user.Id, user)
		if err != nil {
			return err
		}
	}

	return nil
}

func setRememberIfUnset(user *domain.User) error {
	if user.Remember != "" {
		return nil
	}
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

// Validating the RememberToken and hashing
func (s *userService) hmacRememberToken(user *domain.User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = s.hmac.Hashing(user.Remember)
	return nil
}

// bcryptPassword will hash the user passsword with
// salt and pepper and bcrypt the password
func bcryptPassword(user *domain.User) error {
	pwPepper := []byte(user.Password + UserPwPepper) // Salt + Pepper
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(pwPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes) // store the hashedBytes in the struct
	user.Password = ""                      // Don't store Password

	// Look at the token if it is empty then create it.
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	return nil
}

func passwordHashRequired(user *domain.User) error {
	if user.PasswordHash == "" {
		return errors.New("Hashing Password is required")
	}
	return nil
}
