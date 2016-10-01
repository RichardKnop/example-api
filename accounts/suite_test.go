package accounts_test

import (
	"log"
	"os"
	"testing"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/oauth"
	"github.com/RichardKnop/example-api/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	accountsMocks "github.com/RichardKnop/example-api/accounts/mocks"
	emailMocks "github.com/RichardKnop/example-api/email/mocks"
)

var (
	testDbUser = "example_api"
	testDbName = "example_api_accounts_test"

	testFixtures = []string{
		"./oauth/fixtures/scopes.yml",
		"./oauth/fixtures/roles.yml",
		"./oauth/fixtures/test_clients.yml",
		"./oauth/fixtures/test_users.yml",
		"./oauth/fixtures/test_access_tokens.yml",
		"./accounts/fixtures/test_accounts.yml",
		"./accounts/fixtures/test_users.yml",
	}

	testMigrations = []func(*gorm.DB) error{
		oauth.MigrateAll,
		accounts.MigrateAll,
	}
)

func init() {
	if err := os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}

// AccountsTestSuite needs to be exported so the tests run
type AccountsTestSuite struct {
	suite.Suite
	cnf          *config.Config
	db           *gorm.DB
	emailService *emailMocks.ServiceInterface
	emailFactory *accountsMocks.EmailFactoryInterface
	service      *accounts.Service
	accounts     []*accounts.Account
	users        []*accounts.User
	router       *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *AccountsTestSuite) SetupSuite() {
	// Initialise the config
	suite.cnf = config.NewConfig(false, false)

	// Create the test database
	db, err := testutil.CreateTestDatabasePostgres(
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
	suite.emailService = new(emailMocks.ServiceInterface)
	suite.emailFactory = new(accountsMocks.EmailFactoryInterface)

	// Initialise the service
	suite.service = accounts.NewService(
		suite.cnf,
		suite.db,
		oauth.NewService(suite.cnf, suite.db),
		suite.emailService,
		suite.emailFactory,
	)

	// Register routes
	suite.router = mux.NewRouter()
	suite.service.RegisterRoutes(suite.router, "/v1")
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *AccountsTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *AccountsTestSuite) SetupTest() {
	suite.db.Unscoped().Delete(new(accounts.Confirmation))
	suite.db.Unscoped().Delete(new(accounts.Invitation))
	suite.db.Unscoped().Delete(new(accounts.PasswordReset))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(accounts.User))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(accounts.Account))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(oauth.AccessToken))
	suite.db.Unscoped().Delete(new(oauth.RefreshToken))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(oauth.User))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.Client))

	suite.resetMocks()
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

// Reset mocks
func (suite *AccountsTestSuite) resetMocks() {
	suite.emailService.ExpectedCalls = suite.emailService.ExpectedCalls[:0]
	suite.emailService.Calls = suite.emailService.Calls[:0]
	suite.emailFactory.ExpectedCalls = suite.emailFactory.ExpectedCalls[:0]
	suite.emailFactory.Calls = suite.emailFactory.Calls[:0]
}

// Checks that the mock object expectations were met
func (suite *AccountsTestSuite) assertMockExpectations() {
	suite.emailService.AssertExpectations(suite.T())
	suite.emailFactory.AssertExpectations(suite.T())
	suite.resetMocks()
}

// Mock sending confirmation email
func (suite *AccountsTestSuite) mockConfirmationEmail() {
	messageMock := new(email.Message)
	suite.emailFactory.On("NewConfirmationEmail", mock.AnythingOfType("*accounts.Confirmation")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

// Mock sending invitation email
func (suite *AccountsTestSuite) mockInvitationEmail() {
	messageMock := new(email.Message)
	suite.emailFactory.On("NewInvitationEmail", mock.AnythingOfType("*accounts.Invitation")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

// Mock sending password reset email
func (suite *AccountsTestSuite) mockPasswordResetEmail() {
	messageMock := new(email.Message)
	suite.emailFactory.On("NewPasswordResetEmail", mock.AnythingOfType("*accounts.PasswordReset")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

func (suite *AccountsTestSuite) insertTestUser(email, password, firstName, lastName string) (*accounts.User, *oauth.AccessToken, error) {
	var (
		testOauthUser   *oauth.User
		testUser        *accounts.User
		testAccessToken *oauth.AccessToken
		err             error
	)

	// Insert a test user
	testOauthUser, err = suite.service.GetOauthService().CreateUser(
		roles.User,
		email,
		password,
	)
	if err != nil {
		return nil, nil, err
	}

	testUser, err = accounts.NewUser(
		suite.accounts[0],
		testOauthUser,
		"",    //facebook ID
		false, // confirmed
		&accounts.UserRequest{
			FirstName: firstName,
			LastName:  lastName,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	err = suite.db.Create(testUser).Error
	testUser.Account = suite.accounts[0]
	testUser.OauthUser = testOauthUser
	if err != nil {
		return nil, nil, err
	}

	// Login the test user
	testAccessToken, _, err = suite.service.GetOauthService().Login(
		suite.accounts[0].OauthClient,
		testUser.OauthUser,
		"read_write", // scope
	)
	if err != nil {
		return nil, nil, err
	}

	return testUser, testAccessToken, nil
}
