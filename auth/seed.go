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

	db, err := utils.GetDb()
	if err != nil {
		log.Println("Could not get database")
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		defaultUser := models.User{
			Name:         "MD Rashid Hussain",
			Email:        "m3rashid.hussain@gmail.com",
			Password:     hashedPassword,
			OTP:          helpers.GenerateOTP(),
			UserVerified: true,
		}

		if err := tx.Create(&defaultUser).Error; err != nil {
			log.Println("Could not create user")
			return err
		}

		client := models.Client{
			ClientID:           "client-id",
			ClientSecret:       "sample-client-secret",
			Scopes:             helpers.JSONB{"contacts": []string{"read", "create"}},
			SuccessRedirectUri: "http://localhost:5001/success",
			FailureRedirectUri: "http://localhost:5001/failure",
			AppName:            "Campaigns",
			AppLogoUrl:         "http://localhost:5000/public/icons/favicon.ico",
			CreatedByUserID:    defaultUser.ID,
		}

		if err := tx.Create(&client).Error; err != nil {
			log.Println("Could not create client")
			return err
		}

		log.Println("Database Seed Successful")
		return nil
	})
}
