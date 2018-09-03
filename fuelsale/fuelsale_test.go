package fuelsale

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/stretchr/testify/assert"
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
}

// TestConfig method
func (suite *UnitSuite) TestConfig() {
	assert.NotEqual(suite.T(), "", suite.c.AWSRegion, "Expected AWSRegion to be populated")
}

// TestCreateReport method
func (suite *UnitSuite) TestCreateReport() {
	url, err := suite.report.Create()
	suite.NoError(err)
	suite.NotEqual("", url, "Expected signed url to be populated")
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
