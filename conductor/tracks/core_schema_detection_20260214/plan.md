# Implementation Plan for Core Schema Detection

This plan outlines the tasks required to implement the core schema detection logic.

## Phase 1: Basic Schema Inference and Data Type Mapping

- [x] Task: Implement the basic file parsing and JSON object reading functionality. [4f8faab]
    - [x] Write tests for file reading and parsing.
    - [x] Implement file reading and parsing.
- [ ] Task: Implement initial data type mapping for primitive types (String, Integer, Float, Boolean).
    - [ ] Write tests for primitive type mapping.
    - [ ] Implement primitive type mapping.
- [ ] Task: Conductor - User Manual Verification 'Phase 1: Basic Schema Inference and Data Type Mapping' (Protocol in workflow.md)

## Phase 2: Handling Nested and Repeated Structures

- [ ] Task: Implement support for nested `RECORD` (JSON object) types.
    - [ ] Write tests for nested record inference.
    - [ ] Implement nested record inference.
- [ ] Task: Implement support for `REPEATED` primitive types.
    - [ ] Write tests for repeated primitive type inference.
    - [ ] Implement repeated primitive type inference.
- [ ] Task: Implement support for `REPEATED RECORD` types (arrays of objects).
    - [ ] Write tests for repeated record inference, ensuring the schema is merged from all objects in the array.
    - [ ] Implement repeated record inference.
- [ ] Task: Conductor - User Manual Verification 'Phase 2: Handling Nested and Repeated Structures' (Protocol in workflow.md)

## Phase 3: Advanced Rules and Error Handling

- [ ] Task: Implement the logic for `REQUIRED` vs. `NULLABLE` fields.
    - [ ] Write tests to verify the `REQUIRED`/`NULLABLE` logic.
    - [ ] Implement the `REQUIRED`/`NULLABLE` logic.
- [ ] Task: Implement handling of the JSON string literal "null".
    - [ ] Write tests for handling "null" values.
    - [ ] Implement the "null" value handling.
- [ ] Task: Implement robust error handling with clear, actionable error messages.
    - [ ] Write tests for various error conditions (e.g., malformed JSON).
    - [ ] Implement the error handling mechanisms.
- [ ] Task: Conductor - User Manual Verification 'Phase 3: Advanced Rules and Error Handling' (Protocol in workflow.md)
