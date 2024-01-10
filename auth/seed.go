package main

import (
	"internal/helpers"
	"log"

	"github.com/m3rashid/central/auth/models"
	"github.com/m3rashid/central/auth/utils"
	"gorm.io/gorm"
)

func initialSeedDatabase() error {
	hashedPassword, err := helpers.HashPassword("central123")
	if err != nil {
		log.Println("Hash password failed")
		return err
	}

	// create default user
	adminUser := models.User{
		Name:         "MD Rashid Hussain",
		Email:        "m3rashid.hussain@gmail.com",
		Password:     hashedPassword,
		OTP:          helpers.GenerateOTP(),
		UserVerified: true,
	}

	// create default permissions
	permissions := []*models.Permission{
		{Name: "read"}, {Name: "create"}, {Name: "update"}, {Name: "create"},
	}

	// create default scope
	scopes := []*models.Scope{
		{Name: "openid", Description: "", Permissions: permissions},
		{Name: "contacts", Description: "", Permissions: permissions},
	}

	client := models.Client{
		ClientID:           "client-id",
		ClientSecret:       "sample-client-secret",
		Scopes:             scopes,
		SuccessRedirectUri: "http://localhost:5001/success",
		FailureRedirectUri: "http://localhost:5001/failure",
		AppName:            "Contacts",
		AppLogoUrl:         "http://localhost:5000/public/icons/favicon.ico",
	}

	db, err := utils.GetDb()
	if err != nil {
		log.Println("Could not get database")
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&adminUser).Error; err != nil {
			log.Println("Could not create user")
			return err
		}

		if err := tx.Create(permissions).Error; err != nil {
			log.Println("Could not create permissions")
			return err
		}

		if err := tx.Create(scopes).Error; err != nil {
			log.Println("Could not create scopes")
			return err
		}

		if err := tx.Create(&client).Error; err != nil {
			log.Println("Could not create client")
			return err
		}

		log.Println("Database Seed Successful")
		return nil
	})
}
