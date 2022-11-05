package cacheupdater

import "context"

func (c *CacheUpdater) WarmupCelebrityFriendsList(ctx context.Context, username string) {
	go func() {
		select {
		case <-c.stopCh:
			return
		case c.warmupFriendsQueue <- username:
		}
	}()
}

func (c *CacheUpdater) runCelebrityFriendsWarmupWorker(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				c.logger.Debug().Msg("stop friends warmup worker")
				return

			case username := <-c.warmupFriendsQueue:
				c.warmupCelebrityFriendsList(ctx, username)
			}
		}
	}()
}

func (c *CacheUpdater) warmupCelebrityFriendsList(ctx context.Context, username string) {
	celebrityFriends, err := c.durableRepo.GetCelebrityFriends(ctx, username)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't retrieve celebrity friends from storage for user %s", username)
		return
	}

	err = c.cacheRepo.SetCelebrityFriends(ctx, username, celebrityFriends)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't save celebrity friends in cache for user %s", username)
		return
	}
}
