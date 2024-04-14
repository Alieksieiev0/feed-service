package services

import (
	"context"

	"github.com/Alieksieiev0/feed-service/internal/models"
	"github.com/Alieksieiev0/feed-service/internal/types"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type FeedService interface {
	GetPosts(c context.Context, params ...Param) ([]types.Post, error)
	Subscribe(c context.Context, id string, subId string) error
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
	u := &models.User{}
	if err := fs.db.First(u, "id = ?", subId).Error; err != nil {
		return err
	}

	return fs.db.Model(&models.User{Base: models.Base{ID: id}}).
		Association("Subscribers").
		Append(u)
}

func (fs *feedService) Post(c context.Context, id string, post *models.Post) error {
	return fs.db.Model(&models.User{Base: models.Base{ID: id}}).Association("Posts").Append(post)
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
