package fuelsale

import (
	"path"

	"github.com/pulpfree/gdps-fs-dwnld/awsservices"
	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/graphql"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/pulpfree/gdps-fs-dwnld/xlsx"

	log "github.com/sirupsen/logrus"
)

// ReportName constant
const (
	reportFileName = "StationReport"
	timeFrmt       = "2006-01"
)

// Report struct
type Report struct {
	authToken string
	cfg       *config.Config
	request   *model.Request
	file      *xlsx.XLSX
	filenm    string
}

// New function
func New(req *model.Request, cfg *config.Config, authToken string) (r *Report, err error) {
	r = &Report{
		authToken: authToken,
		cfg:       cfg,
		request:   req,
	}
	return r, err
}

// Create method
func (r *Report) Create() (err error) {

	// Init graphql and xlsx packages
	client := graphql.New(r.request, r.cfg, r.authToken)
	r.file, err = xlsx.NewFile()
	if err != nil {
		return err
	}

	// Fetch and create Fuel Sales
	fs, err := client.FuelSales()
	if err != nil {
		log.Errorf("Error fetching FuelSales: %s", err)
		return err
	}
	// Now that we have the station name, we can set
	// would be nice to do earlier, but requires separate query
	r.setFileName(fs.Station.Name)

	err = r.file.FuelSales(fs)
	if err != nil {
		return err
	}

	// Fetch and create NL and DSL Fuel Sales by Station
	fsl, err := client.FuelSalesList()
	if err != nil {
		log.Errorf("Error fetching FuelSalesList: %s", err)
		return err
	}
	err = r.file.FuelSalesListNL(fsl)
	if err != nil {
		log.Errorf("Error creating FuelSalesListNL: %s", err)
		return err
	}

	err = r.file.FuelSalesListDSL(fsl)
	if err != nil {
		log.Errorf("Error creating FuelSalesListDSL: %s", err)
		return err
	}

	// Fetch and create Fuel Delivery
	fd, err := client.FuelDelivery()
	if err != nil {
		log.Errorf("Error fetching FuelDelivery: %s", err)
		return err
	}
	err = r.file.FuelDelivery(fd)
	if err != nil {
		return err
	}

	// Fetch and create monthly overshort
	osm, err := client.OverShortMonth()
	if err != nil {
		log.Errorf("Error fetching FuelSales: %s", err)
		return err
	}
	err = r.file.OverShortMonth(osm)
	if err != nil {
		return err
	}

	// Fetch and create annual overshort
	osa, err := client.OverShortAnnual()
	if err != nil {
		log.Errorf("Error fetching FuelSales: %s", err)
		return err
	}
	err = r.file.OverShortAnnual(osa)
	if err != nil {
		return err
	}

	return err
}

// SaveToDisk method
func (r *Report) SaveToDisk(dir string) (fp string, err error) {

	filePath := path.Join(dir, r.getFileName())
	fp, err = r.file.OutputToDisk(filePath)
	if err != nil {
		return "", err
	}
	return fp, err
}

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {

	output, err := r.file.OutputFile()
	if err != nil {
		return "", err
	}

	s3Serv, err := awsservices.NewS3(r.cfg)
	filePrefix := path.Join(r.cfg.S3FilePrefix, r.getFileName())

	return s3Serv.GetSignedURL(filePrefix, &output)
}

//
// ======================== Helper Functions =============================== //
//

func (r *Report) setFileName(stationName string) {
	r.filenm = stationName + "_" + reportFileName + "_" + r.request.Date.Format(timeFrmt) + ".xlsx"
}

func (r *Report) getFileName() string {
	return r.filenm
}
