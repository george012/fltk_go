#!/bin/bash
set -e

product_name=$(grep ProjectName ../config/config.go | awk -F '"' '{print $2}' | sed 's/\"//g')
Product_version_key="ProjectVersion"
VersionFile=../config/config.go
CURRENT_VERSION=$(grep ${Product_version_key} $VersionFile | awk -F '"' '{print $2}' | sed 's/\"//g')
build_path=./build
RUN_MODE=release
UPLOAD_TMP_DIR=upload_tmp_dir

extend_name_agent=""


OS_TYPE="Unknown"
GetOSType() {
    uNames=$(uname -s)
    osName=${uNames: 0: 4}
    if [ "$osName" == "Darw" ]; then
        OS_TYPE="Darwin"
    elif [ "$osName" == "Linu" ]; then
        OS_TYPE="Linux"
    elif [ "$osName" == "MING" ]; then
        OS_TYPE="Windows"
    else
        OS_TYPE="Unknown"
    fi
}
GetOSType

function getLibVersion() {
    local lib=$1
    version=$(grep "${lib}" go.mod | awk '{print $2}')
    if [ -z "$version" ]; then
        echo "Version not found for ${lib} in go.mod!"
        exit 1
    fi
    echo "$version"
}

function toBuild() {
    rm -rf ${build_path}/${RUN_MODE}
    mkdir -p ${build_path}/${RUN_MODE}

    mkdir -p ${build_path}/${UPLOAD_TMP_DIR}

    go_version=$(go version | awk '{print $3}')
    commit_hash=$(git show -s --format=%H)
    commit_date=$(git show -s --format="%ci")

    if [[ "$OS_TYPE" == "Darwin" ]]; then
        formatted_time=$(date -u -j -f "%Y-%m-%d %H:%M:%S %z" "${commit_date}" "+%Y-%m-%d_%H:%M:%S")
    else
        formatted_time=$(date -u -d "${commit_date}" "+%Y-%m-%d_%H:%M:%S")
    fi

    build_time=$(date -u +"%Y-%m-%d_%H:%M:%S")

    ld_flag_master="-X main.mGitCommitHash=${commit_hash} -X main.mGitCommitTime=${formatted_time} -X main.mGoVersion=${go_version} -X main.mPackageOS=${OS_TYPE} -X main.mPackageTime=${build_time} -X main.mRunMode=${RUN_MODE} -s -w"

    if [[ "$OS_TYPE" == "Darwin" ]]; then

        mkdir -p ${build_path}/${RUN_MODE}/darwin/amd64
        mkdir -p ${build_path}/${RUN_MODE}/darwin/arm64
        mkdir -p ${build_path}/${RUN_MODE}/darwin/universal

        # 处理图标文件
        create_mac_resource

        # Build for macOS x64
        CGO_LDFLAGS="-lpthread -framework OpenGL"
        CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CGO_LDFLAGS=$CGO_LDFLAGS go build -o "${build_path}/${RUN_MODE}/darwin/amd64/${product_name}" -trimpath -ldflags "${ld_flag_master}" main.go
        chmod a+x "${build_path}/${RUN_MODE}/darwin/amd64/${product_name}"
#        package_macos_app "${build_path}/${RUN_MODE}/darwin/amd64" "amd64"

        # Build for macOS arm64
        CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CGO_LDFLAGS=$CGO_LDFLAGS go build -o "${build_path}/${RUN_MODE}/darwin/arm64/${product_name}" -trimpath -ldflags "${ld_flag_master}" main.go
        chmod a+x "${build_path}/${RUN_MODE}/darwin/arm64/${product_name}"
#        package_macos_app "${build_path}/${RUN_MODE}/darwin/arm64" "arm64"

        # 合并二进制文件
        echo "merge ${product_name} darwin amd64 and arm64 to universal"
        lipo -create -output ${build_path}/${RUN_MODE}/darwin/universal/${product_name} ${build_path}/${RUN_MODE}/darwin/amd64/${product_name} ${build_path}/${RUN_MODE}/darwin/arm64/${product_name}
        chmod a+x ${build_path}/${RUN_MODE}/darwin/universal/${product_name}

        package_macos_app "${build_path}/${RUN_MODE}/darwin/universal" "universal"

        rm -rf ${build_path}/${RUN_MODE}/darwin/arm64
        rm -rf ${build_path}/${RUN_MODE}/darwin/amd64
        rm -rf ${build_path}/${RUN_MODE}/darwin/AppIcon.icns

        if [[ -n "$extend_name_agent" ]]; then
            echo "build ${product_name}_${extend_name_agent}"
            CC=x86_64-linux-musl-gcc GOARCH=amd64 GOOS=linux CGO_ENABLED=1 CGO_LDFLAGS="-static" go build -o ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent}/${product_name}_${extend_name_agent} -trimpath -ldflags "${ld_flag_master}" ./${product_name}_${extend_name_agent}/${product_name}_${extend_name_agent}.go \
            && chmod a+x ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent}/${product_name}_${extend_name_agent} \
            && cp ./example_files/${product_name}_${extend_name_agent}.service ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent} \
            && cp ./example_files/install_${product_name}_${extend_name_agent}.sh ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent} \
            && mkdir -p ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent}/conf \
            && cp ./example_files/config_${extend_name_agent}.example.json ${build_path}/${RUN_MODE}/${product_name}_${extend_name_agent}/conf/config_${extend_name_agent}.json

            package_linux_binary_files
        fi
    elif [[ "$OS_TYPE" == "Windows" ]]; then
        mkdir -p ${build_path}/${RUN_MODE}/windows/amd64

        if [ ! -f ./favicon.ico ]; then
          echo "Generating favicon.ico"
          magick ./resources/imgs/Icon.png -strip -depth 8 -type TrueColor -compress None -define icon:auto-resize=256,128,64,32,16 ./favicon.ico
        fi

        generate_windows_package_file

        windres -i main.rc -o main.syso -O coff

        CGO_LDFLAGS="-static -static-libgcc -static-libstdc++ -lglu32 -lopengl32 -lgdiplus -lole32 -luuid -lcomctl32 -lws2_32 -lmsvcrt"

        CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CGO_LDFLAGS=$CGO_LDFLAGS go build -a -trimpath -ldflags "${ld_flag_master} -H windowsgui -w -s" -o ${build_path}/${RUN_MODE}/windows/amd64/${product_name}.exe
        chmod a+x ${build_path}/${RUN_MODE}/windows/amd64/${product_name}.exe

        rm -rf ./main.rc
        rm -rf ./main.syso

        package_windows_files "amd64"
    fi
}

function package_linux_binary_files(){

    local BUILD_OS_TYPE=$(echo "$OS_TYPE" | tr '[:upper:]' '[:lower:]')

    cd ${build_path}/${RUN_MODE} \
    && echo "package ${product_name}_${extend_name_agent}" \
    && zip -r ./${product_name}_${extend_name_agent}_${RUN_MODE}_${CURRENT_VERSION}_linux_amd64.zip ./${product_name}_${extend_name_agent} \
    && mkdir -p ../${UPLOAD_TMP_DIR} \
    && mv *.zip ../${UPLOAD_TMP_DIR} \
    && cd ../../ \
    && echo current dir with $PWD
}


function generate_windows_package_file() {
    # 动态生成 main.rc 文件
    cat <<EOL > main.rc
#include "winver.h"

1 ICON "favicon.ico"

VS_VERSION_INFO VERSIONINFO
FILEVERSION 1,0,0,0
PRODUCTVERSION 1,0,0,0
FILEFLAGSMASK 0x3fL
#ifdef _DEBUG
FILEFLAGS 0x1L
#else
FILEFLAGS 0x0L
#endif
FILEOS 0x4L
FILETYPE 0x1L
FILESUBTYPE 0x0L
BEGIN
    BLOCK "StringFileInfo"
    BEGIN
        BLOCK "040904E4"
        BEGIN
            VALUE "CompanyName", "Free"
            VALUE "FileDescription", "${product_name} Application"
            VALUE "FileVersion", "${CURRENT_VERSION}"
            VALUE "InternalName", "${product_name}"
            VALUE "LegalCopyright", "Free. All rights reserved."
            VALUE "OriginalFilename", "${product_name}.exe"
            VALUE "ProductName", "${product_name}"
            VALUE "ProductVersion", "${CURRENT_VERSION}"
        END
    END
    BLOCK "VarFileInfo"
    BEGIN
        VALUE "Translation", 0x409, 1252
    END
END
EOL
}
function package_windows_files() {

    if [[ "$OS_TYPE" == "Windows" ]]; then
        cd ${build_path}/${RUN_MODE}/windows/amd64
        mkdir -p ${product_name}
        mv ${product_name}.exe ./${product_name}/
#        zip -r ./${product_name}_${RUN_MODE}_${CURRENT_VERSION}_linux_amd64.zip ./${product_name}
        7z a ./${product_name}_${RUN_MODE}_${CURRENT_VERSION}_windows_amd64.zip ./${product_name} >/dev/null 2>&1
        mv ./${product_name}_${RUN_MODE}_${CURRENT_VERSION}_windows_amd64.zip ../../../${UPLOAD_TMP_DIR}
        cd ../../../../
    else
        return
    fi

}

function package_macos_app() {
    local build_dir="$1"
    local arch="$2"
    local app_name="${product_name}.app"
    local app_dir="${build_dir}/${app_name}"
    local contents_dir="${app_dir}/Contents"
    local macos_dir="${contents_dir}/MacOS"
    local resources_dir="${contents_dir}/Resources"

    mkdir -p ${macos_dir}
    mkdir -p ${resources_dir}

    # 将可执行文件移动到MacOS目录
    mv "${build_dir}/${product_name}" ${macos_dir}/

    # 创建Info.plist文件
    cat > ${contents_dir}/Info.plist <<EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>${product_name}</string>
    <key>CFBundleIdentifier</key>
    <string>com.example.${product_name}</string>
    <key>CFBundleName</key>
    <string>${product_name}</string>
    <key>CFBundleVersion</key>
    <string>${CURRENT_VERSION}</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleSignature</key>
    <string>????</string>
    <key>CFBundleInfoDictionaryVersion</key>
    <string>6.0</string>
    <key>CFBundleIconFile</key>
    <string>AppIcon.icns</string>
</dict>
</plist>
EOL

    # 将配置文件移动到Resources目录
    if [ -d "${build_path}/${RUN_MODE}/conf" ]; then
        mv "${build_path}/${RUN_MODE}/conf" ${resources_dir}/
    fi


    # 创建最终的.app包
    cd ${build_dir}
    echo "---------current path $(pwd)-----------"
    cp ../AppIcon.icns ./${product_name}.app/Contents/Resources/AppIcon.icns

    # 打包成 DMG
    create-dmg \
      --volname "${product_name}_${RUN_MODE}_${CURRENT_VERSION}_darwin_${arch}_Installer" \
      --volicon "../AppIcon.icns" \
      --background "../../../../resources/imgs/dmg_bg.png" \
      --window-pos 435 160 \
      --window-size 512 490 \
      --icon-size 128 \
      --icon "${product_name}.app" 190 128 \
      --hide-extension "${product_name}.app" \
      --app-drop-link 382 128 \
      "${product_name}_${RUN_MODE}_${CURRENT_VERSION}_darwin_${arch}.dmg" \
      "${app_name}"

    mv ${product_name}_${RUN_MODE}_${CURRENT_VERSION}_darwin_${arch}.dmg ../../../${UPLOAD_TMP_DIR}/${product_name}_${RUN_MODE}_${CURRENT_VERSION}_darwin_${arch}.dmg

    cd ../../../../
    echo "---------current path $(pwd)-----------"
}

function create_mac_resource() {
    echo "--------------------create_mac_resource-------------------------------------"

    local app_icons=${build_path}/${RUN_MODE}/AppIcon.iconset
    mkdir -p ${app_icons}
    sips -z 16 16     ./resources/imgs/icon.png --out ${app_icons}/icon_16x16.png
    sips -z 32 32     ./resources/imgs/icon.png --out ${app_icons}/icon_16x16@2x.png
    sips -z 32 32     ./resources/imgs/icon.png --out ${app_icons}/icon_32x32.png
    sips -z 64 64     ./resources/imgs/icon.png --out ${app_icons}/icon_32x32@2x.png
    sips -z 128 128   ./resources/imgs/icon.png --out ${app_icons}/icon_128x128.png
    sips -z 256 256   ./resources/imgs/icon.png --out ${app_icons}/icon_128x128@2x.png
    sips -z 256 256   ./resources/imgs/icon.png --out ${app_icons}/icon_256x256.png
    sips -z 512 512   ./resources/imgs/icon.png --out ${app_icons}/icon_256x256@2x.png
    sips -z 512 512   ./resources/imgs/icon.png --out ${app_icons}/icon_512x512.png
    sips -z 1024 1024 ./resources/imgs/icon.png --out ${app_icons}/icon_512x512@2x.png

    iconutil -c icns ${app_icons} -o ${build_path}/${RUN_MODE}/darwin/AppIcon.icns

    rm -rf ${app_icons}
}

function handlerunMode() {
    if [[ "$1" == "release" || "$1" == "" ]]; then
        RUN_MODE=release
    elif [[ "$1" == "test" ]]; then
        RUN_MODE=test
    elif [[ "$1" == "debug" ]]; then
        RUN_MODE=debug
    else
        echo "Usage: bash build.sh [release|test|debug], default is release"
        exit 1
    fi
}

handlerunMode "$1" && toBuild
