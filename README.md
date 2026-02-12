# 🎓 UniSwap – Campus Marketplace for Students

> **Buy, Sell, and Trade: Books, Gadgets, Furniture & More.**  
> A secure, campus-exclusive marketplace designed to connect students.

---

## Project Description

University students frequently buy and sell items such as textbooks, electronics, furniture, and daily-use items. Currently, they rely on scattered platforms like group chats, social media, or generic marketplaces that are not campus-specific, lack trust, and provide poor search and communication experiences.

**UniSwap** is a full-stack web application designed exclusively for university students to buy and sell items within their campus community.

The platform allows authenticated users to:
- Create, edit, and delete item listings with images
- Search listings using filters such as category, price, and keywords
- Communicate with buyers and sellers using real-time chat
- Save listings to a wishlist
- Report inappropriate listings

The system is developed using **Go, PostgreSQL, React, and TypeScript**, and built incrementally using **Agile sprints**.

---

## Team Members

| Name | Role |
|------|------|
| Samarth vinayaka| Frontend Engineer |
| Bhumi Jain| Frontend Engineer |
| Shubhank Chandak | Backend Engineer |
| Nandini Agrawal | Backend Engineer |

---

## Tech Stack

### Frontend
- React  
- TypeScript  
- REST API integration  
- WebSockets for real-time chat  

### Backend
- Go (Golang)  
- RESTful APIs  
- PostgreSQL  
- JWT authentication  
- WebSockets  
- Image upload (local or S3-compatible storage)  

---

## DB Schemas

`users`
- `id` BIGSERIAL PRIMARY KEY
- `full_name` VARCHAR(120) NOT NULL
- `email` VARCHAR(255) NOT NULL UNIQUE
- `password_hash` TEXT NOT NULL
- `university` VARCHAR(150)
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

`listings`
- `id` BIGSERIAL PRIMARY KEY
- `seller_id` BIGINT NOT NULL REFERENCES `users(id)` ON DELETE CASCADE
- `title` VARCHAR(150) NOT NULL
- `description` TEXT NOT NULL
- `category` VARCHAR(80) NOT NULL
- `price` NUMERIC(10,2) NOT NULL CHECK (`price >= 0`)
- `status` VARCHAR(20) NOT NULL DEFAULT `'active'` CHECK (`status IN ('active','sold','hidden')`)
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()
- `updated_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

`listing_images`
- `id` BIGSERIAL PRIMARY KEY
- `listing_id` BIGINT NOT NULL REFERENCES `listings(id)` ON DELETE CASCADE
- `image_url` TEXT NOT NULL
- `is_primary` BOOLEAN NOT NULL DEFAULT FALSE
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

`reports`
- `id` BIGSERIAL PRIMARY KEY
- `listing_id` BIGINT NOT NULL REFERENCES `listings(id)` ON DELETE CASCADE
- `reporter_id` BIGINT NOT NULL REFERENCES `users(id)` ON DELETE CASCADE
- `reason` TEXT NOT NULL
- `created_at` TIMESTAMPTZ NOT NULL DEFAULT NOW()

Indexes currently defined:
- `idx_listings_seller_id` on `listings(seller_id)`
- `idx_listings_category` on `listings(category)`
- `idx_listings_created_at` on `listings(created_at DESC)`
- `idx_listing_images_listing_id` on `listing_images(listing_id)`
- `idx_reports_listing_id` on `reports(listing_id)`
- `idx_reports_reporter_id` on `reports(reporter_id)`
