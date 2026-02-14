# Product Definition

## Initial Concept
The initial concept for this project is to derive the bigQuery schema from ndjson files.

## Vision

## Goals


## Key Features
Key Design Principles:
- Performance through concurrency (using golang channel primitives)
- Big data support through splitted input files (processed in parallel & sampling)
- Enforce consistency checks upstream (avoid things fail downstream)
- Fail fast in case of an error
- File position traceback to enable debugging
- Architect for reusability (generate other destination schema for AlloyDB, Spanner, ...)

Constraints:
- A process will generate one schema. To enable generation of different schemas seperate processes need to be instantiated.

Features:
- Recursively traverse nested and repeated fields
- Core functionality is structured in a reusable packages
- Command line interface (CLI) with enabling configurability
- Single binary executable. Should play well together with Google CLI tools like gcloud and bq.
- Basic test case automation

Not supported yet:
- API integration with Google Cloud (Cloud Storage, Bigquery)
- Constraint relaxation for the "Required" attribute
- No handling for the JSON string literal "null"
- Distribution or clustered deployment
- Schema inference by parsing through repeated records (now the nested schema is based on the first element)
- Incomplete support for string-wrapped types: Timestamp, Time, Date, Geo-types, ...
- Deployment on Dataflow

## Target Audience
- Developers looking for a cli tool to genrate a BIgQuery schema from ndjson data files

