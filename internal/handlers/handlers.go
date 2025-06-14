// Package handlers
package handlers

import (
	"errors"
	"log"
	"playcorner-be/internal/auth"
	"playcorner-be/internal/database"
	"playcorner-be/internal/models"
	"playcorner-be/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Login handles user login
func Login(c *fiber.Ctx) error {
	var body models.LoginBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code: 400, Status: "BAD_REQUEST", Data: models.ErrorData{ErrorMsg: "Cannot parse JSON"},
		})
	}

	var user models.User
	if err := database.DB.Where("id = ?", body.Identifier).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Code: 401, Status: "UNAUTHORIZED", Data: models.ErrorData{ErrorMsg: "Invalid identifier or password"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Database error"},
		})
	}

	if !utils.CheckPasswordHash(body.Password, user.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code: 401, Status: "UNAUTHORIZED", Data: models.ErrorData{ErrorMsg: "Invalid identifier or password"},
		})
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Could not generate tokens"},
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		HTTPOnly: true,
		Secure:   false, // Set true in production with HTTPS
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data: models.TokenCarrier{
			AuthToken:  accessToken,
			UserID:     user.ID,
			ExpireDate: time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		},
	})
}

// RefreshToken handles generating a new access token
func RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code: 401, Status: "UNAUTHORIZED", Data: models.ErrorData{ErrorMsg: "Refresh token not found"},
		})
	}

	claims, err := auth.ValidateToken(refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code: 401, Status: "UNAUTHORIZED", Data: models.ErrorData{ErrorMsg: "Invalid or expired refresh token"},
		})
	}

	accessToken, _, err := auth.GenerateTokens(claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Could not generate access token"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data: models.TokenCarrier{
			AuthToken:  accessToken,
			UserID:     claims.UserID,
			ExpireDate: time.Now().Add(15 * time.Minute).Format(time.RFC3339),
		},
	})
}

// GetUser retrieves a user's details
func GetUser(c *fiber.Ctx) error {
	userID := c.Params("userId")
	var user models.User

	if err := database.DB.Select("id", "name", "faculty", "major", "credit_score", "profile_pict_url").First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Code: 404, Status: "NOT_FOUND", Data: models.ErrorData{ErrorMsg: "User not found"},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Database error"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}

// GetUserHistories retrieves a user's reservation history
func GetUserHistories(c *fiber.Ctx) error {
	userID := c.Params("userId")
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil {
		offset = 0
	}

	var reservations []models.Reservation
	var total int64

	if err := database.DB.Model(&models.Reservation{}).Where("borrower_id = ?", userID).Count(&total).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:   500,
			Status: "SERVER_ERROR",
			Data:   models.ErrorData{ErrorMsg: "Could not count user histories"},
		})
	}

	if err := database.DB.Where("borrower_id = ?", userID).Order("created_at desc").Limit(int(limit)).Offset(int(offset)).Find(&reservations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:   500,
			Status: "SERVER_ERROR",
			Data:   models.ErrorData{ErrorMsg: "Could not fetch user histories"},
		})
	}

	histories := []models.History{}
	for _, r := range reservations {
		histories = append(histories, models.History{
			ID:                  r.ID,
			TVID:                r.TVID,
			ReservationDateTime: r.TimeSlot,
			TVPictURL:           "https://placehold.co/600x400/?text=TV+" + strconv.Itoa(r.TVID),
		})
	}

	pagedData := models.PagedData{
		Offset: offset,
		Limit:  limit,
		Total:  total,
		Data:   histories,
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data:   pagedData,
	})
}

// GetAllTVs retrieves all available TVs and Games
func GetAllTVs(c *fiber.Ctx) error {
	var tvs []models.TVInfo

	if err := database.DB.Preload("Games").Find(&tvs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:   500,
			Status: "SERVER_ERROR",
			Data:   models.ErrorData{ErrorMsg: "Could not fetch TV list with games"},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data:   tvs,
	})
}

// GetTVReservations retrieves the availability status for a specific TV
func GetTVReservations(c *fiber.Ctx) error {
	tvID := c.Params("tvId")

	var tvInfo models.TVInfo
	if err := database.DB.First(&tvInfo, tvID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
			Code: 404, Status: "NOT_FOUND", Data: models.ErrorData{ErrorMsg: "TV not found"},
		})
	}

	var reservations []models.Reservation
	today := time.Now().Format("2006-01-02")
	database.DB.Where("tv_id = ? AND time_slot LIKE ?", tvID, today+"%").Find(&reservations)

	reservedSlots := make(map[string]bool)
	for _, r := range reservations {
		reservedSlots[r.TimeSlot] = true
	}

	var timeSlots []models.TimeSlot
	for i := 9; i <= 17; i++ {
		startTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), i, 0, 0, 0, time.UTC)
		endTime := startTime.Add(1 * time.Hour)

		slotString := startTime.Format(time.RFC3339)
		availability := "available"
		if reservedSlots[slotString] {
			availability = "unavailable"
		}

		timeSlots = append(timeSlots, models.TimeSlot{
			StartTime:    slotString,
			EndTime:      endTime.Format(time.RFC3339),
			Availability: availability,
		})
	}

	tvStatus := models.TV{
		ID:          tvInfo.ID,
		ConsoleType: tvInfo.ConsoleType,
		TimeSlots:   timeSlots,
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		Code:   200,
		Status: "OK",
		Data:   tvStatus,
	})
}

func CreateReservation(c *fiber.Ctx) error {
	var body models.ReservationBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code: 400, Status: "BAD_REQUEST", Data: models.ErrorData{ErrorMsg: "Cannot parse JSON"},
		})
	}

	// Ambil ID pengguna yang sudah diautentikasi dari context.
	// Nilai ini diatur oleh AuthMiddleware dan merupakan sumber kebenaran.
	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		// Ini seharusnya tidak terjadi jika middleware berjalan, tapi ini adalah penjaga yang baik.
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code: 401, Status: "UNAUTHORIZED", Data: models.ErrorData{ErrorMsg: "User identity not found in token"},
		})
	}

	// Cek apakah slot sudah dipesan
	var existingReservation models.Reservation
	if err := database.DB.Where("tv_id = ? AND time_slot = ?", body.TVID, body.Timeslot).First(&existingReservation).Error; err != nil {
		// Pastikan error BUKAN karena data tidak ditemukan
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Database error during availability check"},
			})
		}
	} else {
		// Jika tidak ada error, berarti data ditemukan dan slot sudah dipesan
		return c.Status(fiber.StatusConflict).JSON(models.ErrorResponse{
			Code: 409, Status: "CONFLICT", Data: models.ErrorData{ErrorMsg: "Timeslot is already booked"},
		})
	}

	// Buat record reservasi baru
	newReservation := models.Reservation{
		TVID: body.TVID,
		// Gunakan userID dari token, bukan dari body request.
		BorrowerID: userID,
		TimeSlot:   body.Timeslot,
	}

	// Simpan ke database
	if err := database.DB.Create(&newReservation).Error; err != nil {
		// Log error yang sebenarnya untuk debugging di server
		log.Printf("DATABASE ERROR on CreateReservation: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code: 500, Status: "SERVER_ERROR", Data: models.ErrorData{ErrorMsg: "Could not create reservation"},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.Response{
		Code:   201,
		Status: "CREATED",
		Data:   nil,
	})
}
