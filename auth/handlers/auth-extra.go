package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func VerifyEmail(ctx *fiber.Ctx) error {
	type VerifyEmail struct {
		Email string `json:"email" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}

	var verifyEmail VerifyEmail
	if err := ctx.BodyParser(&verifyEmail); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return err
	}

	var user models.User
	if err := db.Where("email = ?", verifyEmail.Email).First(&user).Error; err != nil {
		return err
	}

	if user.OTP != verifyEmail.OTP {
		return fiber.ErrUnauthorized
	}

	// TODO

	return nil
}

func ResendVerificationEmail(ctx *fiber.Ctx) error {
	type ResendVerificationEmail struct {
		Email string `json:"email" validate:"required"`
	}

	var resendVerificationEmail ResendVerificationEmail
	if err := ctx.BodyParser(&resendVerificationEmail); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return nil
	}

	var user models.User
	if err := db.Where("email = ?", resendVerificationEmail.Email).First(&user).Error; err != nil {
		return err
	}

	if user.UserVerified {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "User already verified",
		})
	}

	// send the otp to the user's email

	return nil
}

func ForgotPassword(ctx *fiber.Ctx) error {
	type ForgotPassword struct {
		Email string `json:"email" validate:"required"`
	}

	var forgotPassword ForgotPassword
	if err := ctx.BodyParser(&forgotPassword); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return nil
	}

	var user models.User
	if err := db.Where("email = ?", forgotPassword.Email).First(&user).Error; err != nil {
		return err
	}

	// create a new OTP and send it to the user's email
	otp := utils.GenerateOTP()
	db.Model(&user).Update("otp", otp)

	// send the otp to the user's email

	return nil
}

func ResetPassword(ctx *fiber.Ctx) error {
	type ResetPassword struct {
		Email string `json:"email" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}

	var resetPassword ResetPassword
	if err := ctx.BodyParser(&resetPassword); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return err
	}

	var user models.User
	if err := db.Where("email = ?", resetPassword.Email).First(&user).Error; err != nil {
		return err
	}

	if user.OTP != resetPassword.OTP {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid OTP",
		})
	}

	// TODO

	return nil
}
