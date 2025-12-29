# 🤝 Contributing to Sirine Go

Terima kasih telah tertarik untuk berkontribusi pada pengembangan Sirine Go App! Dokumen ini berisi panduan untuk membantu Anda memulai.

## 📋 Development Workflow

1. **Clone Repository**
   ```bash
   git clone https://github.com/username/sirine-go.git
   cd sirine-go
   ```

2. **Setup Environment**
   Ikuti langkah-langkah di [Getting Started Guide](./01-getting-started/quickstart.md).

3. **Create Branch**
   Buat branch baru untuk fitur atau fix Anda. Gunakan naming convention:
   - Feature: `feat/nama-fitur` (contoh: `feat/user-profile`)
   - Fix: `fix/nama-bug` (contoh: `fix/login-error`)
   - Refactor: `refactor/nama-modul`

   ```bash
   git checkout -b feat/my-new-feature
   ```

4. **Code Standards**
   - **Backend (Go)**: Ikuti standar `gofmt` dan best practices Go.
   - **Frontend (Vue)**: Gunakan ESLint/Prettier default project.
   - **Commits**: Gunakan Conventional Commits (contoh: `feat: add login page`, `fix: handle null token`).

5. **Testing**
   Pastikan kode Anda lolos testing sebelum commit.
   ```bash
   # Backend test
   cd backend && go test ./...
   
   # Frontend test
   cd frontend && yarn test
   ```

6. **Pull Request**
   - Push branch ke repository.
   - Buat Pull Request (PR) ke branch `main` atau `develop`.
   - Jelaskan perubahan yang Anda buat di deskripsi PR.

## 🏗️ Project Structure

Pahami struktur project di [Folder Structure](./02-architecture/folder-structure.md) sebelum memulai.

## 🛠️ Adding New Features

### Backend
1. Buat model di `backend/models/`.
2. Register model di `backend/database/models_registry.go`.
3. Buat service di `backend/services/`.
4. Buat handler di `backend/handlers/`.
5. Register route di `backend/routes/`.

### Frontend
1. Buat component atau view di `frontend/src/`.
2. Jika perlu state global, update Pinia store.
3. Update router di `frontend/src/router/`.

## 🐛 Reporting Bugs

Jika menemukan bug, silakan buat issue dengan detail:
- Langkah-langkah untuk reproduksi (Steps to reproduce).
- Expected behavior vs Actual behavior.
- Screenshots atau logs error.

## 📄 License

Project ini bersifat Private & Proprietary. Lihat `LICENSE` file untuk detail.
