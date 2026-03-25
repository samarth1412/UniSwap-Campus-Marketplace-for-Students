
# 🎓 UniSwap – Sprint 2 Plan

## Team Members
- **Samarth Vinayaka** – Frontend Engineer  
- **Bhumi Jain** – Frontend Engineer  
- **Shubhank Chandak** – Backend Engineer  
- **Nandini Agrawal** – Backend Engineer  

Tech Stack: **Go, PostgreSQL, React, TypeScript, JWT, REST APIs**

---

# 🚀 Sprint 2 Goals

Sprint 2 focuses on expanding the functionality of UniSwap by improving user interaction, connecting the frontend with backend APIs, and adding core marketplace features.

The main goals for this sprint are:

1. Complete frontend-backend integration for listings
2. Implement wishlist functionality
3. Add advanced search and filtering
4. Implement pagination for listings
5. Allow users to edit and delete their listings
6. Implement a user profile section with “My Listings”
7. Improve backend validation, authorization, and API documentation

These improvements will move UniSwap closer to a fully functional campus marketplace.

---

# ⚠️ Carryover Issues from Sprint 1

Two frontend issues were not completed in Sprint 1 and will be prioritized at the beginning of Sprint 2:

- **[FE-12] Connect create listing form to backend API**
- **[FE-13] Connect listing feed to backend API**

---

# 👨‍💻 Samarth Vinayaka – Frontend Tasks

### FE-16 – Improve API Service Layer
Create a centralized API service module that handles:
- Listings
- Wishlist
- Profile APIs
- Automatic JWT token handling

### FE-17 – Connect Listing Detail Page to API
Integrate frontend with:
GET /listings/:id

### FE-18 – Edit Listing UI
Create the frontend interface for editing a listing with fields:
- Title
- Description
- Price
- Category

### FE-19 – Connect Edit Listing API
Integrate edit form with:
PUT /listings/:id

### FE-20 – Delete Listing Button
Add a delete button to listing pages with confirmation modal.

### FE-21 – Connect Delete Listing API
Connect delete button to:
DELETE /listings/:id

### FE-22 – My Listings Page
Create page showing listings posted by the logged-in user.

Route:
/my-listings

---

# 👩‍💻 Bhumi Jain – Frontend Tasks

### FE-12 – Connect Create Listing Form to Backend API (Carryover)
Connect create listing form to:
POST /listings

### FE-13 – Connect Listing Feed to Backend API (Carryover)
Load listings dynamically from:
GET /listings

### FE-23 – Wishlist Button
Add heart/save icon to listing cards.

### FE-24 – Wishlist Page
Create page displaying user's saved listings.

Route:
/wishlist

### FE-25 – Connect Wishlist APIs
Integrate frontend with:
- POST /wishlist
- DELETE /wishlist/:id
- GET /wishlist

### FE-26 – Advanced Search UI
Add filters for:
- Category
- Price range
- Keyword search

### FE-27 – Pagination UI
Add:
- Next / Previous buttons
- Page indicator

---

# 👨‍💻 Shubhank Chandak – Backend Tasks

### BE-16 – Add Pagination to Listings API
Update endpoint:
GET /listings?page=1&limit=10

### BE-17 – Add Advanced Search Filters
Allow filtering by:
- category
- min_price
- max_price
- keyword

### BE-18 – Create Wishlist Table

Schema:

wishlist
id
user_id
listing_id
created_at

### BE-19 – Add Wishlist Endpoint
POST /wishlist

Adds listing to user's wishlist.

### BE-20 – Remove Wishlist Endpoint
DELETE /wishlist/:id

### BE-21 – Get Wishlist Endpoint
GET /wishlist

Returns all saved listings for a user.

### BE-22 – Get Listings by User
GET /users/:id/listings

Used for "My Listings" page.

---

# 👩‍💻 Nandini Agrawal – Backend Tasks

### BE-23 – Update Listing Endpoint
PUT /listings/:id

Allows editing listing details.

### BE-24 – Delete Listing Endpoint
DELETE /listings/:id

Removes listing from database.

### BE-25 – Authorization Middleware
Ensure users can only:
- Edit their own listings
- Delete their own listings

### BE-26 – Improve Image Upload
Allow multiple images per listing.

### BE-27 – Create Listing Images Table

listing_images
id
listing_id
image_url

### BE-28 – Standardize Error Handling
Return consistent API responses across endpoints.

### BE-29 – API Documentation
Create README documentation including:
- Endpoint descriptions
- Request bodies
- Example responses

---

# 📊 Sprint 2 Work Distribution

| Team Member | Role | Number of Tasks |
|--------------|------|----------------|
| Samarth | Frontend | 7 |
| Bhumi | Frontend | 7 |
| Shubhank | Backend | 7 |
| Nandini | Backend | 7 |

Total Tasks: **28**

---

# 🎯 Expected Outcome After Sprint 2

By the end of Sprint 2, UniSwap should support:

- Creating listings
- Browsing listings
- Editing and deleting listings
- Advanced search and filtering
- Wishlist functionality
- Viewing personal listings
- Pagination for large listing sets

This will significantly improve the usability and completeness of the UniSwap platform.

---
