package storage

import "github.com/Aamir-Lone/students-API/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	DeleteStudentById(id int64) error
	UpdateStudentById(id int64, updatedStudent types.Student) error
}
