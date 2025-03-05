package repository

import (
	_ "errors"
	"github.com/Zaman-R/otp-validator/cmd/otp"
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

func (r *OTPRepository) SaveOTP(otp *otp.OTP) error {
	otp.ID = uuid.New()
	otp.CreatedAt = time.Now()
	otp.UpdatedAt = time.Now()
	return r.db.Create(otp).Error
}

func (r *OTPRepository) GetValidOTPByPurpose(mobileOrEmail, purpose string) (*otp.OTP, error) {
	var otpInstance otp.OTP
	err := r.db.Where("(mobile_number = ? OR email = ?) AND purpose = ? AND status = ?",
		mobileOrEmail, mobileOrEmail, purpose, otp.OTPStatusPending).
		First(&otpInstance).Error
	if err != nil {
		return nil, err
	}
	return &otpInstance, nil
}

func (r *OTPRepository) IncrementRetryCount(id uuid.UUID) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error
}

func (r *OTPRepository) ExpireOTP(id uuid.UUID) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     otp.OTPStatusExpired,
			"updated_at": time.Now(),
		}).Error
}

func (r *OTPRepository) UpdateRetryLimit(id uuid.UUID) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + ?", 1)).Error
}

func (r *OTPRepository) MarkOTPAsUsed(id uuid.UUID) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", id).
		Update("status", otp.OTPStatusUsed).Error
}

func (r *OTPRepository) GetOTPByID(otpID uuid.UUID) (*otp.OTP, error) {
	var otpInstance otp.OTP
	if err := r.db.Where("id = ?", otpID).First(&otpInstance).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("OTP not found")
		}
		return nil, err
	}
	return &otpInstance, nil
}

func (r *OTPRepository) UpdateOTPStatus(otpID uuid.UUID, status string) error {
	return r.db.Model(&otp.OTP{}).Where("id = ?", otpID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
