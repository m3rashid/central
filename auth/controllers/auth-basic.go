package controllers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/m3rashid/central/auth/components"
	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
)

func RenderLoginScreen(ctx *fiber.Ctx) error {
	ctx.Set("Content-Type", "text/html")
	component := components.Login()
	return component.Render(ctx.Context(), ctx.Response().BodyWriter())
}

func LoginWithPassword(ctx *fiber.Ctx) error {
	type Login struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}
	var login Login
	if err := ctx.BodyParser(&login); err != nil {
		return err
	}

	validate := validator.New()
	err := validate.Struct(login)
	if err != nil {
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

	ctx.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "", // TODO
		Domain:  "localhost",
		Expires: time.Now().Add(5 * time.Hour),
	})

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
