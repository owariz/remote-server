package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/owariz/remote-server/internal/config"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cfg := config.New()

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Missing authorization header")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
			}

			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid token: "+err.Error())
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("userId", claims["sub"])
			c.Locals("role", claims["role"])

			// ตรวจสอบสิทธิ์เฉพาะอาจทำที่นี่
			// ตัวอย่าง: ถ้า endpoint ต้องการสิทธิ์ admin แต่ผู้ใช้ไม่ใช่ admin
			/*
				if c.Path() == "/api/v1/admin" && claims["role"] != "admin" {
					return fiber.NewError(fiber.StatusForbidden, "Insufficient permissions")
				}
			*/

			return c.Next()
		}

		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
}
