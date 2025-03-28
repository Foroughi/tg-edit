#!/bin/bash

set -e  # Exit if any command fails

APP_NAME="tg-edit"
PLUGINS_DIR="plugins"

echo "🔹 Cleaning old builds..."
rm -f $APP_NAME $PLUGINS_DIR/*.so

echo "🔹 Building main app..."
go build -o $APP_NAME

echo "🔹 Building plugins..."
for dir in $PLUGINS_DIR/*/; do
    PLUGIN_NAME=$(basename "$dir")
    echo "🔹 Building plugin: $PLUGIN_NAME"
    go build -buildmode=plugin -o "$PLUGINS_DIR/$PLUGIN_NAME.so" "$dir/$PLUGIN_NAME.go"
done

# Add execute permissions on all .so files (plugin files)
echo "🔹 Adding execute permissions to plugin files..."
find $PLUGINS_DIR -type f -name "*.so" -exec chmod +x {} \;

echo "✅ Build completed!"

echo "🚀 Running $APP_NAME..."
./$APP_NAME
