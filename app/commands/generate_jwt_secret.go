package commands

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var GenerateJWTSecretCmd = &cobra.Command{
	Use:   "generate-jwt-secret",
	Short: "Generate and set JWT_SECRET in .env file",
	Long: `Generate a secure random JWT secret and add/update it in the .env file.
If JWT_SECRET doesn't exist, it will be created. If it exists, it will be updated.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := generateAndSetJWTSecret(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func generateAndSetJWTSecret() error {
	// Generate secure random secret (32 bytes = 256 bits)
	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return fmt.Errorf("failed to generate random secret: %w", err)
	}

	// Convert to base64 string (standard format)
	secret := base64.URLEncoding.EncodeToString(secretBytes)

	// Read .env file
	envPath := ".env"
	envContent, err := os.ReadFile(envPath)
	if err != nil {
		// If .env doesn't exist, create it
		if os.IsNotExist(err) {
			envContent = []byte{}
		} else {
			return fmt.Errorf("failed to read .env file: %w", err)
		}
	}

	// Update or add JWT_SECRET
	updatedContent := updateEnvKey(string(envContent), "JWT_SECRET", secret)

	// Write back to .env file
	if err := os.WriteFile(envPath, []byte(updatedContent), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	fmt.Printf("✅ JWT_SECRET generated and saved to .env file\n")
	fmt.Printf("   Key: JWT_SECRET\n")
	fmt.Printf("   Value: %s\n", secret)
	fmt.Printf("\n⚠️  Keep this secret secure and never commit it to version control!\n")

	return nil
}

// updateEnvKey updates or adds a key-value pair in .env content
func updateEnvKey(envContent, key, value string) string {
	lines := strings.Split(envContent, "\n")
	keyFound := false
	var updatedLines []string

	// Process existing lines
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines and comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			updatedLines = append(updatedLines, line)
			continue
		}

		// Check if this line contains the key
		if strings.HasPrefix(trimmed, key+"=") {
			// Update existing key with quotes
			updatedLines = append(updatedLines, fmt.Sprintf("%s=\"%s\"", key, value))
			keyFound = true
		} else {
			// Keep other lines as is
			updatedLines = append(updatedLines, line)
		}
	}

	// If key not found, add it at the end
	if !keyFound {
		// Add newline if content doesn't end with one
		if len(envContent) > 0 && !strings.HasSuffix(envContent, "\n") {
			updatedLines = append(updatedLines, "")
		}
		updatedLines = append(updatedLines, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	return strings.Join(updatedLines, "\n")
}
