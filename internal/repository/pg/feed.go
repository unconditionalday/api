package pg

import (
	"database/sql"
	"time"

	"github.com/unconditionalday/server/internal/app"
)

type FeedRepository struct {
	db *sql.DB
}

func NewFeedRepository(db *sql.DB) *FeedRepository {
	return &FeedRepository{
		db: db,
	}
}

func (f *FeedRepository) Find(query string) ([]app.Feed, error) {
	rows, err := f.db.Query(`
	SELECT
    feeds.title,
    feeds.link,
    feeds.source,
    feeds.language,
    feeds.summary,
    feeds.date,
    COALESCE(feeds.image_url, ''),
    COALESCE(feeds.image_title, '')
FROM 
    feeds, 
    to_tsvector(feeds.title || feeds.summary) document,
    to_tsquery($1) query,
    NULLIF(ts_rank(to_tsvector(feeds.title), query), 0) rank_title,
    NULLIF(ts_rank(to_tsvector(feeds.summary), query), 0) rank_description,
    SIMILARITY($1, feeds.title || feeds.summary) similarity
WHERE 
    (query @@ document OR similarity > 0)
    AND rank_title IS NOT NULL
ORDER BY rank_title, rank_description, similarity DESC;
`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []app.Feed
	for rows.Next() {
		var feed app.Feed
		var dateStr string
		var imageURL, imageTitle sql.NullString

		err := rows.Scan(
			&feed.Title,
			&feed.Link,
			&feed.Source,
			&feed.Language,
			&feed.Summary,
			&dateStr,
			&imageURL,
			&imageTitle,
		)
		if err != nil {
			return nil, err
		}

		feed.Date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			return nil, err
		}

		if imageURL.Valid && imageTitle.Valid {
			feed.Image = &app.Image{
				URL:   imageURL.String,
				Title: imageTitle.String,
			}
		}

		feeds = append(feeds, feed)
	}

	return feeds, nil

}

func (f *FeedRepository) Save(feed app.Feed) error {
	_, err := f.db.Exec(`
		INSERT INTO feeds (title, link, source, language, summary, date, image_url, image_title)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, feed.Title, feed.Link, feed.Source, feed.Language, feed.Summary, feed.Date.Format(time.RFC3339), feed.Image.URL, feed.Image.Title)
	return err
}

func (f *FeedRepository) Update(feed app.Feed) error {
	_, err := f.db.Exec(`
        INSERT INTO feeds (title, link, source, language, summary, date, image_url, image_title)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        ON CONFLICT (link) DO UPDATE
        SET title = $1, source = $3, language = $4, summary = $5, date = $6, image_url = $7, image_title = $8
	`, feed.Title, feed.Link, feed.Source, feed.Language, feed.Summary, feed.Date.Format(time.RFC3339), feed.Image.URL, feed.Image.Title)
	return err
}

func (f *FeedRepository) Count() uint64 {
	var count uint64

	f.db.QueryRow("SELECT COUNT(*) FROM feeds").Scan(&count)
	return count

}

func (f *FeedRepository) Delete(feed app.Feed) error {
	_, err := f.db.Exec("DELETE FROM feeds WHERE link = $1", feed.Link)
	return err
}
