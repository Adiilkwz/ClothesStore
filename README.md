# Gentleman's Wardrobe KZ

> A specialized e-commerce platform for men's clothing in Kazakhstan, built with Go and PostgreSQL.

## The Team
* **Member 1 (Lead/Auth):** [Adil] - Database, Auth, User Management
* **Member 2 (Frontend/Store):** [Nurassyl] - Products, UI/UX, Catalog
* **Member 3 (Orders/Docs):** [Yerassyl] - Shopping Cart, Checkout, Admin Panel

---

## For Developers: How to Start

**Follow these steps exactly to run the project on your local machine.**

### 1. Prerequisites
Ensure you have the following installed:
* **Go** (v1.21 or newer)
* **PostgreSQL** (v14 or newer) + **pgAdmin 4**
* **Git**

### 2. Database Setup (Crucial!)
You cannot run the app without the database.
1. Open **pgAdmin 4**.
2. Create a new database named **`clothes_store`**.
3. Right-click the database -> **Query Tool**.
4. Open the file `schema.sql` from this repository, copy the content, and run it in pgAdmin.
5. (Optional) Create a dedicated user:
   ```sql
   CREATE USER store_admin WITH PASSWORD '12345';
   ALTER DATABASE clothing_store OWNER TO store_admin;