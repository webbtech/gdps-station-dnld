package fuelsale

import (
	"github.com/pulpfree/gdps-fs-dwnld/awsservices"
	"github.com/pulpfree/gdps-fs-dwnld/config"
	"github.com/pulpfree/gdps-fs-dwnld/graphql"
	"github.com/pulpfree/gdps-fs-dwnld/model"
	"github.com/pulpfree/gdps-fs-dwnld/model/dynamo"
	"github.com/pulpfree/gdps-fs-dwnld/xlsx"

	log "github.com/sirupsen/logrus"
)

// Report struct
type Report struct {
	cfg     *config.Config
	db      *dynamo.Dynamo
	request *model.Request
}

// New function
func New(req *model.Request, cfg *config.Config) (r *Report, err error) {

	// Set DynamoDB connection
	dynamo, err := dynamo.NewDB(cfg.Dynamo)
	if err != nil {
		log.Errorf("Error connecting to dynamo: %s", err)
		return r, err
	}

	r = &Report{
		cfg:     cfg,
		db:      dynamo,
		request: req,
	}
	return r, err
}

// Create method
func (r *Report) Create() (signedURL string, err error) {

	client := graphql.New(r.request, r.cfg)
	fs, err := client.FuelSales()
	if err != nil {
		log.Errorf("Error fetching FuelSales: %s", err)
		return "", err
	}

	os, err := client.OverShortMonth()
	if err != nil {
		log.Errorf("Error fetching FuelSales: %s", err)
		return "", err
	}

	file, err := xlsx.NewFile()
	err = file.FuelSales(fs)
	if err != nil {
		return "", err
	}

	err = file.OverShortMonth(os)
	if err != nil {
		return "", err
	}

	output, err := file.OutputFile()
	if err != nil {
		return "", err
	}

	s3Serv, err := awsservices.NewS3(r.cfg)
	filePrefix := "tankfiles/testFile.xslx"

	signedURL, err = s3Serv.GetSignedURL(filePrefix, &output)
	if err != nil {
		return "", err
	}

	return signedURL, err
}
