# ⚔️ Server-Authoritative Idle RPG

A modern, high-performance web-based Idle RPG built with a **Server-Authoritative** architecture. This project demonstrates complex state synchronization between a **Golang Game Engine** and a **Next.js Frontend** using **gRPC**.

## 🏗️ Architecture Overview

The system is split into three main tiers to ensure security, scalability, and smooth user experience:

1.  **Backend (Golang Engine):** The "Source of Truth". Handles all combat math, loot generation, and offline progression simulation.
2.  **BFF (Next.js Route Handlers):** Acts as a bridge, forwarding client requests to the gRPC server and formatting data for the UI.
3.  **Frontend (Next.js + Zustand):** A high-performance UI running a 60FPS game loop for visual feedback, reconciled by periodic server syncs.

### Tech Stack
- **Frontend:** Next.js 14, TypeScript, Zustand, TailwindCSS.
- **Backend:** Golang, gRPC, Protobuf.
- **Database:** PostgreSQL (via Prisma).
- **Communication:** gRPC (Binary protocol for high-speed sync).

## 🎮 Core Mechanics

### 1. Server-Authoritative Combat
All damage calculation happens on the Go server.
- **Formula:** `Final Damage = max(1, (Base Attack - Target Defense) * Multiplier)`
- **Elemental Triangle:** Fire > Wood > Water > Fire (1.5x Advantage / 0.8x Disadvantage).

### 2. Offline Progression (Time-to-Kill)
The server calculates progress based on the time elapsed since the last sync.
- It determines **TTK (Time to Kill)** based on player/monster stats.
- Simulates total kills and generates tiered loot automatically.

### 3. Smart Inventory & Auto-Salvage
To prevent database bloat, the server implements an auto-salvage queue:
- **High-Tier Items & Gems:** Prioritized and kept in inventory.
- **Low-Tier Items:** Automatically converted to Gold if the inventory is full.

### 4. Hybrid Sync Loop
- **Client Side:** Runs a 60FPS loop for smooth HP bars and animations.
- **Server Sync:** Every 60 seconds, the client sends its state to the server. The server calculates the "true" state and **forcefully overwrites** the client's local state to prevent cheating.

## 🚀 Getting Started

### Prerequisites
- [Go](https://golang.org/doc/install) (1.20+)
- [Node.js](https://nodejs.org/) (18+)
- [PostgreSQL](https://www.postgresql.org/)

### 1. Backend Setup
```bash
# Install Go dependencies
go mod tidy

# Run the gRPC Game Server
go run main.go
```

### 2. Frontend Setup
```bash
# Install Node dependencies
npm install

# Run the Next.js dev server
npm run dev
```

## 📁 Project Structure

```text
├── engine/             # Golang: Core Game Logic (Combat, Loot, Progression)
├── proto/              # Protobuf definitions (gRPC contracts)
├── store/              # Zustand: Client-side state management
├── hooks/              # Game Loop (requestAnimationFrame) logic
├── app/api/sync/       # Next.js BFF: gRPC Client bridge
├── components/         # React UI Components
└── schema.prisma       # Database Schema
```

## 🛡️ Security Features
- **Cheat Prevention:** Client cannot "give" themselves gold or items; everything is verified by the Go engine.
- **Resource Efficiency:** Heavy calculations are handled by Go's efficient concurrency model.
- **No-Database Bloat:** Low-value drops are salvaged before they ever hit the database.

---
Created with ❤️ by Gemini CLI Senior Architect.
