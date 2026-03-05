# 🎓 UniSwap – Sprint 1 Report

Team Members:

* Samarth Vinayaka – Frontend Engineer
* Bhumi Jain – Frontend Engineer
* Shubhank Chandak – Backend Engineer
* Nandini Agrawal – Backend Engineer

Tech Stack: Go, PostgreSQL, TypeScript, JWT, REST APIs

---

# 1️⃣ User Stories

**US1 – Authentication**
As a student, I want to register and login using my credentials so that only authenticated campus users can access UniSwap.

**US2 – Protected Access**
As a logged‑in user, I want my session to persist so that I don’t need to login repeatedly.

**US3 – Create Listing**
As a student seller, I want to create a listing with title, description, price, category, and image so that I can sell items.

**US4 – Browse Listings**
As a student buyer, I want to view all listings so that I can browse items available on campus.

**US5 – Listing Details**
As a user, I want to see full listing details so that I can understand item information before contacting the seller.

**US6 – Upload Image**
As a seller, I want to upload an image with my listing so that buyers can see the item clearly.

**US7 – Basic Search**
As a user, I want to search listings by keyword or category so that I can find items faster.

**US8 – Report Listing**
As a user, I want to report inappropriate listings so that unsafe or irrelevant posts can be flagged.

---

# 2️⃣ Issues Planned for Sprint 1

## Frontend Issues

1. Setup React + TypeScript project
2. Implement App Routing
3. Login Page UI
4. Register Page UI
5. JWT Token Storage
6. Protected Route Component
7. API Service Layer
8. Listing Feed Page UI
9. Listing Detail Page UI
10. Create Listing Page UI
11. Image Upload Input Component
12. Connect Create Listing to Backend
13. Connect Listing Feed to Backend
14. Basic Search Bar UI
15. Report Listing Button UI

## Backend Issues

16. Setup Go Project Structure
17. PostgreSQL Setup + Connection
18. Create Users Table
19. Register Endpoint
20. Login Endpoint (JWT)
21. JWT Middleware
22. Create Listings Table
23. Create Listing Endpoint
24. Get All Listings Endpoint
25. Get Listing By ID Endpoint
26. Basic Image Upload Endpoint
27. Reports Table
28. Report Listing Endpoint
29. Basic Search Query Support
30. Postman Collection for Demo

---

# 3️⃣ Issues Successfully Completed

✅ 28 out of 30 planned issues were successfully completed.

The team implemented authentication, listing creation APIs, listing browsing APIs, search UI, reporting functionality, database schema, and REST APIs. Backend endpoints were fully tested using Postman and worked correctly.

The Sprint 1 demo successfully showed:

* User registration and login
* JWT authentication and protected routes
* Create listing flow (backend working)
* View listings feed and details (backend working)
* Image upload functionality
* Basic search functionality
* Report listing functionality

Two frontend-integration issues are still open and will be completed early in Sprint 2.

---

# 4️⃣ Issues Not Completed and Why

Two frontend integration issues are still open:

* [FE-13] Connect listing feed to backend API
* [FE-12] Connect create listing form to backend API

These issues remained open because frontend engineer Bhumi Jain was not able to complete the API integration within the sprint timeline. The backend endpoints were working correctly (verified using Postman), but final UI-to-API connection requires additional debugging and testing.

These tasks are planned as high-priority items at the start of Sprint 2.

---

# 5️⃣ Notes

The team followed Agile workflow with GitHub issues, PR reviews, and collaborative debugging. Sprint 1 established a strong foundation for Sprint 2 features such as wishlist, advanced search, and real‑time chat.

---

# 6️⃣ Team Collaboration & Scrum Process

The team conducted regular stand‑up (Scrum) meetings throughout Sprint 1. In these meetings, we:

* Discussed progress on current issues
* Identified blockers and resolved them collaboratively
* Planned next issues and task ownership
* Ensured frontend and backend integration stayed aligned

These regular check‑ins helped us stay on schedule, maintain clarity in responsibilities, and successfully complete all planned Sprint 1 tasks.