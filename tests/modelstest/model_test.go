package modelstest

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/amaraliou/apetitoso/handlers"
	"github.com/amaraliou/apetitoso/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = handlers.Server{}
var studentInstance = models.Student{}
var adminInstance = models.Admin{}
var shopInstance = models.Shop{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		fmt.Printf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))
	server.DB, err = gorm.Open("postgres", DBURL)
	if err != nil {
		fmt.Print("Cannot connect to Postgres database\n")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Print("We are connected to the Postgres database\n")
	}
}

func refreshStudentTable() error {
	err := server.DB.DropTableIfExists(&models.Student{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Student{}).Error
	if err != nil {
		return err
	}
	return nil
}

func refreshAdminTable() error {
	err := server.DB.DropTableIfExists(&models.Admin{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Admin{}).Error
	if err != nil {
		return err
	}
	return nil
}

func refreshShopTable() error {
	err := server.DB.DropTableIfExists(&models.Shop{}).Error
	if err != nil {
		return err
	}

	err = server.DB.AutoMigrate(&models.Shop{}).Error
	if err != nil {
		return err
	}
	return nil
}

func seedOneStudent() (models.Student, error) {

	refreshStudentTable()

	student := models.Student{
		User: models.User{
			Email:      "email@email.com",
			Password:   "password",
			IsVerified: true,
		},
		IsStudent:      true,
		FirstName:      "Donald",
		LastName:       "Trump",
		BirthDate:      "09/09/1956",
		University:     "Temple University",
		MobileNumber:   "07547775660",
		CountryCode:    "GB",
		GraduationYear: 2021,
		Points:         0,
	}

	err := server.DB.Model(&models.Student{}).Create(&student).Error
	if err != nil {
		return models.Student{}, err
	}
	return student, nil
}

func seedStudents() error {

	var students = []models.Student{
		models.Student{
			User: models.User{
				Email:      "email1@email.com",
				Password:   "password",
				IsVerified: true,
			},
			IsStudent:      true,
			FirstName:      "Donald",
			LastName:       "Trump",
			BirthDate:      "09/09/1956",
			University:     "Temple University",
			MobileNumber:   "07547775660",
			CountryCode:    "GB",
			GraduationYear: 2021,
			Points:         0,
		},
		models.Student{
			User: models.User{
				Email:      "email2@email.com",
				Password:   "password",
				IsVerified: true,
			},
			IsStudent:      true,
			FirstName:      "Bernie",
			LastName:       "Sanders",
			BirthDate:      "09/09/1956",
			University:     "University of Vermont",
			MobileNumber:   "07547775660",
			CountryCode:    "GB",
			GraduationYear: 2021,
			Points:         0,
		},
	}

	for i := range students {
		err := server.DB.Model(&models.Student{}).Create(&students[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func seedOneAdmin() (models.Admin, error) {

	refreshAdminTable()

	admin := models.Admin{
		User: models.User{
			Email:      "email@email.com",
			Password:   "password",
			IsVerified: true,
		},
		FirstName: "Donald WW3",
		LastName:  "Trump",
	}

	err := server.DB.Model(&models.Admin{}).Create(&admin).Error
	if err != nil {
		return models.Admin{}, err
	}
	return admin, nil
}

func seedAdmins() error {

	refreshAdminTable()

	admins := []models.Admin{
		models.Admin{
			User: models.User{
				Email:      "email1@email.com",
				Password:   "password",
				IsVerified: true,
			},
			FirstName: "Donald WW3",
			LastName:  "Trump",
		},
		models.Admin{
			User: models.User{
				Email:      "email2@email.com",
				Password:   "password",
				IsVerified: true,
			},
			FirstName: "Kim Jong",
			LastName:  "Un",
		},
	}

	for i := range admins {
		err := server.DB.Model(&models.Admin{}).Create(&admins[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func seedOneShop() (models.Shop, error) {

	refreshShopTable()

	shop := models.Shop{
		Name:        "Starbucks",
		Logo:        "https://starbucks.com",
		Description: "blah blah blah",
		Latitude:    45.5,
		Longitude:   45.5,
		ShopAddress: models.ShopAddress{
			Postcode:      "G12 8BG",
			AddressNumber: 10,
			AddressLine1:  "Bruh Street",
			TownOrCity:    "Bruh Town",
		},
	}

	err := server.DB.Model(&models.Shop{}).Create(&shop).Error
	if err != nil {
		return models.Shop{}, err
	}
	return shop, nil
}

func seedShops() error {

	refreshShopTable()

	shops := []models.Shop{
		models.Shop{
			Name:        "Starbucks 1",
			Logo:        "https://starbucks-one.com",
			Description: "blah blah blah",
			Latitude:    45.5,
			Longitude:   45.5,
			ShopAddress: models.ShopAddress{
				Postcode:      "G12 8BG",
				AddressNumber: 10,
				AddressLine1:  "Bruh Street",
				TownOrCity:    "Bruh Town",
			},
		},
		models.Shop{
			Name:        "Starbucks 2",
			Logo:        "https://starbucks-two.com",
			Description: "blah blah blah",
			Latitude:    55.5,
			Longitude:   55.5,
			ShopAddress: models.ShopAddress{
				Postcode:      "G12 8BG",
				AddressNumber: 10,
				AddressLine1:  "Bruh Avenue",
				TownOrCity:    "Bruh City",
			},
		},
	}

	for i := range shops {
		err := server.DB.Model(&models.Shop{}).Create(&shops[i]).Error
		if err != nil {
			return err
		}
	}

	return nil
}
