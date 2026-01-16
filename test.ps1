# Test script with colored output
# Runs Go tests and displays results in colored format

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "RQ_MobCounter - Running Tests" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Run tests with JSON output
$testOutput = & go test ./... -json 2>&1

# Check if go command failed
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: Failed to run tests" -ForegroundColor Red
    exit 1
}

# Parse JSON output and display colored results
$testResults = $testOutput | ConvertFrom-Json

$passed = 0
$failed = 0
$totalTime = 0

foreach ($result in $testResults) {
    if ($result.Action -eq "pass" -and $result.Test) {
        $time = [math]::Round($result.Elapsed * 1000, 0)
        Write-Host "$($result.Test) ... " -NoNewline
        Write-Host "PASS" -ForegroundColor Green -NoNewline
        Write-Host " ($time ms)"
        $passed++
        $totalTime += $result.Elapsed
    }
    elseif ($result.Action -eq "fail" -and $result.Test) {
        Write-Host "$($result.Test) ... " -NoNewline
        Write-Host "FAIL" -ForegroundColor Red
        $failed++

        # Show failure output if available
        if ($result.Output) {
            Write-Host "  $($result.Output)" -ForegroundColor Red
        }
    }
    elseif ($result.Action -eq "output" -and $result.Test -and $result.Output -match "FAIL") {
        Write-Host "  $($result.Output)" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Test Results:" -ForegroundColor Cyan
Write-Host "  Passed: $passed" -ForegroundColor Green
Write-Host "  Failed: $failed" -ForegroundColor Red
Write-Host ("  Total time: {0:F2}s" -f $totalTime) -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Cyan

if ($failed -gt 0) {
    exit 1
}

Write-Host ""
Write-Host "All tests passed!" -ForegroundColor Green
Write-Host ""