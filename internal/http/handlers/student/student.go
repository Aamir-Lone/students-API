package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Aamir-Lone/students-API/internal/storage"
	"github.com/Aamir-Lone/students-API/internal/types"
	"github.com/Aamir-Lone/students-API/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating a student")
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {

			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}
		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)
		slog.Info("user created succesfully", slog.String("userId", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return

		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("getting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("error getting user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return

		}
		response.WriteJson(w, http.StatusOK, student)

	}

}

func GetList(storage storage.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return

		}

		response.WriteJson(w, http.StatusOK, students)

	}

}

func Delete(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("deleting a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid ID format: %w", err)))
			return
		}

		err = storage.DeleteStudentById(intId)
		if err != nil {
			if errors.Is(err, errors.New("record not found")) {
				slog.Error("student not found", slog.String("id", id))
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student with ID %d not found", intId)))
				return
			}
			slog.Error("error deleting student", slog.String("id", id), slog.Any("error", err))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to delete student: %w", err)))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student deleted successfully"})
	}
}
func Update(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		slog.Info("updating a student", slog.String("id", id))
		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid ID format: %w", err)))
			return
		}

		var updatedStudent types.Student
		if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid JSON body: %w", err)))
			return
		}

		updatedStudent.Id = intId // Ensure the ID is set correctly in the struct

		err = storage.UpdateStudentById(intId, updatedStudent)
		if err != nil {
			if errors.Is(err, errors.New("record not found")) {
				slog.Error("student not found for update", slog.String("id", id))
				response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student with ID %d not found", intId)))
				return
			}
			slog.Error("error updating student", slog.String("id", id), slog.Any("error", err))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to update student: %w", err)))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student updated successfully"})
	}
}
