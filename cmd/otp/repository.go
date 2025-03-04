package otp

import (
	_ "errors"
	"github.com/pkg/errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OTPRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) *OTPRepository {
	return &OTPRepository{db: db}
}

func (r *OTPRepository) SaveOTP(otp *OTP) error {
	otp.ID = uuid.New()
	otp.CreatedAt = time.Now()
	otp.UpdatedAt = time.Now()
	return r.db.Create(otp).Error
}

func (r *OTPRepository) GetValidOTPByPurpose(mobileOrEmail, purpose string) (*OTP, error) {
	var otp OTP
	err := r.db.Where("(mobile_number = ? OR email = ?) AND purpose = ? AND status = ?",
		mobileOrEmail, mobileOrEmail, purpose, OTPStatusPending).
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) IncrementRetryCount(id uuid.UUID) error {
	return r.db.Model(&OTP{}).Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

func (r *OTPRepository) ExpireOTP(id uuid.UUID) error {
	return r.db.Model(&OTP{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     OTPStatusExpired,
			"updated_at": time.Now(),
		}).Error
}

func (r *OTPRepository) UpdateRetryLimit(id uuid.UUID) error {
	return r.db.Model(&OTP{}).Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r *OTPRepository) MarkOTPAsUsed(id uuid.UUID) error {
	return r.db.Model(&OTP{}).Where("id = ?", id).
		Update("status", OTPStatusUsed).Error
}

func (r *OTPRepository) GetOTPByID(otpID uuid.UUID) (*OTP, error) {
	var otp OTP
	if err := r.db.Where("id = ?", otpID).First(&otp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("OTP not found")
		}
		return nil, err
	}
	return &otp, nil
}

func (r *OTPRepository) UpdateOTPStatus(otpID uuid.UUID, status string) error {
	return r.db.Model(&OTP{}).Where("id = ?", otpID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
