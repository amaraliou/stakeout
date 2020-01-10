package utils

import (
	"log"

	"github.com/amaraliou/apetitoso/models"
	"github.com/jinzhu/gorm"
)

var students = []models.Student{
	models.Student{
		User: models.User{
			Email:      "2310549a@student.gla.ac.uk",
			Password:   "password",
			IsVerified: true,
		},
		IsStudent:      true,
		FirstName:      "Aliou",
		LastName:       "Amar",
		BirthDate:      "10/09/1997",
		University:     "University of Glasgow",
		MobileNumber:   "07547775660",
		CountryCode:    "GB",
		GraduationYear: 2021,
		Points:         0,
	},
}

var admins = []models.Admin{
	models.Admin{
		User: models.User{
			Email:      "admin@gmail.com",
			Password:   "password",
			IsVerified: true,
		},
		FirstName: "Admin",
		LastName:  "Admin",
	},
}

// Load ... making my linter happy
func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Student{}, &models.Admin{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.Student{}, &models.Admin{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	for i := range students {
		err = db.Debug().Model(&models.Student{}).Create(&students[i]).Error
		if err != nil {
			log.Fatalf("cannot seed students table: %v", err)
		}
	}

	for i := range admins {
		err = db.Debug().Model(&models.Admin{}).Create(&admins[i]).Error
		if err != nil {
			log.Fatalf("cannot seed admins table: %v", err)
		}
	}
}
