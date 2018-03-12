package resource

import (
	"github.com/sul-dlss-labs/taco/db"
	"github.com/sul-dlss-labs/taco/generated/restapi/operations"
)

func Create(database *db.Database, id string, params operations.DepositResourceParams) error {
	err := db.Insert(database, loadParams(id, params))
	return err
}

func loadParams(id string, params operations.DepositResourceParams) interface{} {
	// NOTE: This section will be replaced by DataUtils
	return map[string]interface{}{
		"id":        id,
		"attype":    params.Payload.AtType,
		"atcontext": params.Payload.AtContext,
		"access":    params.Payload.Access,
		"label":     params.Payload.Label,
		"preserve":  params.Payload.Preserve,
		"publish":   params.Payload.Publish,
		"sourceid":  params.Payload.SourceID,
	}
}
