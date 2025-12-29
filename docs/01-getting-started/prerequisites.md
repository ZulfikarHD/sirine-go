# 📋 Prerequisites - System Requirements

Pastikan semua software requirements terinstall sebelum memulai setup Sirine Go App.

---

## 🖥️ System Requirements

### **Minimum Requirements:**
- **OS:** Ubuntu 20.04+ / Debian 11+ / macOS 12+ / Windows 10+ (dengan WSL2)
- **RAM:** 4GB
- **Storage:** 2GB free space
- **Network:** Internet connection untuk download dependencies

### **Recommended Requirements:**
- **OS:** Ubuntu 22.04 LTS / macOS 13+
- **RAM:** 8GB+
- **Storage:** 5GB+ free space
- **CPU:** 2+ cores

---

## 📦 Required Software

### **1. Go 1.24+**

Go adalah bahasa programming untuk backend aplikasi.

#### **Check if installed:**
```bash
go version
```

Expected output: `go version go1.24.x linux/amd64` (atau platform Anda)

#### **Install Go (Ubuntu/Debian):**

```bash
# Download Go
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Extract ke /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### **Install Go (macOS):**

```bash
# Using Homebrew
brew install go

# Or download from: https://go.dev/dl/

# Verify
go version
```

#### **Install Go (Windows dengan WSL2):**

```bash
# Download
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Extract
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Add to PATH in ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version
```

---

### **2. Node.js 18+ & Yarn**

Node.js untuk frontend tooling, Yarn untuk package management.

#### **Check if installed:**
```bash
node --version
yarn --version
```

Expected output:
```
v20.x.x  (atau v18.x.x+)
1.22.x
```

#### **Install Node.js (Ubuntu/Debian):**

```bash
# Install Node.js 20.x LTS
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify
node --version
npm --version
```

#### **Install Node.js (macOS):**

```bash
# Using Homebrew
brew install node

# Or download from: https://nodejs.org/

# Verify
node --version
npm --version
```

#### **Install Yarn:**

```bash
# Install Yarn globally via npm
npm install -g yarn

# Verify
yarn --version
```

**Alternative - Install Yarn via package manager:**

```bash
# Ubuntu/Debian
curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add -
echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list
sudo apt update && sudo apt install yarn

# macOS
brew install yarn
```

---

### **3. MySQL 8.0+**

MySQL adalah database yang digunakan untuk persist data.

#### **Check if installed:**
```bash
mysql --version
```

Expected output: `mysql  Ver 8.0.x`

#### **Install MySQL (Ubuntu/Debian):**

```bash
# Update package list
sudo apt update

# Install MySQL Server
sudo apt install mysql-server -y

# Verify installation
mysql --version

# Start MySQL service
sudo systemctl start mysql

# Enable auto-start on boot
sudo systemctl enable mysql

# Check status
sudo systemctl status mysql
```

#### **Secure MySQL Installation:**

```bash
# Run secure installation wizard
sudo mysql_secure_installation

# Follow prompts:
# 1. Set root password (if not set)
# 2. Remove anonymous users? Yes
# 3. Disallow root login remotely? Yes (for development, can choose No)
# 4. Remove test database? Yes
# 5. Reload privilege tables? Yes
```

#### **Install MySQL (macOS):**

```bash
# Using Homebrew
brew install mysql

# Start MySQL service
brew services start mysql

# Secure installation
mysql_secure_installation

# Verify
mysql --version
```

#### **Test MySQL Connection:**

```bash
# Login to MySQL
mysql -u root -p
# Enter your password

# You should see MySQL prompt:
mysql>

# Test query
mysql> SELECT VERSION();

# Exit
mysql> EXIT;
```

---

### **4. Git**

Git untuk version control dan cloning repository.

#### **Check if installed:**
```bash
git --version
```

Expected output: `git version 2.x.x`

#### **Install Git (Ubuntu/Debian):**

```bash
sudo apt update
sudo apt install git -y

# Verify
git --version
```

#### **Install Git (macOS):**

```bash
# Using Homebrew
brew install git

# Or install Xcode Command Line Tools
xcode-select --install

# Verify
git --version
```

#### **Configure Git:**

```bash
# Set your name and email
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"

# Verify configuration
git config --list
```

---

## ✅ Verification Checklist

Setelah install semua software, verifikasi dengan commands berikut:

```bash
# Run all verification commands
go version      && \
node --version  && \
yarn --version  && \
mysql --version && \
git --version   && \
echo "✅ All prerequisites installed!"
```

**Expected Output:**

```
go version go1.24.0 linux/amd64
v20.11.0
1.22.19
mysql  Ver 8.0.35 for Linux on x86_64 (MySQL Community Server - GPL)
git version 2.34.1
✅ All prerequisites installed!
```

### **Individual Checks:**

| Software | Command | Expected Output |
|----------|---------|-----------------|
| **Go** | `go version` | `go version go1.24.x` |
| **Node.js** | `node --version` | `v18.x.x` atau higher |
| **Yarn** | `yarn --version` | `1.22.x` |
| **MySQL** | `mysql --version` | `mysql Ver 8.0.x` |
| **Git** | `git --version` | `git version 2.x.x` |

---

## 🛠️ Optional Tools

Tools ini tidak required, tapi sangat helpful untuk development:

### **1. Make**

Make untuk menjalankan Makefile commands.

```bash
# Ubuntu/Debian
sudo apt install make

# macOS (included with Xcode Command Line Tools)
xcode-select --install

# Verify
make --version
```

### **2. cURL**

cURL untuk testing API endpoints.

```bash
# Ubuntu/Debian
sudo apt install curl

# macOS (already included)

# Verify
curl --version
```

### **3. jq**

jq untuk pretty-print JSON responses.

```bash
# Ubuntu/Debian
sudo apt install jq

# macOS
brew install jq

# Verify
echo '{"name":"test"}' | jq
```

### **4. Visual Studio Code**

Recommended code editor dengan Go & Vue extensions.

**Download:** https://code.visualstudio.com/

**Recommended Extensions:**
- Go (golang.go)
- Vue - Official (vue.volar)
- Tailwind CSS IntelliSense (bradlc.vscode-tailwindcss)
- ESLint (dbaeumer.vscode-eslint)

---

## 🐛 Troubleshooting

### **Problem: Go command not found after installation**

**Solution:**
```bash
# Add Go to PATH manually
export PATH=$PATH:/usr/local/go/bin

# Make it permanent
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### **Problem: MySQL service won't start**

**Solution:**
```bash
# Check MySQL status
sudo systemctl status mysql

# View logs
sudo journalctl -u mysql -n 50

# Reset MySQL data directory (CAUTION: Deletes all data!)
sudo rm -rf /var/lib/mysql/*
sudo mysqld --initialize --user=mysql
sudo systemctl start mysql
```

### **Problem: Permission denied when installing npm packages globally**

**Solution:**
```bash
# Option 1: Use sudo (not recommended)
sudo npm install -g yarn

# Option 2: Fix npm permissions (recommended)
mkdir ~/.npm-global
npm config set prefix '~/.npm-global'
echo 'export PATH=~/.npm-global/bin:$PATH' >> ~/.bashrc
source ~/.bashrc

# Now install without sudo
npm install -g yarn
```

---

## 🎯 Next Steps

Setelah semua prerequisites terinstall:

1. ✅ Verify semua software dengan checklist di atas
2. ➡️ Continue to [database-setup.md](./database-setup.md) untuk setup database
3. ➡️ Atau lihat [installation.md](./installation.md) untuk complete setup guide

---

## 📚 Related Documentation

- [installation.md](./installation.md) - Complete installation guide
- [database-setup.md](./database-setup.md) - Database configuration
- [quickstart.md](./quickstart.md) - 5-minute quick setup

---

## 📞 Need Help?

Jika ada masalah saat install prerequisites:
- **Developer:** Zulfikar Hidayatullah
- **Phone:** +62 857-1583-8733
- **Timezone:** Asia/Jakarta (WIB)

---

**Last Updated:** 28 Desember 2025  
**Version:** 2.0.0 (Phase 2 Restructure)
