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

# Install necessary dependencies for OpenCV and GoCV
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    pkg-config \
    libgtk-3-dev \
    libavcodec-dev \
    libavformat-dev \
    libswscale-dev \
    libv4l-dev \
    libxvidcore-dev \
    libx264-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    gfortran \
    openexr \
    libatlas-base-dev \
    python3-dev \
    python3-numpy \
    libtbb2 \
    libtbb-dev \
    libdc1394-22-dev \
    libopenexr-dev \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer1.0-dev \
    wget \
    unzip

# Set the Go environment variables
ENV GO111MODULE=on

# Download and install OpenCV
WORKDIR /root
RUN wget -O opencv.zip https://github.com/opencv/opencv/archive/4.5.1.zip && \
    wget -O opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/4.5.1.zip && \
    unzip opencv.zip && \
    unzip opencv_contrib.zip && \
    mkdir -p opencv-4.5.1/build && \
    cd opencv-4.5.1/build && \
    cmake -D CMAKE_BUILD_TYPE=RELEASE \
          -D CMAKE_INSTALL_PREFIX=/usr/local \
          -D OPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-4.5.1/modules \
          -D WITH_CUDA=OFF \
          -D BUILD_EXAMPLES=OFF .. && \
    make -j$(nproc) && \
    make install && \
    ldconfig

# Install GoCV
WORKDIR /go/src/gocv.io/x/gocv
RUN make install

# Set up your Go workspace
WORKDIR /go/src/app

# Copy the rest of the application code
COPY . .

# Expose port
EXPOSE 8080

# Run the Go app when the container launches
CMD ["air", "-c", ".air.toml"]