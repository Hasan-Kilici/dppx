struct Field {
    ValueType type;

    std::string key;

    std::vector<uint8_t> value;
};

struct Document {
    std::vector<Field> fields;
};