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
        "security": [
          {
            "RemoteUser": []
          }
        ],
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
        "security": [
          {
            "RemoteUser": []
          }
        ],
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
        "security": [
          {
            "RemoteUser": []
          }
        ],
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
          },
          {
            "type": "string",
            "description": "The version of the requested resource",
            "name": "Version",
            "in": "query"
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
      "delete": {
        "description": "Deletes a TACO resource (Collection, Digital Repository Object, File resource (metadata) and File binary, or subclass of those).",
        "produces": [
          "application/json"
        ],
        "summary": "Delete a TACO Resource.",
        "operationId": "deleteResource",
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
          "204": {
            "description": "TACO resource metadata delete."
          },
          "401": {
            "description": "You are not authorized to delete a resource in TACO."
          },
          "404": {
            "description": "Invalid ID supplied"
          },
          "500": {
            "description": "This resource could not be deleted at this time by TACO."
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
        "security": [
          {
            "RemoteUser": []
          }
        ],
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
        "tacoIdentifier": {
          "type": "string",
          "example": "oo000oo0001"
        }
      },
      "example": {
        "tacoIdentifier": "oo000oo0001"
      }
    },
    "Resource": {
      "type": "object"
    },
    "ResourceResponse": {
      "type": "object"
    }
  },
  "securityDefinitions": {
    "RemoteUser": {
      "type": "apiKey",
      "name": "On-Behalf-Of",
      "in": "header"
    }
  }
}`))
}
