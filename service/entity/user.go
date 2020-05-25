package entity

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// 用户
type User struct {
	Entity

	Email string `gorm:"unique_index"`      //用户邮箱
	Phone string                            // 用户手机号
	Password string                         // 用户密码
	NickName string                         // 用户昵称
	Avatar string                           // 用户头像
	IsAdmin bool                            // 判断是否超级管理员
}

// 校验用户密码
func (user *User) VerificationPassword(password string) bool {
	if user.Password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// 保存用户
func (user *User) Save() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("密码加密出错: " + err.Error())
	}
	user.Password = string(hash)
	return store[UserStore].Create(user).Error
}

// 根据 ID 更新用户（用户名称、用户头像）
func (user *User) Update() error {
	return store[UserStore].Model(user).Updates(map[string]interface{}{
		"nick_name": user.NickName,
		"avatar": user.Avatar,
	}).Error
}

// 根据 ID 加载用户信息
func (user *User) Load() error {
	return store[UserStore].First(user).Error
}

// 根据用户邮箱号加载用户
func (user *User) LoadByEmail() error {
	return store[UserStore].First(user, "email = ?", user.Email).Error
}

// 判断系统是否有超级管理员
func HasAdmin() bool {
	var count int
	store[UserStore].Model(&User{}).Where("is_admin = ?", true).Count(&count)
	return count > 0
}

