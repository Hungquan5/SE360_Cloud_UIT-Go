# Kiến trúc tổng quan UIT-Go

## 🎯 Mục tiêu

UIT-Go là một nền tảng gọi xe giả tưởng. Ở giai đoạn “Bộ xương”, hệ thống được triển khai dưới dạng ba **microservice độc lập**, mỗi service có cơ sở dữ liệu riêng và được đóng gói bằng Docker.  
Mục tiêu của kiến trúc:

- Cung cấp API đơn giản cho người dùng và tài xế (đăng ký, đăng nhập, yêu cầu chuyến đi, theo dõi vị trí tài xế).  
- Cho phép mở rộng dễ dàng, chịu tải cao, hỗ trợ event-driven bằng NATS.  
- Hỗ trợ triển khai hạ tầng bằng Terraform để có thể tái tạo môi trường nhanh chóng.

---

## 🧩 Các thành phần chính

### 1. User Service
- **Trách nhiệm:** Quản lý người dùng (hành khách & tài xế), đăng ký, đăng nhập, JWT Auth.  
- **Công nghệ:** Go, PostgreSQL, JWT, Docker.  
- **API chính:**  
  - `POST /users` – Đăng ký tài khoản  
  - `POST /sessions` – Đăng nhập  
  - `GET /users/me` – Thông tin người dùng  
  - `POST /drivers/apply` – Đăng ký tài xế  

---

### 2. Trip Service
- **Trách nhiệm:** Quản lý và điều phối chuyến đi, xử lý trạng thái, tính giá cước, phát sự kiện.  
- **Công nghệ:** Go, PostgreSQL, NATS, WebSocket.  
- **API chính:**  
  - `POST /trips` – Tạo chuyến  
  - `POST /trips/{id}/cancel` – Hủy chuyến  
  - `POST /trips/{id}/complete` – Hoàn tất chuyến  
  - `POST /trips/{id}/rate` – Đánh giá tài xế  
  - `GET /trips/{id}/ws` – Kênh WebSocket theo dõi tài xế  

---

### 3. Driver Service
- **Trách nhiệm:** Cập nhật vị trí tài xế theo thời gian thực, quản lý trạng thái và tìm kiếm tài xế gần nhất.  
- **Công nghệ:** Go, Redis (Geospatial), Docker.  
- **API chính:**  
  - `POST /drivers/{id}/location` – Cập nhật vị trí  
  - `POST /drivers/{id}/status/{state}` – Cập nhật trạng thái  
  - `GET /drivers/nearby` – Tìm tài xế gần điểm đón  

---

## ⚙️ Luồng dữ liệu chuẩn

1. **Đăng nhập & xác thực:**  
   Người dùng đăng ký và đăng nhập qua `User Service`, JWT trả về được gửi kèm trong mọi request.  

2. **Tạo chuyến:**  
   - Hành khách gửi yêu cầu qua `Trip Service`.  
   - Service tính giá cước (dựa vào Haversine) → ghi vào DB → phát sự kiện `trip.requested` qua **NATS**.  

3. **Tìm tài xế:**  
   - `Driver Service` lưu vị trí tài xế vào Redis `GEOADD`.  
   - Khi có yêu cầu, service sử dụng `GEOSEARCH` để tìm tài xế gần nhất.  
   - Khi tài xế nhận chuyến, phát sự kiện `driver.accepted`.  

4. **Theo dõi chuyến:**  
   - `Trip Service` nhận sự kiện `driver.location.updated` và đẩy qua **WebSocket** đến hành khách.  
   - Khi hoàn tất, cập nhật DB và phát sự kiện `trip.completed`.  

---

## ☁️ Hạ tầng triển khai

- **Docker Compose**: Mỗi service có `app.compose.yml` riêng (chạy Postgres/Redis/NATS).  
- **Terraform**: Triển khai AWS EC2, VPC, Security Groups, NATS JetStream, tagging chi phí.  
- **NATS JetStream**: Message bus trung gian giữa các service.  
- **CI/CD**: Có thể mở rộng dùng GitHub Actions (Module E).  

---

## 🚀 Module E – Automation & FinOps

- Thiết kế pipeline **CI/CD** (GitHub Actions).  
- Cấu trúc lại **Terraform** theo module có thể tái sử dụng.  
- Gắn **tag AWS** để theo dõi chi phí từng service.  
- Tạo **AWS Budget Alerts**.  
- Khuyến nghị dùng **Spot Instances** hoặc **Graviton** để tiết kiệm chi phí.  

---

## 🔮 Định hướng mở rộng

- **Module A (Scalability):** Load Testing (k6, JMeter), Auto-Scaling, SQS async comm.  
- **Module B (Reliability):** Multi-AZ, Load Balancer, Chaos Engineering.  
- **Module C (Security):** Zero-Trust VPC, IAM Least Privilege, Secrets Manager.  
- **Module D (Observability):** Prometheus, CloudWatch, X-Ray, Dashboard SLO.  

---

## 🧭 Kết luận

UIT-Go hiện là hệ thống **microservices hoàn chỉnh**, có khả năng mở rộng theo chiều sâu (Module E) và chiều ngang (Module A-D).  
Hệ thống sẵn sàng được triển khai thực tế hoặc làm nền tảng cho đồ án mở rộng.

---
