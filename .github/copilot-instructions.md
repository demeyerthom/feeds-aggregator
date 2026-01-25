# COPILOT EDITS OPERATIONAL GUIDELINES

## PRIME DIRECTIVE

	Avoid working on more than one file at a time.
	Multiple simultaneous edits to a file will cause corruption.
	Be chatty and teach about what you are doing while coding.

## LARGE FILE & COMPLEX CHANGE PROTOCOL

### MANDATORY PLANNING PHASE

	When working with large files (>300 lines) or complex changes:
		1. ALWAYS start by creating a detailed plan BEFORE making any edits
            2. Your plan MUST include:
                   - All functions/sections that need modification
                   - The order in which changes should be applied
                   - Dependencies between changes
                   - Estimated number of separate edits required
                
            3. Format your plan as:

## PROPOSED EDIT PLAN

	Working with: [filename]
	Total planned edits: [number]

### MAKING EDITS

	- Focus on one conceptual change at a time
	- Show clear "before" and "after" snippets when proposing changes
	- Include concise explanations of what changed and why
	- Always check if the edit maintains the project's coding style

### Edit sequence:

	1. [First specific change] - Purpose: [why]
	2. [Second specific change] - Purpose: [why]
	3. Do you approve this plan? I'll proceed with Edit [number] after your confirmation.
	4. WAIT for explicit user confirmation before making ANY edits when user ok edit [number]

### EXECUTION PHASE

	- After each individual edit, clearly indicate progress:
		"✅ Completed edit [#] of [total]. Ready for next edit?"
	- If you discover additional needed changes during editing:
	- STOP and update the plan
	- Get approval before continuing

### REFACTORING GUIDANCE

	When refactoring large files:
	- Break work into logical, independently functional chunks
	- Ensure each intermediate state maintains functionality
	- Consider temporary duplication as a valid interim step
	- Always indicate the refactoring pattern being applied

### RATE LIMIT AVOIDANCE

	- For very large files, suggest splitting changes across multiple sessions
	- Prioritize changes that are logically complete units
	- Always provide clear stopping points

## General Requirements

	Use modern technologies as described below for all code suggestions. Prioritize clean, maintainable code with appropriate comments.

## Go Requirements

	- **Target Version**: Go 1.25 or higher
	- **Features to Use**:
	- Generics for type-safe data structures and functions
    - Write tests using the `testing` package and `testify` for assertions for every function.
    - Use Go modules for dependency management.
    - Build binaries with `go build` and manage dependencies with `go mod`. Binaries should be placed in the `bin/` directory.

## Folder Structure

	Follow this structured directory layout:

		feeds-aggregator/
		├── .github/                            # GitHub-specific files (workflows, issue templates)
		├── bin/                                # Compiled binaries and executables
        ├── cmd/                                # Main applications for this project
        │   └── {application-name}/main.go      # Entry point per application
        ├── configs/                            # Configuration files and templates
        │   └── {application-name}/*            # Config files for each application
        ├── docker /                            # Docker-related files
        │   └── {application-name}/Dockerfile   # Dockerfile for each application
        ├── internal/                           # Private application and library code

## Documentation Requirements

	- Include comments for all functions and complex code sections.
	- Document complex functions with clear examples.
	- Maintain concise Markdown documentation.
	- Minimum docblock info: `param`, `return`, `throws`, `author`

## Security Considerations

	- Sanitize all user inputs thoroughly.
	- Parameterize database queries.
	- Enforce strong Content Security Policies (CSP).
	- Use CSRF protection where applicable.
	- Ensure secure cookies (`HttpOnly`, `Secure`, `SameSite=Strict`).
	- Limit privileges and enforce role-based access control.
	- Implement detailed internal logging and monitoring.