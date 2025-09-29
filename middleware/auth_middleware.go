package middleware

import (
	// "database/sql"
	"hello-fiber/app/model"
	"strings"
	"time"
	"strconv"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key") // In production, use environment variable

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	RoleID int    `json:"role_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for authenticated user
func GenerateJWT(user model.User) (string, error) {
	claims := Claims{
		UserID: user.ID,
		Email:  user.Email,
		RoleID: user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// JWTMiddleware validates JWT token
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Extract claims and store in context
		if claims, ok := token.Claims.(*Claims); ok {
			c.Locals("user_id", claims.UserID)
			c.Locals("email", claims.Email)
			c.Locals("role_id", claims.RoleID)
		}

		return c.Next()
	}
}

// AdminOnlyMiddleware restricts access to admin users only
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleID := c.Locals("role_id")
		if roleID == nil || roleID.(int) != 1 { // 1 = admin role
			return c.Status(403).JSON(fiber.Map{
				"error": "Access denied. Admin role required",
			})
		}
		return c.Next()
	}
}

func PekerjaanOwnerMiddleware(db *sql.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        userID, ok := c.Locals("user_id").(int)
        if !ok {
            return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
        }

        idParam := c.Params("id")
        pekerjaanID, err := strconv.Atoi(idParam)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid pekerjaan id"})
        }

        // Ambil alumni_id dari pekerjaan_alumni
        var pekerjaanAlumniID int
        err = db.QueryRow("SELECT alumni_id FROM pekerjaan_alumni WHERE id = $1", pekerjaanID).Scan(&pekerjaanAlumniID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Tidak dapat pekerjaan alumni"})
        }

        // Ambil alumni_id dari users berdasarkan user yang login
        var userAlumniID int
        err = db.QueryRow("SELECT alumni_id FROM users WHERE id = $1", userID).Scan(&userAlumniID)
        if err != nil {
            return c.Status(404).JSON(fiber.Map{"error": "Tidak dapat alumni dari users"})
        }

        // Bandingkan alumni_id dari users dan pekerjaan_alumni
        if pekerjaanAlumniID != userAlumniID {
            return c.Status(403).JSON(fiber.Map{"error": "Forbidden: You can only update your own data"})
        }

        return c.Next()
    }
}