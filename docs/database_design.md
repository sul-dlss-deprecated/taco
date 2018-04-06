# Database design

## Key schema
### Primary Key
The primary key is `tacoIdentifier`. This is a UUID generated for every record and is unique.

### Composite Keys
  * externalIdentifier + Version - This is composed of a DRUID (for Object or Collection) or a UUID (for Files and Filesets) and the version of the resource. This is indexed as `ResourceByExternalIDAndVersion`

### Unique Keys
  * Catkey - TODO: Is this truly unique? How does it interact with ID + Version?

## Dynamodb concerns
DynamoDB only allows top level attributes to be indexed. So many of the fields that we originally wished to put in the `identification` subschema have been moved to the root level.

>  Every attribute in the index key schema must be a top-level attribute of type String, Number, or Binary. Other data types, including documents and sets, are not allowed.

https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/SecondaryIndexes.html

*1MB limit* DynamoDB queries are limited to 1MB of data. If our records are
going to exceede that threshold we need to update `db/retrieve_latest.go` and `db/retrieve_version.go`, which use a query to return the full record.  Instead
we could query for the primary key and then do a `GetItem` request.
