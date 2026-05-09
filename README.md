# Assets Management Backend

## 🚀 Deploy trên Railway

### Chuẩn bị
1. Push code lên GitHub
2. Tạo tài khoản Railway: https://railway.app

### Deploy Steps

#### 1. Tạo Project
- New Project → Deploy from GitHub repo
- Chọn repo này

#### 2. Thêm PostgreSQL
- Click "+ New" → Database → PostgreSQL

#### 3. Environment Variables
```env
PORT=3001
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_NAME=${{Postgres.PGDATABASE}}
JWT_SECRET=HDwdhCOP4endFG3cFVkJu06B1F7vcJRp8vT7EU8EFk1
FE_URL=https://your-frontend-url.vercel.app
CLOUDINARY_CLOUD_NAME=ddg7y4e24
CLOUDINARY_API_KEY=915124992879912
CLOUDINARY_API_SECRET=r-nR94Zrj0XBo5WLujo2-ENlsmw
GIN_MODE=release
```

#### 4. Generate Domain
Settings → Domains → Generate Domain

#### 5. Tạo Admin User
Sau khi deploy, tạo user admin:

```bash
# Cách 1: Qua API
curl -X POST https://your-railway-url.up.railway.app/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","role":"admin"}'

# Cách 2: Qua Railway Query
# Vào PostgreSQL service → Data → Query
# Chạy script tạo hash password trước:
cd scripts
go run hash_password.go admin123
# Copy hash và chạy SQL
```

## 🔧 Local Development

### Requirements
- Go 1.21+
- PostgreSQL 14+

### Setup
```bash
# Install dependencies
go mod download

# Copy .env
cp .env.example .env

# Edit .env với config local

# Run
go run cmd/api/main.go
```

### Tạo Admin User Local
```bash
# Hash password
cd scripts
go run hash_password.go your_password

# Hoặc dùng API
curl -X POST http://localhost:3001/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123","role":"admin"}'
```

## 📚 API Documentation

### Auth
- `POST /api/auth/register` - Đăng ký user mới
- `POST /api/auth/login` - Đăng nhập
- `GET /api/me` - Lấy thông tin user hiện tại

### Loan Slips
- `GET /api/loan-slips` - Danh sách phiếu mượn
- `GET /api/loan-slips/:id` - Chi tiết phiếu mượn
- `POST /api/loan-slips` - Tạo phiếu mượn mới (Admin/IT)
- `PUT /api/loan-slips/:id` - Cập nhật phiếu mượn (Admin/IT)
- `PATCH /api/loan-slips/:id/status` - Cập nhật trạng thái (Admin/IT)
- `DELETE /api/loan-slips/:id` - Xóa phiếu mượn (Admin/IT)

### Notifications
- `GET /api/notifications` - Danh sách thông báo (Admin/IT)
- `PUT /api/notifications/:id` - Đánh dấu đã đọc (Admin/IT)
- `GET /api/notifications/unread/count` - Số thông báo chưa đọc (Admin/IT)

### Dashboard
- `GET /api/dashboard/loan-metrics` - Thống kê (Admin/IT)

## 🔐 Roles
- `admin` - Toàn quyền
- `it` - Quản lý phiếu mượn, thông báo
- `user` - Xem phiếu mượn của mình

## 📝 Database Schema

### Users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

### Loan Slips
```sql
CREATE TABLE loan_slips (
    id SERIAL PRIMARY KEY,
    borrower_name VARCHAR(255) NOT NULL,
    department VARCHAR(255),
    asset_name VARCHAR(255) NOT NULL,
    asset_code VARCHAR(255),
    quantity INTEGER NOT NULL,
    borrow_date TIMESTAMP NOT NULL,
    expected_return_date TIMESTAMP,
    actual_return_date TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    reason TEXT,
    notes TEXT,
    image_url TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## 🛠️ Scripts

### Hash Password
```bash
cd scripts
go run hash_password.go <password>
```

Output sẽ cho bạn:
- Hashed password
- SQL command để insert user
