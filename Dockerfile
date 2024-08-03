FROM golang:latest

WORKDIR /app

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev
ENV TESSDATA_PREFIX=/usr/share/tesseract-ocr/5/tessdata/
RUN apt-get install -y -qq \
  tesseract-ocr-eng \
  tesseract-ocr-deu \
  tesseract-ocr-jpn

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
# RUN go build -o main .

# Expose port
EXPOSE 8080

# Run the Go app when the container launches
# CMD ["./main"]
CMD ["air", "-c", ".air.toml"]