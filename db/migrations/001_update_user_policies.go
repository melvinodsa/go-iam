package migrations

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/melvinodsa/go-iam/db"
	"github.com/melvinodsa/go-iam/db/models"
	"github.com/melvinodsa/go-iam/services/policy/system"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// OldUser represents the user structure before migration
type OldUser struct {
	Id         string                         `bson:"id"`
	ProjectId  string                         `bson:"project_id"`
	Name       string                         `bson:"name"`
	Email      string                         `bson:"email"`
	Phone      string                         `bson:"phone"`
	Enabled    bool                           `bson:"enabled"`
	ProfilePic string                         `bson:"profile_pic"`
	Expiry     *time.Time                     `bson:"expiry"`
	Roles      map[string]models.UserRoles    `bson:"roles"`
	Resources  map[string]models.UserResource `bson:"resources"`
	Policies   map[string]string              `bson:"policies"` // Old format: map[string]string
	CreatedAt  *time.Time                     `bson:"created_at"`
	CreatedBy  string                         `bson:"created_by"`
	UpdatedAt  *time.Time                     `bson:"updated_at"`
	UpdatedBy  string                         `bson:"updated_by"`
}

func init() {
	db.RegisterMigration(db.MigrationInfo{
		Version:     "001",
		Name:        "update_user_policies",
		Description: "Remove old policies field and add NewAccessToCreatedResource policy to all users",
		Up:          updateUserPoliciesUp,
		Down:        updateUserPoliciesDown,
	})
}

func updateUserPoliciesUp(ctx context.Context, dbConn db.DB) error {
	userModel := models.GetUserModel()

	// Get the NewAccessToCreatedResource policy
	accessPolicy := system.NewAccessToCreatedResource(nil)
	newPolicyData := map[string]models.UserPolicy{
		accessPolicy.ID(): {
			Name: accessPolicy.Name(),
		},
	}

	log.Info("Starting user policies migration...")

	// Process users in batches of 50
	batchSize := int64(50)
	skip := int64(0)
	totalProcessed := int64(0)

	for {
		// Find users in current batch
		findOpts := options.Find().SetLimit(batchSize).SetSkip(skip)
		cursor, err := dbConn.Find(ctx, userModel, bson.M{}, findOpts)
		if err != nil {
			return fmt.Errorf("failed to find users: %w", err)
		}
		defer func() {
			if err := cursor.Close(context.Background()); err != nil {
				log.Errorf("failed to close cursor: %w", err)
			}
		}()

		var users []OldUser
		if err := cursor.All(ctx, &users); err != nil {
			return fmt.Errorf("failed to decode users: %w", err)
		}

		// Break if no more users
		if len(users) == 0 {
			break
		}

		log.Infof("Processing batch of %d users (skip: %d)", len(users), skip)

		// Update each user in the batch
		for _, user := range users {
			now := time.Now()

			// Prepare update - remove old policies field and add new policies structure
			update := bson.M{
				"$set": bson.M{
					"policies":   newPolicyData, // Add new policy structure
					"updated_at": &now,
					"updated_by": "migration_001",
				},
			}

			filter := bson.M{userModel.IdKey: user.Id}

			result, err := dbConn.UpdateOne(ctx, userModel, filter, update)
			if err != nil {
				return fmt.Errorf("failed to update user %s: %w", user.Id, err)
			}

			if result.ModifiedCount == 0 {
				log.Warnf("User %s was not updated (may already be migrated)", user.Id)
			} else {
				log.Debugf("Successfully updated user %s", user.Id)
			}
		}

		totalProcessed += int64(len(users))
		skip += batchSize

		log.Infof("Processed %d users so far", totalProcessed)
	}

	log.Infof("Successfully migrated %d users", totalProcessed)
	return nil
}

func updateUserPoliciesDown(ctx context.Context, dbConn db.DB) error {
	userModel := models.GetUserModel()

	log.Info("Rolling back user policies migration...")

	// Process users in batches of 50
	batchSize := int64(50)
	skip := int64(0)
	totalProcessed := int64(0)

	for {
		// Find users in current batch
		findOpts := options.Find().SetLimit(batchSize).SetSkip(skip)
		cursor, err := dbConn.Find(ctx, userModel, bson.M{}, findOpts)
		if err != nil {
			return fmt.Errorf("failed to find users: %w", err)
		}
		defer func() {
			if err := cursor.Close(context.Background()); err != nil {
				log.Errorf("failed to close cursor: %w", err)
			}
		}()

		var users []models.User
		if err := cursor.All(ctx, &users); err != nil {
			return fmt.Errorf("failed to decode users: %w", err)
		}

		// Break if no more users
		if len(users) == 0 {
			break
		}

		log.Infof("Rolling back batch of %d users (skip: %d)", len(users), skip)

		// Update each user in the batch
		for _, user := range users {
			now := time.Now()

			// Convert new policies structure back to old format (empty map for rollback)
			oldPolicies := make(map[string]string)

			// Prepare update - remove new policies field and add old policies structure
			update := bson.M{
				"$set": bson.M{
					"policies":   oldPolicies, // Set old policy structure (empty)
					"updated_at": &now,
					"updated_by": "migration_001_rollback",
				},
			}

			filter := bson.M{userModel.IdKey: user.Id}

			result, err := dbConn.UpdateOne(ctx, userModel, filter, update)
			if err != nil {
				return fmt.Errorf("failed to rollback user %s: %w", user.Id, err)
			}

			if result.ModifiedCount == 0 {
				log.Warnf("User %s was not rolled back", user.Id)
			} else {
				log.Debugf("Successfully rolled back user %s", user.Id)
			}
		}

		totalProcessed += int64(len(users))
		skip += batchSize

		log.Infof("Rolled back %d users so far", totalProcessed)
	}

	log.Infof("Successfully rolled back %d users", totalProcessed)
	return nil
}
