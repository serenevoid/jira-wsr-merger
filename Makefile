BINARY_NAME=WSR.exe

all: build

build:
	go build -o ${BINARY_NAME} .
	mv ${BINARY_NAME} .\WSR_Generator
	7z a WSR_Generator .\WSR_Generator\*

clean:
	go clean
	rm .\WSR_Generator\${BINARY_NAME} WSR.csv WSR_Generator.7z
