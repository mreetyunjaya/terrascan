package runtime

import (
	"os"

	"github.com/accurics/terrascan/pkg/utils"
	"go.uber.org/zap"

	cloudProvider "github.com/accurics/terrascan/pkg/cloud-providers"
	iacProvider "github.com/accurics/terrascan/pkg/iac-providers"
	"github.com/accurics/terrascan/pkg/iac-providers/output"
)

// Executor object
type Executor struct {
	filePath      string
	dirPath       string
	cloudType     string
	iacType       string
	iacVersion    string
	iacProvider   iacProvider.IacProvider
	cloudProvider cloudProvider.CloudProvider
}

// NewExecutor creates a runtime object
func NewExecutor(iacType, iacVersion, cloudType, filePath, dirPath string) (e *Executor, err error) {
	e = &Executor{
		filePath:   filePath,
		dirPath:    dirPath,
		cloudType:  cloudType,
		iacType:    iacType,
		iacVersion: iacVersion,
	}

	// initialized executor
	if err = e.Init(); err != nil {
		return e, err
	}

	return e, nil
}

// Init validates input and initializes iac and cloud providers
func (e *Executor) Init() error {

	// validate inputs
	err := e.ValidateInputs()
	if err != nil {
		return err
	}

	// create new IacProvider
	e.iacProvider, err = iacProvider.NewIacProvider(e.iacType, e.iacVersion)
	if err != nil {
		zap.S().Errorf("failed to create a new IacProvider for iacType '%s'. error: '%s'", e.iacType, err)
		return err
	}

	// create new CloudProvider
	e.cloudProvider, err = cloudProvider.NewCloudProvider(e.cloudType)
	if err != nil {
		zap.S().Errorf("failed to create a new CloudProvider for cloudType '%s'. error: '%s'", e.cloudType, err)
		return err
	}

	return nil
}

// Execute validates the inputs, processes the IaC, creates json output
func (e *Executor) Execute() error {

	// load iac config
	var (
		iacOut output.AllResourceConfigs
		err    error
	)
	if e.dirPath != "" {
		iacOut, err = e.iacProvider.LoadIacDir(e.dirPath)
	} else {
		// create config from IaC
		iacOut, err = e.iacProvider.LoadIacFile(e.filePath)
	}
	if err != nil {
		return err
	}

	// create normalized json
	normalized, err := e.cloudProvider.CreateNormalizedJSON(iacOut)
	if err != nil {
		return err
	}
	utils.PrintJSON(normalized, os.Stdout)

	// write output

	// successful
	return nil
}