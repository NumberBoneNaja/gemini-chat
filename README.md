# Four Project (Go 1.23)

ระบบนี้พัฒนาโดยใช้ภาษา Go เวอร์ชัน 1.23 และมีการใช้ไลบรารี (libs) สำหรับเชื่อมต่อฐานข้อมูล, สร้าง REST API, และใช้งาน Google GenAI API

## 📦 ไลบรารีที่ใช้งาน

- `gorm.io/gorm` - ORM สำหรับ Go
- `gorm.io/driver/postgres` - PostgreSQL Driver
- `github.com/gin-gonic/gin` - Web Framework
- `google.golang.org/genai` - ไลบรารีสำหรับใช้งาน Google GenAI API
- `github.com/joho/godotenv` - โหลด `.env` ไฟล์สำหรับตั้งค่า

## 🧰 ก่อนติดตั้ง

### 1. ติดตั้ง Go (แนะนำใช้เวอร์ชัน 1.23 หรือสูงกว่า)

ดาวน์โหลด Go ได้จาก [https://go.dev/dl/](https://go.dev/dl/)


### 2. ติดตั้ง PostgreSQL และเตรียมฐานข้อมูล
รัน Docker compose
```bash
docker-compose up
```
เข้า pgadmin
- Host = postgres
- Port = 5432
- Maintenancee database = mydatabase
- Username = myuser
- password mypassword
