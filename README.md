### รายละเอียดแต่ละไฟล์/โฟลเดอร์

- **main.go**  
  ไฟล์หลักสำหรับรันเซิร์ฟเวอร์ด้วย Gin Framework มีการเชื่อมต่อฐานข้อมูล PostgreSQL และกำหนด API สำหรับจัดการข้อมูลผู้ใช้และแผนก

- **go.mod / go.sum**  
  ใช้สำหรับจัดการ dependencies ของโปรเจกต์และระบุเวอร์ชันของไลบรารีที่ใช้งาน

- **types/**  
  โฟลเดอร์นี้ใช้เก็บไฟล์ที่เกี่ยวกับโครงสร้างข้อมูล (Struct) เช่น `Department` และ `User` เพื่อใช้ในการ map ข้อมูลกับฐานข้อมูล

- **README.md**  
  ไฟล์อธิบายรายละเอียดโปรเจกต์และโครงสร้าง เพื่อให้ผู้ใช้งานหรือผู้พัฒนาคนอื่นเข้าใจภาพรวมของโปรเจกต์

## ขั้นตอนการรันโปรเจกต์

1. **ตั้งค่าการเชื่อมต่อฐานข้อมูล**  
   ตรวจสอบและแก้ไขค่าต่อไปนี้ในไฟล์ `main.go` ให้ตรงกับเครื่องของคุณ

   - host: `localhost`
   - port: `5432`
   - user: `postgres`
   - password: `123456`
   - dbname: `Test-TNS`
   - sslmode: `disable`

2. **ติดตั้ง dependencies**  
    เปิด terminal ที่โฟลเดอร์ backend แล้วรันคำสั่ง:

   ```
   go mod tidy
   go get github.com/lib/pq
   go get github.com/gin-gonic/gin
   go get gorm.io/gorm
   go get gorm.io/driver/postgres

   ```

3. **รันโปรเจกต์**  
   ใช้คำสั่งนี้ใน terminal:

   ```
   go run main.go
   ```

   ระบบจะสร้างตารางในฐานข้อมูลและเพิ่มข้อมูล department เริ่มต้นให้อัตโนมัติ

4. **ทดสอบ API**  
   สามารถเรียกใช้งาน API ผ่าน Postman หรือ Frontend ที่เชื่อมต่อกับ backend นี้  
   ตัวอย่าง endpoint:
   - GET `/users` : ดึงข้อมูลผู้ใช้ทั้งหมด
   - POST `/users/addUser` : เพิ่มผู้ใช้ใหม่
   - PUT `/users/:id` : แก้ไขข้อมูลผู้ใช้
   - DELETE `/users/:id` : ลบผู้ใช้
   - GET `/department` : ดึงข้อมูลแผนกทั้งหมด
