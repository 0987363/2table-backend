package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 保留原有的正则表达式和常量定义
var RegexpUserNameLogin = regexp.MustCompile(`^[A-Za-z0-9\+\-_\.\@]{3,30}$`)
var RegexpPwd = regexp.MustCompile(`^[A-Za-z0-9]{64}$`)
var RegexpEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var RegexpPhone = regexp.MustCompile(`^[0-9\+\-][0-9\+\- ]{1,18}[0-9]$`)

const (
	RoleRoot    = "root"
	RoleManager = "manager"
	// ... 其他角色定义保持不变
)

// 用户结构体调整为适合BadgerDB的格式
type User struct {
	ID             string     `json:"id"`
	UserName       string     `json:"username"`
	FullName       string     `json:"full_name"`
	Password       string     `json:"password,omitempty"`
	OldPassword    string     `json:"old_password,omitempty"`
	HashedPassword string     `json:"hashed_password,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	DeletedBy      string     `json:"deleted_by,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	Address        string     `json:"address,omitempty"`
	Phone          string     `json:"phone,omitempty"`
	Email          string     `json:"email,omitempty"`
	Token          string     `json:"token,omitempty"`
	Expiry         *time.Time `json:"expiry,omitempty"`
	Role           string     `json:"role,omitempty"`
	CreatorID      string     `json:"creator_id,omitempty"`
	CompanyID      string     `json:"company_id,omitempty"`
	CompanyName    string     `json:"company_name,omitempty"`
	Description    string     `json:"description,omitempty"`
	AgencyID       string     `json:"agency_id,omitempty"`
	AgencyName     string     `json:"agency_name,omitempty"`
}

// 用户存储结构
type UserStore struct {
	db *badger.DB
}

func NewUserStore(db *badger.DB) *UserStore {
	return &UserStore{db: db}
}

// 生成存储键名
func (s *UserStore) userKey(id string) []byte {
	return []byte("user:" + id)
}

func (s *UserStore) usernameIndexKey(username string) []byte {
	return []byte("user_index:username:" + username)
}

func (s *UserStore) emailIndexKey(email string) []byte {
	return []byte("user_index:email:" + email)
}

// 创建用户
func (s *UserStore) CreateUser(user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// 密码处理
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %v", err)
		}
		user.HashedPassword = string(hashed)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		// 检查用户名唯一性
		if _, err := txn.Get(s.usernameIndexKey(user.UserName)); err == nil {
			return fmt.Errorf("用户名已存在")
		}

		// 序列化用户数据
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		// 存储主数据
		if err := txn.Set(s.userKey(user.ID), data); err != nil {
			return err
		}

		// 创建索引
		if err := txn.Set(s.usernameIndexKey(user.UserName), []byte(user.ID)); err != nil {
			return err
		}

		if user.Email != "" {
			if err := txn.Set(s.emailIndexKey(user.Email), []byte(user.ID)); err != nil {
				return err
			}
		}

		return nil
	})
}

// 通过用户名获取用户
func (s *UserStore) GetUserByUsername(username string) (*User, error) {
	var userID string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.usernameIndexKey(username))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			userID = string(val)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return s.GetUserByID(userID)
}

// 通过ID获取用户
func (s *UserStore) GetUserByID(id string) (*User, error) {
	var user User
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(s.userKey(id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})
	return &user, err
}

// 更新用户信息
func (s *UserStore) UpdateUser(user *User) error {
	return s.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		return txn.Set(s.userKey(user.ID), data)
	})
}

// 软删除用户
func (s *UserStore) SoftDelete(userID string, deletedBy string) error {
	now := time.Now()
	return s.db.Update(func(txn *badger.Txn) error {
		user, err := s.GetUserByID(userID)
		if err != nil {
			return err
		}

		user.DeletedAt = &now
		user.DeletedBy = deletedBy
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return txn.Set(s.userKey(userID), data)
	})
}

// 保持原有的验证方法
func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}

// 保持原有的工具方法
func MD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
