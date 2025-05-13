#!/bin/bash
set -e

# Install MSYS2 using winget
echo "Installing MSYS2 using winget..."
winget install -e --id MSYS2.MSYS2

# Wait for installation to complete
echo "Waiting for MSYS2 installation to complete..."
sleep 30

# Update MSYS2 packages using ucrt64 shell
echo "Updating MSYS2 packages using ucrt64 shell..."
/c/msys64/usr/bin/bash -lc "pacman -Sy --noconfirm"

# Install MinGW-w64 toolchain using mingw64 shell
echo "Installing MinGW-w64 toolchain using mingw64 shell..."
/c/msys64/usr/bin/bash -lc "pacman -S --noconfirm mingw-w64-x86_64-toolchain"

# Set MINGW64_HOME environment variable in Windows
echo "Setting MINGW64_HOME environment variable..."
powershell -Command "[System.Environment]::SetEnvironmentVariable('MINGW64_HOME', 'C:\\msys64\\mingw64', [System.EnvironmentVariableTarget]::Machine)"

# Add MINGW64_HOME\bin to PATH
echo "Adding MINGW64_HOME\\bin to PATH..."
powershell -Command "\$newPath = [System.Environment]::GetEnvironmentVariable('PATH', [System.EnvironmentVariableTarget]::Machine) + ';C:\\msys64\\mingw64\\bin'; [System.Environment]::SetEnvironmentVariable('PATH', \$newPath, [System.EnvironmentVariableTarget]::Machine)"

echo "MSYS2 and MinGW-w64 toolchain installation complete."
