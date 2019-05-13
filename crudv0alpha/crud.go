package crudv0alpha

import (
	"fmt"
	"net/http"
	"time"

	"github.com/arundo/data-sdk-go/ginutils"
	"github.com/arundo/data-sdk-go/utils"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Skeleton represents metadata of a skeleton
type Skeleton struct {
	ID        int64      `json:"-" gorm:"type:bigserial;primary_key"`
	GUID      string     `json:"guid" gorm:"not null;unique;index:skeleton_guid"`
	BoneCount int        `json:"boneCount,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"type:timestamp"`
	CreatedBy string     `json:"createBy,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"type:timestamp"`
	UpdatedBy string     `json:"updatedBy,omitempty"`
}

// BeforeCreate executes operations before create
func (skeleton *Skeleton) BeforeCreate() error {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	formattedTime, err := time.Parse(time.RFC3339, timeNow)
	if err != nil {
		return fmt.Errorf("Failed to parse time. %v", err)
	}
	skeleton.UpdatedAt = &formattedTime
	skeleton.CreatedAt = &formattedTime
	return nil
}

// BeforeUpdate executes operations before update
func (skeleton *Skeleton) BeforeUpdate() error {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	formattedTime, err := time.Parse(time.RFC3339, timeNow)
	if err != nil {
		return fmt.Errorf("Failed to parse time. %v", err)
	}
	skeleton.UpdatedAt = &formattedTime
	return nil
}

// AfterUpdate executes operations after update. This is needed to display the object's timestamp in the same format as createdAt
func (skeleton *Skeleton) AfterUpdate() error {
	timeNow := time.Now().UTC().Format(time.RFC3339)
	formattedTime, err := time.Parse(time.RFC3339, timeNow)
	if err != nil {
		return fmt.Errorf("Failed to parse time. %v", err)
	}
	skeleton.UpdatedAt = &formattedTime
	return nil
}

// HealthCheck returns http.StatusNoContent if service is running
func HealthCheck(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		ginutils.FormatResponse(c, http.StatusNoContent, "", "", "", utils.WhereAmI())
	}
}

func getErrorStatusCode(err error) int {
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		case "23505":
			return http.StatusConflict
		case "22023":
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	} else {
		return http.StatusInternalServerError
	}
}

// CreateSkeleton creates skeleton's data in the database
func CreateSkeleton(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		newRecord := Skeleton{}

		if err := c.ShouldBindJSON(&newRecord); err != nil {
			ginutils.FormatResponse(c, http.StatusInternalServerError, "about:blank", http.StatusText(http.StatusInternalServerError), err.Error(), utils.WhereAmI())
			return
		}

		newRecord.GUID = uuid.NewV4().String()
		username, _ := c.Get("username")
		newRecord.CreatedBy = username.(string)
		newRecord.UpdatedBy = username.(string)

		if err := db.Create(&newRecord).Error; err != nil {
			ginutils.FormatResponse(c, getErrorStatusCode(err), "about:blank", http.StatusText(getErrorStatusCode(err)), err.Error(), utils.WhereAmI()) // TODO: this error is not end Employee friendly
			return
		}

		ginutils.FormatResponse(c, http.StatusCreated, "", "", newRecord, utils.WhereAmI())
		return
	}
}

// GetSkeleton gets skeleton's data from the database
func GetSkeleton(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if uuid, uuidFound := c.Params.Get("skeletonId"); uuidFound {
			skeleton := Skeleton{}

			if err := db.Where("guid = ?", uuid).First(&skeleton).Error; err != nil {
				ginutils.FormatResponse(c, http.StatusNotFound, "about:blank", http.StatusText(http.StatusNotFound), err.Error(), utils.WhereAmI())
				return
			}

			ginutils.FormatResponse(c, http.StatusOK, "", "", skeleton, utils.WhereAmI())
		} else {
			ginutils.FormatResponse(c, http.StatusBadRequest, "about:blank", http.StatusText(http.StatusBadRequest), "SkeletonID not provided", utils.WhereAmI())
		}
		return
	}
}

// UpdateSkeleton updates skeleton's data in the database
func UpdateSkeleton(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if uuid, uuidFound := c.Params.Get("skeletonId"); uuidFound {
			requestedUpdates := Skeleton{}

			if err := c.ShouldBindJSON(&requestedUpdates); err != nil {
				ginutils.FormatResponse(c, http.StatusInternalServerError, "about:blank", http.StatusText(http.StatusInternalServerError), err.Error(), utils.WhereAmI())
				return
			}

			skeleton := Skeleton{}

			if err := db.Where("guid = ?", uuid).First(&skeleton).Error; err != nil {
				ginutils.FormatResponse(c, http.StatusNotFound, "about:blank", http.StatusText(http.StatusNotFound), err.Error(), utils.WhereAmI())
				return
			}

			username, _ := c.Get("username")
			requestedUpdates.UpdatedBy = username.(string)

			db.Model(&skeleton).Omit("guid", "created_at", "created_by").Updates(&requestedUpdates)
			ginutils.FormatResponse(c, http.StatusOK, "", "", skeleton, utils.WhereAmI())
		} else {
			ginutils.FormatResponse(c, http.StatusBadRequest, "about:blank", http.StatusText(http.StatusBadRequest), "SkeletonID not provided", utils.WhereAmI())
		}
	}
}

// DeleteSkeleton deletes skeleton's data from the database
func DeleteSkeleton(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		if uuid, uuidFound := c.Params.Get("skeletonId"); uuidFound {
			skeleton := Skeleton{}

			if err := db.Where("guid = ?", uuid).First(&skeleton).Error; err != nil {
				ginutils.FormatResponse(c, http.StatusNotFound, "about:blank", http.StatusText(http.StatusNotFound), err.Error(), utils.WhereAmI())
				return
			}

			db.Delete(&skeleton)
			ginutils.FormatResponse(c, http.StatusNoContent, "", "", "", utils.WhereAmI())
		} else {
			ginutils.FormatResponse(c, http.StatusBadRequest, "about:blank", http.StatusText(http.StatusBadRequest), "SkeletonID not provided", utils.WhereAmI())
		}
	}
}
