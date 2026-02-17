# Sprint 1 — UniSwap Campus Marketplace

## 1) User Stories

| ID | User Story |
|----|------------|
| US1 | As a student, I want to register/login so that only authenticated users can use UniSwap. |
| US2 | As a logged-in user, I want to stay logged in (token persistence) so that I don't need to sign in repeatedly. |
| US3 | As a user, I want to create a listing with title/price/category/description so that I can sell items. |
| US4 | As a user, I want to view a feed of listings so that I can browse items. |
| US5 | As a user, I want to view listing details so that I can see full description and contact seller. |
| US6 | As a user, I want to upload at least one image for my listing so that buyers can see it clearly. |
| US7 | As a user, I want to search/filter listings by keyword/category so that I can find items faster. |
| US8 | As a user, I want to report a listing so that inappropriate posts can be flagged. |

---

## 2) Issues Planned for Sprint 1

### Frontend

| Issue | Title | Owner |
|-------|--------|--------|
| FE-1 | Setup React + TypeScript project | Samarth |
| FE-2 | Implement App Routing | Samarth |
| FE-3 | Login Page UI | Samarth |
| FE-4 | Register Page UI | Samarth |
| FE-5 | JWT Token Storage | Samarth |
| FE-6 | Protected Route Component | Samarth |
| FE-7 | API Service Layer | Samarth |
| FE-8 | Listing Feed Page UI (Mocked Data Initially) | Bhumi |
| FE-9 | Listing Detail Page UI | Bhumi |
| FE-10 | Create Listing Page UI | Bhumi |
| FE-11 | Image Upload Input Component | Bhumi |
| FE-12 | Connect Create Listing to Backend | Bhumi |
| FE-13 | Connect Listing Feed to Backend | Bhumi |
| FE-14 | Basic Search Bar UI | Bhumi |
| FE-15 | Report Listing Button UI | Bhumi |

### Backend

| Issue | Title | Owner |
|-------|--------|--------|
| BE-1 | Setup Go Project Structure | Shubhank |
| BE-2 | PostgreSQL Setup + Connection | Shubhank |
| BE-3 | Create Users Table | Shubhank |
| BE-4 | Register Endpoint | Shubhank |
| BE-5 | Login Endpoint (JWT) | Shubhank |
| BE-6 | JWT Middleware | Shubhank |
| BE-7 | Create Listings Table | Nandini |
| BE-8 | Create Listing Endpoint | Nandini |
| BE-9 | Get All Listings Endpoint | Nandini |
| BE-10 | Get Listing By ID Endpoint | Nandini |
| BE-11 | Basic Image Upload Endpoint (Local Storage) | Nandini |
| BE-12 | Reports Table | Nandini |
| BE-13 | Report Listing Endpoint | Nandini |
| BE-14 | Basic Search Query Support | Nandini |
| BE-15 | Postman Collection for Demo | Nandini |

---

## 3) Issues Completed

*(Update as you close issues. Use ✅ and the issue number.)*

### Frontend

- ✅ FE-8 — Listing Feed Page UI (mocked & connected to API)
- ✅ FE-9 — Listing Detail Page UI
- ✅ FE-10 — Create Listing Page UI with image upload
- ✅ FE-11 — Image Upload Input Component
- ✅ FE-12 — Connect Create Listing to Backend
- ✅ FE-13 — Connect Listing Feed to Backend
- ✅ FE-14 — Basic Search Bar UI
- ✅ FE-15 — Report Listing Button UI

### Backend

- (Add when done, e.g. ✅ BE-1 — Setup Go Project Structure)

---

## 4) Issues Not Completed + Why

*(Fill only for issues that were planned but not completed.)*

| Issue | Reason (blocked / underestimated / scope cut / etc.) |
|-------|------------------------------------------------------|
| (e.g. FE-15) | (e.g. Ran out of time; will pick in Sprint 2) |

---

## 5) Demo Notes (optional)

**Frontend demo flow:**  
Login → View listing feed → Open listing detail → Create listing (with image) → Search → Report listing.

**Backend demo flow (Postman/curl):**  
Register → Login (get JWT) → Create listing (with auth header) → Get listings → Report listing.

**Known limitations:**  
*(e.g. Image upload uses local storage only; search is frontend-only until BE-14 is done.)*
