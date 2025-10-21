package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"
	"hello-fiber/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// JWTMiddleware validates JWT token and stores claims in context
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header required",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format",
			})
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		token, err := jwt.ParseWithClaims(tokenString, &utils.Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return utils.GetJWTSecret(), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token", "detail": err.Error()})
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role_id", claims.RoleID) // Now storing as string (hex format)

		return c.Next()
	}
}

// AdminOnlyMiddleware restricts access to admin users only
func AdminOnlyMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleIDStr, ok := c.Locals("role_id").(string)
		if !ok {
			fmt.Printf("[DEBUG] role_id not found in locals\n")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Admin role required"})
		}

		fmt.Printf("[DEBUG] roleIDStr from JWT: %s\n", roleIDStr)

		roleID, err := bson.ObjectIDFromHex(roleIDStr)
		if err != nil {
			fmt.Printf("[DEBUG] Failed to convert roleIDStr to ObjectID: %v\n", err)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Invalid role"})
		}

		// Get admin role from database
		adminRoleColl := database.MongoDB.Collection("roles")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := adminRoleColl.Find(ctx, bson.M{})
		if err == nil {
			defer cursor.Close(ctx)
			var allRoles []model.Role
			if err := cursor.All(ctx, &allRoles); err == nil {
				fmt.Printf("[DEBUG] All roles in database:\n")
				for _, r := range allRoles {
					fmt.Printf("  - ID: %s, Role: %s\n", r.ID.Hex(), r.Role)
				}
			}
		}

		var adminRole model.Role
		err = adminRoleColl.FindOne(ctx, bson.M{"role": "admin"}).Decode(&adminRole)
		if err != nil {
			fmt.Printf("[DEBUG] Admin role not found with query {role: 'admin'}: %v\n", err)
			// Try alternative query
			err = adminRoleColl.FindOne(ctx, bson.M{"role": "Admin"}).Decode(&adminRole)
			if err != nil {
				fmt.Printf("[DEBUG] Admin role not found with query {role: 'Admin'}: %v\n", err)
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Admin role not found in database"})
			}
		}

		fmt.Printf("[DEBUG] User roleID: %s, Admin roleID: %s\n", roleID.Hex(), adminRole.ID.Hex())

		if roleID != adminRole.ID {
			fmt.Printf("[DEBUG] Role mismatch - User has: %s, Admin is: %s\n", roleID.Hex(), adminRole.ID.Hex())
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Admin role required"})
		}

		fmt.Printf("[DEBUG] Admin verification passed!\n")
		return c.Next()
	}
}

// PekerjaanOwnerMiddlewareMongo memastikan user hanya boleh akses pekerjaan miliknya
// Verifikasi dilakukan dengan membandingkan id, email, dan roleID dari JWT
func PekerjaanOwnerMiddlewareMongo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr, ok := c.Locals("user_id").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		email, ok := c.Locals("email").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		roleIDStr, ok := c.Locals("role_id").(string) // Now receiving as string
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		pekerjaanIDStr := c.Params("id")
		if pekerjaanIDStr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pekerjaan id"})
		}

		pekerjaanID, err := bson.ObjectIDFromHex(pekerjaanIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid pekerjaan id format"})
		}

		userID, err := bson.ObjectIDFromHex(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user id"})
		}

		roleID, err := bson.ObjectIDFromHex(roleIDStr) // Convert role_id string to ObjectID
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid role id"})
		}

		coll := database.MongoDB.Collection("pekerjaan_alumni")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{"_id": pekerjaanID}
		var pekerjaan struct {
			AlumniID bson.ObjectID `bson:"alumni_id"`
		}
		if err := coll.FindOne(ctx, filter).Decode(&pekerjaan); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pekerjaan alumni tidak ditemukan"})
		}

		usersColl := database.MongoDB.Collection("users")
		var userDoc model.User
		if err := usersColl.FindOne(ctx, bson.M{"_id": userID}).Decode(&userDoc); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		// 1. User ID from JWT matches user in database
		// 2. Email from JWT matches user email in database
		// 3. RoleID from JWT matches user role in database
		// 4. Alumni ID from pekerjaan matches user's alumni_id
		if userDoc.ID != userID || userDoc.Email != email || userDoc.RoleID != roleID { // Now comparing ObjectID to ObjectID
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: JWT claims do not match user data"})
		}

		if pekerjaan.AlumniID != userDoc.ID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden: You can only access your own data"})
		}

		return c.Next()
	}
}
