echo "Building..."
go build ../handlers
go build

echo "Installing..."
go install

echo "Running..."
sudo env "PATH=$PATH" dodec
