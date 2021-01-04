#! /bin/bash

echo "Copying static files to server..."
rm -rf ../server/build/
cp ./build/ -r ../server/build/
echo "Frontend static files built by React for server" > ../server/build/README.md