package oauth_test

import (
	"log"
	"os"
	"testing"

	"github.com/RichardKnop/example-api/config"
	"github.com/RichardKnop/example-api/models"
	"github.com/RichardKnop/example-api/services/oauth"
	"github.com/RichardKnop/example-api/test-util"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

var (
	testDbUser = "example_api"
	testDbName = "example_api_oauth_test"

	testFixtures = []string{
		"./oauth/fixtures/scopes.yml",
		"./oauth/fixtures/roles.yml",
		"./oauth/fixtures/test_clients.yml",
		"./oauth/fixtures/test_users.yml",
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

// OauthTestSuite needs to be exported so the tests run
type OauthTestSuite struct {
	suite.Suite
	cnf     *config.Config
	db      *gorm.DB
	service *oauth.Service
	clients []*models.OauthClient
	users   []*models.OauthUser
	router  *mux.Router
}

// The SetupSuite method will be run by testify once, at the very
// start of the testing suite, before any tests are run.
func (suite *OauthTestSuite) SetupSuite() {
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

	// Fetch test client
	suite.clients = make([]*models.OauthClient, 0)
	if err := suite.db.Order("id").Find(&suite.clients).Error; err != nil {
		log.Fatal(err)
	}

	// Fetch test users
	suite.users = make([]*models.OauthUser, 0)
	if err := suite.db.Order("id").Find(&suite.users).Error; err != nil {
		log.Fatal(err)
	}

	// Initialise the service
	suite.service = oauth.NewService(suite.cnf, suite.db)

	// Register routes
	suite.router = mux.NewRouter()
	suite.service.RegisterRoutes(suite.router, "/v1/oauth")
}

// The TearDownSuite method will be run by testify once, at the very
// end of the testing suite, after all tests have been run.
func (suite *OauthTestSuite) TearDownSuite() {
	//
}

// The SetupTest method will be run before every test in the suite.
func (suite *OauthTestSuite) SetupTest() {
	//
}

// The TearDownTest method will be run after every test in the suite.
func (suite *OauthTestSuite) TearDownTest() {
	// Scopes are static, populated from fixtures,
	// so there is no need to clear them after running a test
	suite.db.Unscoped().Delete(new(models.OauthAuthorizationCode))
	suite.db.Unscoped().Delete(new(models.OauthRefreshToken))
	suite.db.Unscoped().Delete(new(models.OauthAccessToken))
	suite.db.Unscoped().Not("id", []int64{1, 2}).Delete(new(models.OauthUser))
	suite.db.Unscoped().Not("id", []int64{1, 2, 3}).Delete(new(models.OauthClient))
}

// TestOauthTestSuite ...
// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestOauthTestSuite(t *testing.T) {
	suite.Run(t, new(OauthTestSuite))
}
