#include "payload.hpp"
#include <cstring>

static void appendPrimitive(std::vector<uint8_t>& buffer, const void* value, size_t size) {
    const auto* bytes = reinterpret_cast<const uint8_t*>(value);
    buffer.insert(buffer.end(), bytes, bytes + size);
}

std::vector<uint8_t> serializeField(const Field& field) {
    std::vector<uint8_t> buffer;
    uint8_t type = static_cast<uint8_t>(field.type);
    uint16_t keySize = static_cast<uint16_t>(field.key.size());
    uint32_t valueSize = static_cast<uint32_t>(field.value.size());

    appendPrimitive(buffer, &type, sizeof(type));
    appendPrimitive(buffer, &keySize, sizeof(keySize));
    buffer.insert(buffer.end(), field.key.begin(), field.key.end());
    appendPrimitive(buffer, &valueSize, sizeof(valueSize));
    buffer.insert(buffer.end(), field.value.begin(), field.value.end());
    return buffer;
}

std::vector<uint8_t> serializePayload(const PayloadDocument& document) {
    std::vector<uint8_t> buffer;
    appendPrimitive(buffer, &document.id, sizeof(document.id));
    uint16_t fieldCount = static_cast<uint16_t>(document.fields.size());
    appendPrimitive(buffer, &fieldCount, sizeof(fieldCount));

    for (const auto& field : document.fields) {
        auto fieldBytes = serializeField(field);
        buffer.insert(buffer.end(), fieldBytes.begin(), fieldBytes.end());
    }
    return buffer;
}
