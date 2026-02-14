# Specification for Core Schema Detection Logic

## 1. Overview
This document outlines the technical specifications for implementing the core schema detection logic for converting Newline Delimited JSON (NDJSON) files into a BigQuery schema. The implementation will address current limitations and enhance the robustness of the schema inference process.

## 2. Functional Requirements

### 2.1. Core Schema Inference
- The system must parse NDJSON files line by line, with each line representing a distinct JSON object.
- It must recursively traverse nested JSON objects and arrays to build a hierarchical schema structure.
- Data types for each field must be inferred based on the values encountered in the NDJSON data.

### 2.2. Handling of "null" literals
- The system must correctly handle the JSON string literal "null". When "null" is encountered for a field, it should not affect the inferred data type of that field. If a field contains only "null" values across all records, it should be treated as a `STRING` type by default.

### 2.3. Constraint Relaxation for "Required" Fields
- The `REQUIRED` mode for BigQuery fields will be inferred based on the presence of a field across all processed JSON objects.
- A field will be marked as `REQUIRED` only if it is present in every single JSON object in the NDJSON file.
- If a field is missing from any object, it will be marked as `NULLABLE`.

### 2.4. Improved Inference for Nested Schemas
- For fields that are arrays of objects (repeated records), the schema will be inferred by analyzing all elements of the array across multiple records, not just the first element of the first record encountered.
- The system will merge the schemas of all objects within a repeated field to create a comprehensive schema that accounts for all possible fields.

### 2.5. Data Type Mapping
The system will map JSON data types to BigQuery data types as follows:
- JSON String -> BigQuery `STRING`
- JSON Number (integer) -> BigQuery `INTEGER`
- JSON Number (floating-point) -> BigQuery `FLOAT`
- JSON Boolean -> BigQuery `BOOLEAN`
- JSON Object -> BigQuery `RECORD` (nested schema)
- JSON Array of primitive types -> BigQuery `REPEATED` primitive type
- JSON Array of objects -> BigQuery `REPEATED RECORD`

## 3. Non-Functional Requirements

### 3.1. Performance through Concurrency
- The implementation should leverage Go's concurrency features to process large NDJSON files efficiently.
- The scanning through the ndjson file should be parallelized. Example file are provided in the ndjson/benchmarks folder.
- The scanning (file i/o) and inference processing should communicate via a golang channel. Following a producer consumer pattern with multiple producers and multiple consumers.
- The resulting schema schema should be shared accross consumer goroutines and protected via a RWLock
- File I/O and JSON parsing should be optimized to minimize processing time.

### 3.2. Error Handling
- The system must provide clear and informative error messages, including the file and line number where an error occurred.
- In case of a fatal error (e.g., malformed JSON), the process should fail fast.

### 3.3. Reusability
- The core schema detection logic should be encapsulated in a reusable Go package called "core" that can be easily integrated into other applications or services.

## 4. Out of Scope
- API integration with Google Cloud services (Cloud Storage, BigQuery).
- Distributed or clustered deployment.
- Support for string-wrapped data types (e.g., `TIMESTAMP`, `DATE`).
- Deployment on Google Cloud Dataflow.
