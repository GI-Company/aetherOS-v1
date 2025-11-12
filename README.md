# Aether: A Browser-Native OS with an AI Core

Aether is an experimental project to build a lightweight, browser-native operating system. It uses a Go-based microservices architecture to provide a runtime for WebAssembly (Wasm) applications, with a dedicated AI core powered by Google Gemini.

## Architecture

The Aether ecosystem is composed of three main components working in concert:

### 1. The Frontend: Aether OS Shell

*   **Location:** `/frontend`
*   **Technology:** React.js
*   **Description:** This is the user-facing interface for Aether. It acts as the "desktop environment" for the browser-based OS, allowing users to launch and interact with Aether applications.

### 2. The Backend: Aether Kernel

*   **Location:** `/backend`
*   **Description:** This is the heart of the operating system. It is a Go-based server that runs Wasm modules, providing them with a secure sandbox and a set of "system calls" for interacting with the outside world.
*   **Features:**
    *   **Wasm Runtime:** Uses `wazero` to execute Wasm applications.
    *   **System Services:** Provides APIs for authentication, key-value storage, and service lifecycle management.
    *   **IPC:** Uses a WebSocket-based message bus for real-time communication between the frontend and running Wasm modules.

### 3. The AI Service: Aether AI Core

*   **Location:** `/` (Root)
*   **Description:** A specialized microservice that functions as the dedicated "AI Processor" for the entire Aether ecosystem.
*   **Features:**
    *   **Gemini Integration:** Exposes a simple HTTP endpoint (`/v1/ai/generate`) to provide access to Google's Gemini models.
    *   **Decoupled Design:** By running as a separate service, the AI Core can be scaled and updated independently from the main Kernel.

## Getting Started

To run the full Aether ecosystem, you need to start the two main backend services and the frontend development server.

### Prerequisites

*   Go (1.23 or later)
*   Node.js and npm
*   A valid `GEMINI_API_KEY` set as an environment variable.

### 1. Run the AI Core

In the root directory, start the AI service:

```bash
go run main.go
```
It will be available at `http://localhost:8080`.

### 2. Run the Aether Kernel

In a separate terminal, navigate to the `backend` directory and start the Kernel:

```bash
cd backend
go run main.go
```
This server will run on port `8080` by default and also serve the frontend. To avoid conflicts with the AI Core, you may need to configure it to use a different port.

### 3. Run the Frontend

In a third terminal, navigate to the `frontend` directory, install dependencies, and start the Vite development server:

```bash
cd frontend
npm install
npm run dev
```

## Usage

You can interact with the AI Core directly using `curl`:

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"prompt": "Tell me a story about a brave robot."}' \
  http://localhost:8080/v1/ai/generate
```
