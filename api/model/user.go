package model

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/skanehira/vue-go-oauth2/api/common"
)

// User user info
type User struct {
	ID          string     `gorm:"primary_key;not null" json:"id_str"`
	ScreenName  string     `gorm:"not null" json:"screen_name"`
	Name        string     `gorm:"not null" json:"name"`
	URL         string     `gorm:"not null" json:"url"`
	Description string     `gorm:"null" json:"description"`
	IsSignedIn  bool       `gorm:"not null" json:"is_signed_in"`
	CreatedAt   time.Time  `gorm:"null" json:"create_at"`
	UpdatedAt   time.Time  `gorm:"null" json:"update_at"`
	DeletedAt   *time.Time `gorm:"null" json:"-"`
}

// SaveUser save user info
func (m *Model) SaveUser(user User) error {
	user.UpdatedAt = common.GetTime()
	user.CreatedAt = common.GetTime()

	if err := m.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// GetUser get user info
func (m *Model) GetUser(id string) (User, error) {
	user := User{ScreenName: id}

	// get user info
	if err := m.db.Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return user, common.ErrNotFoundUserInfo
		}

		log.Println("error: " + err.Error())
		return user, err
	}

	return user, nil
}

// UpdateUser update user info
func (m *Model) UpdateUser(user User) (User, error) {
	// if user no exist
	oldUser, err := m.GetUser(user.ScreenName)
	if err != nil {
		return user, err
	}

	// update date time
	newTime := common.GetTime()
	user.UpdatedAt = newTime
	user.CreatedAt = oldUser.CreatedAt

	if err := m.db.Model(&user).Updates(user).Error; err != nil {
		return user, err
	}

	newUser, err := m.GetUser(user.ScreenName)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

// DeleteUser delete user
func (m *Model) DeleteUser(id string) error {

	// if user not exist
	user, err := m.GetUser(id)
	if err != nil {
		return err
	}

	// start transaction
	db := m.db.Begin()

	// delete user
	if err := db.Delete(&user).Error; err != nil {
		log.Println("error: " + err.Error())
		if err := db.Rollback().Error; err != nil {
			return err
		}

		return err
	}

	// db commit
	if err := db.Commit().Error; err != nil {
		return err
	}

	return nil
}
