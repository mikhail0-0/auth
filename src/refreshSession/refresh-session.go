package refreshSession

import (
	"auth/src/common"
	"auth/src/user"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshSession struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	User        user.User `gorm:"constraint"`
	UserId      string    `gorm:"type:varchar"`
	RefreshHash string    `gorm:"type:varchar"`
	Ip          string    `gorm:"type:varchar"`
	ExpiresIn   time.Time `gorm:"type:timestamptz"`
}

var db *gorm.DB

func Init(refDB *gorm.DB) {
	db = refDB
	if !db.Migrator().HasTable(&RefreshSession{}) {
		db.Migrator().CreateTable(&RefreshSession{})
	}
}

func Create(
	userId, refreshString, ip string,
	expiresIn int64,
) (*RefreshSession, error) {
	refreshHash, err := common.GetHash(refreshString)
	if err != nil {
		return nil, err
	}

	refreshSession := RefreshSession{
		UserId:      userId,
		RefreshHash: refreshHash,
		Ip:          ip,
		ExpiresIn:   time.Unix(expiresIn, 0),
	}
	result := db.Create(&refreshSession)

	if result.Error != nil {
		return nil, result.Error
	}

	return &refreshSession, nil
}

func FindByUserId(userId string) (*RefreshSession, error) {
	var refreshSession RefreshSession
	result := db.First(&refreshSession, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}

	return &refreshSession, nil
}

func Delete(rs *RefreshSession) error {
	result := db.Delete(rs)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func CheckRefresh(guid, refreshToken string) (*RefreshSession, error) {
	targetSession, err := FindByUserId(guid)
	if err != nil {
		return nil, err
	}

	err = common.CompareHashAndString(targetSession.RefreshHash, refreshToken)
	if err != nil {
		return nil, common.ErrWrongRefreshToken
	}

	if time.Now().Unix() > targetSession.ExpiresIn.Unix() {
		return nil, common.ErrRefreshTokenExpired
	}

	return targetSession, nil
}
