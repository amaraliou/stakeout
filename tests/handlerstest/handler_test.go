package handlerstest

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

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
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
	log.Printf("Successfully refreshed table")
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
		log.Fatalf("cannot seed students table: %v", err)
	}
	return student, nil
}

func seedStudents() ([]models.Student, error) {

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
			return []models.Student{}, err
		}
	}

	return students, nil
}
