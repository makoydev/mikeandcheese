# mikeandcheese

A minimalist blog site for exploring ideas about technology, AI, rationality, and interesting things.

## Tech Stack

- **Frontend**: React + TypeScript + Vite + Tailwind CSS v4
- **Backend**: Go + Chi router + SQLite

## Getting Started

### Backend

```bash
cd backend
go mod tidy
go run .
```

The API server starts on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

The dev server starts on `http://localhost:5173` and proxies API requests to the backend.

## Design

Inspired by [LessWrong](https://www.lesswrong.com/) and [Dario Amodei's blog](https://www.darioamodei.com/) — clean, serif-forward, editorial aesthetic with dark/light mode support.
