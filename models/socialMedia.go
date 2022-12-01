package models

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	GormModel
	Name           string `json:"name" form:"name" valid:"required~Name is required"`
	SocialMediaUrl string `json:"social_media_url" form:"social_media_url" valid:"required~Url of your social media is required"`
	UserID         uint
	User           *User
}

func (s *SocialMedia) BeforeCreate(tx *gorm.DB) (err error) {
	_, errCreate := govalidator.ValidateStruct(s)

	if errCreate != nil {
		err = errCreate
		return
	}

	err = nil
	return

}
