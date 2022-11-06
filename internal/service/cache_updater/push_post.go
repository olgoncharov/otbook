package cacheupdater

import (
	"context"

	"github.com/olgoncharov/otbook/internal/entity"
)

func (c *CacheUpdater) PushPostToFeed(ctx context.Context, post entity.Post) {
	go func() {
		select {
		case <-c.stopCh:
			return
		case c.pushPostQueue <- post:
		}
	}()
}

func (c *CacheUpdater) runPostPushWorker(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				c.logger.Debug().Msg("stop post push worker")
				return

			case post := <-c.pushPostQueue:
				c.pushPostToFeed(ctx, post)
			}
		}
	}()
}

func (c *CacheUpdater) pushPostToFeed(ctx context.Context, post entity.Post) {
	authorProfile, err := c.durableRepo.GetProfileByUsername(ctx, post.Author)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't retrieve profile of user %s", post.Author)
		return
	}
	if authorProfile.IsCelebrity {
		// do not cache posts from celebrity to avoid Lady Gaga effect
		return
	}

	followers, err := c.durableRepo.GetFollowersOfUser(ctx, post.Author)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't retrieve followrs of user %s", post.Author)
		return
	}

	if err = c.cacheRepo.PushPostToFeed(ctx, post, followers, uint(c.cfg.PostFeedLimit())); err != nil {
		c.logger.Error().Err(err).Msg("can't push post to feed")
		return
	}
}
