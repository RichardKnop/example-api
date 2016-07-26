package accounts_test

import (
	"log"
	"os"
	"testing"

	"github.com/RichardKnop/example-api/accounts"
	"github.com/RichardKnop/example-api/accounts/roles"
	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/database"
	"github.com/RichardKnop/example-api/email"
	"github.com/RichardKnop/example-api/oauth"
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
)

var testFixtures = []string{
	"./oauth/fixtures/test_clients.yml",
	"./oauth/fixtures/test_users.yml",
	"./oauth/fixtures/test_access_tokens.yml",
	"./accounts/fixtures/roles.yml",
	"./accounts/fixtures/test_accounts.yml",
	"./accounts/fixtures/test_users.yml",
}

// db migrations needed for tests
var testMigrations = []func(*gorm.DB) error{
	oauth.MigrateAll,
	accounts.MigrateAll,
}

func init() {
	if err := os.Chdir("../"); err != nil {
		log.Fatal(err)
	}
}

// AccountsTestSuite needs to be exported so the tests run
type AccountsTestSuite struct {
	suite.Suite
	cnf           *config.Config
	db            *gorm.DB
	emailService  *emailMocks.ServiceInterface
	emailFactory  *accountsMocks.EmailFactoryInterface
	service       *accounts.Service
	accounts      []*accounts.Account
	users         []*accounts.User
	superuserRole *accounts.Role
	userRole      *accounts.Role
	router        *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *AccountsTestSuite) SetupSuite() {

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

	// Fetch test roles
	suite.superuserRole = new(accounts.Role)
	err = suite.db.Where("id = ?", roles.Superuser).First(&suite.superuserRole).Error
	if err != nil {
		log.Fatal(err)
	}
	suite.userRole = new(accounts.Role)
	err = suite.db.Where("id = ?", roles.User).First(&suite.userRole).Error
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
	accounts.RegisterRoutes(suite.router, suite.service)
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
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(oauth.AccessToken))
	suite.db.Unscoped().Delete(new(oauth.RefreshToken))

	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(accounts.User))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(accounts.Account))

	// Service.CreateUser also creates a new oauth.User instance
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(oauth.User))

	// Service.CreateAccount also creates a new oauth.Client instance
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.Client))

	// Reset mocks
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
	emailMock := new(email.Email)
	suite.emailFactory.On(
		"NewConfirmationEmail",
		mock.AnythingOfType("*accounts.Confirmation"),
	).Return(emailMock, nil)
	suite.emailService.On("Send", emailMock).Return(nil)
}

// Mock sending invitation email
func (suite *AccountsTestSuite) mockInvitationEmail() {
	emailMock := new(email.Email)
	suite.emailFactory.On(
		"NewInvitationEmail",
		mock.AnythingOfType("*accounts.Invitation"),
	).Return(emailMock, nil)
	suite.emailService.On("Send", emailMock).Return(nil)
}

// Mock sending password reset email
func (suite *AccountsTestSuite) mockPasswordResetEmail() {
	emailMock := new(email.Email)
	suite.emailFactory.On(
		"NewPasswordResetEmail",
		mock.AnythingOfType("*accounts.PasswordReset"),
	).Return(emailMock, nil)
	suite.emailService.On("Send", emailMock).Return(nil)
}
