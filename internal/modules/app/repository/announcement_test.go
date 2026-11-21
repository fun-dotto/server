package repository

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fun-dotto/server/internal/shared/microcms"
)

func TestToDomainAnnouncementID(t *testing.T) {
	a := toDomainAnnouncement(microcms.Article{ID: "abc"}, "https://fun-dotto.com")
	b := toDomainAnnouncement(microcms.Article{ID: "abc"}, "https://fun-dotto.com")
	c := toDomainAnnouncement(microcms.Article{ID: "xyz"}, "https://fun-dotto.com")

	_, err := uuid.Parse(a.ID)
	require.NoError(t, err, "ID は UUID 形式")
	assert.Equal(t, a.ID, b.ID, "同じ記事なら同じ ID")
	assert.NotEqual(t, a.ID, c.ID, "異なる記事なら異なる ID")
	assert.Equal(t, "https://fun-dotto.com/articles/abc", a.URL)
}
