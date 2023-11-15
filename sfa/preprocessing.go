package sfa

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	"github.com/jraddaoui/preprocessing/preprocessing/activities"
	"github.com/jraddaoui/preprocessing/spec"
)

var _ spec.Preprocessing = (*Preprocessing)(nil)

type Preprocessing struct{}

func (s *Preprocessing) Activities() []spec.Activity {
	return []spec.Activity{
		{Name: activities.ExtractPackageName, Execute: activities.NewExtractPackage()},
		{Name: activities.CheckSipStructureName, Execute: activities.NewCheckSipStructure()},
		{Name: activities.AllowedFileFormatsName, Execute: activities.NewAllowedFileFormatsActivity()},
		{Name: activities.MetadataValidationName, Execute: activities.NewMetadataValidationActivity()},
		{Name: activities.SipCreationName, Execute: activities.NewSipCreationActivity()},
	}
}

func (s *Preprocessing) Execute(ctx workflow.Context, params spec.Params) (spec.Result, error) {
	result := spec.Result{}
	preProcCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2,
			MaximumInterval:    time.Second * 10,
			MaximumAttempts:    2,
			NonRetryableErrorTypes: []string{
				"TemporalTimeout:StartToClose",
			},
		},
	})

	// Extract package.
	var extractPackageRes activities.ExtractPackageResult
	err := workflow.ExecuteActivity(preProcCtx, activities.ExtractPackageName, &activities.ExtractPackageParams{
		Path: params.Path,
		Key:  "path_should_include_the_proper_extension.zip",
	}).Get(ctx, &extractPackageRes)
	if err != nil {
		return result, err
	}

	// Validate SIP structure.
	var checkStructureRes activities.CheckSipStructureResult
	err = workflow.ExecuteActivity(preProcCtx, activities.CheckSipStructureName, &activities.CheckSipStructureParams{SipPath: extractPackageRes.Path}).Get(ctx, &checkStructureRes)
	if err != nil {
		return result, err
	}

	var allowedFileFormats activities.AllowedFileFormatsResult
	err = workflow.ExecuteActivity(preProcCtx, activities.AllowedFileFormatsName, &activities.AllowedFileFormatsParams{SipPath: extractPackageRes.Path}).Get(ctx, &allowedFileFormats)
	if err != nil {
		return result, err
	}

	// Validate metadata.xsd.
	var metadataValidation activities.MetadataValidationResult
	err = workflow.ExecuteActivity(preProcCtx, activities.MetadataValidationName, &activities.MetadataValidationParams{SipPath: extractPackageRes.Path}).Get(ctx, &metadataValidation)
	if err != nil {
		return result, err
	}

	// Repackage SFA Sip into a Bag.
	var sipCreation activities.SipCreationResult
	err = workflow.ExecuteActivity(preProcCtx, activities.SipCreationName, &activities.SipCreationParams{SipPath: extractPackageRes.Path}).Get(ctx, &sipCreation)
	if err != nil {
		return result, err
	}

	// We do this so that the code above only stops when a non-bussines error is found.
	if !allowedFileFormats.Ok {
		return result, activities.ErrIlegalFileFormat
	}
	if !checkStructureRes.Ok {
		return result, activities.ErrInvaliSipStructure
	}

	return spec.Result{Path: sipCreation.NewSipPath}, nil
}
