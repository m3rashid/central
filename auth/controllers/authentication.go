package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func PasswordLogin(ctx *fiber.Ctx) error {
	type Login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var login Login
	if err := ctx.BodyParser(&login); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return err
	}

	var user models.User
	if err := db.Where("email = ?", login.Email).First(&user).Error; err != nil {
		return err
	}

	if err := utils.ComparePassword(user.Password, login.Password); err != nil {
		return err
	}

	// TODO

	return nil
}

func Register(ctx *fiber.Ctx) error {
	type Register struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	var register Register
	if err := ctx.BodyParser(&register); err != nil {
		return err
	}

	db, err := utils.GetDb()
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(register.Password)
	if err != nil {
		return err
	}

	// create a new OTP and send it to the user's email
	otp := utils.GenerateOTP()

	db.Create(&models.User{
		OTP:      otp,
		Name:     register.Name,
		Email:    register.Email,
		Password: hashedPassword,
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

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
