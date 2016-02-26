package accounts

import (
	"log"
	"testing"

	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/database"
	"github.com/RichardKnop/recall/oauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	// sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

var testDbPath = "/tmp/accounts_testdb.sqlite"

var testFixtures = []string{
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../oauth/fixtures/test_access_tokens.yml",
	"fixtures/roles.yml",
	"fixtures/test_accounts.yml",
	"fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	MigrateAll,
}

// AccountsTestSuite needs to be exported so the tests run
type AccountsTestSuite struct {
	suite.Suite
	cnf      *config.Config
	db       *gorm.DB
	service  *Service
	accounts []*Account
	users    []*User
	router   *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *AccountsTestSuite) SetupSuite() {

	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := database.CreateTestDatabase(testDbPath, testMigrations, testFixtures)
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	// Fetch test accounts
	suite.accounts = make([]*Account, 0)
	if suite.db.Preload("OauthClient").Find(&suite.accounts).Error != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*User, 0)
	err = suite.db.Preload("Account").Preload("OauthUser").Preload("Role").
		Find(&suite.users).Error
	if err != nil {
		log.Fatal(err)
	}

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		oauth.NewService(suite.cnf, suite.db),
	)

	// Register routes
	suite.router = mux.NewRouter()
	RegisterRoutes(suite.router, suite.service)
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *AccountsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *AccountsTestSuite) SetupTest() {
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(User))

	// Service.CreateUser also creates a new oauth.User instance
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.User))

	// Service.CreateAccount also creates a new oauth.Client instance
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.Client))
}

// The TearDownTest method will be run after every test in the suite.
func (suite *AccountsTestSuite) TearDownTest() {
	//
}

// TestAccountsTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
