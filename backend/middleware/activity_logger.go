package middleware

import (
	"encoding/json"
	"sirine-go/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ActivityLogger merupakan middleware untuk auto-log critical actions
// yang mencakup CREATE, UPDATE, DELETE operations pada entities
func ActivityLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request first
		c.Next()

		// Check apakah ada activity yang perlu di-log
		action, exists := c.Get("activity_action")
		if !exists {
			return
		}

		// Get current user
		currentUser, exists := c.Get("user")
		if !exists {
			return
		}
		user := currentUser.(*models.User)

		// Get activity details dari context
		entityType, _ := c.Get("activity_entity_type")
		entityID, _ := c.Get("activity_entity_id")
		changesBefore, _ := c.Get("activity_changes_before")
		changesAfter, _ := c.Get("activity_changes_after")

		// Create activity log
		activityLog := models.ActivityLog{
			UserID:     user.ID,
			Action:     action.(models.ActivityAction),
			EntityType: entityType.(string),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.GetHeader("User-Agent"),
		}

		// Set entity ID jika ada
		if entityID != nil {
			if id, ok := entityID.(uint64); ok {
				activityLog.EntityID = &id
			}
		}

		// Set changes (before/after)
		if changesBefore != nil || changesAfter != nil {
			changeData := models.ChangeData{
				Before: changesBefore,
				After:  changesAfter,
			}

			changesJSON, err := json.Marshal(changeData)
			if err == nil {
				activityLog.Changes = changesJSON
			}
		}

		// Save activity log ke database
		// Gunakan background goroutine untuk tidak block response
		go func() {
			db.Create(&activityLog)
		}()
	}
}
