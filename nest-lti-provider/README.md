## Cấu hình Environment Variables (.env)

Để chạy ứng dụng LTI Provider, bạn cần tạo file `.env` trong thư mục gốc của dự án với các biến môi trường sau:

### Cấu hình cơ sở dữ liệu
```env
# Tên cơ sở dữ liệu MongoDB
DATABASE_NAME=lti_provider

# Tên đăng nhập MongoDB
DATABASE_USERNAME=admin

# Mật khẩu MongoDB
DATABASE_PASSWORD=admin

# Cổng kết nối MongoDB
DATABASE_PORT=27017

# URI kết nối MongoDB
DATABASE_URI=mongodb://mongo:27017
```

### Hướng dẫn thiết lập

## 1. **Kích hoạt các plugin cần thiết trong Moodle**

- Truy cập vào Moodle (http://localhost:8888 => file /moodle/docker-compose.yaml expose moodle bằng port 8888 với giao thức http)
- Với vai trò quản trị viên Moodle (tài khoản trong file /moodle/docker-compose.yaml), truy cập **Site administration > Plugins > Authentication > Manage authentication** và đảm bảo **LTI authentication plugin** đã được kích hoạt.
- Quan trọng nhất để kết nối với LTI provider bên ngoài, truy cập **Site administration > Plugins > Activity modules > External tool > Manage tools** và đảm bảo plugin này đã được kích hoạt để thêm các công cụ LTI như các hoạt động.

## 2. **Thêm và cấu hình External LTI Tool trong Moodle**

- Điều hướng đến **Site administration > Plugins > Activity modules > External tool > Manage tools**.
- Nhấp vào **Configure a tool manually** hoặc **Add LTI Advantage**.
- Nhập thông tin sau từ LTI provider của bạn (công cụ bên ngoài mà bạn muốn kết nối):
    - **Tool name:** Tên cho công cụ (hiển thị cho người tạo khóa học).
    - **Tool URL:** URL khởi chạy của LTI provider (http://localhost:3000).
    - **LTI version:** Chọn **LTI 1.3** (chưa rõ Moodle có version nào nên test thử đã).
    - **Public key type:** chọn **Keyset URL**.
    - **Public keyset URL:** URL JWKS từ LTI provider (http://localhost:3000/lti/keys). Để Moodle xác minh JWT từ LTI Provider
    - **Initiate login URL:** URL đăng nhập được cung cấp bởi LTI provider.
    - **Redirection URL:** URL chuyển hướng sau khi đăng nhập, từ provider.
- Cấu hình các thiết lập bổ sung như (bật full hết để test đã):
    - Kích hoạt **Deep Linking (Content-Item Message)** nếu được hỗ trợ để lựa chọn nội dung.
    - Cài đặt quyền riêng tư (ví dụ: chia sẻ tên/email người khởi chạy, chấp nhận điểm số).
    - Tùy chọn sử dụng cấu hình công cụ (ví dụ: hiển thị trong activity chooser).
    - Đánh dấu **Force SSL** nếu trang web của bạn sử dụng HTTPS.

## 3. **Lấy thông tin cấu hình từ Moodle**

Sau khi thiết lập Moodle xong, tại `http://localhost:8888/mod/lti/toolconfigure.php` bấm vào LTI Provider vừa tạo (cái icon details ấy - 3 dấu chấm rồi 3 gạch ngang), nó sẽ hiển thị thông tin có dạng như sau:

```
Platform ID: http://localhost:8888
Client ID: DUInSyShV52fzQr
Deployment ID: 1
Public keyset URL: http://localhost:8888/mod/lti/certs.php
Access token URL: http://localhost:8888/mod/lti/token.php
Authentication request URL: http://localhost:8888/mod/lti/auth.php
```

## 4. **Cấu hình file .env**

Từ những thông tin trên, cập nhật phần cấu hình tích hợp với Moodle trong file `.env`:

```env
# Cổng chạy ứng dụng
PORT=3000

# URL endpoint cho LTI key (Public keyset lúc điền thông tin tạo external tool đó)
LTI_KEY=http://localhost:3000/lti/keys

# Tên của LTI provider (lúc điền thông tin tạo external tool)
LTI_NAME=TQTos-Provider

# URL của platform Moodle (Platform ID)
LTI_PLATFORM_URL=http://localhost:8888

# Client ID được Moodle sử dụng để nhận diện LTI provider
LTI_CLIENT_ID=DUInSyShV52fzQr

# URL chứa public keyset (JWKS) của Moodle để xác thực JWT
LTI_PUBLIC_KEYSET_URL=http://localhost:8888/mod/lti/certs.php

# URL để LTI Provider lấy access token OAuth 2.0 từ Moodle
LTI_ACCESS_TOKEN_URL=http://localhost:8888/mod/lti/token.php

# URL để bắt đầu quá trình xác thực với Moodle
LTI_AUTHENTICATION_URL=http://localhost:8888/mod/lti/auth.php
```

## 5. **Chạy ứng dụng**

1. Mở terminal và chuyển đến thư mục `nest-lti-provider`:
   ```bash
   cd nest-lti-provider
   ```

2. Khởi chạy ứng dụng bằng Docker:
   ```bash
   docker-compose up -d
   ```

3. Sau khi chạy xong, kiểm tra ứng dụng bằng cách truy cập:
   ```
   http://localhost:3000/lti/ping
   ```

**Lưu ý:** Đảm bảo rằng cả Moodle (port 8888) và MongoDB đã được khởi chạy trước khi chạy LTI Provider.