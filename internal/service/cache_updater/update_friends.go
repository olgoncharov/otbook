package cacheupdater

import (
	"context"

	"github.com/olgoncharov/otbook/internal/repository/dto"
)

func (c *CacheUpdater) AddFriend(ctx context.Context, user string, friend string) {
	go func() {
		select {
		case <-c.stopCh:
			return
		case c.updateFriendsQueue <- updateFriendsMsg{user: user, addedFriend: friend}:
		}
	}()
}

func (c *CacheUpdater) DeleteFriend(ctx context.Context, user string, friend string) {
	go func() {
		select {
		case <-c.stopCh:
			return
		case c.updateFriendsQueue <- updateFriendsMsg{user: user, deletedFriend: friend}:
		}
	}()
}

func (c *CacheUpdater) runFriendsUpdateWorker(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				c.logger.Debug().Msg("stop friends update worker")
				return

			case msg := <-c.updateFriendsQueue:
				if msg.addedFriend != "" {
					c.updateFriends(ctx, msg.user, msg.addedFriend, c.cacheRepo.AddCelebrityFriend)
					c.warmupFeed(ctx, msg.user) // full feed update
				}

				if msg.deletedFriend != "" {
					c.updateFriends(ctx, msg.user, msg.deletedFriend, c.cacheRepo.DeleteCelebrityFriend)
					c.removePostsOfFriendFromFeed(ctx, msg.user, msg.deletedFriend)
				}
			}
		}
	}()
}

func (c *CacheUpdater) updateFriends(ctx context.Context, user string, friend string, updateFn func(context.Context, string, string) error) {
	friendProfile, err := c.durableRepo.GetProfileByUsername(ctx, friend)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't retrieve profile of user %s", friend)
		return
	}

	if friendProfile.IsCelebrity {
		if err = updateFn(ctx, user, friend); err != nil {
			c.logger.Error().Err(err).Msg("can't update celebrity friends")
		}
	}
}

func (c *CacheUpdater) removePostsOfFriendFromFeed(ctx context.Context, user string, exFriend string) {
	postFeed, err := c.cacheRepo.GetPostFeed(ctx, user, uint(c.cfg.PostFeedLimit()))
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't get post feed from user %s", user)
		return
	}

	postsForRemoval := make([]dto.PostShortInfo, 0)
	for _, post := range postFeed {
		if post.Author == exFriend {
			postsForRemoval = append(postsForRemoval, post)
		}
	}

	if len(postsForRemoval) == 0 {
		return
	}

	err = c.cacheRepo.DeletePostsFromFeed(ctx, user, postsForRemoval)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't delete posts from feed for user %s", user)
	}
}
