# Script để push Backend lên GitHub

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Push Backend to GitHub" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Kiểm tra git
try {
    git --version | Out-Null
} catch {
    Write-Host "❌ Git chưa cài đặt. Tải tại: https://git-scm.com/download/win" -ForegroundColor Red
    exit
}

# Lấy GitHub username
$username = Read-Host "Nhập GitHub username của bạn"
if ([string]::IsNullOrWhiteSpace($username)) {
    Write-Host "❌ Username không được để trống" -ForegroundColor Red
    exit
}

Write-Host ""
Write-Host "📝 Chuẩn bị push code..." -ForegroundColor Yellow
Write-Host ""

# Check git status
$status = git status --porcelain
if ($status) {
    Write-Host "Có thay đổi chưa commit:" -ForegroundColor Yellow
    git status --short
    Write-Host ""

    # Add files
    Write-Host "Thêm files..." -ForegroundColor Yellow
    git add .

    # Commit
    $commitMsg = Read-Host "Nhập commit message (Enter để dùng 'Prepare for deployment')"
    if ([string]::IsNullOrWhiteSpace($commitMsg)) {
        $commitMsg = "Prepare for deployment"
    }
    git commit -m "$commitMsg"
    Write-Host "✓ Đã commit" -ForegroundColor Green
} else {
    Write-Host "✓ Không có thay đổi mới" -ForegroundColor Green
}

# Check remote
$remoteExists = git remote | Select-String "origin"
if ($remoteExists) {
    $currentRemote = git remote get-url origin
    Write-Host "Remote hiện tại: $currentRemote" -ForegroundColor Cyan
    Write-Host "Giữ nguyên remote này? (Y/n)" -ForegroundColor Yellow
    $keep = Read-Host
    if ($keep -eq "n" -or $keep -eq "N") {
        git remote remove origin
        $remoteExists = $false
    }
}

# Add remote nếu chưa có
if (-not $remoteExists) {
    $repoName = "assets-management-be"
    git remote add origin "https://github.com/$username/$repoName.git"
    Write-Host "✓ Đã thêm remote: https://github.com/$username/$repoName.git" -ForegroundColor Green
}

# Set branch
git branch -M main

# Push
Write-Host ""
Write-Host "🚀 Đang push lên GitHub..." -ForegroundColor Yellow
Write-Host ""

try {
    git push -u origin main
    Write-Host ""
    Write-Host "==================================" -ForegroundColor Cyan
    Write-Host "✅ Push thành công!" -ForegroundColor Green
    Write-Host "==================================" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Repository URL:" -ForegroundColor Yellow
    Write-Host "https://github.com/$username/$repoName" -ForegroundColor Green
    Write-Host ""
    Write-Host "Bước tiếp theo:" -ForegroundColor Yellow
    Write-Host "1. Vào https://railway.app" -ForegroundColor Cyan
    Write-Host "2. New Project → Deploy from GitHub repo" -ForegroundColor Cyan
    Write-Host "3. Chọn repository này" -ForegroundColor Cyan
    Write-Host "4. Thêm PostgreSQL database" -ForegroundColor Cyan
    Write-Host "5. Cấu hình Environment Variables (xem README.md)" -ForegroundColor Cyan
    Write-Host "6. Generate Domain" -ForegroundColor Cyan
} catch {
    Write-Host ""
    Write-Host "❌ Lỗi khi push" -ForegroundColor Red
    Write-Host ""
    Write-Host "Có thể bạn cần:" -ForegroundColor Yellow
    Write-Host "1. Tạo repository trên GitHub trước: https://github.com/new" -ForegroundColor Cyan
    Write-Host "   Tên repo: assets-management-be" -ForegroundColor Cyan
    Write-Host "2. Đăng nhập GitHub nếu chưa" -ForegroundColor Cyan
    Write-Host "3. Thử lại script này" -ForegroundColor Cyan
}

Write-Host ""
