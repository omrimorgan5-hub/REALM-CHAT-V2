# Realm-Chat V2

A real-time chat application built with a high-performance Go backend and a modern React + TypeScript (TSX) frontend.

## 🛠 Tech Stack
- Backend: Go (Golang)
- Frontend: React, TypeScript, Vite
- Tooling: Portable Node.js & Go (No-Admin Environment)

---

## 🚀 Getting Started

### 1. Environment Setup
Since this project uses portable binaries in a restricted environment, you must add the tools to your path for the current terminal session. 

**Windows (PowerShell):**
$env:Path = "C:\Users\<USER>\LocalDev\node-v24.13.0-win-x64;C:\Users\<USER>\LocalDev\go\bin;" + $env:Path

### 2. Backend Setup (Go)
The backend handles WebSocket connections and API requests.
```bash
cd backend
go mod tidy
go run main.go
```

*Server runs on: http://localhost:8080*

### 3. Frontend Setup (React/TSX)
The frontend uses Vite with a built-in proxy to communicate with the Go server.

```bash
cd frontend
npm install
npm run dev
```

*UI runs on: http://localhost:5173*

---

## 📁 Project Structure
- /backend: Go source code, handlers, and modules.
- /frontend: React components (.tsx), styles, and Vite config.
- vite.config.ts: Set up with a /api proxy to prevent CORS errors during development.

## 🏗 Planned Features
- [ ] Real-time messaging via WebSockets.
- [ ] User session management.
- [ ] Chat history persistence.
- [ ] Responsive UI for mobile and desktop.

---
*Created as part of the Realm-Chat evolution.*