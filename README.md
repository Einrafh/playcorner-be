# PlayCorner Backend API

![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8.svg?style=flat-square)
![Docker](https://img.shields.io/badge/Docker-24.0-2496ED.svg?style=flat-square)
![Nginx](https://img.shields.io/badge/Nginx-1.25-009639.svg?style=flat-square)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791.svg?style=flat-square)

Selamat datang di repositori backend untuk **PlayCorner**, sebuah layanan API yang dirancang untuk mengelola sistem peminjaman Game Corner di Fakultas Ilmu Komputer (FILKOM).

Proyek ini dibangun menggunakan Go (dengan framework Fiber) dan berjalan di dalam lingkungan Docker yang siap produksi dengan Nginx sebagai *reverse proxy*.

## ✨ Fitur Utama
- **Otentikasi Pengguna**: Sistem login berbasis JWT dengan *access token* dan *refresh token* (disimpan di HttpOnly cookie).
- **Manajemen User**: Mengambil data profil dan riwayat peminjaman pengguna.
- **Informasi Game Corner**: Mendapatkan daftar TV dan game yang tersedia, dengan relasi spesifik untuk setiap TV.
- **Reservasi Real-time**: Mengecek ketersediaan slot waktu dan membuat reservasi baru.
- **Siap Produksi**: Dikonfigurasi untuk berjalan dengan Docker dan Nginx, lengkap dengan penanganan SSL/TLS.
- **Dokumentasi API**: Dokumentasi lengkap dan interaktif yang dibuat secara otomatis menggunakan **Zudoku**.

## 📖 Dokumentasi API
Dokumentasi API lengkap yang menjelaskan setiap *endpoint*, skema, dan contoh penggunaan dapat diakses di sini:

**[Lihat Dokumentasi API PlayCorner](https://api.playcorner.einrafh.com/docs/introduction)**

*(Catatan: Ganti dengan URL domain produksi Zudoku Anda jika berbeda)*

## 🏗️ Struktur Proyek
Proyek ini mengikuti struktur layout standar Go untuk skalabilitas dan keterbacaan.
```
/
├── cmd/api/             # Main package aplikasi Go
├── docs/                # Proyek dokumentasi Zudoku
├── internal/            # Semua logika bisnis, model, dan handler
│   ├── auth/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   └── utils/
├── .dockerignore        # File yang diabaikan oleh Docker
├── .env                 # (LOKAL) File variabel lingkungan (JANGAN DI-COMMIT)
├── .env.example         # Contoh file environment
├── docker-compose.yml   # Konfigurasi layanan Docker
├── Dockerfile           # Instruksi untuk membangun image aplikasi
├── go.mod               # Dependensi proyek Go
├── nginx.conf           # Konfigurasi Nginx
└── README.md            # Anda sedang membacanya
```

## 🚀 Memulai (Development Lokal)
Untuk menjalankan proyek ini di lingkungan lokal, Anda hanya memerlukan Docker dan Docker Compose.

**1. Clone Repositori**
```bash
git clone <URL_REPOSITORI_ANDA>
cd playcorner-be
```

**2. Konfigurasi Environment**
Salin file `.env.example` menjadi `.env`. File ini akan digunakan oleh Docker Compose untuk mengkonfigurasi semua layanan.
```bash
cp .env.example .env
```
Anda bisa mengubah nilai di dalam `.env` jika perlu, tetapi pengaturan default sudah cukup untuk development.

**3. Jalankan dengan Docker Compose**
Perintah ini akan membangun *image* aplikasi, mengunduh *image* PostgreSQL & Nginx, lalu menjalankan semua layanan secara bersamaan.
```bash
docker-compose up --build
```
* Opsi `--build` hanya diperlukan saat pertama kali atau jika ada perubahan pada kode Go atau `Dockerfile`.
* Untuk menjalankan di latar belakang, gunakan `docker-compose up -d`.

**4. Selesai!**
Server API Anda sekarang berjalan dan dapat diakses di `http://localhost:3000`.

## 🔧 Variabel Lingkungan
Variabel lingkungan dikelola di dalam file `.env`. Pastikan semua variabel ini terisi dengan benar.

| Variabel               | Deskripsi                                                        | Contoh Nilai                               |
| ---------------------- | ---------------------------------------------------------------- | ------------------------------------------ |
| `SERVER_PORT`          | Port internal yang digunakan oleh aplikasi Go.                   | `3000`                                     |
| `DB_HOST`              | Hostname layanan database. **Harus `db`** saat di Docker.        | `db`                                       |
| `DB_PORT`              | Port internal database PostgreSQL.                               | `5432`                                     |
| `DB_USER`              | Username untuk database.                                         | `postgres`                                 |
| `DB_PASSWORD`          | Password untuk database.                                         | `mysecretpassword`                         |
| `DB_NAME`              | Nama database yang akan dibuat.                                  | `playcorner_db`                            |
| `JWT_SECRET`           | Kunci rahasia yang sangat panjang dan acak untuk menandatangani JWT. | `your_super_secret_key_...`                |
| `CORS_ALLOWED_ORIGINS` | Daftar URL frontend yang diizinkan (pisahkan dengan koma).         | `http://localhost:5173,https://app.com`    |

## 📜 Lisensi
Hak Cipta &copy; 2025 **Muhammad Rafly Ash Shiddiqi**.

Semua hak dilindungi undang-undang. Dilarang memperbanyak, mendistribusikan, atau mentransmisikan bagian mana pun dari perangkat lunak ini dalam bentuk apa pun atau dengan cara apa pun tanpa izin tertulis sebelumnya dari pemilik hak cipta.

---
Dibuat dengan ❤️ di FILKOM.

