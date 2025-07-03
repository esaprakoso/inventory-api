# Functional Requirements Document (FRD) - POS API

Dokumen ini menjelaskan persyaratan fungsional untuk sistem Point of Sale (POS) API.

## 1. Manajemen Pengguna (User Management)

### 1.1. Registrasi Pengguna
- **Deskripsi:** Sistem harus memungkinkan pengguna baru untuk mendaftar.
- **Persyaratan:**
    - Pengguna mendaftar dengan memberikan `nama`, `email`, dan `password`.
    - Email harus unik di seluruh sistem.
    - Secara default, pengguna baru akan memiliki peran (role) "user".
    - API Endpoint: `POST /api/auth/register`

### 1.2. Login Pengguna
- **Deskripsi:** Sistem harus memungkinkan pengguna terdaftar untuk masuk.
- **Persyaratan:**
    - Pengguna login menggunakan `email` dan `password`.
    - Setelah berhasil, sistem akan memberikan JWT (Access Token dan Refresh Token).
    - API Endpoint: `POST /api/auth/login`

### 1.3. Refresh Token
- **Deskripsi:** Sistem harus bisa memperbarui access token yang sudah kedaluwarsa menggunakan refresh token.
- **Persyaratan:**
    - Pengguna mengirimkan refresh token yang valid.
    - Sistem akan memberikan access token baru.
    - API Endpoint: `POST /api/auth/refresh`

### 1.4. Pengelolaan Pengguna oleh Admin
- **Deskripsi:** Pengguna dengan peran "admin" dapat mengelola data semua pengguna.
- **Persyaratan:**
    - Admin dapat melihat daftar semua pengguna (`GET /api/users`).
    - Admin dapat melihat detail satu pengguna berdasarkan ID (`GET /api/users/:id`).
    - Admin dapat memperbarui data pengguna (`PATCH /api/users/:id`).
    - Admin dapat menghapus pengguna (`DELETE /api/users/:id`).
    - Akses ke endpoint ini hanya diizinkan untuk admin.

## 2. Manajemen Profil Pengguna (User Profile)

### 2.1. Lihat Profil
- **Deskripsi:** Pengguna yang sudah login dapat melihat informasi profilnya sendiri.
- **Persyaratan:**
    - Sistem menampilkan data pengguna yang sedang login.
    - API Endpoint: `GET /api/profile`

### 2.2. Perbarui Profil
- **Deskripsi:** Pengguna dapat memperbarui data profilnya sendiri (misal: nama, email).
- **Persyaratan:**
    - Pengguna hanya bisa mengubah data miliknya sendiri.
    - API Endpoint: `PATCH /api/profile`

### 2.3. Perbarui Password
- **Deskripsi:** Pengguna dapat mengubah password akunnya.
- **Persyaratan:**
    - Pengguna harus menyediakan password lama dan password baru.
    - API Endpoint: `PATCH /api/profile/password`

## 3. Manajemen Produk dan Kategori

### 3.1. Pengelolaan Kategori
- **Deskripsi:** Admin dapat mengelola kategori produk.
- **Persyaratan:**
    - Semua pengguna dapat melihat daftar kategori (`GET /api/categories`) dan detail kategori (`GET /api/categories/:id`).
    - Hanya admin yang dapat membuat (`POST /api/categories`), memperbarui (`PUT /api/categories/:id`), dan menghapus (`DELETE /api/categories/:id`) kategori.

### 3.2. Pengelolaan Produk
- **Deskripsi:** Admin dapat mengelola produk.
- **Persyaratan:**
    - Semua pengguna dapat melihat daftar produk (`GET /api/products`) dan detail produk (`GET /api/products/:id`).
    - Hanya admin yang dapat membuat (`POST /api/products`), memperbarui (`PUT /api/products/:id`), dan menghapus (`DELETE /api/products/:id`) produk.

### 3.3. Manajemen Stok Produk
- **Deskripsi:** Admin dapat menyesuaikan jumlah stok untuk sebuah produk.
- **Persyaratan:**
    - Admin dapat menambah atau mengurangi stok produk tertentu.
    - Setiap perubahan stok harus tercatat.
    - API Endpoint: `PATCH /api/products/:id/stock`

## 4. Manajemen Promosi

### 4.1. Promosi Produk (Product Promotions)
- **Deskripsi:** Admin dapat membuat promosi yang berlaku untuk produk tertentu (misal: diskon harga).
- **Persyaratan:**
    - Semua pengguna dapat melihat promosi produk yang aktif.
    - Hanya admin yang dapat membuat, memperbarui, dan menghapus promosi produk.
    - Endpoints: `GET /api/product-promotions`, `POST /api/product-promotions`, `PUT /api/product-promotions/:id`, `DELETE /api/product-promotions/:id`.

### 4.2. Promosi Keranjang (Cart Promotions)
- **Deskripsi:** Admin dapat membuat promosi yang berlaku untuk total belanja di keranjang (misal: diskon jika total belanja > Rp 100.000).
- **Persyaratan:**
    - Semua pengguna dapat melihat promosi keranjang yang aktif.
    - Hanya admin yang dapat membuat, memperbarui, dan menghapus promosi keranjang.
    - Endpoints: `GET /api/cart-promotions`, `POST /api/cart-promotions`, `PUT /api/cart-promotions/:id`, `DELETE /api/cart-promotions/:id`.

## 5. Manajemen Pesanan (Order Management)

### 5.1. Membuat Pesanan
- **Deskripsi:** Pengguna yang login dapat membuat pesanan baru dari produk yang tersedia.
- **Persyaratan:**
    - Sistem akan menghitung total harga, menerapkan promosi yang berlaku, dan mengurangi stok produk.
    - API Endpoint: `POST /api/orders`

### 5.2. Melihat Riwayat Pesanan
- **Deskripsi:** Pengguna dapat melihat riwayat pesanannya sendiri. Admin dapat melihat semua pesanan.
- **Persyaratan:**
    - Pengguna biasa hanya melihat pesanannya sendiri.
    - Admin dapat melihat semua pesanan yang masuk ke sistem.
    - Endpoints: `GET /api/orders`, `GET /api/orders/:id`.
