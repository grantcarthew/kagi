# p-001: Project Initialization

- Status: Pending
- Started: -

## Overview

This is the initialization project for a new repository using the Documentation Driven Development (DDD) template. The goal is to transition this generic template into a specific, active software project by gathering requirements from the user and populating the core documentation.

## Goals

1. Understand the user's intent for this repository.
2. Remove generic template placeholders from root documentation.
3. Define the first concrete work package (p-002).

## Scope

In Scope:

- Gathering project requirements from user
- Updating root documentation (README.md, AGENTS.md)
- Creating p-002 for the first milestone

Out of Scope:

- Actual implementation work (covered by p-002)
- Design decisions (covered during p-002)

## Success Criteria

- User has provided project name, description, and tech stack
- `README.md` accurately describes the new project
- `AGENTS.md` accurately describes the new project
- `.ai/projects/p-002-*.md` exists and is well-formed
- `AGENTS.md` points to p-002 as the Active Project

## Deliverables

- Updated `README.md`
- Updated `AGENTS.md`
- New `.ai/projects/p-002-[milestone-name].md`

## Agent Instructions

As the AI Agent working on this project, your goal is to "interview" the user and then set up the project. Follow these steps:

1. Context Gathering (The Interview):
    - Ask the user for the Project Name
    - Ask for a High-Level Description (What are we building?)
    - Ask about the Target Tech Stack (Languages, Frameworks)
    - Ask for the First Major Milestone (What is the first tangible thing to build?)

2. Documentation Updates:
    - Update `README.md`: Replace the template intro with the new Project Name and Description
    - Update `AGENTS.md`: Update the `[Project Name]` and Description placeholders

3. Project Planning:
    - Based on the First Major Milestone, create p-002
    - Use the standard project template (`.ai/projects/p-writing-guide.md`)
    - Ensure p-002 has clear goals, success criteria, and deliverables

4. Handover:
    - Move this file (p-001) to `.ai/projects/completed/`
    - Update `AGENTS.md` to set `Active Project: .ai/projects/p-002-....md`
    - Inform the user that the project is initialized and ready for p-002
