# UniSwap - Sprint 4

## Team Members

| Name | Role |
| --- | --- |
| Samarth Vinayaka | Frontend Engineer |
| Bhumi Jain | Frontend Engineer |
| Shubhank Chandak | Backend Engineer |
| Nandini Agrawal | Backend Engineer |

Tech stack: Go, PostgreSQL, React, TypeScript, JWT, REST APIs, Vitest, Cypress.

## Sprint 4 Goals

Sprint 4 focused on finishing Sprint 3 carry-over work, adding buyer-to-seller contact functionality, strengthening automated tests, and preparing the project for final demonstration.

Primary goals:

1. Complete unimplemented or unstable Sprint 3 functionality.
2. Add buyer/seller contact request flow.
3. Add unit tests for new backend and frontend behavior.
4. Document application setup, usage, test commands, and API behavior.
5. Prepare final video material showing the finished project, new functionality, and test results.

## Work Completed in Sprint 4

### Buyer and Seller Contact Requests

- Added backend support for buyers to contact sellers through a listing.
- Added validation so users cannot contact themselves about their own listings.
- Added persistent contact request storage with buyer, seller, listing, message, status, and timestamp fields.
- Added seller-facing endpoint to view received contact requests.
- Added frontend contact seller form on the listing detail page.
- Added frontend contact request inbox page at `/messages`.
- Added loading, empty, error, retry, and success states for contact request UI.

### Listing Flow Improvements

- Continued stabilizing create, view, edit, delete, and image upload flows from Sprint 3.
- Preserved paginated listing browsing with keyword, category, and price filters.
- Kept owner-only listing edit/delete behavior protected behind authentication.
- Improved user feedback after listing create, edit, delete, report, and contact actions.

### Backend API and Data Model Updates

- Added `contact_requests` database table.
- Added contact request model, repository, service, and handler layers.
- Updated listing routes to support `/api/listings/{id}/contact`.
- Updated application wiring in `backend/main.go`.
- Added response handling for unauthorized, forbidden, not found, validation, and server error cases.

### Documentation and Final Submission Prep

- Documented setup and run requirements in this Sprint 4 file for backend, frontend, database, and tests.
- Listed frontend unit tests, Cypress tests, and backend unit tests.
- Added backend API documentation for all current REST endpoints.
- Added a final presentation checklist for the narrated submission video.

## Requirements to Run the Application

### Required Tools

- Go 1.26 or compatible local Go version
- Node.js and npm
- PostgreSQL
- Git

### Backend Setup

1. Create a PostgreSQL database for UniSwap.
2. Configure backend environment values in `backend/resources/app.env`.
3. Apply the database schema from `backend/db/schema.sql`.
4. Start the backend:

```bash
cd backend
go run .
```

Default backend URL: `http://localhost:8080`

### Frontend Setup

1. Install frontend dependencies:

```bash
cd frontend
npm install
```

2. Start the Vite development server:

```bash
npm run dev
```

Default frontend URL: `http://localhost:5173`

### Environment Notes

- Frontend API base URL is controlled by `VITE_API_URL`; if unset, it defaults to `http://localhost:8080/api`.
- Backend CORS currently allows `http://localhost:5173` and `http://127.0.0.1:5173`.
- Authenticated API calls require `Authorization: Bearer <jwt>`.

## How to Use the Application

1. Register a student account with name, email, password, and university.
2. Log in to receive a JWT-backed session.
3. Browse listings on the home page.
4. Search and filter listings by keyword, category, price range, and page.
5. Create a listing with title, description, price, category, and images.
6. View a listing detail page.
7. Contact a seller from a listing detail page.
8. Review received buyer contact requests from `/messages`.
9. Save listings to the wishlist.
10. Edit or delete your own listings from the listing detail or my listings flow.
11. Report inappropriate listings.

## New Functionality Demonstrated in Sprint 4

- Buyer opens a listing detail page and selects "Contact Seller".
- Buyer submits a message about the listing.
- Backend creates a pending contact request tied to the listing, buyer, and seller.
- Seller opens the contact requests page and sees the buyer message with listing title, buyer name, buyer email, message, status, and timestamp.
- UI handles empty inbox, API errors, retry, and successful request display.

## Frontend Unit Tests

Command:

```bash
cd frontend
npm test
```

Current local result from this workspace: 7 test files passed, 24 tests passed.

Test files:

| File | Coverage |
| --- | --- |
| `src/components/ProtectedRoute.test.tsx` | Redirects unauthenticated users to login |
| `src/pages/HomePage.test.tsx` | Paginated listings, debounced search, pagination |
| `src/pages/CreateListingPage.test.tsx` | Create listing with image upload and success state |
| `src/pages/EditListingPage.test.tsx` | Form prefill, replacement image upload, redirect state |
| `src/pages/MyListingsPage.test.tsx` | Empty state for current user's listings |
| `src/pages/ListingDetailPage.test.tsx` | Contact seller, owner restrictions, validation, auth errors, delete flow |
| `src/pages/ContactRequestsPage.test.tsx` | Loading, empty, received requests, backend error, retry |

## Cypress Tests

Command:

```bash
cd frontend
npm run cypress:run
```

Current Cypress test file:

| File | Coverage |
| --- | --- |
| `cypress/e2e/login-form.cy.ts` | Login page renders, credentials can be entered, sign-in button is enabled |

## Backend Unit Tests

Command:

```bash
cd backend
go test ./...
```

Current local result from this workspace: backend package tests are mostly passing, but `go test ./...` currently fails in `main_test.go` at `TestWithCORSOptions` because the test expects `Access-Control-Allow-Origin: http://localhost:5173` on an OPTIONS request that does not include an `Origin` header. Package results from the latest run:

| Package | Latest Result |
| --- | --- |
| `uniswap-campus-marketplace` | Fails: `TestWithCORSOptions` |
| `uniswap-campus-marketplace/apiresponse` | Pass |
| `uniswap-campus-marketplace/config` | Pass |
| `uniswap-campus-marketplace/handlers` | Pass |
| `uniswap-campus-marketplace/middleware` | Pass |
| `uniswap-campus-marketplace/models` | Pass |
| `uniswap-campus-marketplace/repository` | Pass |
| `uniswap-campus-marketplace/services` | Pass |

Backend test files:

| File | Coverage |
| --- | --- |
| `backend/main_test.go` | Health check and CORS behavior |
| `backend/apiresponse/response_test.go` | Success and error response helpers |
| `backend/config/config_test.go` | Config loading, env precedence, database DSN |
| `backend/config/database_test.go` | Database open failure paths |
| `backend/middleware/jwt_middleware_test.go` | Missing, invalid, valid JWT auth handling |
| `backend/models/models_test.go` | JSON tags and password hash omission |
| `backend/repository/repository_test.go` | Repository constructors and closed DB error paths |
| `backend/handlers/auth_handler_test.go` | Register, login, me endpoint behavior |
| `backend/handlers/listing_handler_test.go` | Listing route parsing, update, delete, contact request handler cases |
| `backend/handlers/listing_query_test.go` | Listing query parsing and validation |
| `backend/handlers/upload_handler_test.go` | Listing image upload path parsing and handler cases |
| `backend/handlers/wishlist_handler_test.go` | Wishlist create, get, delete, conflict, auth cases |
| `backend/handlers/user_handler_test.go` | User listing ownership and authorization |
| `backend/handlers/contact_request_handler_test.go` | Received contact request endpoint cases |
| `backend/services/auth_service_test.go` | Auth validation, register, login, token parsing, user lookup |
| `backend/services/listing_service_test.go` | Create validation, update auth, delete behavior |
| `backend/services/listing_image_service_test.go` | Image validation, authorization, not found, success |
| `backend/services/report_service_test.go` | Report validation, not found, success |
| `backend/services/wishlist_service_test.go` | Wishlist validation, create, delete, list |
| `backend/services/contact_request_service_test.go` | Contact validation, self-contact blocking, success, seller inbox |

## Backend API Documentation

All JSON API responses use a consistent wrapper:

```json
{
  "success": true,
  "data": {}
}
```

Error responses use:

```json
{
  "success": false,
  "error": "message"
}
```

### Health

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| GET | `/health` | No | Backend health check |

### Authentication

| Method | Endpoint | Auth | Body | Description |
| --- | --- | --- | --- | --- |
| POST | `/api/auth/register` | No | `full_name`, `email`, `password`, `university` | Register a new user |
| POST | `/api/auth/login` | No | `email`, `password` | Log in and receive JWT plus user |
| GET | `/api/auth/me` | Yes | None | Get the current authenticated user |

### Listings

| Method | Endpoint | Auth | Body or Query | Description |
| --- | --- | --- | --- | --- |
| GET | `/api/listings` | No | `keyword` or `search`, `category`, `min_price`, `max_price`, `page`, `limit` | List paginated and filtered listings |
| POST | `/api/listings` | Yes | `title`, `description`, `price`, `category` | Create a listing |
| GET | `/api/listings/{id}` | No | None | Get listing details |
| PUT | `/api/listings/{id}` | Yes, owner only | `title`, `description`, `price`, `category` | Update a listing |
| DELETE | `/api/listings/{id}` | Yes, owner only | None | Delete a listing |
| POST | `/api/listings/{id}/report` | Yes | `reason` | Report a listing |
| POST | `/api/listings/{id}/contact` | Yes, non-owner | `message` | Contact seller about a listing |
| POST | `/api/listings/{id}/images` | Yes, owner only | multipart `files` | Upload one or more listing images |

### Uploads

| Method | Endpoint | Auth | Body | Description |
| --- | --- | --- | --- | --- |
| POST | `/api/uploads/image` | Yes | multipart `file` | Upload a standalone image and receive a URL |

### Wishlist

| Method | Endpoint | Auth | Body | Description |
| --- | --- | --- | --- | --- |
| GET | `/api/wishlist` | Yes | None | Get current user's saved listings |
| POST | `/api/wishlist` | Yes | `listing_id` | Save a listing |
| DELETE | `/api/wishlist/{wishlist_id}` | Yes | None | Remove a saved listing |

### Users

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| GET | `/api/users/{user_id}/listings` | Yes, same user only | Get listings owned by a user |

### Contact Requests

| Method | Endpoint | Auth | Description |
| --- | --- | --- | --- |
| GET | `/api/contact-requests/received` | Yes | Get contact requests received by the authenticated seller |

### Common Status Codes

| Status | Meaning |
| --- | --- |
| 200 | Successful read or update |
| 201 | Resource created |
| 204 | Successful delete with no response body |
| 400 | Invalid request body or validation error |
| 401 | Missing or invalid authentication |
| 403 | Authenticated user is not allowed to perform the action |
| 404 | Resource not found |
| 405 | Method not allowed |
| 409 | Duplicate resource conflict |
| 500 | Unexpected server error |

## Final Video Presentation Checklist

Each team member should narrate a portion of the final video. The video should include:

1. New Sprint 4 functionality: buyer contacts seller and seller views received contact request.
2. Frontend walkthrough: register, login, browse, search/filter, listing details, create, edit, image upload, wishlist, report, contact requests, delete.
3. Backend API walkthrough: auth, listings, uploads, wishlist, user listings, reports, contact requests.
4. Test results: frontend Vitest, Cypress, backend Go tests, including Sprint 3 tests.
5. Project pitch: explain UniSwap as a campus marketplace for trusted student-to-student buying and selling.

## Known Follow-Up

- Fix or update `TestWithCORSOptions` so the expected CORS behavior matches the current implementation requiring an `Origin` header.
- Run Cypress in the final demo environment after starting the Vite dev server.
- Keep the root `README.md` aligned with the setup and usage details documented above.
