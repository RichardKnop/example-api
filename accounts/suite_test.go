package accounts

import (
	"log"
	"testing"

	"github.com/RichardKnop/recall/accounts/roles"
	"github.com/RichardKnop/recall/config"
	"github.com/RichardKnop/recall/database"
	"github.com/RichardKnop/recall/email"
	"github.com/RichardKnop/recall/oauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	testDbUser = "recall"
	testDbName = "recall_accounts_test"
)

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
	cnf              *config.Config
	db               *gorm.DB
	emailServiceMock *email.ServiceMock
	emailFactoryMock *EmailFactoryMock
	service          *Service
	accounts         []*Account
	users            []*User
	superuserRole    *Role
	userRole         *Role
	router           *mux.Router
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
	suite.accounts = make([]*Account, 0)
	err = suite.db.Preload("OauthClient").Order("id").Find(&suite.accounts).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*User, 0)
	err = suite.db.Preload("Account").Preload("OauthUser").Preload("Role").
		Order("id").Find(&suite.users).Error
	if err != nil {
		log.Fatal(err)
	}

	// Fetch test roles
	suite.superuserRole = new(Role)
	err = suite.db.Where("id = ?", roles.Superuser).First(&suite.superuserRole).Error
	if err != nil {
		log.Fatal(err)
	}
	suite.userRole = new(Role)
	err = suite.db.Where("id = ?", roles.User).First(&suite.userRole).Error
	if err != nil {
		log.Fatal(err)
	}

	// Initialise mocks
	suite.emailServiceMock = new(email.ServiceMock)
	suite.emailFactoryMock = new(EmailFactoryMock)

	// Initialise the service
	suite.service = NewService(
		suite.cnf,
		suite.db,
		oauth.NewService(suite.cnf, suite.db),
		suite.emailServiceMock,
		suite.emailFactoryMock,
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
	suite.db.Unscoped().Delete(new(Confirmation))
	suite.db.Unscoped().Delete(new(Invitation))
	suite.db.Unscoped().Delete(new(PasswordReset))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3, 4}).Delete(new(oauth.AccessToken))
	suite.db.Unscoped().Delete(new(oauth.RefreshToken))

	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(User))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(Account))

	// Service.CreateUser also creates a new oauth.User instance
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(oauth.User))

	// Service.CreateAccount also creates a new oauth.Client instance
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(oauth.Client))

	// Reset mocks
	suite.emailServiceMock.ExpectedCalls = suite.emailServiceMock.ExpectedCalls[:0]
	suite.emailServiceMock.Calls = suite.emailServiceMock.Calls[:0]
	suite.emailFactoryMock.ExpectedCalls = suite.emailFactoryMock.ExpectedCalls[:0]
	suite.emailFactoryMock.Calls = suite.emailFactoryMock.Calls[:0]
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

// Checks that the mock object expectations were met
func (suite *AccountsTestSuite) assertMockExpectations() {
	suite.emailServiceMock.AssertExpectations(suite.T())
	suite.emailFactoryMock.AssertExpectations(suite.T())
}

// Mock sending confirmation email
func (suite *AccountsTestSuite) mockConfirmationEmail() {
	emailMock := new(email.Email)
	suite.emailFactoryMock.On(
		"NewConfirmationEmail",
		mock.AnythingOfType("*accounts.Confirmation"),
	).Return(emailMock)
	suite.emailServiceMock.On("Send", emailMock).Return(nil)
}

// Mock sending invitation email
func (suite *AccountsTestSuite) mockInvitationEmail() {
	emailMock := new(email.Email)
	suite.emailFactoryMock.On(
		"NewInvitationEmail",
		mock.AnythingOfType("*accounts.Invitation"),
	).Return(emailMock)
	suite.emailServiceMock.On("Send", emailMock).Return(nil)
}

// Mock sending password reset email
func (suite *AccountsTestSuite) mockPasswordResetEmail() {
	emailMock := new(email.Email)
	suite.emailFactoryMock.On(
		"NewPasswordResetEmail",
		mock.AnythingOfType("*accounts.PasswordReset"),
	).Return(emailMock)
	suite.emailServiceMock.On("Send", emailMock).Return(nil)
}
