#pragma once

#include <cstdint>
#include <fstream>
#include <string>

enum class WalOp : uint8_t {
    Insert = 1,
    Delete = 2,
};

struct WalRecordHeader {
    uint8_t op;
    uint32_t length;
};

class WalWriter {
public:
    explicit WalWriter(const std::string& path);
    ~WalWriter();

    void appendInsert(uint64_t id, const float* vector, int dim);
    void appendDelete(uint64_t id);
    void flush();
    void reopen();

private:
    std::string path_;
    std::ofstream stream_;
};
