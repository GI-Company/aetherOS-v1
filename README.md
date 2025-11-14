# Aether: A Browser-Native OS with an AI Core

Aether is an experimental project to build a lightweight, browser-native operating system. It uses a Go-based microservices architecture to provide a runtime for applications, with a dedicated AI core powered by Google Gemini.

## Key Features

*   **AI-Powered Core:** Integrates with Google's Gemini models to provide powerful generative AI capabilities to the entire OS.
*   **Browser-Native OS:** Runs entirely in the browser, providing a lightweight and portable operating system experience.
*   **Microservices Architecture:** Built on a modular, Go-based backend with services for AI, authentication, and file management.
*   **Real-Time Communication:** Uses WebSockets to enable real-time, bidirectional communication between the frontend and the backend.
*   **User Authentication:** A complete authentication system that allows users to sign in and personalize their Aether experience.
*   **Virtual File System (VFS):** A flexible VFS that proxies file operations to the frontend, which is responsible for persistence (e.g., using IndexedDB).

## Architecture

The Aether ecosystem is composed of two main components working in concert:

### 1. The Frontend: Aether OS Shell

*   **Location:** `/frontend`
*   **Technology:** React.js
*   **Description:** This is the user-facing interface for Aether. It acts as the "desktop environment" for the browser-based OS, allowing users to launch and interact with Aether applications and services. It communicates with the backend kernel via a WebSocket connection.

### 2. The Backend: Aether Kernel

*   **Location:** `/` (Root Project)
*   **Technology:** Go
*   **Description:** This is the heart of the operating system. It's a Go-based application that runs as a single server. It is built on a microservice-style architecture where different kernel services (like VFS and AI) communicate over a central, in-process message bus.

## Getting Started

To run the full Aether ecosystem, you need to start the backend kernel and the frontend development server.

### Prerequisites

*   Go (1.23 or later)
*   Node.js and npm
*   A valid `GEMINI_API_KEY` set as an environment variable.

### 1. Run the Aether Kernel

In your terminal, from the root directory of the project, run the following command:

```bash
go run main.go
```

The kernel will start and will be available on the port configured in `config.yaml` (default: `8080`). It handles both HTTP requests and WebSocket connections.

### 2. Run the Frontend

In a separate terminal, navigate to the `frontend` directory, install dependencies, and start the Vite development server:

```bash
cd frontend
npm install
npm run dev
```

The Aether OS shell will now be accessible in your browser, typically at `http://localhost:5173`.

## Pricing

This project includes both a free, open-source core and a commercial edition with advanced features.

### GAiNOS / AetherOS Pricing Overview

| Tier                    | Price       | Users / Devices  | Key Features                                                                                                    |
| ----------------------- | ----------- | ---------------- | --------------------------------------------------------------------------------------------------------------- |
| **Free Tier**           | $0          | 1 device         | Open-source core OS, default apps, read-only SDK, Apache 2.0 license                                            |
| **Developer Tier**      | $199 / year | 1–3 devices      | Full SDK access, Marketplace app submission, experimental AI modules                                            |
| **Small Business Tier** | $899 / year | Up to 20 devices | All Developer features + multi-device support, API access, standard support                                     |
| **Enterprise Tier**     | Custom      | 20+ devices      | All Small Business features + bulk deployment, priority updates, premium enterprise support, kernel integration |

### Add-On Modules (Commercial Edition Only)

| Feature               | Price                  | Description                          |
| --------------------- | ---------------------- | ------------------------------------ |
| Advanced AI Analytics | $199 / year / device   | Enhanced AI Core functionality       |
| Developer Pro Toolkit | $99 / year / developer | Advanced debugging & emulation tools |
| Marketplace Promotion | $299 / app / year      | Featured app listing in Marketplace  |

### Legend / Notes

* **Free Tier**: v1 repo remains open-source; Apache 2.0 licensed.
* **Paid tiers**: GAiNOS v2 / commercial edition, governed by Proprietary License & EULA.
* **Revenue Split**: Paid apps: 70% developer / 30% Global Intent Company.
* **Enterprise**: Custom pricing based on deployment size and support requirements.

## Contributing

Aether is an open and experimental project. Contributions are welcome! If you'd like to get involved, please check out the open issues on GitHub. Feel free to fork the repository, make your changes, and submit a pull request.

## Legal

This project is governed by a set of legal documents, policies, and agreements. All legal information, including the proprietary license, SDK license, EULA, and other related documents, can be found in the [`LEGAL`](./LEGAL) directory. Please review these documents to understand your rights and obligations when using or contributing to this project.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
