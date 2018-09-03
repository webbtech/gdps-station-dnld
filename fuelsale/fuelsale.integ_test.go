package fuelsale

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date             = "2018-05-01"
	defaultsFilePath = "../config/defaults.yaml"
	stationID        = "d03224a7-f1df-4863-bcaa-5c6e61af11fc"
	timeFormat       = "2006-01-02"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	c       *config.Config
	request *model.Request
	report  *Report
}

// SetupTest method
func (suite *UnitSuite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	suite.request = &model.Request{
		Date:      dte,
		StationID: stationID,
	}
	suite.c = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.c.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.c)

	suite.report, err = New(suite.request, suite.c)
	suite.NoError(err)
	suite.IsType(new(Report), suite.report)

	err = suite.report.Create()
	suite.NoError(err)
}

// TestConfig method
func (suite *UnitSuite) TestConfig() {
	suite.NotEqual("", suite.c.AWSRegion, "Expected AWSRegion to be populated")
}

// TestSaveToDisk method
func (suite *UnitSuite) TestSaveToDisk() {
	err := suite.report.Create()
	suite.NoError(err)

	fp, err := suite.report.SaveToDisk("../tmp")
	suite.NoError(err)
	suite.NotEqual("", fp, "Expected file path to be populated")
}

// TestCreateSignedURL method
func (suite *UnitSuite) TestCreateSignedURL() {
	err := suite.report.Create()
	suite.NoError(err)

	url, err := suite.report.CreateSignedURL()
	suite.NoError(err)
	suite.NotEqual("", url, "Expected url to be populated")

	response, err := http.Get(url)
	suite.NoError(err)
	defer response.Body.Close()
	suite.Equal(200, response.StatusCode, "Expect response code to be 200")
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
