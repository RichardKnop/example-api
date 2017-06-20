package accounts_test

import (
	"log"
	"os"
	"testing"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/accounts"
	"github.com/RichardKnop/example-api/services/email"
	"github.com/RichardKnop/example-api/services/oauth"
	"github.com/RichardKnop/example-api/services/oauth/roles"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	accountsMocks "github.com/RichardKnop/example-api/services/accounts/mocks"
	emailMocks "github.com/RichardKnop/example-api/services/email/mocks"
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
		"./accounts/fixtures/test_users.yml",
	}

	testMigrations = []func(*gorm.DB) error{
		models.MigrateAll,
	}
)

func init() {
	if err := os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}

// oauth clientsTestSuite needs to be exported so the tests run
type AccountsTestSuite struct {
	suite.Suite
	cnf          *config.Config
	db           *gorm.DB
	emailService *emailMocks.ServiceInterface
	emailFactory *accountsMocks.EmailFactoryInterface
	service      *accounts.Service
	clients      []*models.OauthClient
	users        []*models.User
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

	// Fetch test clients
	suite.clients = make([]*models.OauthClient, 0)
	err = suite.db.Order("id").Find(&suite.clients).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*models.User, 0)
	err = models.UserPreload(suite.db).Order("id").Find(&suite.users).Error
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
	suite.db.Unscoped().Delete(new(models.Confirmation))
	suite.db.Unscoped().Delete(new(models.Invitation))
	suite.db.Unscoped().Delete(new(models.PasswordReset))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(models.User))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(models.OauthAccessToken))
	suite.db.Unscoped().Delete(new(models.OauthRefreshToken))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(models.OauthUser))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(models.OauthClient))

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
	suite.emailFactory.On("NewConfirmationEmail", mock.AnythingOfType("*models.Confirmation")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

// Mock sending invitation email
func (suite *AccountsTestSuite) mockInvitationEmail() {
	messageMock := new(email.Message)
	suite.emailFactory.On("NewInvitationEmail", mock.AnythingOfType("*models.Invitation")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

// Mock sending password reset email
func (suite *AccountsTestSuite) mockPasswordResetEmail() {
	messageMock := new(email.Message)
	suite.emailFactory.On("NewPasswordResetEmail", mock.AnythingOfType("*models.PasswordReset")).
		Return(messageMock, nil)
	suite.emailService.On("Send", messageMock).Return(nil)
}

func (suite *AccountsTestSuite) insertTestUser(email, password, firstName, lastName string) (*models.User, *models.OauthAccessToken, error) {
	var (
		testOauthUser   *models.OauthUser
		testUser        *models.User
		testAccessToken *models.OauthAccessToken
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

	testUser, err = models.NewUser(
		suite.clients[0],
		testOauthUser,
		"", //facebook ID
		firstName,
		lastName,
		"",    // picture
		false, // confirmed
	)
	if err != nil {
		return nil, nil, err
	}
	err = suite.db.Create(testUser).Error
	testUser.OauthClient = suite.clients[0]
	testUser.OauthUser = testOauthUser
	if err != nil {
		return nil, nil, err
	}

	// Login the test user
	testAccessToken, _, err = suite.service.GetOauthService().Login(
		suite.clients[0],
		testUser.OauthUser,
		"read_write", // scope
	)
	if err != nil {
		return nil, nil, err
	}

	return testUser, testAccessToken, nil
}
