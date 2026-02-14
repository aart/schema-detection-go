# Specification for Concurrency Refactoring

## 1. Overview
This document outlines the technical specifications for refactoring the schema detection implementation to improve performance by introducing concurrency.

## 2. Functional Requirements

### 2.1. Parallel File Processing
- The file reading logic should be refactored to process multiple files in parallel.
- The system should use a pool of workers to read and parse lines from the input NDJSON files concurrently.
- A channel should be used to send the parsed JSON objects from the workers to the schema inference logic.

## 3. Non-Functional Requirements

### 3.1. Performance
- The schema detection process should be significantly faster for large NDJSON files.
- The degree of parallelism (e.g., the number of workers) should be configurable.

## 4. Acceptance Criteria
- The refactored implementation passes all existing unit tests.
- Benchmark tests show a significant performance improvement (e.g., at least 2x faster) for large files.

## 5. Out of Scope
- Refactoring of the schema inference logic itself.
- Changes to the public API of the `core` package.
