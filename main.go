package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq" //driver postgress

	"github.com/gin-gonic/gin"
)

type Student struct {
	Student_id       uint64 `json:"student_id" binding:"required"`
	Student_name     string `json:"student_name" binding:"required"`
	Student_age      uint64 `json:"student_age" binding:"required"`
	Student_address  string `json:"student_address" binding:"required"`
	Student_phone_no string `json:"student_phone_no" binding:"required"`
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
func rowToStruct(rows *sql.Rows, dest *[]Student) error {
	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.Student_id, &student.Student_name, &student.Student_age, &student.Student_address, &student.Student_phone_no); err != nil {
			return err
		}
		*dest = append(*dest, student)
	}
	return nil
}
func postHandler(c *gin.Context, db *sql.DB) {
	var student Student

	if c.Bind(&student) == nil {
		_, err := db.Exec("insert into students values ($1, $2, $3, $4, $5)", student.Student_id, student.Student_name, student.Student_age, student.Student_address, student.Student_phone_no)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "success create"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"message": "something was error"})
}

func getHandlerByParam(c *gin.Context, db *sql.DB) {
	var student []Student

	studentId := c.Param("id")
	row, err := db.Query("select * from students where student_id=$1", studentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowToStruct(row, &student)
	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})
}
func getAllHandler(c *gin.Context, db *sql.DB) {
	var student []Student

	row, err := db.Query("select * from students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowToStruct(row, &student)
	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": student})

}

func updateByParam(c *gin.Context, db *sql.DB) {
	var student Student
	studentId := c.Param("id")

	if c.Bind(&student) == nil {
		_, err := db.Exec("update students set student_name=$1 where student_id=$2", student.Student_name, studentId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "success updated"})
	}

}

func deleteHandler(c *gin.Context, db *sql.DB) {
	studentId := c.Param("id")

	_, err := db.Exec("delete from students where student_id=$1", studentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success deleted"})
}

func setupRouter() *gin.Engine {

	conn := "user=postgres password='semogaberkah' dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	r.POST("/student", func(ctx *gin.Context) {
		postHandler(ctx, db)
	})
	r.GET("/student", func(ctx *gin.Context) {
		getAllHandler(ctx, db)
	})
	r.GET("/student/:id", func(ctx *gin.Context) {
		getHandlerByParam(ctx, db)
	})
	r.PUT("/student/:id", func(ctx *gin.Context) {
		updateByParam(ctx, db)
	})
	r.DELETE("/student/:id", func(ctx *gin.Context) {
		deleteHandler(ctx, db)
	})

	return r

}

func main() {
	server := setupRouter()

	server.Run(":8080")
}
