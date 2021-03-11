package basic

import (
	"context"
	"github.com/agui2200/wingman-store/examples/basic/store"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	ctx := context.Background()
	ctx, err := store.New(ctx, "sqlite4", "file:ent?mode=memory&cache=shared&_fk=1", false)
	require.Nil(t, err)
	c := store.WithContext(ctx)
	assert.NotNil(t, c)
}

func TestCloseDatabase(t *testing.T) {
	ctx := context.Background()
	ctx, err := store.New(ctx, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", true)
	require.Nil(t, err)
	err = store.CloseDatabase(ctx)
	require.Nil(t, err)
}

func TestMigration(t *testing.T) {
	ctx := context.Background()
	ctx, err := store.New(ctx, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1", true)
	require.Nil(t, err)
	err = store.Migration(ctx, true)
	require.Nil(t, err)
}
