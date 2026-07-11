# Custom Command Configuration: /ticket

## Usage Syntax
/ticket <github-issue-url-or-number>

## Execution Workflow
When this command is invoked, the AI must strictly execute the following pipeline:

1. Data Extraction (GitHub CLI):
   - Parse the issue number/URL provided by the user.
   - Run the following terminal command to extract *every single detail* (body, comments, labels, milestones, assignees):
     `gh issue view <issue-id-or-url> --comments`
   - Read and digest the entire command output. Do not truncate the comments; parse all discussion threads for crucial constraints, shifting requirements, or edge cases.

2. Codebase Discovery:
   - Perform a structural scan of the local repository.
   - Map the extracted issue requirements to specific files, core logic components, or configurations in the current project.

3. Deliverable Generation:
   Produce a comprehensive, highly technical **Implementation Plan** structured exactly as follows:

   ## 🎯 Complete Issue Breakdown
   - **Summary:** [2-3 sentence technical summary of the objective]
   - **Metadata:** [Labels, Milestones, or Assignees if relevant]
   - **Discovered Constraints:** [Crucial details or direction changes extracted from the comment history]

   ## 🔍 Key Requirements & Edge Cases
   - [ ] **Requirement 1:** [Core item from description]
   - [ ] **Edge Case:** [Discovered system limitation or edge case]

   ## 🗺️ Codebase Impact Map
   - `path/to/affected/file.ext`: [Line-by-line or conceptual change required]

   ## 🛠️ Step-by-Step Implementation Blueprint
   1. **Phase 1: [Name]** - Detailed architectural changes or scaffolding.
   2. **Phase 2: [Name]** - Logic implementation and integration.

   ## 🧪 Verification & Testing Strategy
   - [Specific automated test cases or manual verification steps required to close this ticket]