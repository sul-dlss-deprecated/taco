// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

// SwaggerJSON embedded version of the swagger document used at generation time
var SwaggerJSON json.RawMessage

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "TACO, the Stanford Digital Repository (SDR) Management Layer API",
    "title": "taco",
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "0.1.0"
  },
  "host": "sdr.dlss.stanford.edu",
  "basePath": "/v1",
  "paths": {
    "/file": {
      "post": {
        "description": "Deposits a new File (binary) into SDR. Will return the SDR identifier for the File resource (aka the metadata object generated and persisted for management of the provided binary).",
        "consumes": [
          "multipart/form-data"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Deposit New File (binary).",
        "operationId": "depositFile",
        "parameters": [
          {
            "type": "file",
            "description": "Binary to be added to an Object in TACO.",
            "name": "upload",
            "in": "formData",
            "required": true
          }
        ],
        "responses": {
          "201": {
            "description": "TACO binary ingested, File management metadata created, \u0026 File processing started.",
            "schema": {
              "$ref": "#/definitions/ResourceResponse"
            }
          },
          "401": {
            "description": "You are not authorized to ingest a File into TACO."
          },
          "415": {
            "description": "Unsupported file type provided."
          },
          "500": {
            "description": "This file could be ingested at this time by TACO."
          }
        }
      }
    },
    "/healthcheck": {
      "get": {
        "description": "The healthcheck endpoint provides information about the health of the service.",
        "summary": "Health Check",
        "operationId": "healthCheck",
        "responses": {
          "200": {
            "description": "The service is functioning nominally",
            "schema": {
              "$ref": "#/definitions/HealthCheckResponse"
            }
          },
          "503": {
            "description": "The service is not working correctly",
            "schema": {
              "$ref": "#/definitions/HealthCheckResponse"
            }
          }
        }
      }
    },
    "/resource": {
      "post": {
        "description": "Deposits a new resource (Collection, Digital Repository Object, File [metadata only] or subclass of those) into SDR. Will return the SDR identifier for the resource.",
        "consumes": [
          "application/json",
          "application/json+ld"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Deposit New TACO Resource.",
        "operationId": "depositResource",
        "parameters": [
          {
            "description": "JSON-LD representation of the resource metadata going into SDR. Needs to fit the SDR 3.0 MAP requirements.",
            "name": "payload",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Resource"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "TACO resource created \u0026 processing started.",
            "schema": {
              "$ref": "#/definitions/ResourceResponse"
            }
          },
          "401": {
            "description": "You are not authorized to create a resource in TACO."
          },
          "415": {
            "description": "Unsupported resource type provided. TACO resources should be handed over as JSON or JSON-LD."
          },
          "422": {
            "description": "The resource JSON provided had an unspecified or unsupported field, or is otherwise unprocessable by TACO.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "This resource could be created at this time by TACO."
          }
        }
      }
    },
    "/resource/{ID}": {
      "get": {
        "description": "Retrieves the metadata (as JSON-LD following our SDR3 MAP v.1) for an existing TACO resource (Collection, Digital Repository Object, File metadata object [not binary] or subclass of those). The resource is identified by the TACO identifier.",
        "produces": [
          "application/json"
        ],
        "summary": "Retrieve TACO Resource Metadata.",
        "operationId": "retrieveResource",
        "parameters": [
          {
            "type": "string",
            "description": "TACO Resource Identifier.",
            "name": "ID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Resource metadata retrieved.",
            "schema": {
              "$ref": "#/definitions/Resource"
            }
          },
          "401": {
            "description": "You are not authorized to view this resource in TACO."
          },
          "404": {
            "description": "Resource not found. Please check your provided TACO identifier."
          },
          "500": {
            "description": "The resource could not be retrieved by TACO at this time."
          }
        }
      },
      "patch": {
        "description": "Updates the metadata for an existing TACO resource (Collection, Digital Repository Object, File metadata object [not binary] or subclass of those). Only include the required fields and the fields you wish to have changed. Will return the TACO resource identifier.",
        "consumes": [
          "application/json",
          "application/json+ld"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Update TACO Resource.",
        "operationId": "updateResource",
        "parameters": [
          {
            "type": "string",
            "description": "SDR Identifier for the Resource.",
            "name": "ID",
            "in": "path",
            "required": true
          },
          {
            "description": "JSON-LD Representation of the resource metadata required fields and only the fields you wish to update (identified via its TACO identifier). Needs to fit the SDR 3.0 MAP requirements.",
            "name": "payload",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Resource"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "TACO resource metadata updated \u0026 processing started.",
            "schema": {
              "$ref": "#/definitions/ResourceResponse"
            }
          },
          "400": {
            "description": "Invalid ID supplied"
          },
          "401": {
            "description": "You are not authorized to update a resource in TACO."
          },
          "415": {
            "description": "Unsupported resource type provided. TACO resources should be handed over as JSON or JSON-LD."
          },
          "422": {
            "description": "The resource JSON provided had an unspecified or unsupported field, or is otherwise unprocessable by TACO.",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "500": {
            "description": "This resource could be updated at this time by TACO."
          }
        }
      }
    },
    "/status/{ID}": {
      "get": {
        "description": "Get the processing status and history for a resource.",
        "produces": [
          "application/json"
        ],
        "summary": "Resource Processing Status.",
        "operationId": "getProcessStatus",
        "parameters": [
          {
            "type": "string",
            "description": "SDR Identifier for the Resource.",
            "name": "ID",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Processing status for the TACO resource.",
            "schema": {
              "$ref": "#/definitions/ProcessResponse"
            }
          },
          "401": {
            "description": "You are not authorized to view this resource's processing status in TACO."
          },
          "404": {
            "description": "Resource not found. Please check your provided TACO identifier."
          },
          "500": {
            "description": "This resource's processing status could be retrieved at this time by TACO."
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "properties": {
        "detail": {
          "description": "a human-readable explanation specific to this occurrence of the problem.",
          "type": "string",
          "example": "Title must contain at least three characters."
        },
        "source": {
          "type": "object",
          "properties": {
            "pointer": {
              "type": "string",
              "example": "/data/attributes/title"
            }
          }
        },
        "title": {
          "description": "a short, human-readable summary of the problem that SHOULD NOT change from occurrence to occurrence of the problem.",
          "type": "string",
          "example": "Invalid Attribute"
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "errors": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Error"
          }
        }
      }
    },
    "HealthCheckResponse": {
      "type": "object",
      "properties": {
        "status": {
          "description": "The status of the service",
          "type": "string"
        }
      },
      "example": {
        "status": "OK"
      }
    },
    "ProcessResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "oo000oo0001"
        }
      },
      "example": {
        "id": "oo000oo0001"
      }
    },
    "Resource": {
      "type": "object",
      "required": [
        "@context",
        "@type",
        "access",
        "label",
        "preserve",
        "publish"
      ],
      "properties": {
        "@context": {
          "description": "URI for the JSON-LD context definitions",
          "type": "string",
          "format": "uri",
          "pattern": "http://sdr\\.sul\\.stanford\\.edu/contexts/taco-base\\.jsonld",
          "example": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld"
        },
        "@type": {
          "description": "URI for the resource type",
          "type": "string",
          "format": "uri",
          "pattern": "http://sdr\\.sul\\.stanford\\.edu/models/sdr3-(object|collection|file)\\.jsonld",
          "example": "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld"
        },
        "access": {
          "description": "What groups should be able to access (view) the resource in Access environments",
          "type": "string",
          "enum": [
            "world",
            "stanford",
            "location-based",
            "citation-only",
            "dark"
          ]
        },
        "contained-by": {
          "description": "The parent resource(s) of this resource.",
          "type": "array",
          "items": {
            "type": "string",
            "format": "uri"
          }
        },
        "contains": {
          "description": "The child resource(s) of this resource.",
          "type": "array",
          "items": {
            "type": "string",
            "format": "uri"
          }
        },
        "id": {
          "description": "The TACO identifier for the resource. Usually DRUID-derived.",
          "type": "string",
          "example": "oo000oo0001"
        },
        "label": {
          "description": "The label or processing title for the resource.",
          "type": "string",
          "example": "Label for this resource"
        },
        "preserve": {
          "description": "Should the resource be released to Preservation environments",
          "type": "boolean"
        },
        "publish": {
          "description": "Should the resource's metadata be released to Access environments",
          "type": "boolean"
        },
        "sourceId": {
          "description": "The source identifier (bib id, archival id) for the resource that was digitized or derived from to create the TACO resource.",
          "type": "string",
          "example": "bib12345678"
        }
      },
      "example": {
        "@context": "http://sdr.sul.stanford.edu/contexts/taco-base.jsonld",
        "@type": "http://sdr.sul.stanford.edu/models/sdr3-object.jsonld",
        "access": "world",
        "id": "oo000oo0001",
        "label": "My SDR3 resource",
        "preserve": true,
        "publish": true,
        "sourceId": "bib12345678"
      }
    },
    "ResourceResponse": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "example": "oo000oo0001"
        }
      },
      "example": {
        "id": "oo000oo0001"
      }
    }
  }
}`))
}
