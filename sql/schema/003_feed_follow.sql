CREATE TABLE
    feed_follow (
        id UUID PRIMARY KEY,
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW (),
        user_id UUID,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
        feed_id UUID,
        FOREIGN KEY (feed_id) REFERENCES feeds (id) ON DELETE CASCADE,
        UNIQUE (user_id, feed_id)
    );