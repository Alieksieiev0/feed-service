package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/types"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type FeedService interface {
	GetPosts(c context.Context, params ...Param) ([]types.Post, error)
	Subscribe(c context.Context, id string, subId string) error
	Unsubscribe(c context.Context, id string, subId string) error
	Post(c context.Context, id string, post *models.Post) error
}

func NewFeedService(db *gorm.DB) FeedService {
	return &feedService{
		db: db,
	}
}

type feedService struct {
	db *gorm.DB
}

func (fs *feedService) GetPosts(c context.Context, params ...Param) ([]types.Post, error) {
	posts := []models.Post{}
	if err := ApplyParams(fs.db, params...).Find(&posts).Error; err != nil {
		return nil, err
	}

	usersPosts := fs.mapToOwners(posts)
	users := []models.User{}
	if err := fs.db.Find(&users, "Id IN ?", maps.Keys(usersPosts)).Error; err != nil {
		return nil, err
	}

	return fs.mapPosts(users, usersPosts), nil
}

func (fs *feedService) Subscribe(c context.Context, id string, subId string) error {
	if err := fs.userExists(&models.User{}, id); err != nil {
		return err
	}
	u := &models.User{}
	if err := fs.userExists(u, id); err != nil {
		return err
	}

	return fs.db.Model(&models.User{Base: models.Base{ID: id}}).
		Association("Subscribers").
		Append(u)
}

func (fs *feedService) Unsubscribe(c context.Context, id string, subId string) error {
	if err := fs.userExists(&models.User{}, id); err != nil {
		return err
	}
	u := &models.User{}
	if err := fs.userExists(u, id); err != nil {
		return err
	}

	return fs.db.Model(&models.User{Base: models.Base{ID: id}}).
		Association("Subscribers").
		Delete(u)
}

func (fs *feedService) Post(c context.Context, id string, post *models.Post) error {
	return fs.db.Model(&models.User{Base: models.Base{ID: id}}).Association("Posts").Append(post)
}

func (fs *feedService) userExists(user *models.User, id string) error {
	err := fs.db.First(user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user with provided id does not exist")
		}
		return err
	}
	return nil
}

func (fs *feedService) mapToOwners(posts []models.Post) map[string][]models.Post {
	ownerIds := map[string][]models.Post{}
	for _, p := range posts {
		values, ok := ownerIds[p.UserId]
		if ok {
			ownerIds[p.UserId] = append(values, p)
			continue
		}
		ownerIds[p.UserId] = []models.Post{p}
	}

	return ownerIds
}

func (fs *feedService) mapPosts(
	users []models.User,
	usersPosts map[string][]models.Post,
) []types.Post {
	mps := []types.Post{}
	for _, u := range users {
		for _, p := range usersPosts[u.ID] {
			mp := types.Post{
				Id:        p.ID,
				CreatedAt: p.CreatedAt,
				Title:     p.Title,
				Body:      p.Body,
				OwnerName: u.Username,
				OwnerId:   u.ID,
			}
			mps = append(mps, mp)
		}

	}

	return mps
}
