# 🎓 UniSwap – Sprint 3 Plan

## Team Members
- **Samarth Vinayaka** – Frontend Engineer  
- **Bhumi Jain** – Frontend Engineer  
- **Shubhank Chandak** – Backend Engineer  
- **Nandini Agrawal** – Backend Engineer  

Tech Stack: **Go, PostgreSQL, React, TypeScript, JWT, REST APIs**

---

# 🚀 Sprint 3 Goals

Sprint 3 focuses on **integration completion and stabilization**. The goal is to ensure all previously built features work end‑to‑end with real backend integration and eliminate breakages before moving into testing-heavy Sprint 4.

The main goals for this sprint are:

1. Run baseline verification (frontend + backend)
2. Fix build/test failures
3. Complete listing filters + pagination support
4. Align frontend completely with backend APIs
5. Implement end-to-end image upload flow
6. Ensure all core flows work seamlessly
7. Added all the Unit tests in place

---

# ⚠️ Context from Previous Sprints

From Sprint 1 and Sprint 2 :
- Core features (auth, listings, wishlist, edit/delete) are implemented
- Advanced search + pagination partially implemented
- Some frontend-backend mismatches and mock fallbacks remain
- Image upload flow not fully wired end-to-end

Sprint 3 focuses on completing and stabilizing these areas.

---

# 👨‍💻 Shubhank Chandak – Backend Tasks

### BE-30 – Fix Backend Test & Build Failures
- Resolve interface mismatches (e.g., GetByUserID changes)
- Ensure Go test suite compiles successfully
- Increase overall test coverage

### BE-31 – Complete Listing Query Support
Enhance:
GET /listings

Add full support for:
- category
- min_price
- max_price
- keyword
- page
- limit

### BE-32 – Validate API Response Consistency
- Standardize response format
- Ensure proper error handling

---

# 👩‍💻 Nandini Agrawal – Backend Tasks

### BE-33 – Core Listings Experience (Create/Edit/Manage/Browse)
- Complete end-to-end listings flow with proper loading/error/empty states and role-based actions.

### BE-34 – Frontend UI/UX Standardization and Route Protection
- Make styling consistent/responsive/accessibility-friendly and ensure protected pages redirect/authorize correctly

### BE-35 – Testing Consolidation (Unit + E2E)
- Establish complete test coverage strategy for critical user flows.

### BE-36 - Platform Stability and Configuration Hardening
- Harden env/config handling, stabilize dependencies/scripts, and align local + CI behavior.


---

# 👨‍💻 Samarth Vinayaka – Frontend Tasks

### FE-28 – Align Listing APIs with Backend
- Update API calls for filters & pagination
- Match request/response formats

### FE-29 – Remove Mock Data
- Replace dummy data with real API data
- Add loading + error states

### FE-30 – Search & Filter Integration
- Connect UI filters to backend query params

---

# 👩‍💻 Bhumi Jain – Frontend Tasks

### FE-31 – Image Upload UI
- Create listing flow
- Edit listing flow
- Add image preview support

### FE-32 – UI Fixes & Flow Corrections
- Fix listing display issues
- Handle empty/error states

### FE-33 – End-to-End Flow Validation
- Create → View → Edit → Delete listing
- Verify image rendering

---

# 📊 Sprint 3 Work Distribution

| Team Member | Role | Number of Tasks |
|--------------|------|----------------|
| Samarth | Frontend | 3 |
| Bhumi | Frontend | 3 |
| Shubhank | Backend | 3 |
| Nandini | Backend | 3 |

Total Tasks: **12**

---

# 🎯 Project status After Sprint 3

We have achieved

- Have fully working frontend-backend integration
- Support filtering and pagination correctly
- Have complete image upload functionality
- Have no dependency on mock data
- Successfully run core flows end-to-end

---

# 📌 Notes

- Buyer seller interaction mode feature for Sprint 4
- TBD : Payment Gateway integration
- Focus is on fixing, completing, and integrating
- Prepares system for testing-heavy Sprint 4

---

