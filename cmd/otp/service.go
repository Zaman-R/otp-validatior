package otp

import (
	"errors"
	"fmt"
	"github.com/Zaman-R/otp-validator/cmd/client"
	"time"

	"github.com/Zaman-R/otp-validator/cmd/utils"
	"github.com/google/uuid"
)

// OTPService handles OTP generation, validation, and sending.
type OTPService struct {
	repo          *OTPRepository
	smsProvider   client.SMSProvider
	emailProvider client.EmailProvider
}

// NewOTPService initializes a new OTPService.
func NewOTPService(repo *OTPRepository, smsProvider client.SMSProvider, emailProvider client.EmailProvider) *OTPService {
	return &OTPService{
		repo:          repo,
		smsProvider:   smsProvider,
		emailProvider: emailProvider,
	}
}

func (s *OTPService) IsOTPExpired(otp OTP) bool {
	return time.Now().After(otp.ExpiresAt)
}

func (s *OTPService) GenerateOTP(email, phone, purpose string, retryLimit, expiryMinutes int, transactionPayload map[string]interface{}) (*OTP, string, error) {
	rawOTP, err := utils.GenerateSecureOTP()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate OTP: %v", err)
	}

	hashedOTP, err := utils.HashOTP(rawOTP)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash OTP: %v", err)
	}

	deliveryMethod := utils.DetermineDeliveryMethod(email, phone)
	if deliveryMethod == "" {
		return nil, "", errors.New("no valid delivery method provided (email or phone required)")
	}

	var encodedPayload string
	if purpose == "transaction" && transactionPayload != nil {
		encodedPayload, err = utils.EncodeBase64(transactionPayload)
		if err != nil {
			return nil, "", fmt.Errorf("failed to encode transaction payload: %v", err)
		}
	}

	otp := &OTP{
		ID:                 uuid.New(),
		Purpose:            purpose,
		HashedOTP:          hashedOTP,
		Delivery:           deliveryMethod,
		MobileNumber:       phone,
		Email:              email,
		RetryLimit:         retryLimit,
		ExpiresAt:          time.Now().Add(time.Duration(expiryMinutes) * time.Minute),
		Status:             OTPStatusPending,
		TransactionPayload: encodedPayload,
	}

	if err := s.repo.SaveOTP(otp); err != nil {
		return nil, "", fmt.Errorf("failed to save OTP: %v", err)
	}

	return otp, rawOTP, nil
}

func (s *OTPService) CanRetry(otp OTP) bool {
	return otp.RetryCount < otp.RetryLimit
}

type SendOTPRequest struct {
	FromAccount  string
	Payload      map[string]interface{}
	Length       int
	RetryLimit   int
	Expiration   time.Duration
	MobileNumber *string
	Email        *string
	SMSBody      *string
	EmailSubject *string
	EmailBody    *string
}

func (s *OTPService) SendOTPFromParams(params map[string]interface{}) (string, error) {
	request := SendOTPRequest{
		FromAccount:  utils.GetString(params, "from_account"),
		Payload:      utils.GetMap(params, "payload"),
		Length:       utils.GetInt(params, "length", 6),
		RetryLimit:   utils.GetInt(params, "retry_limit", 3),
		Expiration:   utils.GetDuration(params, "expiration", 5*time.Minute),
		MobileNumber: utils.GetStringPtr(params, "mobile_number"),
		Email:        utils.GetStringPtr(params, "email"),
		SMSBody:      utils.GetStringPtr(params, "sms_body"),
		EmailSubject: utils.GetStringPtr(params, "email_subject"),
		EmailBody:    utils.GetStringPtr(params, "email_body"),
	}

	return s.SendOTP(request)
}

func (s *OTPService) SendOTP(req SendOTPRequest) (string, error) {
	if req.MobileNumber == nil && req.Email == nil {
		return "", errors.New("please provide a valid mobile_number, email, or both")
	}
	if req.MobileNumber != nil && (req.SMSBody == nil || !utils.Contains(*req.SMSBody, "<otp>")) {
		return "", errors.New("invalid SMS body format, missing `<otp>` placeholder")
	}
	if req.Email != nil && (req.EmailBody == nil || !utils.Contains(*req.EmailBody, "<otp>")) {
		return "", errors.New("invalid Email body format, missing `<otp>` placeholder")
	}

	otp, rawOTP, err := s.GenerateOTP(
		utils.GetStringValue(req.Email),
		utils.GetStringValue(req.MobileNumber),
		req.FromAccount,
		req.RetryLimit,
		int(req.Expiration.Minutes()),
		req.Payload,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate OTP: %v", err)
	}

	smsBody := utils.ReplacePlaceholder(req.SMSBody, rawOTP)
	emailBody := utils.ReplacePlaceholder(req.EmailBody, rawOTP)

	payload := map[string]interface{}{"otp_ref": otp.ID}
	token, err := utils.GenerateToken(payload, int(req.Expiration.Seconds()))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	if req.MobileNumber != nil && s.smsProvider != nil {
		err = s.smsProvider.SendSMS(*req.MobileNumber, smsBody)
		if err != nil {
			fmt.Printf("Failed to send SMS: %v\n", err)
		}
	}

	if req.Email != nil && s.emailProvider != nil {
		err = s.emailProvider.SendEmail(*req.Email, emailBody)
		if err != nil {
			fmt.Printf("Failed to send Email: %v\n", err)
		}
	}

	return token, nil
}

func (s *OTPService) ValidateOTP(otpCode string, payloadToken string) (map[string]interface{}, error) {
	payload, err := utils.ValidateToken(payloadToken)
	if err != nil {
		return nil, fmt.Errorf("invalid token provided: %v", err)
	}
	otpRefStr, ok := payload["otp_ref"].(string)
	if !ok {
		return nil, errors.New("invalid otp_ref format")
	}

	otpRef, err := uuid.Parse(otpRefStr)
	if err != nil {
		return nil, errors.New("invalid UUID format for otp_ref")
	}

	otpInstance, err := s.repo.GetOTPByID(otpRef)
	if err != nil || otpInstance == nil {
		return nil, errors.New("OTP not found")
	}

	if otpInstance.Status != OTPStatusPending {
		return nil, errors.New("OTP is no longer valid")
	}

	if otpInstance.RetryCount >= otpInstance.RetryLimit {
		_ = s.repo.ExpireOTP(otpRef)
		return nil, errors.New("maximum retry attempts reached, OTP expired")
	}

	if time.Now().After(otpInstance.ExpiresAt) {
		_ = s.repo.UpdateRetryLimit(otpRef)
		return nil, errors.New("OTP expired")
	}

	if !utils.ValidateOTP(otpCode, otpInstance.HashedOTP, otpInstance.CreatedAt, otpInstance.ExpiresAt) {
		_ = s.repo.UpdateRetryLimit(otpRef)
		return nil, errors.New("invalid OTP provided")
	}

	var sanitizedPayload map[string]interface{}
	if otpInstance.TransactionPayload != "" {
		transactionPayload, err := utils.DecodeBase64(otpInstance.TransactionPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to decode transaction payload: %v", err)
		}

		sanitizedPayload, err = utils.SanitizePayload(transactionPayload)
		if err != nil {
			return nil, fmt.Errorf("failed to sanitize transaction payload: %v", err)
		}
	}

	err = s.repo.UpdateOTPStatus(otpRef, OTPStatusVerified)
	if err != nil {
		return nil, fmt.Errorf("failed to update OTP status: %v", err)
	}

	if otpInstance.Purpose == "login" || otpInstance.Purpose == "register" {
		return map[string]interface{}{"status": "OTP verified"}, nil
	}

	return sanitizedPayload, nil
}
