#include "wal.hpp"
#include <filesystem>

WalWriter::WalWriter(const std::string& path)
    : path_(path) {
    std::filesystem::create_directories(std::filesystem::path(path).parent_path());
    stream_.open(path_, std::ios::binary | std::ios::app);
}

WalWriter::~WalWriter() {
    flush();
    if (stream_.is_open()) {
        stream_.close();
    }
}

void WalWriter::appendInsert(uint64_t id, const float* vector, int dim) {
    if (!stream_.is_open() || !vector) {
        return;
    }

    WalRecordHeader header;
    header.op = static_cast<uint8_t>(WalOp::Insert);
    header.length = sizeof(id) + sizeof(float) * dim;

    stream_.write(reinterpret_cast<const char*>(&header), sizeof(header));
    stream_.write(reinterpret_cast<const char*>(&id), sizeof(id));
    stream_.write(reinterpret_cast<const char*>(vector), sizeof(float) * dim);
}

void WalWriter::appendDelete(uint64_t id) {
    if (!stream_.is_open()) {
        return;
    }

    WalRecordHeader header;
    header.op = static_cast<uint8_t>(WalOp::Delete);
    header.length = sizeof(id);

    stream_.write(reinterpret_cast<const char*>(&header), sizeof(header));
    stream_.write(reinterpret_cast<const char*>(&id), sizeof(id));
}

void WalWriter::flush() {
    if (stream_.is_open()) {
        stream_.flush();
    }
}

void WalWriter::reopen() {
    if (stream_.is_open()) {
        stream_.close();
    }

    stream_.open(path_, std::ios::binary | std::ios::app);
}
