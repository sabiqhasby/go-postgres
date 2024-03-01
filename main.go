package main

import (
	"fmt"
	"learn2/auth"
	"learn2/middleware"
	"log"
	"net/http"

	_ "github.com/lib/pq" //driver postgress
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type Student struct {
	Student_id       uint64 `json:"student_id"`
	Student_name     string `json:"student_name"`
	Student_age      uint64 `json:"student_age"`
	Student_address  string `json:"student_address"`
	Student_phone_no string `json:"student_phone_no"`
}

// func rowToStruct(rows *sql.Rows, dest interface{}) error {
// 	destv := reflect.ValueOf(dest).Elem()

// 	args := make([]interface{}, destv.Type().Elem().NumField())

// 	for rows.Next() {
// 		rowp := reflect.New(destv.Type().Elem())
// 		rowv := rowp.Elem()

//			for i := 0; i < rowv.NumField(); i++ {
//				args[i] = rowv.Field(i).Addr().Interface()
//			}
//			if err := rows.Scan(args...); err != nil {
//				return err
//			}
//			destv.Set(reflect.Append(destv, rowv))
//		}
//		return nil
//	}
// func rowToStruct(rows *sql.Rows, dest *[]Student) error {
// 	for rows.Next() {
// 		var student Student
// 		if err := rows.Scan(&student.Student_id, &student.Student_name, &student.Student_age, &student.Student_address, &student.Student_phone_no); err != nil {
// 			return err
// 		}
// 		*dest = append(*dest, student)
// 	}
// 	return nil
// }

func postHandler(c *gin.Context, db *gorm.DB) {
	var student Student

	// if c.Bind(&student) == nil {
	// 	_, err := db.Exec("insert into students values ($1, $2, $3, $4, $5)", student.Student_id, student.Student_name, student.Student_age, student.Student_address, student.Student_phone_no)
	// 	if err != nil {
	// 		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(http.StatusOK, gin.H{"message": "success create"})
	// 	return
	// }

	// c.JSON(http.StatusBadRequest, gin.H{"message": "something was error"})

	c.Bind(&student)
	db.Create(&student)
	c.JSON(http.StatusOK, gin.H{"message": "success created", "data": student})

}

func getHandlerByParam(c *gin.Context, db *gorm.DB) {
	// var student []Student

	// studentId := c.Param("id")
	// row, err := db.Query("select * from students where student_id=$1", studentId)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// rowToStruct(row, &student)
	// if student == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"data": student})
	// ==================================================================

	// db.Where("student_id = ?", studentId).First(&student)

	var student Student

	studentId := c.Param("id")

	if err := db.First(&student, "student_id = ?", studentId).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "found by id", "data": student})

}

func getAllHandler(c *gin.Context, db *gorm.DB) {
	var student []Student

	// row, err := db.Query("select * from students")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	// rowToStruct(row, &student)
	// if student == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{"data": student})
	// ==========================================
	db.Find(&student)
	c.JSON(http.StatusOK, gin.H{"message": "show all data", "data": student})

}

func updateByParam(c *gin.Context, db *gorm.DB) {
	// var student Student
	// studentId := c.Param("id")

	// 	if c.Bind(&student) == nil {
	// 		_, err := db.Exec("update students set student_name=$1 where student_id=$2", student.Student_name, studentId)
	// 		if err != nil {
	// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		}

	// 		c.JSON(http.StatusOK, gin.H{"message": "success updated"})
	// 	}

	var student Student
	bodyReq := Student{}
	studentId := c.Param("id")

	//check dahulu apakah id tersebut ada atau tidak
	if err := db.First(&student, studentId).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.Bind(&bodyReq)
	db.Model(&student).Where("student_id = ?", studentId).Updates(bodyReq)
	c.JSON(http.StatusOK, gin.H{"message": "success updated", "data": student})

}

func deleteHandler(c *gin.Context, db *gorm.DB) {
	var student Student
	studentId := c.Param("id")

	// 	_, err := db.Exec("delete from students where student_id=$1", studentId)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// ---------------------------------------------------------

	// db.Delete(&student, "student_id = ?", studentId)
	// sama saja dengan
	db.Delete(&student, studentId)
	c.JSON(http.StatusOK, gin.H{"message": "success deleted"})

}

func setupRouter() *gin.Engine {

	conn := "user=postgres password='semogaberkah' dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	Migrate(db)

	r.POST("/login", auth.LoginHandler)

	r.POST("/students", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})
	r.GET("/students", middleware.AuthValid, func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})
	r.GET("/students/:id", func(ctx *gin.Context) {
		getHandlerByParam(ctx, db)
	})
	r.PUT("/students/:id", func(ctx *gin.Context) {
		updateByParam(ctx, db)
	})
	r.DELETE("/students/:id", func(ctx *gin.Context) {
		deleteHandler(ctx, db)
	})

	return r

}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Student{})

	data := Student{}

	// err := db.First(&data).Error
	// if errors.Is(err, gorm.ErrRecordNotFound) {
	// 	fmt.Println("======seeder running=====")
	// 	seederStudent(db)
	// }

	if err := db.First(&data).Error; err == gorm.ErrRecordNotFound {
		fmt.Println("======seeder running=====")
		seederStudent(db)

	}
	fmt.Println("======database is not empty=======")
}

func seederStudent(db *gorm.DB) {
	data := Student{
		Student_id:       1,
		Student_name:     "student1",
		Student_age:      20,
		Student_address:  "Jalan kenangan",
		Student_phone_no: "0909099",
	}
	db.Create(&data)
}

func main() {
	server := setupRouter()

	server.Run(":8080")
}
