package facebook_test

import (
	"log"
	"testing"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/database"
	"github.com/RichardKnop/example-api/facebook"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/gorilla/mux"
	fb "github.com/huandu/facebook"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"

	accountsMocks "github.com/RichardKnop/example-api/accounts/mocks"
	emailMocks "github.com/RichardKnop/example-api/email/mocks"
)

var (
	testDbUser = "example_api"
	testDbName = "example_api_facebook_test"
)

var testFixtures = []string{
	"../oauth/fixtures/scopes.yml",
	"../oauth/fixtures/roles.yml",
	"../oauth/fixtures/test_clients.yml",
	"../oauth/fixtures/test_users.yml",
	"../accounts/fixtures/test_accounts.yml",
	"../accounts/fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
}

// FacebookTestSuite needs to be exported so the tests run
type FacebookTestSuite struct {
	suite.Suite
	cnf         *config.Config
	db          *gorm.DB
	adapterMock *facebook.AdapterMock
	service     *facebook.Service
	router      *mux.Router
	accounts    []*accounts.Account
	users       []*accounts.User
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *FacebookTestSuite) SetupSuite() {
	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := database.CreateTestDatabasePostgres(
		testDbUser,
		testDbName,
		testMigrations,
		testFixtures,
	)
	if err != nil {
		log.Fatal(err)
	}
	suite.db = db

	// Fetch test accounts
	suite.accounts = make([]*accounts.Account, 0)
	err = accounts.AccountPreload(suite.db).Order("id").Find(&suite.accounts).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*accounts.User, 0)
	err = accounts.UserPreload(suite.db).Order("id").Find(&suite.users).Error
	if err != nil {
		log.Fatal(err)
	}

	// Initialise mocks
	suite.adapterMock = new(facebook.AdapterMock)

	// Initialise the service
	suite.service = facebook.NewService(
		suite.cnf,
		suite.db,
		accounts.NewService(
			suite.cnf,
			suite.db,
			oauth.NewService(suite.cnf, suite.db),
			new(emailMocks.ServiceInterface),
			new(accountsMocks.EmailFactoryInterface),
		),
		suite.adapterMock,
	)

	// Register routes
	suite.router = mux.NewRouter()
	suite.service.RegisterRoutes(suite.router, "/v1/facebook")
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *FacebookTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *FacebookTestSuite) SetupTest() {
	// loginHandler also creates a new user and oauth tokens
	suite.db.Unscoped().Delete(new(oauth.RefreshToken))
	suite.db.Unscoped().Delete(new(oauth.AccessToken))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(accounts.User))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(oauth.User))

	// Reset mocks
	suite.adapterMock.ExpectedCalls = suite.adapterMock.ExpectedCalls[:0]
	suite.adapterMock.Calls = suite.adapterMock.Calls[:0]
}

// The TearDownTest method will be run after every test in the suite.
func (suite *FacebookTestSuite) TearDownTest() {
	//
}

// TestFacebookTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFacebookTestSuite(t *testing.T) {
	suite.Run(t, new(FacebookTestSuite))
}

// Mock fetching of facebook profile data
func (suite *FacebookTestSuite) mockFacebookGetMe(result fb.Result, err error) {
	suite.adapterMock.On("GetMe", "facebook_token").Return(result, err)
}
