# World Cup Scoreboard

A simple in-memory Go library for managing live football matches.

## Features

- Start a match with a 0–0 score
- Update match scores
- Finish and remove matches
- Get a summary sorted by total score (desc), then by most recently added match

## Run Tests

```bash
go test ./...
go test -race ./...
```

## Design Decisions

- **Map** for fast match lookup and updates.
- **RWMutex** for safe concurrent access.
- **Sequence number** to track match creation order.
- No database or external dependencies.

## Assumptions

- A team can only participate in one active match.
- Team names are compared case-insensitively.
- Scores cannot be negative.
- Updating a score does not change match order.
