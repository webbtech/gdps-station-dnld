package xlsx

import (
	"os"
	"testing"
	"time"

	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/graphql"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/stretchr/testify/suite"
)

const (
	date             = "2018-05-01"
	defaultsFilePath = "../config/defaults.yaml"
	filePath         = "../tmp/testfile.xlsx"
	stationID        = "d03224a7-f1df-4863-bcaa-5c6e61af11fc"
	timeFormat       = "2006-01-02"
)

// Suite struct
type Suite struct {
	suite.Suite
	cfg     *config.Config
	request *model.Request
	file    *XLSX
	graphql *graphql.Client
}

// SetupTest method
func (suite *Suite) SetupTest() {
	os.Setenv("Stage", "test")
	dte, err := time.Parse(timeFormat, date)
	suite.request = &model.Request{
		Date:      dte,
		StationID: stationID,
	}
	suite.cfg = &config.Config{DefaultsFilePath: defaultsFilePath}
	err = suite.cfg.Load()
	suite.NoError(err)
	suite.IsType(new(config.Config), suite.cfg)

	suite.file, err = NewFile()
	suite.NoError(err)
	suite.IsType(new(XLSX), suite.file)

	suite.graphql = graphql.New(suite.request, suite.cfg)
	suite.IsType(new(graphql.Client), suite.graphql)
}

// TestOutput method
func (suite *Suite) TestOutput() {

	fs, err := suite.graphql.FuelSales()
	suite.NoError(err)
	suite.IsType(new(model.FuelSales), fs)

	fd, err := suite.graphql.FuelDelivery()
	suite.NoError(err)
	suite.IsType(new(model.FuelDelivery), fd)

	osm, err := suite.graphql.OverShortMonth()
	suite.NoError(err)
	suite.IsType(new(model.OverShortMonth), osm)

	osa, err := suite.graphql.OverShortAnnual()
	suite.NoError(err)
	suite.IsType(new(model.OverShortAnnual), osa)

	err = suite.file.FuelSales(fs)
	suite.NoError(err)

	err = suite.file.FuelDelivery(fd)
	suite.NoError(err)

	err = suite.file.OverShortMonth(osm)
	suite.NoError(err)

	err = suite.file.OverShortAnnual(osa)
	suite.NoError(err)

	err = suite.file.OutputToDisk(filePath)
	suite.NoError(err)
}

// TestXLSXSuite function
func TestXLSXSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
