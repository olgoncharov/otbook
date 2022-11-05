package cacheupdater

import "context"

func (c *CacheUpdater) WarmupFeed(ctx context.Context, username string) {
	go func() {
		select {
		case <-c.stopCh:
			return
		case c.warmupFeedQueue <- username:
		}
	}()
}

func (c *CacheUpdater) runFeedWarmupWorker(ctx context.Context) {
	c.wg.Add(1)

	go func() {
		defer c.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				c.logger.Debug().Msg("stop feed warmup worker")
				return

			case username := <-c.warmupFeedQueue:
				c.warmupFeed(ctx, username)
			}
		}
	}()
}

func (c *CacheUpdater) warmupFeed(ctx context.Context, username string) {
	posts, err := c.durableRepo.GetPostFeedWithoutCelebrities(ctx, username, uint(c.cfg.PostFeedLimit()))
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't retrieve post feed from storage for user %s", username)
		return
	}

	err = c.cacheRepo.SetPostFeed(ctx, username, posts)
	if err != nil {
		c.logger.Error().Err(err).Msgf("can't save post feed in cache for user %s", username)
	}
}
