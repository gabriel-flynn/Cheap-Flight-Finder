protoc -I=. --go_out=. --go-grpc_out=. flight_scraping.proto
python -m grpc_tools.protoc -I. --python_out=./scraper --grpc_python_out=./scraper flight_scraping.proto
