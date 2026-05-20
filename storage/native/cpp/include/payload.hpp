#pragma once

#include <cstdint>
#include <string>
#include <vector>

enum class FieldType : uint8_t {
    Int32 = 1,
    Int64 = 2,
    Float32 = 3,
    Float64 = 4,
    Bool = 5,
    String = 6,
    Timestamp = 7,
    Vector = 8,
};

struct Field {
    FieldType type;
    std::string key;
    std::vector<uint8_t> value;
};

struct PayloadDocument {
    uint64_t id;
    std::vector<Field> fields;
};

std::vector<uint8_t> serializeField(const Field& field);
std::vector<uint8_t> serializePayload(const PayloadDocument& document);
