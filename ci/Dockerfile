FROM ubuntu:22.04

# 環境変数
ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="/usr/local/go/bin:$PATH"

# 必要なパッケージと LLVM リポジトリ追加
RUN apt-get update && \
    apt-get install -y wget gnupg software-properties-common && \
    wget https://apt.llvm.org/llvm-snapshot.gpg.key && \
    apt-key add llvm-snapshot.gpg.key && \
    add-apt-repository "deb http://apt.llvm.org/jammy/ llvm-toolchain-jammy-16 main" && \
    apt-get update && \
    apt-get install -y \
    curl \
    git \
    build-essential \
    llvm-16 \
    llvm-16-dev \
    llvm-16-runtime \
    llvm-16-tools \
    clang-16 \
    lldb-16 \
    lld-16 \
    vim-common 

RUN update-alternatives --install /usr/bin/clang clang /usr/bin/clang-16 100 && \
    update-alternatives --install /usr/bin/llvm-config llvm-config /usr/bin/llvm-config-16 100 && \
    update-alternatives --install /usr/bin/lli lli /usr/bin/lli-16 100 && \
    rm -rf /var/lib/apt/lists/*

# Go のインストール（arm64）
RUN wget https://go.dev/dl/go1.23.2.linux-arm64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-arm64.tar.gz && \
    rm go1.23.2.linux-arm64.tar.gz

# Go と LLVM の確認
RUN go version
RUN llvm-config --version
RUN lli --version

# 作業ディレクトリ
WORKDIR /work

