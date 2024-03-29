package model

import (
	"TikTokServer/pkg/errorcode"

	"gorm.io/gorm"
)

type UserFavoriteVideos struct {
	UserID  int64 `gorm:"primaryKey;"`
	VideoID int64 `gorm:"primaryKey;"`
}

func Favorite(userID, videoID int64) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		user := new(User)
		video := new(Video)
		if err := tx.First(video, videoID).Error; err != nil {
			return err
		}

		if err := tx.First(user, userID).Error; err != nil {
			return err
		}

		// tlog.Debugf("video: %v, user: %v", video, user)
		if err := tx.Model(user).Association("FavoriteVideos").Append(video); err != nil {
			return err
		}
		// videoCnt := tx.Model(video).Association("FavoriteVideos").Count()
		var videoCnt int64
		tx.Model(&UserFavoriteVideos{}).Where("video_id = ?", videoID).Count(&videoCnt)

		// tlog.Debugf("video: %v, user: %v", video, user)
		// tlog.Debugf("videoCnt: %v", videoCnt)
		res := tx.Model(video).Update("favorite_count", videoCnt)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			errCode := errorcode.ErrHttpDatabase
			errCode.SetMsg("更新字段数量大于1")
			return errCode
		}

		// TODO: 这两个字段重要性不高，可以后续使用定时任务从 redis 缓存里更新
		// res = tx.Model(user).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }

		// if res.RowsAffected != 1 {
		// 	errCode := errorcode.ErrHttpDatabase
		// 	errCode.SetMsg("更新字段数量大于1")
		// 	return errCode
		// }

		// // 更新视频作者的 TotalFavorited 字段
		// res = tx.Model(User{}).Where("id = ?", video.AuthorID).Update("total_favorited", gorm.Expr("total_favorited + ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }
		// if res.RowsAffected != 1 {
		// 	errCode := errorcode.ErrHttpDatabase
		// 	errCode.SetMsg("更新字段数量大于1")
		// 	return errCode
		// }

		return nil
	})
	return err
}

func DisFavorite(userID, videoID int64) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		user := new(User)
		video := new(Video)

		if err := tx.First(video, videoID).Error; err != nil {
			return err
		}

		if err := tx.First(user, userID).Error; err != nil {
			return err
		}

		err := tx.Unscoped().Model(user).Association("FavoriteVideos").Delete(video)
		if err != nil {
			return err
		}

		var videoCnt int64
		tx.Model(&UserFavoriteVideos{}).Where("video_id = ?", videoID).Count(&videoCnt)
		res := tx.Model(video).Update("favorite_count", videoCnt)

		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			errCode := errorcode.ErrHttpDatabase
			errCode.SetError(res.Error)
			return errCode
		}

		// res = tx.Model(user).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }

		// if res.RowsAffected != 1 {
		// 	errCode := errorcode.ErrHttpDatabase
		// 	errCode.SetMsg("更新字段数量大于1")
		// 	return errCode
		// }

		// // 更新视频作者的 TotalFavorited 字段
		// res = tx.Model(User{}).Where("id = ?", video.AuthorID).Update("total_favorited", gorm.Expr("total_favorited - ?", 1))
		// if res.Error != nil {
		// 	return res.Error
		// }

		// if res.RowsAffected != 1 {
		// 	errCode := errorcode.ErrHttpDatabase
		// 	errCode.SetMsg("更新字段数量大于1")
		// 	return errCode
		// }

		return nil
	})
	return err
}

func GetFavoriteList(userID int64) ([]*Video, error) {
	videos := []*Video{}

	user := &User{}
	if err := db.First(user, userID).Error; err != nil {
		return nil, err
	}

	if err := db.Model(user).Association("FavoriteVideos").Find(&videos); err != nil {
		return nil, err
	}

	for i, v := range videos {
		author, err := GetUserByID(v.AuthorID)
		if err != nil {
			return videos, err
		}

		videos[i].Author = *author
	}

	return videos, nil
}
