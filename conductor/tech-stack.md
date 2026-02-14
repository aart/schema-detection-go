# Tech Stack

## Overview
This project is built using go, leveraging its performance and concurrency features.

## Core Technologies

### Programming Language
-   **Go**: A programming language focused on simplicity, performance, and concurrency.

### Key Libraries and Frameworks
-   **github.com/spf13/cobra** for implementing commands for the cli
-   **cloud.google.com/go/bigquery** import and use the BigQuery schema and field types
-	**encoding/json** for parsing the json in the ndjson files
-	**sync"** for protecting the shared data structures from concurrent access

## Architecture

### Key Design Principles:
- Performance through concurrency (using golang channel primitives)
- Big data support through splitted input files (processed in parallel & sampling)
- Enforce consistency checks upstream (avoid things fail downstream)
- Fail fast in case of an error
- File position traceback to enable debugging
- Architect for reusability (generate other destination schema for AlloyDB, Spanner, ...)
