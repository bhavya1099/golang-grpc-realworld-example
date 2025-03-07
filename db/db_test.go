package db

import (
	errors "errors"
	fmt "fmt"
	os "os"
	testing "testing"
	sync "sync"
	sql "database/sql"
	ioutil "io/ioutil"
	gosqlmock "github.com/DATA-DOG/go-sqlmock"
	gorm "github.com/jinzhu/gorm"
	model "github.com/raahii/golang-grpc-realworld-example/model"
	debug "runtime/debug"
	time "time"
	assert "github.com/stretchr/testify/assert"
	uuid "github.com/google/uuid"
	godotenv "github.com/joho/godotenv"
	gotxdb "github.com/DATA-DOG/go-txdb"
	bytes "bytes"
	toml "github.com/BurntSushi/toml"
)



var invalidDSN = "invalid_dsn"
var openFunc = func(dialect string, args ...interface{}) (*gorm.DB, error) {
	if args[0] == validDSN {
		db, _, _ := gosqlmock.New()
		return &gorm.DB{DB: db}, nil
	}
	return nil, errors.New("failed to connect")
}
var validDSN = "valid_dsn"




/*
ROOST_METHOD_HASH=dsn_e202d1c4f9
ROOST_METHOD_SIG_HASH=dsn_b336e03d64

FUNCTION_DEF=func dsn() (string, error) 

*/
func TestDsn(t *testing.T) {
	type testScenario struct {
		name          string
		envSetup      func()
		expectedDSN   string
		expectedError string
	}

	tests := []testScenario{
		{
			name: "Scenario 1: Successful DSN Generation",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
			},
			expectedDSN:   "user:password@(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			expectedError: "",
		},
		{
			name: "Scenario 2: Missing DB_HOST Environment Variable",
			envSetup: func() {
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_HOST")
			},
			expectedDSN:   "",
			expectedError: "$DB_HOST is not set",
		},
		{
			name: "Scenario 3: Missing DB_USER Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_USER")
			},
			expectedDSN:   "",
			expectedError: "$DB_USER is not set",
		},
		{
			name: "Scenario 4: Missing DB_PASSWORD Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_NAME", "dbname")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_PASSWORD")
			},
			expectedDSN:   "",
			expectedError: "$DB_PASSWORD is not set",
		},
		{
			name: "Scenario 5: Missing DB_NAME Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_PORT", "3306")
				os.Unsetenv("DB_NAME")
			},
			expectedDSN:   "",
			expectedError: "$DB_NAME is not set",
		},
		{
			name: "Scenario 6: Missing DB_PORT Environment Variable",
			envSetup: func() {
				os.Setenv("DB_HOST", "localhost")
				os.Setenv("DB_USER", "user")
				os.Setenv("DB_PASSWORD", "password")
				os.Setenv("DB_NAME", "dbname")
				os.Unsetenv("DB_PORT")
			},
			expectedDSN:   "",
			expectedError: "$DB_PORT is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			tt.envSetup()

			dsn, err := dsn()

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("Expected error: %s, but got: %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if dsn != tt.expectedDSN {
					t.Errorf("Expected DSN: %s, but got: %s", tt.expectedDSN, dsn)
				}
			}
			t.Logf("Success: %s", tt.name)
		})
	}
}


/*
ROOST_METHOD_HASH=AutoMigrate_68b002ff76
ROOST_METHOD_SIG_HASH=AutoMigrate_19557fe223

FUNCTION_DEF=func AutoMigrate(db *gorm.DB) error // AutoMigrate is a wrapper of (*gorm.DB).AutoMigrate


*/
func TestAutoMigrate(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "Successful Migration of All Models",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Database Connection Error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("connection error"))
			},
			expectedError: errors.New("connection error"),
		},
		{
			name: "Partial Migration Failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("syntax error"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("syntax error"),
		},
		{
			name: "Concurrent Migration Requests",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "Unsupported Database Dialect",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("unsupported dialect"))
			},
			expectedError: errors.New("unsupported dialect"),
		},
		{
			name: "Invalid Model Definitions",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("CREATE TABLE").WillReturnError(errors.New("invalid model"))
				mock.ExpectRollback()
			},
			expectedError: errors.New("invalid model"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' occurred when opening a stub database connection", err)
			}
			defer db.Close()

			gormDB, err := gorm.Open("postgres", db)
			if err != nil {
				t.Fatalf("an error '%s' occurred when initializing gorm", err)
			}

			tt.setupMock(mock)

			err = AutoMigrate(gormDB)
			if err != nil && tt.expectedError == nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error %v, but got none", tt.expectedError)
			}
			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, but got %v", tt.expectedError, err)
			}

		})
	}
}


/*
ROOST_METHOD_HASH=DropTestDB_a4e5029e5a
ROOST_METHOD_SIG_HASH=DropTestDB_362a330faf

FUNCTION_DEF=func DropTestDB(d *gorm.DB) error // DropTestDB close connection


*/
func TestDropTestDb(t *testing.T) {
	type testScenario struct {
		description string
		setup       func() *gorm.DB
		expectedErr error
	}

	scenarios := []testScenario{
		{
			description: "Successful Database Connection Closure",
			setup: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Handling of Nil Database Connection",
			setup: func() *gorm.DB {
				return nil
			},
			expectedErr: nil,
		},
		{
			description: "Database Connection Already Closed",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose()
				gormDB.Close()
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Concurrent Database Connection Closures",
			setup: func() *gorm.DB {
				db, _, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				return gormDB
			},
			expectedErr: nil,
		},
		{
			description: "Mocked Database Connection with Error on Close",
			setup: func() *gorm.DB {
				db, mock, _ := sqlmock.New()
				gormDB, _ := gorm.Open("sqlmock", db)
				mock.ExpectClose().WillReturnError(fmt.Errorf("close error"))
				return gormDB
			},
			expectedErr: nil,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			db := scenario.setup()

			var wg sync.WaitGroup
			if scenario.description == "Concurrent Database Connection Closures" {
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						err := DropTestDB(db)
						assert.Equal(t, scenario.expectedErr, err)
					}()
				}
				wg.Wait()
			} else {
				err := DropTestDB(db)
				assert.Equal(t, scenario.expectedErr, err)
			}

			t.Logf("Scenario '%s' passed.", scenario.description)
		})
	}
}


/*
ROOST_METHOD_HASH=New_9360a2bab3
ROOST_METHOD_SIG_HASH=New_9365800ea3

FUNCTION_DEF=func New() (*gorm.DB, error) // New return mysql connection


*/
func TestNew(t *testing.T) {
	type testCase struct {
		name        string
		dsn         func() (string, error)
		open        func(dialect string, args ...interface{}) (*gorm.DB, error)
		expectedErr bool
		expectedNil bool
	}

	tests := []testCase{
		{
			name: "Successful Database Connection",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Failed Database Connection Due to Invalid DSN",
			dsn: func() (string, error) {
				return invalidDSN, nil
			},
			open:        openFunc,
			expectedErr: true,
			expectedNil: true,
		},
		{
			name: "Retry Logic on Temporary Connection Failure",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				staticCounter := 0
				return func(dialect string, args ...interface{}) (*gorm.DB, error) {
					if staticCounter < 3 {
						staticCounter++
						return nil, errors.New("temporary error")
					}
					db, _, _ := gosqlmock.New()
					return &gorm.DB{DB: db}, nil
				}(dialect, args...)
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Max Retry Limit Reached",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				return nil, errors.New("persistent error")
			},
			expectedErr: true,
			expectedNil: true,
		},
		{
			name: "Check Database Connection Pool Configuration",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
			},
			expectedErr: false,
			expectedNil: false,
		},
		{
			name: "Log Mode Configuration",
			dsn: func() (string, error) {
				return validDSN, nil
			},
			open: func(dialect string, args ...interface{}) (*gorm.DB, error) {
				db, _, _ := gosqlmock.New()
				return &gorm.DB{DB: db}, nil
			},
			expectedErr: false,
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			dsn = tt.dsn
			openFunc = tt.open

			db, err := New()

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedNil {
				assert.Nil(t, db)
			} else {
				assert.NotNil(t, db)
			}
		})
	}
}


/*
ROOST_METHOD_HASH=NewTestDB_4ca8dd9747
ROOST_METHOD_SIG_HASH=NewTestDB_ab2b828b16

FUNCTION_DEF=func NewTestDB() (*gorm.DB, error) // NewTestDB return mysql connection wrapped txdb


*/
func TestNewTestDb(t *testing.T) {
	tests := []struct {
		name           string
		envFileExists  bool
		mockDSNSuccess bool
		txdbError      bool
		expectError    bool
	}{
		{
			name:           "Successful Database Connection Initialization",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Missing Environment File",
			envFileExists:  false,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    true,
		},
		{
			name:           "Invalid DSN String",
			envFileExists:  true,
			mockDSNSuccess: false,
			txdbError:      false,
			expectError:    true,
		},
		{
			name:           "Concurrency and Mutex Locking",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Maximum Idle Connections Setting",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      false,
			expectError:    false,
		},
		{
			name:           "Error During txdb Registration",
			envFileExists:  true,
			mockDSNSuccess: true,
			txdbError:      true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()

			if tt.envFileExists {
				_ = godotenv.Load("../env/test.env")
			} else {
				_ = os.Remove("../env/test.env")
			}

			dsn = mockDSN(tt.mockDSNSuccess)

			if tt.txdbError {
				txdb.Register = func(driverName, dataSourceName, dsn string) error {
					return errors.New("txdb registration failed")
				}
			}

			db, err := NewTestDB()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else {
					t.Logf("Received expected error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if db == nil {
					t.Error("Expected non-nil gorm.DB object")
				} else {
					t.Log("Received valid gorm.DB object")
				}
			}

			if tt.name == "Maximum Idle Connections Setting" && db != nil {
				if db.DB().Stats().Idle > 3 {
					t.Errorf("Expected maximum idle connections to be 3, got: %v", db.DB().Stats().Idle)
				}
			}
		})
	}
}

func mockDSN(success bool) func() (string, error) {
	return func() (string, error) {
		if success {
			return "valid_dsn", nil
		}
		return "", errors.New("invalid DSN")
	}
}


/*
ROOST_METHOD_HASH=Seed_a96d5d4dee
ROOST_METHOD_SIG_HASH=Seed_878933cebc

FUNCTION_DEF=func Seed(db *gorm.DB) error 

*/
func TestSeed(t *testing.T) {

	txdb.Register("txdb", "mysql", "dsn")

	tests := []struct {
		name       string
		setup      func() (*gorm.DB, func())
		wantErr    bool
		assertFunc func(*testing.T, *gorm.DB)
	}{
		{
			name: "Successful Seeding of Users from a Valid TOML File",
			setup: func() (*gorm.DB, func()) {
				db, mock, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnResult(gosqlmock.NewResult(1, 1))
				mock.ExpectCommit()

				tomlData := `
				[[Users]]
				Username = "john_doe"
				Email = "john@example.com"
				Password = "password123"
				Bio = "A regular user"
				Image = "http://example.com/image.png"
				`
				ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)

				return gdb, func() {
					db.Close()
					os.Remove("db/seed/users.toml")
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var user model.User
				if err := db.First(&user, "username = ?", "john_doe").Error; err != nil {
					t.Fatalf("expected user to be found, got error: %v", err)
				}
				if user.Email != "john@example.com" {
					t.Errorf("expected email to be 'john@example.com', got '%s'", user.Email)
				}
			},
		},
		{
			name: "File Read Error",
			setup: func() (*gorm.DB, func()) {
				db, _, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				return gdb, func() {
					db.Close()
				}
			},
			wantErr: true,
		},
		{
			name: "TOML Decode Error",
			setup: func() (*gorm.DB, func()) {
				db, _, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				ioutil.WriteFile("db/seed/users.toml", []byte("invalid toml data"), 0644)

				return gdb, func() {
					db.Close()
					os.Remove("db/seed/users.toml")
				}
			},
			wantErr: true,
		},
		{
			name: "Database Insertion Error",
			setup: func() (*gorm.DB, func()) {
				db, mock, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `users`").WillReturnError(errors.New("unique constraint violation"))
				mock.ExpectRollback()

				tomlData := `
				[[Users]]
				Username = "john_doe"
				Email = "john@example.com"
				Password = "password123"
				Bio = "A regular user"
				Image = "http://example.com/image.png"
				`
				ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)

				return gdb, func() {
					db.Close()
					os.Remove("db/seed/users.toml")
				}
			},
			wantErr: true,
		},
		{
			name: "Empty User List in TOML File",
			setup: func() (*gorm.DB, func()) {
				db, _, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				tomlData := `
				[[Users]]
				`
				ioutil.WriteFile("db/seed/users.toml", []byte(tomlData), 0644)

				return gdb, func() {
					db.Close()
					os.Remove("db/seed/users.toml")
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var count int
				db.Model(&model.User{}).Count(&count)
				if count != 0 {
					t.Errorf("expected no users to be inserted, found %d", count)
				}
			},
		},
		{
			name: "Large Number of Users in TOML File",
			setup: func() (*gorm.DB, func()) {
				db, mock, err := gosqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}
				gdb, err := gorm.Open("mysql", db)
				if err != nil {
					t.Fatalf("failed to open gorm DB: %v", err)
				}

				mock.ExpectBegin()
				for i := 0; i < 1000; i++ {
					mock.ExpectExec("INSERT INTO `users`").WillReturnResult(gosqlmock.NewResult(1, 1))
				}
				mock.ExpectCommit()

				var tomlData bytes.Buffer
				for i := 0; i < 1000; i++ {
					tomlData.WriteString(fmt.Sprintf(`
					[[Users]]
					Username = "user%d"
					Email = "user%d@example.com"
					Password = "password123"
					Bio = "A regular user"
					Image = "http://example.com/image%d.png"
					`, i, i, i))
				}
				ioutil.WriteFile("db/seed/users.toml", tomlData.Bytes(), 0644)

				return gdb, func() {
					db.Close()
					os.Remove("db/seed/users.toml")
				}
			},
			wantErr: false,
			assertFunc: func(t *testing.T, db *gorm.DB) {
				var count int
				db.Model(&model.User{}).Count(&count)
				if count != 1000 {
					t.Errorf("expected 1000 users to be inserted, found %d", count)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Panic encountered so failing test. %v\n%s", r, string(debug.Stack()))
					t.Fail()
				}
			}()
			db, cleanup := tt.setup()
			defer cleanup()

			err := Seed(db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Seed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.assertFunc != nil {
				tt.assertFunc(t, db)
			}
		})
	}
}

