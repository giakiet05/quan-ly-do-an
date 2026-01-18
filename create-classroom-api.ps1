# Script tạo classroom qua API
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJsa2ZvcnVtLXVzZXJzIiwiZXhwIjoxNzY4NjgwMTc5LCJpYXQiOjE3Njg2NzY1NzksImlzcyI6ImxrZm9ydW0iLCJqdGkiOiI1ZTQ5ODA3MS1iM2Y3LTQzNjctYTQ3NS1jMDZjMDQ3YTY4MGYiLCJzdWIiOiI2OTZiYzk0MWJlM2I0YWNjNmY4M2VkMWIifQ.OM1PFraaTf8gbG1aD_XCX_MAXFnWhMwctsbsbYeHnfk"

$body = @{
    name = "Đồ án Phát triển ứng dụng Web - K18"
    description = "Môn học về phát triển ứng dụng web full-stack với React, Node.js và MongoDB"
    semester = "HK1"
    year = 2024
    auto_approve = $true
    max_students = 100
    can_student_delete_group = $true
} | ConvertTo-Json

$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

Write-Host "Creating classroom via API..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/classrooms" -Method POST -Headers $headers -Body $body
    Write-Host "✅ SUCCESS!" -ForegroundColor Green
    Write-Host "Classroom ID: $($response.data.id)" -ForegroundColor Cyan
    Write-Host "Invitation Code: $($response.data.invitation_code)" -ForegroundColor Cyan
} catch {
    Write-Host "❌ ERROR:" -ForegroundColor Red
    Write-Host $_.Exception.Message
}
