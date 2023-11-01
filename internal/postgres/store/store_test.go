package store

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/autostrada-test/store"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var TestDBString = os.Getenv("DB_TEST_DSN")

func setupTest(t testing.TB) (*PostgresStore, func(tb testing.TB)) {
	postgresStore, closer, err := NewPostgresStore(TestDBString, true)
	require.Nil(t, err)
	return postgresStore, func(tb testing.TB) {
		defer closer()
		log.Println("teardown suite")
		err := migrateDb(TestDBString, DirectionDown)
		require.Nil(t, err)
	}
}

func TestPostgresStoreUserInsert(t *testing.T) {
	t.Run("test UserInsert method happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := uuid.New()
		admin := true

		user, err := postgresStore.UserInsert(email, password, id, admin)

		fmt.Println(err)
		require.Nil(t, err)

		require.Equal(t, email, user.Email)
		require.Equal(t, admin, user.Admin)
		require.Equal(t, id, user.ID)
		require.NotNil(t, user.Created)
	})

	t.Run("test UserInsert for duplicates", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		id := uuid.New()
		admin := true

		_, _ = postgresStore.UserInsert(email, password, id, admin)
		_, err := postgresStore.UserInsert(email, password, id, admin)

		require.Error(t, err)
		require.Equal(t, store.ErrUserExists, err)
	})
}

func TestPostgresStoreUserList(t *testing.T) {
	t.Run("test UserList method happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := uuid.New()

		_, _ = postgresStore.UserInsert("a"+email, password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email, password, id, admin)
		_, _ = postgresStore.UserInsert(email+"a", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"b", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"c", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"d", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"e", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"f", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"g", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"h", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"i", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"j", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"k", password, uuid.New(), admin)
		_, _ = postgresStore.UserInsert(email+"l", password, uuid.New(), admin)
		users, err := postgresStore.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 14, users.TotalObjects)
		require.Equal(t, 2, users.TotalPages)

		require.Equal(t, email, users.Data[1].Email)
		require.Equal(t, admin, users.Data[1].Admin)
		require.Equal(t, id, users.Data[1].ID)
		require.NotNil(t, users.Data[1].Created)
	})

	t.Run("test UserList method no data", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		users, err := postgresStore.UserList(1, 10)
		require.Nil(t, err)
		require.NotNil(t, users)

		require.Equal(t, 0, users.TotalObjects)
		require.Equal(t, 0, users.TotalPages)
	})
}

func TestPostgresStoreUserRetrieve(t *testing.T) {
	t.Run("UserRetrieve happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		admin := true
		id := uuid.New()

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		retrieved, err := postgresStore.UserRetrieve(user.ID)
		t.Log(retrieved)
		require.Nil(t, err)
		require.NotNil(t, retrieved)

		require.Equal(t, user, retrieved)
	})

	t.Run("UserRetrieve not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := uuid.New()

		retrieved, err := postgresStore.UserRetrieve(id)
		require.Nil(t, retrieved)
		require.NotNil(t, err)
		t.Log(err)
	})
}

func TestPostgresStoreUserRetrieveByEmail(t *testing.T) {
	t.Run("UserRetrieveByEmail happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"
		password := "password"
		admin := true
		id := uuid.New()

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		retrieved, err := postgresStore.UserRetrieveByEmail(user.Email)
		t.Log(retrieved)
		require.Nil(t, err)
		require.NotNil(t, retrieved)

		require.Equal(t, user, retrieved)
	})

	t.Run("UserRetrieveByEmail not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im"

		retrieved, err := postgresStore.UserRetrieveByEmail(email)
		require.Nil(t, retrieved)
		require.NotNil(t, err)
		t.Log(err)
	})
}
func TestNewPostgresStoreUserUpdatePassword(t *testing.T) {
	t.Run("UserUpdatePassword happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		newPassword := "newPassword"
		admin := true
		id := uuid.New()

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserUpdatePassword(user.ID, newPassword)
		require.Nil(t, err)
	})
	t.Run("UserUpdatePassword user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := uuid.New()
		newPassword := "newPassword"

		err := postgresStore.UserUpdatePassword(id, newPassword)
		t.Log(err)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
func TestNewPostgresStoreUserUpdateAdmin(t *testing.T) {
	t.Run("UserUpdateAdmin happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := uuid.New()
		newAdminValue := false

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserUpdateAdmin(user.ID, newAdminValue)
		require.Nil(t, err)

		u, err := postgresStore.UserRetrieve(user.ID)
		require.Nil(t, err)
		require.NotNil(t, u)

		require.Equal(t, newAdminValue, u.Admin)

	})
	t.Run("UserUpdateAdmin user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := uuid.New()
		newAdminValue := false

		err := postgresStore.UserUpdateAdmin(id, newAdminValue)
		t.Log(err)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
func TestNewPostgresStoreUserDelete(t *testing.T) {
	t.Run("UserDelete happy path", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		email := "im@parham.im123"
		password := "password"
		admin := true
		id := uuid.New()

		user, err := postgresStore.UserInsert(email, password, id, admin)
		require.Nil(t, err)
		require.NotNil(t, user)

		err = postgresStore.UserDelete(id)
		require.Nil(t, err)
	})
	t.Run("UserUpdateAdmin user not exists", func(t *testing.T) {
		postgresStore, teardownTest := setupTest(t)
		defer teardownTest(t)

		id := uuid.New()

		err := postgresStore.UserDelete(id)
		require.NotNil(t, err)
		require.Equal(t, store.ErrUserNotFound, err)
	})
}
