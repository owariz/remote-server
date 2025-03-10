# Remot-Server

เครื่องมือสำหรับการจัดการและควบคุมเซิร์ฟเวอร์ระยะไกลผ่าน Discord Bot และ API

## 📝 รายละเอียดโปรเจค

Remot-Server เป็นแพลตฟอร์มที่ออกแบบมาเพื่อช่วยให้ผู้ดูแลระบบสามารถจัดการและควบคุมเซิร์ฟเวอร์ได้จากระยะไกล โดยมีฟีเจอร์หลักคือ:

- **Discord Bot Integration**: ควบคุมเซิร์ฟเวอร์ผ่านคำสั่งใน Discord (Node.js + TypeScript)
- **RESTful API**: ช่องทางสำหรับเชื่อมต่อกับแอปพลิเคชันอื่นๆ (Go + Fiber)
- **การจัดการเซิร์ฟเวอร์**: ตรวจสอบสถานะ, รีสตาร์ทบริการ, เรียกดูล็อก และอื่นๆ
- **การแจ้งเตือน**: รับการแจ้งเตือนเมื่อเกิดเหตุการณ์สำคัญบนเซิร์ฟเวอร์

## 🔄 สถาปัตยกรรม

โปรเจคนี้แบ่งเป็นสองส่วนหลัก:

1. **API Server (Go + Fiber)**
   - ทำหน้าที่เชื่อมต่อและควบคุมเซิร์ฟเวอร์โดยตรง
   - สร้าง RESTful API endpoints สำหรับการจัดการเซิร์ฟเวอร์
   - ใช้ Fiber framework เพื่อประสิทธิภาพสูงและการพัฒนาที่รวดเร็ว
   
2. **Discord Bot (Node.js + TypeScript)**
   - ส่วนติดต่อกับผู้ใช้ผ่าน Discord
   - ใช้ Discord.js library
   - เรียกใช้งาน API ที่สร้างด้วย Go เพื่อควบคุมเซิร์ฟเวอร์

## 🚀 เริ่มต้นใช้งาน

### ความต้องการของระบบ
- Go 1.24+
- Node.js 20+ (สำหรับ Discord Bot)
- TypeScript 5+ (สำหรับ Discord Bot)

### การติดตั้ง API Server (Go + Fiber)
```bash
# Clone repository
git clone https://github.com/owariz/remot-server.git
cd remot-server/api

# Download dependencies
go mod download

# Build the project
go build -o remot-server

# Run the server
./remot-server
```

### การติดตั้ง Discord Bot (Node.js + TypeScript)
```bash
# Navigate to bot directory
cd ../bot

# Install dependencies
npm install

# Build TypeScript
npm run build

# Start the bot
npm start
```

### การตั้งค่า
1. สร้างไฟล์ `.env` ตามตัวอย่างใน `.env.example` ในทั้งสองไดเรกทอรี (api และ bot)
2. สำหรับ API Server กำหนดค่าพอร์ตและการเชื่อมต่ออื่นๆ
3. สำหรับ Discord Bot:
   - เพิ่ม Discord Bot Token
   - กำหนด URL ของ API Server

## 🔧 การใช้งาน

### Discord Bot Commands
```
/status - แสดงสถานะเซิร์ฟเวอร์
/restart [service] - รีสตาร์ทบริการที่กำหนด
/logs [service] [lines] - แสดงล็อกของบริการ
/exec [command] - รันคำสั่งบนเซิร์ฟเวอร์ (ต้องมีสิทธิ์)
```

### RESTful API Endpoints (Fiber)
```
GET /api/v1/status - รับสถานะเซิร์ฟเวอร์
POST /api/v1/service/:name/restart - รีสตาร์ทบริการ
GET /api/v1/service/:name/logs?lines=100 - รับล็อกของบริการ
POST /api/v1/exec - รันคำสั่งบนเซิร์ฟเวอร์ (ต้องมีสิทธิ์)
GET /api/v1/metrics - รับข้อมูลประสิทธิภาพของเซิร์ฟเวอร์
```

## 🛠️ เทคโนโลยี

### API Server
- **Go**: ภาษาที่มีประสิทธิภาพสูงสำหรับการพัฒนาแบ็คเอนด์
- **Fiber**: Web framework สำหรับ Go ที่มีความเร็วสูงและใช้งานง่าย
- **JWT**: สำหรับการยืนยันตัวตน
- **Middleware**: สำหรับการจัดการการเข้าถึง

### Discord Bot
- **Node.js**: Runtime สำหรับรันแอปพลิเคชัน JavaScript
- **TypeScript**: Superset ของ JavaScript ที่เพิ่มความปลอดภัยด้วย type
- **Discord.js**: Library สำหรับพัฒนา Discord Bot

## 📂 โครงสร้างโปรเจค
```
remot-server/
├── api/                   # Go Fiber API
│   ├── cmd/               # จุดเริ่มต้นของแอปพลิเคชัน
│   ├── internal/          # โค้ดภายในที่ไม่เปิดเผยภายนอก
│   │   ├── config/        # การตั้งค่าของแอปพลิเคชัน
│   │   ├── handlers/      # API route handlers
│   │   ├── middleware/    # Middleware
│   │   ├── models/        # Data models
│   │   └── services/      # Business logic
│   ├── pkg/               # โค้ดที่สามารถนำไปใช้ซ้ำได้
│   ├── go.mod             # Go module dependencies
│   └── go.sum             # Go module checksums
│
├── bot/                   # Discord Bot (Node.js + TypeScript)
│   ├── src/               # Source code
│   │   ├── commands/      # Bot commands
│   │   ├── events/        # Event handlers
│   │   ├── services/      # Services for API communication
│   │   └── utils/         # Utility functions
│   ├── package.json       # npm dependencies
│   └── tsconfig.json      # TypeScript configuration
│
└── README.md              # Project documentation
```

## 🤝 การมีส่วนร่วม
คำแนะนำและคอนทริบิวชันยินดีรับเสมอ!

1. Fork โปรเจค
2. สร้าง feature branch (`git checkout -b feature/amazing-feature`)
3. Commit การเปลี่ยนแปลง (`git commit -m 'Add some amazing feature'`)
4. Push ไปยัง branch (`git push origin feature/amazing-feature`)
5. เปิด Pull Request

## 📄 ลิขสิทธิ์
โครงการนี้ได้รับอนุญาตภายใต้ใบอนุญาต MIT ดูไฟล์ LICENSE สำหรับรายละเอียดเพิ่มเติม