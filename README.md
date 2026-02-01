# 🎓 UniSwap – Campus Marketplace for Students

> **Buy, Sell, and Trade: Books, Gadgets, Furniture & More.**  
> A secure, campus-exclusive marketplace designed to connect students.

---

## 📖 Project Overview
**UniSwap** is a full-stack web application tailored for university campuses. It addresses the need for a trusted platform where students can easily trade items without the friction and safety concerns of open marketplaces. By restricting access to the campus community and offering real-time communication, UniSwap ensures a safer and more relevant trading experience.

## ✨ Key Features

### 🛍️ Marketplace Core
- **Rich Listings**: Create detailed product listings with multiple images, descriptions, and condition tags.
- **Advanced Search & Filtering**: Find exactly what you need by category (Books, Electronics, Furniture), price range, or condition.
- **Wishlist**: Save items for later and get notified of price drops.

### 💬 Social & Interaction
- **Real-time Chat**: Integrated WebSocket-based chat allows buyers and sellers to negotiate instantly without sharing personal phone numbers.
- **User Profiles**: View seller ratings, history, and verified student status.

### 🛡️ Administration
- **Moderation Panel**: Admin dashboard to review flagged listings, ban users, and manage content quality.
- **Safety**: Image content moderation and university email verification (planned).

---

## 🛠️ Technology Stack

### Backend Engineering (Go)
- **Language**: Go (Golang)
- **API Style**: RESTful API
- **Database**: PostgreSQL (Relational Data)
- **Real-time**: WebSockets (Gorilla/Melody) for Chat
- **Storage**: S3-compatible object storage (or local file system) for image uploads.

### Frontend Engineering (React)
- **Framework**: React + TypeScript
- **Styling**: TailwindCSS / Vanilla CSS
- **State Management**: React Context / Zustand
- **Build Tool**: Vite

---

## 📅 Sprint Roadmap (8 Weeks)

This project is executed in 4 strategic sprints:

| Sprint | Focus Area | Key Deliverables |
| :--- | :--- | :--- |
| **Sprint 1** | **Foundation & Auth** | Repo setup, DB Schema, User Authentication (JWT), Basic UI Layout. |
| **Sprint 2** | **Marketplace Core** | Listing CRUD, Image Uploads, Search & Filtering Logic. |
| **Sprint 3** | **Interaction** | Real-time Chat (WebSockets), Wishlist, User Profiles. |
| **Sprint 4** | **Admin & Polish** | Moderation Dashboard, End-to-End Testing, UI Refinement. |

---

## 👥 Team Members

### Frontend Engineers
- **[Name]** - Role: UI/UX, Component Architecture
- **[Name]** - Role: State Management, API Integration

### Backend Engineers
- **[Name]** - Role: API Design, Database Architecture
- **[Name]** - Role: Real-time Services, Infrastructure

---

## 🚀 Getting Started

### Prerequisites
- Node.js (v18+)
- Go (v1.21+)
- PostgreSQL

### Installation
*Instructions to be added in Sprint 1...*
